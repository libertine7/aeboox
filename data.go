package main

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
)

func AddEventsV1(events []*EbooxBookEvent) {
	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: "http://10.133.19.227:8086"})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{Database: "aaa", Precision: "us"})

	for _, event := range events {
		tags := map[string]string{
			"EventType":   event.EventType,
			"Screen":      event.Screen,
			"UserAgent":   event.UserAgent,
			"Platform":    event.Platform,
			"OsVersion":   event.OsVersion,
			"AppVersion":  event.AppVersion,
			"DeviceModel": event.DeviceModel,
			"BookName":    event.BookName,
			"BookFormat":  event.BookFormat,
			"City":        event.City,
			"Subdivision": event.Subdivision,
			"Country":     event.Country,
			"TimeZone":    event.TimeZone,
		}

		fields := map[string]interface{}{
			"DeviceId":     event.DeviceId,
			"EventPayload": event.EventPayload,
			"EventTime":    event.EventTime,
		}

		pt, err := client.NewPoint("aevents", tags, fields)
		if err != nil {
			fmt.Println("Error:", err.Error())
			continue
		}
		bp.AddPoint(pt)
	}

	err = c.Write(bp)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
