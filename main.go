package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"strconv"
)

type Centers struct {
	Centers []struct {
		CenterID     int    `json:"center_id"`
		Name         string `json:"name"`
		Address      string `json:"address"`
		StateName    string `json:"state_name"`
		DistrictName string `json:"district_name"`
		BlockName    string `json:"block_name"`
		Pincode      int    `json:"pincode"`
		Lat          int    `json:"lat"`
		Long         int    `json:"long"`
		From         string `json:"from"`
		To           string `json:"to"`
		FeeType      string `json:"fee_type"`
		Sessions     []Session `json:"sessions"`
		VaccineFees []struct {
			Vaccine string `json:"vaccine"`
			Fee     string `json:"fee"`
		} `json:"vaccine_fees,omitempty"`
	} `json:"centers"`
}

type Session struct {
	SessionID              string   `json:"session_id"`
	Date                   string   `json:"date"`
	AvailableCapacity      int      `json:"available_capacity"`
	MinAgeLimit            int      `json:"min_age_limit"`
	Vaccine                string   `json:"vaccine"`
	Slots                  []string `json:"slots"`
	AvailableCapacityDose1 int      `json:"available_capacity_dose1"`
	AvailableCapacityDose2 int      `json:"available_capacity_dose2"`
} 
var vaccine = "COVAXIN"
func main() {
	var dateParam string
	var url string
	spaceClient := http.Client{
		Timeout: time.Second * 1, // Timeout after 2 seconds
	}
	for{
	dateParam = time.Now().Format("02-01-2006")
	url ="https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?district_id=294&date="+dateParam
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	centers := Centers{}
	jsonErr := json.Unmarshal(body, &centers)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	for _, elem := range centers.Centers {
		for _, sess := range elem.Sessions{
			if sess.IsSlotAvaliable() {
				fmt.Println("CenterName: "+elem.Name+", Pincode: "+strconv.Itoa(elem.Pincode)+", Vaccine: "+sess.Vaccine+", Date: "+sess.Date+", dose1: "+strconv.Itoa(sess.AvailableCapacityDose1)+", dose2: "+strconv.Itoa(sess.AvailableCapacityDose2))
			}else{
				fmt.Println("Retrying....")
			}
		}
    }
	time.Sleep(1 * time.Second)
}
	
}

func (sess Session) IsSlotAvaliable() bool { 
	if sess.AvailableCapacity>0 && sess.MinAgeLimit<45 && vaccine == sess.Vaccine{
		return true
	}
	return false
}
