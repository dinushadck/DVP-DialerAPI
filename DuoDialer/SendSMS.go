package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GenerateSMS(smsServerUrl, fromNumber, message, phoneNumber string) string {
	resportUrl := "http://45.55.171.228:9998/reply"

	ru := url.QueryEscape(resportUrl)
	fmt.Println(ru)

	param := fmt.Sprintf("username=foo&password=bar&from=%s&to=%s&content=%s&dlr-url=%s&dlr-level=2", fromNumber, phoneNumber, message, ru)

	request := fmt.Sprintf("http://%s/send?%s", smsServerUrl, param)
	//path := fmt.Sprintf("send?")
	//param := fmt.Sprintf(" %s%s %s", params, furl, data)

	//u, _ := url.Parse(request)
	//u.Path += path
	//u.Path += param

	fmt.Println("request: ", request)
	return request
}

func SendSmsDirect(company, tenant int, message, phoneNumber string) {
	campaignId := "SMS_Direct"
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	serverInfo := GetResourceServerInfo(company, tenant, "SMS1", "sms")
	_, ani, dnis := GetTrunkCode(internalAuthToken, "", phoneNumber)

	IncrConcurrentChannelCount(serverInfo.ResourceServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, 240, "SMS", "DIRECT", "DIALOUT", "1", campaignId, campaignId, "0", dnis, dnis, "start", "dial_start", time.Now().UTC().Format(layout4), serverInfo.ResourceServerId)

	smsRequest := GenerateSMS(serverInfo.Url, ani, message, dnis)
	resp, err := SendSms(smsRequest)
	HandleSmsResponse(resp, err, serverInfo, campaignId, dnis)
}

func SendSms(smsUrl string) (*http.Response, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SendSms", r)
		}
	}()

	resp, err := http.Get(smsUrl)
	return resp, err
}

func HandleSmsResponse(resp *http.Response, err error, server ResourceServerInfo, campaignId, sessionId string) {
	if err != nil {
		DecrConcurrentChannelCount(server.ResourceServerId, campaignId)
		SetSessionInfo(campaignId, sessionId, "Reason", "dial_failed")
		SetSessionInfo(campaignId, sessionId, "DialerStatus", "failed")
		go UploadSessionInfo(campaignId, sessionId)
		fmt.Println(err.Error())
	}

	if resp != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		tmx := string(response[:])
		fmt.Println("response: ", tmx)
		resultInfo := strings.Split(tmx, " ")
		if len(resultInfo) > 0 {
			if resultInfo[0] == "Success" {
				DecrConcurrentChannelCount(server.ResourceServerId, campaignId)

				SetSessionInfo(campaignId, sessionId, "Reason", "dial_success")
				SetSessionInfo(campaignId, sessionId, "DialerStatus", "connected")
			} else {
				SetSessionInfo(campaignId, sessionId, "Reason", "dial_failed")
				SetSessionInfo(campaignId, sessionId, "DialerStatus", "not-connected")
			}

			go UploadSessionInfo(campaignId, sessionId)
		}
	}
}
