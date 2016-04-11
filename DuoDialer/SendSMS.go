package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func SaveShortMessageInformation(sms Sms) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SaveShortMessageInformation", r)
		}
	}()
	smsByte, err := json.Marshal(sms)
	if err != nil {
		fmt.Println(err)
		return
	}
	text := string(smsByte)
	fmt.Println(text)
	//upload to campaign service
	serviceurl := fmt.Sprintf("http://%s/DuoMessageTemplate/MesssageDispatcherService.svc/Json/saveShortMessageInformation", casServerHost)
	//jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	//internalAuthToken := fmt.Sprintf("%s:%s", sessionInfo["TenantId"], sessionInfo["CompanyId"])

	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(smsByte))
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("authorization", jwtToken)
	//req.Header.Set("companyinfo", internalAuthToken)
	fmt.Println("request:", serviceurl)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	body, errb := ioutil.ReadAll(resp.Body)
	//if success remove hashInfo
	if errb != nil {
		fmt.Println(err.Error())
	} else {
		result := string(body)
		fmt.Println("response Body:", result)
	}
}

func GenerateSMS(fromNumber, message, phoneNumber string) Sms {
	numbers := make([]string, 0)
	smsInfoData := SmsInfo{}
	smsData := Sms{}

	numbers = append(numbers, phoneNumber)
	smsInfoData.Date = time.Now().String()
	smsInfoData.FromPhoneNumber = fromNumber
	smsInfoData.GatewayName = 1
	smsInfoData.MessageContent = message
	smsInfoData.PhoneNumbers = numbers

	smsData.shortMessageInfo = smsInfoData
	smsData.securityToken = v5_1SecurityToken

	return smsData
}

func SendSms(company, tenant int, resourceServer ResourceServerInfo, campaignId, camClass, camType, camCategory, fromNumber, message, phoneNumber string) {
	IncrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, 240, camClass, camType, camCategory, "1", campaignId, phoneNumber, phoneNumber, "start", "start", time.Now().UTC().Format(layout4), resourceServer.ResourceServerId)
	smsRequest := GenerateSMS(fromNumber, message, phoneNumber)
	SaveShortMessageInformation(smsRequest)

	DecrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)

	SetSessionInfo(campaignId, phoneNumber, "Reason", "dial_success")
	SetSessionInfo(campaignId, phoneNumber, "DialerStatus", "connected")
	go UploadSessionInfo(campaignId, phoneNumber)
}
