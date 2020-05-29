package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Holiday struct {
	Date        string      `json:"date"`
	LocalName   string      `json:"localName"`
	Name        string      `json:"name"`
	CountryCode string      `json:"countryCode"`
	Fixed       bool        `json:"fixed"`
	Global      bool        `json:"global"`
	Counties    string      `json:"counties"`
	LaunchYear  string      `json:"launchYear"`
	Type        string      `json:"type"`
}

func Weekend(date time.Time) bool {
	if date.Weekday().String() == "Sunday" || date.Weekday().String() == "Saturday" {
		return true
	}
	return false
}

func bigHoliday(day time.Time) (time.Time, bool) {
	big := false
	January := "2020-01-07"
	jan , _ := time.Parse("2006-01-02", January)

	August := "2020-08-24"
	aug , _ := time.Parse("2006-01-02", August)

	var firstDay time.Time
	if day.Equal(jan){
		firstDay = day.AddDate(0,0,-3)
		big = true
	} else if day.Equal(aug){
		firstDay = day.AddDate(0,0, -2)
		big = true
	}
	return firstDay, big
}


func main() {
	var data []Holiday
	day := time.Now()
	url := "https://date.nager.at/api/v2/publicholidays/2020/UA"

	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	getHoliday := body

	error := json.Unmarshal(getHoliday, &data)
	if error != nil {
		panic(err.Error())
	}

	for _, value := range data {
		holiday, err := time.Parse("2006-01-02", value.Date)
		if err != nil {
			panic(err.Error())
		}

		if holiday.After(day){
			firstDay := holiday.AddDate(0, 0, 1)
			var  adjacentFirst bool
			var  adjacentSecond bool
			var lastDay time.Time
			if Weekend(firstDay) {
				lastDay = firstDay.AddDate(0, 0, 1)
				adjacentFirst = true
			} else {
				firstDay, adjacentSecond = bigHoliday(holiday)
			}

			if  adjacentFirst {
				days := (lastDay.Sub(holiday).Hours() / 24) + 1
				fmt.Printf(
					"The next holiday is '%v', %v %v, and the weekend will %v days: %v %v - %v %v.",
					value.Name,
					holiday.Month(),
					holiday.Day(),
					days,
					holiday.Month(),
					holiday.Day(),
					lastDay.Month(),
					lastDay.Day())
			} else if  adjacentSecond{
				days := (holiday.Sub(firstDay).Hours() / 24) + 1
				fmt.Printf(
					"The next holiday is '%v', %v %v, and the weekend will %v days: %v %v - %v %v.",
					value.Name,
					holiday.Month(),
					holiday.Day(),
					days,
					firstDay.Month(),
					firstDay.Day(),
					holiday.Month(),
					holiday.Day())
			} else {
				fmt.Printf("The next holiday is '%v' on %v %v", value.Name, holiday.Month(), holiday.Day())
			}
			break
		}
	}

}

