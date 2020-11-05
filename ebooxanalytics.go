package main

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"strings"
	"time"
)

type EbooxAnalytics struct {
}

type EbooxBookEvent struct {
	DeviceId     string `valid:"required"`
	EventType    string `valid:"required"`
	Screen       string `valid:"required"`
	UserAgent    string `valid:"-"`
	Platform     string `valid:"required"`
	OsVersion    string `valid:"required"`
	AppVersion   string `valid:"required"`
	DeviceModel  string `valid:"required"`
	EventPayload string `valid:"-"`
	EventTime    string `valid:"required"`
	// book
	BookName   string `valid:"-"`
	BookFormat string `valid:"-"`
	// geo data from ip
	City        string `valid:"-"`
	Subdivision string `valid:"-"`
	Country     string `valid:"-"`
	TimeZone    string `valid:"-"`
}

func (c *EbooxAnalytics) AddEventV1(r *http.Request, req *[]*EbooxBookEvent, res *bool) error {
	userIp := strings.Split(r.RemoteAddr, ":")[0]
	// if ngnix use
	if val, ok := r.Header["X-Real-Ip"]; ok {
		userIp = val[0]
	}

	geo := IpToGeo(userIp)
	items := *req
	for _, item := range items {
		item.City = geo.City
		item.Subdivision = geo.Subdivision
		item.Country = geo.Country
		item.TimeZone = geo.TimeZone
		_, err := govalidator.ValidateStruct(item)
		if err != nil {
			*res = false
			return errors.New("data validation error")
		}
	}

	for _, aaa := range items {
		fmt.Println(fmt.Sprintf("* %v", aaa))
	}

	time.Sleep(10 * time.Millisecond) // write emulation
	AddEventsV1(items)                // write to db

	*res = true
	return nil
}
