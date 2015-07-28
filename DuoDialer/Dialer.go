package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetUuid() string {
	resp, err := http.Get(uuidService)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	} else {
		defer resp.Body.Close()
		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		} else {
			tmx := string(response[:])
			fmt.Println(tmx)
			return tmx
		}
	}
}

func GetTrunkCode(authToken, ani, dnis string) (trunkCode, rAni, rDnis string) {
	fmt.Println("Start GetTrunkCode: ", authToken, ": ", ani, ": ", dnis)
	client := &http.Client{}

	//request := fmt.Sprintf("%s/ANI/%s/DNIS/%s", callRuleService, ani, dnis)
	request := fmt.Sprintf("%s?ANI=%s&DNIS=%s", callRuleService, ani, dnis)
	fmt.Println("Start GetTrunkCode request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", ""
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var apiResult CallRuleApiResult
	json.Unmarshal(response, &apiResult)
	if apiResult.IsSuccess == true {
		fmt.Println("callRule: ", apiResult.Result.GatewayCode, "ANI: ", apiResult.Result.ANI, "DNIS: ", apiResult.Result.DNIS)
		return apiResult.Result.GatewayCode, apiResult.Result.ANI, apiResult.Result.DNIS
	} else {
		return "", "", ""
	}
}

func DialNumber(company, tenant int, callServer CallServerInfo, campaignId, uuid, fromNumber, trunkCode, phoneNumber, extention string) {
	fmt.Println("Start DialNumber: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)
	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)
	request := fmt.Sprintf("http://%s", callServer.Url)
	path := fmt.Sprintf("api/originate?")
	param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,return_ring_ready=true,origination_uuid=%s,origination_caller_id_number=%s}sofia/gateway/%s/%s %s xml dialer", subChannelName, campaignId, customCompanyStr, uuid, fromNumber, trunkCode, phoneNumber, extention)

	u, _ := url.Parse(request)
	u.Path += path
	u.Path += param

	fmt.Println(u.String())
	IncrConcurrentChannelCount(callServer.CallServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, "1", campaignId, uuid, "start", "start", time.Now().Format(layout4), callServer.CallServerId)
	resp, err := http.Get(u.String())
	if err != nil {
		DecrConcurrentChannelCount(callServer.CallServerId, campaignId)
		SetSessionInfo(uuid, "reason", "dial_failed")
		SetSessionInfo(uuid, "dialerStatus", "failed")
		go UploadSessionInfo(uuid)
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	if resp != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		tmx := string(response[:])
		fmt.Println(tmx)
		resultInfo := strings.Split(tmx, " ")
		if len(resultInfo) > 0 {
			if resultInfo[0] == "-ERR" {
				DecrConcurrentChannelCount(callServer.CallServerId, campaignId)

				if len(resultInfo) > 1 {
					reason := resultInfo[1]
					if reason == "" {
						SetSessionInfo(uuid, "reason", "not_specified")
					} else {
						SetSessionInfo(uuid, "reason", reason)
					}
				} else {
					SetSessionInfo(uuid, "reason", "not_specified")
				}
				SetSessionInfo(uuid, "dialerStatus", "not_connected")
				go UploadSessionInfo(uuid)
			} else {
				SetSessionInfo(uuid, "reason", "dial_success")
				SetSessionInfo(uuid, "dialerStatus", "connected")
			}
		}
	}
}

func DialNumberFIFO(company, tenant int, callServer CallServerInfo, campaignId, uuid, fromNumber, trunkCode, phoneNumber, extention string) {
	fmt.Println("Start DialNumber: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)
	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)
	request := fmt.Sprintf("http://%s", callServer.Url)
	path := fmt.Sprintf("api/originate?")
	param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,return_ring_ready=false,origination_uuid=%s,origination_caller_id_number=%s}sofia/gateway/%s/%s %s xml dialer", subChannelName, campaignId, customCompanyStr, uuid, fromNumber, trunkCode, phoneNumber, extention)

	u, _ := url.Parse(request)
	u.Path += path
	u.Path += param

	fmt.Println(u.String())
	IncrConcurrentChannelCount(callServer.CallServerId, campaignId)
	InitiateSessionInfo(company, tenant, "1", campaignId, uuid, "start", "start", time.Now().Format(layout4), callServer.CallServerId)
	IncrCampaignDialCount(company, tenant, campaignId)
	resp, err := http.Get(u.String())
	if err != nil {
		DecrConcurrentChannelCount(callServer.CallServerId, campaignId)
		SetSessionInfo(uuid, "reason", "dial_failed")
		SetSessionInfo(uuid, "dialerStatus", "failed")
		go UploadSessionInfo(uuid)
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	if resp != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		tmx := string(response[:])
		fmt.Println(tmx)
		resultInfo := strings.Split(tmx, " ")
		if len(resultInfo) > 0 {
			if resultInfo[0] == "-ERR" {
				DecrConcurrentChannelCount(callServer.CallServerId, campaignId)

				if len(resultInfo) > 1 {
					reason := resultInfo[1]
					if reason == "" {
						SetSessionInfo(uuid, "reason", "not_specified")
					} else {
						SetSessionInfo(uuid, "reason", reason)
					}
				} else {
					SetSessionInfo(uuid, "reason", "not_specified")
				}
				SetSessionInfo(uuid, "dialerStatus", "not_connected")
				go UploadSessionInfo(uuid)
			} else {
				SetSessionInfo(uuid, "reason", "dial_success")
				SetSessionInfo(uuid, "dialerStatus", "connected")
			}
		}
	}
}
