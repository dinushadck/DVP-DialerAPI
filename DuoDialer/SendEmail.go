package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func SaveEmailInformation(_email Email) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SaveShortMessageInformation", r)
		}
	}()
	emailByte, err := json.Marshal(_email)
	if err != nil {
		fmt.Println(err)
		return
	}
	text := string(emailByte)
	fmt.Println(text)
	//upload to campaign service
	serviceurl := fmt.Sprintf("http://%s/DuoMessageTemplate/EmailService.svc/Json/saveEmailInformation", casServerHost)
	//jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	//internalAuthToken := fmt.Sprintf("%s:%s", sessionInfo["TenantId"], sessionInfo["CompanyId"])

	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(emailByte))
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

func GenerateEmail(fromEmail, subject, message, toEmail string) Email {
	toEmailAddresses := make([]string, 0)
	_emailInfo := EmailInformation{}
	_email := Email{}

	toEmailAddresses = append(toEmailAddresses, toEmail)
	_emailInfo.Date = fmt.Sprintf("/Date(%d)/", time.Now().UnixNano())
	_emailInfo.ToEmailAddresses = toEmailAddresses
	_emailInfo.Subject = subject
	_emailInfo.Content = message

	_email.EmailInformation = _emailInfo
	_email.SecurityToken = v5_1SecurityToken

	return _email
}

func SendEmail(company, tenant int, resourceServer ResourceServerInfo, campaignId, camClass, camType, camCategory, fromEmail, subject, message, toEmail string) {
	IncrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, 240, camClass, camType, camCategory, "1", campaignId, toEmail, toEmail, "start", "start", time.Now().UTC().Format(layout4), resourceServer.ResourceServerId)
	emailRequest := GenerateEmail(fromEmail, subject, message, toEmail)
	SaveEmailInformation(emailRequest)

	DecrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)

	SetSessionInfo(campaignId, toEmail, "Reason", "dial_success")
	SetSessionInfo(campaignId, toEmail, "DialerStatus", "connected")
	go UploadSessionInfo(campaignId, toEmail)
}

func RequestEmailInfo(company, tenant int, campaignId string) EmailAdditionalData {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RequestEmailInfo", r)
		}
	}()
	//Request campaign from Campaign Manager service
	additionalInfo := EmailAdditionalData{}

	response := RequestCampaignAddtionalData(company, tenant, campaignId, "EMAIL", "mode1", "BLAST")
	if response != "" {
		var tempAdditionalInfo EmailAdditionalData
		json.Unmarshal([]byte(response), &tempAdditionalInfo)
		additionalInfo = tempAdditionalInfo
	}
	return additionalInfo
}
