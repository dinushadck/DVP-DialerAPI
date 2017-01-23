package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetAppoinmentsForSchedule(internalAuthToken, schedulrId string) Schedule {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetAppoinmentsForSchedule", r)
		}
	}()
	fmt.Println("Start Get Schedule Schedule service")
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	client := &http.Client{}
	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/LimitAPI/Schedule/%s/Appointments/Info", CreateHost(scheduleServiceHost, scheduleServicePort), schedulrId)
	fmt.Println("request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Schedulr API response:: ", string(response))
	var apiResult ScheduleDetails
	json.Unmarshal(response, &apiResult)

	if len(apiResult.Result) > 0 {
		fmt.Println("Schedulr apiResult.Result: ", apiResult.Result[0])
		return apiResult.Result[0]
	} else {
		return Schedule{}
	}
}

func GetTimeZoneFroSchedule(internalAuthToken, schedulrId string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetTimeZoneFroSchedule", r)
		}
	}()
	fmt.Println("Start Get Schedule Schedule service")
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	client := &http.Client{}
	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/LimitAPI/Schedule/%s", CreateHost(scheduleServiceHost, scheduleServicePort), schedulrId)
	fmt.Println("request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var apiResult ScheduleDetails
	json.Unmarshal(response, &apiResult)

	if len(apiResult.Result) > 0 {
		fmt.Println("Schedulr apiResult.Result: ", apiResult.Result[0])
		return apiResult.Result[0].TimeZone
	} else {
		return ""
	}
}

func CheckAppoinments(appoinments []Appoinment, timeNow time.Time, timeZone string) (appinment Appoinment, endTime time.Time) {
	location, _ := time.LoadLocation(timeZone)

	for _, appmnt := range appoinments {
		fmt.Println("CheckAppoinments: ", appmnt.AppointmentName)
		fmt.Println("daysOfWeek: ", appmnt.DaysOfWeek)
		daysOfWeek := strings.Split(appmnt.DaysOfWeek, ",")
		if stringInSlice(timeNow.Weekday().String(), daysOfWeek) {
			fmt.Println("match daysOfWeek: ", timeNow.Weekday().String())

			tempstartDate, _ := time.Parse(layout2, appmnt.StartDate)
			tempendDate, _ := time.Parse(layout2, appmnt.EndDate)

			splitStartTime := strings.Split(appmnt.StartTime, ":")
			splitEndTime := strings.Split(appmnt.EndTime, ":")

			startDate := time.Date(tempstartDate.Year(), tempstartDate.Month(), tempstartDate.Day(), 0, 0, 0, 0, location)
			endDate := time.Date(tempendDate.Year(), tempendDate.Month(), tempendDate.Day(), 0, 0, 0, 0, location)

			fmt.Println("appoinment startDate: ", startDate.String())
			fmt.Println("appoinment endDate: ", endDate.String())

			if startDate.Before(timeNow) && endDate.After(timeNow) {
				tempStartHr := 0
				tempStartMin := 0
				tempEndHr := 0
				tempEndMin := 0

				if len(splitStartTime) == 2 {
					tempStartHr, _ = strconv.Atoi(splitStartTime[0])
					tempStartMin, _ = strconv.Atoi(splitStartTime[1])
				}

				if len(splitEndTime) == 2 {
					tempEndHr, _ = strconv.Atoi(splitEndTime[0])
					tempEndMin, _ = strconv.Atoi(splitEndTime[1])
				}

				localStartTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), tempStartHr, tempStartMin, 0, 0, location)
				localEndTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), tempEndHr, tempEndMin, 0, 0, location)

				fmt.Println("serverTimeLocal: ", timeNow.String())
				fmt.Println("appoinment startTime: ", localStartTime.String())
				fmt.Println("appoinment enendTimedDate: ", localEndTime.String())

				if localStartTime.Before(timeNow) && localEndTime.After(timeNow) {
					fmt.Println("match appoinment date&time: ", timeNow.String())

					endTime = localEndTime
					appinment = appmnt
					return
				}
			}
		}
	}

	return
}

func CheckAppoinmentForCampaign(internalAuthToken, schedulrId string) (appointment Appoinment, timeZone string, campaignEndTime time.Time) {
	schedule := GetAppoinmentsForSchedule(internalAuthToken, schedulrId)
	location, _ := time.LoadLocation(schedule.TimeZone)
	timeNow := time.Now().In(location)
	appointment, campaignEndTime = CheckAppoinments(schedule.Appointment, timeNow, schedule.TimeZone)
	timeZone = schedule.TimeZone

	fmt.Println("CheckAppoinmentForCampaign::: appointment :: ", appointment)
	fmt.Println("CheckAppoinmentForCampaign::: campaignEndTime :: ", campaignEndTime)
	fmt.Println("CheckAppoinmentForCampaign::: timeZone :: ", timeZone)
	return
}

func CheckAppoinmentForCallback(company, tenant int, schedulrId string, timeToCheck time.Time, timeZone string) bool {
	defaultAppoinment := Appoinment{}
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	schedule := GetAppoinmentsForSchedule(internalAuthToken, schedulrId)
	machingAppoinment, _ := CheckAppoinments(schedule.Appointment, timeToCheck, timeZone)
	if machingAppoinment == defaultAppoinment {
		return false
	} else {
		return true
	}
}
