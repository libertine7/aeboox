package main

import (
	"encoding/json"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

type GeoAddress struct {
	City        string
	Subdivision string
	Country     string
	TimeZone    string
}

func (geo *GeoAddress) Json() ([]byte, error) {
	return json.Marshal(geo)
}

func (geo *GeoAddress) FromJson(jsonStr []byte) error {
	return json.Unmarshal(jsonStr, geo)
}

func IpToGeo(ipstr string) (result GeoAddress) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		log.Println("IP address wrong")
		return GeoAddress{}
	}

	cached, err := Cache.Get(ipstr)
	if err == nil {
		result.FromJson(cached)
		return result
	}

	//wget http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
	// and extract it
	//db, err := geoip2.Open("GeoIP2-City.mmdb")
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Println(err)
		return GeoAddress{}
	}
	defer db.Close()

	record, err := db.City(ip)
	if err != nil {
		log.Println(err)
		return GeoAddress{}
	}

	if len(record.City.Names) == 0 || len(record.Subdivisions) == 0 {
		return GeoAddress{}
	}

	result = GeoAddress{City: record.City.Names["en"], Subdivision: record.Subdivisions[0].IsoCode, Country: record.Country.IsoCode, TimeZone: record.Location.TimeZone}
	buf, _ := result.Json()

	Cache.Set(ipstr, buf)

	return result
}

