/*
 Copyright (C) 2023 Michail Angelos Tsiantakis

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"embed"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Weather struct {
	LastWeatherUpdate time.Time
	LastPMUpdate      time.Time
	Temperature       float64
	Humidity          float64
	PM25              float64
	PM10              float64
}

type ResponseData struct {
	Background string
	Width      float64
	Value      float64
}

var (
	//go:embed all:app
	appDir embed.FS

	weatherStorage Weather
)

func main() {

	fs := http.FileServer(http.FS(appDir))
	http.Handle("/", fs)
	http.HandleFunc("/weatherData", handleWeather)

	log.Fatalln(http.ListenAndServe(":3000", nil))
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	updateWeatherData()
	responseData := ResponseData{}
	if weatherStorage.Temperature < 15 {
		responseData.Background = "bg-info"
	} else if weatherStorage.Temperature < 30 {
		responseData.Background = "bg-success"
	} else {
		responseData.Background = "bg-danger"
	}
	responseData.Value = weatherStorage.Temperature
	responseData.Width = (weatherStorage.Temperature / 50) * 100
	response := fmt.Sprintf("<div id='temp-bar' class='progress-bar %s'  style='width: %f%%;'>%.2f°C</div>", responseData.Background, responseData.Width, responseData.Value)
	w.Write([]byte(response))

	//handle time last updated
	if weatherStorage.LastWeatherUpdate.Add(60 * time.Minute).After(time.Now()) {
		w.Write([]byte(fmt.Sprintf("<p id='last-updated' hx-swap-oob='true' class='mx-1 px-1 rounded-4 text-bg-success fs-6 fst-italic'>Τελευταία μέτρηση μετεωρολογικών δεδομένων: %s</p>", weatherStorage.LastWeatherUpdate.Format("02.01.2006 15:04"))))
	} else {
		w.Write([]byte(fmt.Sprintf("<p id='last-updated' hx-swap-oob='true' class='mx-1 px-1 rounded-4 text-bg-danger fs-6 fst-italic'>Τελευταία μέτρηση μετεωρολογικών δεδομένων: %s</p>", weatherStorage.LastWeatherUpdate.Format("02.01.2006 15:04"))))

	}
	if weatherStorage.LastPMUpdate.Add(60 * time.Minute).After(time.Now()) {
		w.Write([]byte(fmt.Sprintf("<p id='last-updated-pm' hx-swap-oob='true' class='mx-1 px-1 rounded-4 text-bg-success fs-6 fst-italic'>Τελευταία μέτρηση αέριων σωματιδίων: %s</p>", weatherStorage.LastPMUpdate.Format("02.01.2006 15:04"))))
	} else {
		w.Write([]byte(fmt.Sprintf("<p id='last-updated-pm' hx-swap-oob='true' class='mx-1 px-1 rounded-4 text-bg-danger fs-6 fst-italic'>Τελευταία μέτρηση αέριων σωματιδίων: %s</p>", weatherStorage.LastPMUpdate.Format("02.01.2006 15:04"))))
	}

	// handle humidity
	if weatherStorage.Humidity < 30 {
		responseData.Background = "bg-info"
	} else if weatherStorage.Humidity < 60 {
		responseData.Background = "bg-success"
	} else {
		responseData.Background = "bg-danger"
	}
	responseData.Value = weatherStorage.Humidity
	responseData.Width = (weatherStorage.Humidity / 100) * 100
	response = fmt.Sprintf("<div id='hum-bar' class='progress-bar %s'  style='width: %f%%;'>%.2f%%</div>", responseData.Background, responseData.Width, responseData.Value)
	w.Write([]byte(response))

	// handle pm25
	if weatherStorage.PM25 < 15 {
		responseData.Background = "bg-info"
	} else if weatherStorage.PM25 < 30 {
		responseData.Background = "bg-success"
	} else if weatherStorage.PM25 < 55 {
		responseData.Background = "bg-warning-subtle"
	} else if weatherStorage.PM25 < 110 {
		responseData.Background = "bg-warning"
	} else {
		responseData.Background = "bg-danger"
	}
	responseData.Value = weatherStorage.PM25
	responseData.Width = (weatherStorage.PM25 / 110) * 100
	response = fmt.Sprintf("<div id='pm25-bar' class='progress-bar %s'  style='width: %f%%;'>%.1f</div>", responseData.Background, responseData.Width, responseData.Value)
	w.Write([]byte(response))

	//handle pm10
	if weatherStorage.PM10 < 25 {
		responseData.Background = "bg-info"
	} else if weatherStorage.PM10 < 50 {
		responseData.Background = "bg-success"
	} else if weatherStorage.PM10 < 90 {
		responseData.Background = "bg-warning-subtle"
	} else if weatherStorage.PM10 < 180 {
		responseData.Background = "bg-warning"
	} else {
		responseData.Background = "bg-danger"
	}
	responseData.Value = weatherStorage.PM10
	responseData.Width = (weatherStorage.PM10 / 240) * 100
	response = fmt.Sprintf("<div id='pm10-bar' class='progress-bar %s'  style='width: %f%%;'>%.1f</div>", responseData.Background, responseData.Width, responseData.Value)
	w.Write([]byte(response))

}
func handleHum(w http.ResponseWriter, r *http.Request) {
	updateWeatherData()
	w.WriteHeader(http.StatusOK)
	responseData := ResponseData{}
	if weatherStorage.Humidity < 30 {
		responseData.Background = "bg-info"
	} else if weatherStorage.Humidity < 60 {
		responseData.Background = "bg-success"
	} else {
		responseData.Background = "bg-danger"
	}
	responseData.Value = weatherStorage.Humidity
	responseData.Width = (weatherStorage.Humidity / 100) * 100
	response := fmt.Sprintf("<div id='hum-bar' class='progress-bar %s'  style='width: %f%%;'>%.2f%%</div>", responseData.Background, responseData.Width, responseData.Value)
	w.Write([]byte(response))
}

func handlePM25(w http.ResponseWriter, r *http.Request) {
	updateWeatherData()
	w.WriteHeader(http.StatusOK)
	responseData := ResponseData{}
	if weatherStorage.PM25 < 15 {
		responseData.Background = "bg-info"
	} else if weatherStorage.PM25 < 30 {
		responseData.Background = "bg-success"
	} else if weatherStorage.PM25 < 55 {
		responseData.Background = "bg-warning-subtle"
	} else if weatherStorage.PM25 < 110 {
		responseData.Background = "bg-warning"
	} else {
		responseData.Background = "bg-danger"
	}
	responseData.Value = weatherStorage.PM25
	responseData.Width = (weatherStorage.PM25 / 110) * 100
	response := fmt.Sprintf("<div id='pm-25' class='progress-bar %s'  style='width: %f%%;'>%.1f</div>", responseData.Background, responseData.Width, responseData.Value)
	w.Write([]byte(response))
}

func handlePM10(w http.ResponseWriter, r *http.Request) {
	updateWeatherData()
	w.WriteHeader(http.StatusOK)
	responseData := ResponseData{}
	if weatherStorage.PM10 < 25 {
		responseData.Background = "bg-info"
	} else if weatherStorage.PM10 < 50 {
		responseData.Background = "bg-success"
	} else if weatherStorage.PM10 < 90 {
		responseData.Background = "bg-warning-subtle"
	} else if weatherStorage.PM10 < 180 {
		responseData.Background = "bg-warning"
	} else {
		responseData.Background = "bg-danger"
	}
	responseData.Value = weatherStorage.PM10
	responseData.Width = (weatherStorage.PM10 / 240) * 100
	response := fmt.Sprintf("<div id='pm-10' class='progress-bar %s'  style='width: %f%%;'>%.1f</div>", responseData.Background, responseData.Width, responseData.Value)
	w.Write([]byte(response))
}

func updateWeatherData() {
	if weatherStorage.LastWeatherUpdate.Add(20 * time.Minute).After(time.Now()) {
		log.Println("Weather data is up to date!")
		return
	} else {
		log.Println("Weather data is outdated! Updating...")
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	csvFile, err := client.Get("https://writethatdown.site/greendigital/meteo48h.csv.txt")
	if err != nil {
		log.Println("Cannot get weather data csv!")
	}
	defer csvFile.Body.Close()

	r := csv.NewReader(csvFile.Body)

	records, _ := r.ReadAll()

	record := records[len(records)-1]

	parsedTime, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		log.Println("Cannot parse time!")
		parsedTime = -999
	}
	unixTime := time.Unix(parsedTime, 0)

	temp, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		log.Println("Cannot parse temperature!")
		temp = -999
	}

	hum, err := strconv.ParseFloat(record[6], 64)
	if err != nil {
		log.Println("Cannot parse humidity!")
		hum = -999
	}

	pmCsvFile, err := client.Get("https://writethatdown.site/greendigital/sen55.csv")
	if err != nil {
		log.Println("Cannot get sen55 data!")
	}
	defer pmCsvFile.Body.Close()

	pmR := csv.NewReader(pmCsvFile.Body)

	pmRecords, _ := pmR.ReadAll()

	pmRecord := pmRecords[len(pmRecords)-1]

	parsedPMTime, err := strconv.ParseInt(pmRecord[0], 10, 64)
	if err != nil {
		log.Println("Cannot parse time!")
		parsedPMTime = -999
	}
	unixPMTime := time.Unix(parsedPMTime, 0)

	pm25, err := strconv.ParseFloat(pmRecord[2], 64)
	if err != nil {
		log.Println("Cannot parse bar!")
		pm25 = -999
	}

	pm10, err := strconv.ParseFloat(pmRecord[4], 64)
	if err != nil {
		log.Println("Cannot parse bar!")
		pm10 = -999
	}

	weather := Weather{
		LastWeatherUpdate: unixTime,
		LastPMUpdate:      unixPMTime,
		Temperature:       temp,
		Humidity:          hum,
		PM25:              pm25,
		PM10:              pm10,
	}

	weatherStorage = weather
}
