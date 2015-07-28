package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetAppoinmentsForSchedule(authToken, schedulrId string) []Appoinment {
	fmt.Println("Start Get Schedule Schedule service")
	client := &http.Client{}
	request := fmt.Sprintf("%s/Schedule/%s/Appointments", scheduleService, schedulrId)
	fmt.Println("request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var apiResult ScheduleDetails
	json.Unmarshal(response, &apiResult)
	fmt.Println("Schedulr apiResult.Result: ", apiResult.Result)
	return apiResult.Result
}

func CheckAppoinments(appoinments []Appoinment) Appoinment {
	for _, appmnt := range appoinments {
		fmt.Println("CheckAppoinments: ", appmnt.AppointmentName)
		timeNow := time.Now().UTC()
		fmt.Println("daysOfWeek: ", appmnt.DaysOfWeek)
		daysOfWeek := strings.Split(appmnt.DaysOfWeek, ",")
		if stringInSlice(timeNow.Weekday().String(), daysOfWeek) {
			fmt.Println("match daysOfWeek: ", timeNow.Weekday().String())
			tempstartDate, _ := time.Parse(layout2, appmnt.StartDate)
			tempendDate, _ := time.Parse(layout2, appmnt.EndDate)

			startDate := time.Date(tempstartDate.Year(), tempstartDate.Month(), tempstartDate.Day(), tempstartDate.Hour(), tempstartDate.Minute(), tempstartDate.Second(), 0, time.UTC)
			endDate := time.Date(tempendDate.Year(), tempendDate.Month(), tempendDate.Day(), tempendDate.Hour(), tempendDate.Minute(), tempendDate.Second(), 0, time.UTC)

			fmt.Println("appoinment startDate: ", startDate.String())
			fmt.Println("appoinment endDate: ", endDate.String())

			if startDate.Before(timeNow) && endDate.After(timeNow) {
				startTime, _ := time.Parse(layout1, appmnt.StartTime)
				endTime, _ := time.Parse(layout1, appmnt.EndTime)

				localStartTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), 0, time.UTC)
				localEndTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), 0, time.UTC)

				fmt.Println("serverTimeUTC: ", timeNow.String())
				fmt.Println("appoinment startTime: ", localStartTime.String())
				fmt.Println("appoinment enendTimedDate: ", localEndTime.String())

				if localStartTime.Before(timeNow) && localEndTime.After(timeNow) {
					fmt.Println("match appoinment date&time: ", timeNow.String())
					appmnt.StartTime = localStartTime.String()
					appmnt.EndTime = localEndTime.String()

					return appmnt
				}
			}
		}
	}

	return Appoinment{}
}

func CheckAppoinmentForCampaign(authToken, schedulrId string) Appoinment {
	appionments := GetAppoinmentsForSchedule(authToken, schedulrId)
	return CheckAppoinments(appionments)
}
