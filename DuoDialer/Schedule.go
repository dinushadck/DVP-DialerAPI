package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetAppoinmentsForSchedule(authToken, schedulrId string) []Appoinment {
	client := &http.Client{}
	request := fmt.Sprintf("%s/Schedule/%s/Appointment", callRuleService, schedulrId)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var apiResult ScheduleDetails
	json.Unmarshal(response, &apiResult)

	return apiResult.Result
}
