package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//Add preview dial request to dialer
func AddPreviewDialRequest(company, tenant int, callServer CallServerInfo, campaignId, dialoutMec, uuid, fromNumber, trunkCode, phoneNumber, numExtraData, tryCount, extention string) {
	fmt.Println("Start AddPreviewDialRequest: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)

	IncrConcurrentChannelCount(callServer.CallServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "PreviewDial", tryCount, campaignId, uuid, phoneNumber, "ards added", "start", time.Now().UTC().Format(layout4), callServer.CallServerId)
	SetSessionInfo(campaignId, uuid, "FromNumber", fromNumber)
	SetSessionInfo(campaignId, uuid, "TrunkCode", trunkCode)
	SetSessionInfo(campaignId, uuid, "Extention", extention)

	//get attribute info from redis ** after put data stucture to cam service
	attributeInfo := make([]string, 0)

	attributeInfo = RequestCampaignAttributeInfo(company, tenant, campaignId)

	reqOtherData := RequestOtherData{}
	reqOtherData.CampaignId = campaignId
	reqOtherData.StrData = numExtraData
	reqOtherData.DialoutMec = dialoutMec
	tmpReqOtherData, _ := json.Marshal(reqOtherData)

	resp, err := AddRequest(company, tenant, uuid, string(tmpReqOtherData), attributeInfo)
	if err != nil {
		DecrConcurrentChannelCount(callServer.CallServerId, campaignId)
		SetSessionInfo(campaignId, uuid, "Reason", "ards_failed")
		SetSessionInfo(campaignId, uuid, "DialerStatus", "failed")
		go UploadSessionInfo(campaignId, uuid)
		fmt.Println(err.Error())
	}

	if resp != "" {
		var ardsRes = ArdsResult{}
		json.Unmarshal([]byte(resp), &ardsRes)
		if ardsRes.IsSuccess == false {
			DecrConcurrentChannelCount(callServer.CallServerId, campaignId)
			SetSessionInfo(campaignId, uuid, "Reason", ardsRes.CustomMessage)
			SetSessionInfo(campaignId, uuid, "DialerStatus", "failed")
			go UploadSessionInfo(campaignId, uuid)
		}
	}
}

//Send data to agent for preview
func SendPreviewDataToAgent(resourceInfo ArdsCallbackInfo, reqOData RequestOtherData) {
	//send call detail to given agent
	refData, _ := json.Marshal(resourceInfo)
	refDataStr := string(refData)
	campaignId := reqOData.CampaignId

	pushD := PushData{}
	pushD.From = campaignId
	pushD.To = resourceInfo.ResourceInfo.ResourceId
	pushD.Direction = "BY"
	pushD.Message = reqOData.StrData
	pushD.Callback = fmt.Sprintf("http://%s/DVP/DialerAPI/PreviewCallBack", CreateHost(lbIpAddress, lbPort))
	pushD.Ref = refDataStr

	jsonData, _ := json.Marshal(pushD)

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%s:%s", resourceInfo.Tenant, resourceInfo.Company)
	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/NotificationService/Notification/initiate", CreateHost(notificationServiceHost, notificationServicePort))
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	fmt.Println("request:", serviceurl)
	fmt.Println(string(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	//if done update Session ExpireTime
	SetSessionInfo(campaignId, resourceInfo.SessionID, "Resource", resourceInfo.ResourceInfo.ResourceId)
	SetSessionInfo(campaignId, resourceInfo.SessionID, "ArdsCategory", resourceInfo.RequestType)
}
