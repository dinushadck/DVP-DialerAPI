package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

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

	reqOtherData := PreviewRequestOtherData{}
	reqOtherData.CampaignId = campaignId
	reqOtherData.PreviewData = numExtraData
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

func SendPreviewDataToAgent(resourceInfo ArdsCallbackInfo) {
	//send call detail to given agent
	var reqOData PreviewRequestOtherData
	json.Unmarshal([]byte(resourceInfo.OtherInfo), &reqOData)

	refData, _ := json.Marshal(resourceInfo)
	refDataStr := string(refData)

	pushD := PushData{}
	pushD.From = "Campaign"
	pushD.To = resourceInfo.ResourceInfo.ResourceId
	pushD.Direction = "BY"
	pushD.Message = reqOData.PreviewData
	pushD.Callback = fmt.Sprintf("http://%s/DVP/DialerAPI/PreviewCallBack", CreateHost(lbIpAddress, lbPort))
	pushD.Ref = refDataStr

	jsonData, _ := json.Marshal(pushD)

	authToken := fmt.Sprintf("%d#%d", resourceInfo.Tenant, resourceInfo.Company)
	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/NotificationService/Notification/initiate", CreateHost(notificationServiceHost, notificationServicePort))
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", authToken)
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
}

func DialPreviewNumber(contactName, domain, contactType, resourceId, company, tenant, campaignId, ardsClass, ardsType, ardsCategory, sessionId string) {
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
		customCompanyStr := fmt.Sprintf("%s_%s", company, tenant)
		param := fmt.Sprintf(" {sip_h_DVP-DESTINATION-TYPE=GATEWAY,DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, sessionId, fromNumber)
		furl := fmt.Sprintf("sofia/gateway/%s/%s ", trunkCode, phoneNumber)
		var data string
		var dial bool
		if contactType == "PRIVATE" {
			dial = true
			data = fmt.Sprintf(" &bridge({sip_h_DVP-DESTINATION-TYPE=PRIVATE_USER,ards_client_uuid=%s,ards_resource_id=%s,tenantid=%s,companyid=%s,ards_class=%s,ards_type=%s,ards_category=%s}user/%s@%s)", sessionId, resourceId, tenant, company, ardsClass, ardsType, ardsCategory, contactName, domain)
		} else if contactType == "PUBLIC" {
			dial = true
			data = fmt.Sprintf(" &bridge({sip_h_DVP-DESTINATION-TYPE=PUBLIC_USER,ards_client_uuid=%s,ards_resource_id=%s,tenantid=%s,companyid=%s,ards_class=%s,ards_type=%s,ards_category=%s}sofia/external/%s@%s)", sessionId, resourceId, tenant, company, ardsClass, ardsType, ardsCategory, contactName, domain)
		} else if contactType == "TRUNK" {
			dial = true
			data = fmt.Sprintf(" &bridge({sip_h_DVP-DESTINATION-TYPE=GATEWAY,ards_client_uuid=%s,ards_resource_id=%s,tenantid=%s,companyid=%s,ards_class=%s,ards_type=%s,ards_category=%s}sofia/gateway/%s/%s)", sessionId, resourceId, tenant, company, ardsClass, ardsType, ardsCategory, domain, contactName)
		} else {
			dial = false
			fmt.Println("Invalied ContactType")
		}

		if dial == true {
			SetSessionInfo(campaignId, sessionId, "Reason", "Dial Number")

			resp, err := Dial(callServer.Url, param, furl, data)
			HandleDialResponse(resp, err, callServer, campaignId, sessionId)
		} else {
			SetSessionInfo(campaignId, sessionId, "Reason", "Invalied ContactType")
			RejectPreviewNumber(campaignId, sessionId, "Invalied ContactType")
		}
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

//-----------------------------------CampaignManager Service-------------------------------------------------

func RequestCampaignAttributeInfo(company, tenant int, campaignId string) []string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RequestCampaignAttributeInfo", r)
		}
	}()
	//Request campaign from Campaign Manager service
	attributeDetails := make([]string, 0)
	authToken := fmt.Sprintf("%d#%d", tenant, company)

	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaign/%s/AdditinalData/PREVIEW/ARDS/ATTRIBUTE", CreateHost(campaignServiceHost, campaignServicePort), campaignId)
	fmt.Println("Start RequestCampaign request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return attributeDetails
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var campaignAdditionalDataResult CampaignAdditionalDataResult
	json.Unmarshal(response, &campaignAdditionalDataResult)
	if campaignAdditionalDataResult.IsSuccess == true {
		var attInfo []string
		json.Unmarshal([]byte(campaignAdditionalDataResult.Result.AdditionalData), &attInfo)
		attributeDetails = attInfo
	}
	return attributeDetails
}
