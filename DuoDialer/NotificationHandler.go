package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

func SendNotificationToRoom(roomName, from, direction, message, ref string, company, tenant int) {
	defer func() {
		if r := recover(); r != nil {
			color.Red(fmt.Sprintf("Recovered in UpdateCampaignStartStatus %+v", r))
		}
	}()
	pushD := PushDataRoom{}
	pushD.From = from
	pushD.Direction = direction
	pushD.message = message
	pushD.Ref = ref

	jsonData, _ := json.Marshal(pushD)

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/NotificationService/Notification/initiate/%s",  CreateHost(notificationServiceHost, notificationServicePort), roomName)
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	req.Header.Set("eventname", ref)
	//DialerLog(fmt.Sprintf("request:%s", serviceurl))
	color.Yellow(fmt.Sprintf("NOTIFICATION SENT - URL : %s, DATA : %v", serviceurl, pushD))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	DialerLog(fmt.Sprintf("response Status:%s", resp.Status))
	DialerLog(fmt.Sprintf("response Headers:%s", resp.Header))
	body, errb := ioutil.ReadAll(resp.Body)
	if errb != nil {
		color.Red(err.Error())
	} else {
		result := string(body)
		DialerLog(fmt.Sprintf("response Body:%s", result))
	}
}