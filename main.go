package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MarineInstituteResponse struct {
	Table struct {
		Rows [][]interface{} `json:"rows"`
	} `json:"table"`
}

type TideData struct {
	Time time.Time
	Tide float64
}

func main() {
	var today = time.Now()
	var url = "https://erddap.marine.ie/erddap/tabledap/IMI-TidePrediction_epa.json?time,longitude,latitude,stationID,sea_surface_height&time%3E=" + today.Format("2006-01-02T15:04:05Z") + "&time%3C=" + today.AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z") + "&stationID=%22IEEABWC140_0000_0300_MODELLED%22"
	//  request, err := http.Get("https://erddap.marine.ie/erddap/tabledap/IMI-TidePrediction_epa.json?time,longitude,latitude,stationID,sea_surface_height&time%3E=2023-06-22T00:00:00Z&time%3C=2023-06-22T23:59:00Z&stationID=%22IESEBWC010_0000_0100_MODELLED%22")
	request, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer request.Body.Close()

	b, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	var response MarineInstituteResponse
	err = json.Unmarshal(b, &response)
	if err != nil {
		panic(err)
	}

	var data []TideData
	for _, row := range response.Table.Rows {
		var tideData TideData
		t := row[0].(string)
		tideData.Time, _ = time.Parse(time.RFC3339, t)
		tideData.Tide = row[4].(float64)
		data = append(data, tideData)
	}

	tides := MinMaxTide(data)

	for _, d := range tides {
		fmt.Println(d.Time, d.Tide)
	}

}

func MinMaxTide(data []TideData) []TideData {
	var ret []TideData
	down := data[1].Tide < data[0].Tide

	for i := 1; i < len(data); i++ {
		currentValue := data[i].Tide
		previousValue := data[i-1].Tide

		if down && currentValue > previousValue {
			ret = append(ret, data[i-1])
			down = false
		} else if !down && currentValue < previousValue {
			ret = append(ret, data[i-1])
			down = true
		}
	}

	return ret
}
