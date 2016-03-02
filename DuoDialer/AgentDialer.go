package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//Add preview dial request to dialer
func AddAgentDialRequest(company, tenant int, callServer CallServerInfo, campaignId, dialoutMec, uuid, fromNumber, trunkCode, phoneNumber, numExtraData, tryCount, extention string) {
	fmt.Println("Start AddPreviewDialRequest: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)

	IncrConcurrentChannelCount(callServer.CallServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "AgentDial", tryCount, campaignId, uuid, phoneNumber, "ards added", "start", time.Now().UTC().Format(layout4), callServer.CallServerId)
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

//Once agent accept, send dial data to server
func DialAgent(contactName, domain, contactType, resourceId, company, tenant, campaignId, ardsServerType, ardsReqType, sessionId string) {
	sessionInfoKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
	if RedisCheckKeyExist(sessionInfoKey) {
		sessionInfo := RedisHashGetAll(sessionInfoKey)
		fromNumber := sessionInfo["FromNumber"]
		trunkCode := sessionInfo["TrunkCode"]
		phoneNumber := sessionInfo["Number"]
		extention := sessionInfo["Extention"]
		callServerId := sessionInfo["ServerId"]

		companyInt, _ := strconv.Atoi(company)
		tenantInt, _ := strconv.Atoi(tenant)
		callServer := GetCallServerInfo(companyInt, tenantInt, callServerId)

		fmt.Println("Start DialPreviewNumber: ", sessionId, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)
		customCompanyStr := fmt.Sprintf("%s_%s", company, tenant)
		param := fmt.Sprintf(" {sip_h_DVP-DESTINATION-TYPE=GATEWAY,DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, sessionId, fromNumber)
		furl := fmt.Sprintf("sofia/gateway/%s/%s ", trunkCode, phoneNumber)
		var data string
		var dial bool
		if contactType == "PRIVATE" {
			dial = true
			data = fmt.Sprintf(" &bridge({sip_h_DVP-DESTINATION-TYPE=PRIVATE_USER,ards_client_uuid=%s,ards_resource_id=%s,tenantid=%s,companyid=%s,ards_server_type=%s,ards_request_type=%s}user/%s@%s)", sessionId, resourceId, tenant, company, ardsServerType, ardsReqType, contactName, domain)
		} else if contactType == "PUBLIC" {
			dial = true
			data = fmt.Sprintf(" &bridge({sip_h_DVP-DESTINATION-TYPE=PUBLIC_USER,ards_client_uuid=%s,ards_resource_id=%s,tenantid=%s,companyid=%s,ards_server_type=%s,ards_request_type=%s}sofia/external/%s@%s)", sessionId, resourceId, tenant, company, ardsServerType, ardsReqType, contactName, domain)
		} else if contactType == "TRUNK" {
			dial = true
			data = fmt.Sprintf(" &bridge({sip_h_DVP-DESTINATION-TYPE=GATEWAY,ards_client_uuid=%s,ards_resource_id=%s,tenantid=%s,companyid=%s,ards_server_type=%s,ards_request_type=%s}sofia/gateway/%s/%s)", sessionId, resourceId, tenant, company, ardsServerType, ardsReqType, domain, contactName)
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
			AgentReject(company, tenant, campaignId, sessionId, ardsReqType, resourceId, "Invalied ContactType")
		}
	}
}

func AgentReject(company, tenant, campaignId, sessionId, requestType, resourceId, rejectReason string) {
	sessionInfoKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
	if RedisCheckKeyExist(sessionInfoKey) {
		callServerId := RedisHashGetField(sessionInfoKey, "ServerId")
		companyInt, _ := strconv.Atoi(company)
		tenantInt, _ := strconv.Atoi(tenant)
		callServer := GetCallServerInfo(companyInt, tenantInt, callServerId)
		DecrConcurrentChannelCount(callServer.CallServerId, campaignId)
		SetSessionInfo(campaignId, sessionId, "Reason", rejectReason)
		SetSessionInfo(campaignId, sessionId, "DialerStatus", "agent_reject")
		ClearResourceSlotWhenReject(company, tenant, requestType, resourceId, sessionId)
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
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)

	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaign/%s/AdditinalData/PREVIEW/ARDS/ATTRIBUTE", CreateHost(campaignServiceHost, campaignServicePort), campaignId)
	fmt.Println("Start RequestCampaign request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
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
