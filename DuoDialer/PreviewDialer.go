package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

func AddPreviewDialRequest(company, tenant int, callServer CallServerInfo, campaignId, dialoutMec, uuid, fromNumber, trunkCode, phoneNumber, tryCount, extention string) {
	fmt.Println("Start AddPreviewDialRequest: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)

	IncrConcurrentChannelCount(callServer.CallServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, 240, tryCount, campaignId, uuid, phoneNumber, "ards added", "start", time.Now().UTC().Format(layout4), callServer.CallServerId)
	SetSessionInfo(campaignId, uuid, "FromNumber", fromNumber)
	SetSessionInfo(campaignId, uuid, "TrunkCode", trunkCode)
	SetSessionInfo(campaignId, uuid, "Extention", extention)

	//get attribute info from redis ** after put data stucture to cam service
	attributeInfo := make([]string, 0)

	resp, err := AddRequest(company, tenant, uuid, dialoutMec, attributeInfo)
	if err != nil {
		DecrConcurrentChannelCount(callServer.CallServerId, campaignId)
		SetSessionInfo(campaignId, uuid, "Reason", "ards_failed")
		SetSessionInfo(campaignId, uuid, "DialerStatus", "failed")
		go UploadSessionInfo(campaignId, uuid)
		fmt.Println(err.Error())
	}

	if resp != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		result := string(response)
		fmt.Println("response Body:", result)

		var ardsRes = ArdsResult{}
		json.Unmarshal(response, &ardsRes)
		if ardsRes.IsSuccess == false {
			DecrConcurrentChannelCount(callServer.CallServerId, campaignId)
			SetSessionInfo(campaignId, uuid, "Reason", ardsRes.CustomMessage)
			SetSessionInfo(campaignId, uuid, "DialerStatus", "failed")
			go UploadSessionInfo(campaignId, uuid)
		}
	}
}

func SendPreviewDataToAgent(resourceInfo ArdsCallback, campaignId, previewData string) {
	//send call detail to given agent
	//if done update Session ExpireTime
}

func DialPreviewNumber(company, tenant int, campaignId, sessionId string) {
	sessionInfoKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
	if RedisCheckKeyExist(sessionInfoKey) {
		sessionInfo := RedisHashGetAll(sessionInfoKey)
		fromNumber := sessionInfo["FromNumber"]
		trunkCode := sessionInfo["TrunkCode"]
		phoneNumber := sessionInfo["Number"]
		extention := sessionInfo["Extention"]
		callServerId := sessionInfo["ServerId"]

		callServer := GetCallServerInfo(callServerId)

		fmt.Println("Start DialPreviewNumber: ", sessionId, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)
		customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)
		param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, sessionId, fromNumber)
		furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, phoneNumber, extention)
		data := " xml dialer"

		SetSessionInfo(campaignId, sessionId, "Reason", "Dial Number")

		resp, err := Dial(callServer.Url, param, furl, data)
		HandleDialResponse(resp, err, callServer, campaignId, sessionId)
	}
}

func RejectPreviewNumber(campaignId, sessionId, rejectReason string) {
	sessionInfoKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
	if RedisCheckKeyExist(sessionInfoKey) {
		callServerId := RedisHashGetField(sessionInfoKey, "ServerId")

		callServer := GetCallServerInfo(callServerId)
		DecrConcurrentChannelCount(callServer.CallServerId, campaignId)
		SetSessionInfo(campaignId, sessionId, "Reason", rejectReason)
		SetSessionInfo(campaignId, sessionId, "DialerStatus", "agent_reject")
		go UploadSessionInfo(campaignId, sessionId)
	}
}
