package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
)

//Add preview dial request to dialer
func AddPreviewDialRequest(company, tenant int, resourceServer ResourceServerInfo, campaignId, scheduleId, campaignName, dialoutMec, uuid, fromNumber, trunkCode, phoneNumber, xGateway, numExtraData, tryCount, extention string, integrationData *IntegrationConfig, contacts *[]Contact, thirdpartyreference, businessUnit string, numberWiseSkills []string) {
	fmt.Println("Start AddPreviewDialRequest: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention, ": ", xGateway)

	strTenant := strconv.Itoa(tenant)
	strCompany := strconv.Itoa(company)
	AddCampaignCallsRealtime(phoneNumber, tryCount, "WAITING", strTenant, strCompany, campaignId, uuid)

	IncrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, -1, "Campaign", "Dialer", "PreviewDial", tryCount, campaignId, scheduleId, campaignName, uuid, phoneNumber, "ards added", "dial_start", time.Now().UTC().Format(layout4), resourceServer.ResourceServerId, integrationData, contacts, numExtraData, thirdpartyreference, businessUnit)
	InitiateAgentSessionInfo(company, tenant, -1, campaignId, campaignName, uuid, phoneNumber, integrationData, thirdpartyreference)
	SetSessionInfo(campaignId, uuid, "FromNumber", fromNumber)
	SetSessionInfo(campaignId, uuid, "TrunkCode", trunkCode)
	SetSessionInfo(campaignId, uuid, "Extention", extention)
	SetSessionInfo(campaignId, uuid, "XGateway", xGateway)
	SetSessionInfo(campaignId, uuid, "ResourceServerId", resourceServer.ResourceServerId)

	//get attribute info from redis ** after put data stucture to cam service
	attributeInfo := make([]string, 0)

	if numberWiseSkills != nil && len(numberWiseSkills) > 0 {
		attributeInfo = numberWiseSkills
	} else {
		attributeInfo = RequestCampaignAttributeInfo(company, tenant, campaignId)
	}

	reqOtherData := RequestOtherData{}
	reqOtherData.CampaignId = campaignId
	reqOtherData.StrData = numExtraData
	reqOtherData.DialoutMec = dialoutMec
	tmpReqOtherData, _ := json.Marshal(reqOtherData)
	tmpAttributeInfo, _ := json.Marshal(attributeInfo)
	SetSessionInfo(campaignId, uuid, "Skills", string(tmpAttributeInfo))

	resp, err := AddRequest(company, tenant, uuid, string(tmpReqOtherData), attributeInfo)
	if err != nil {
		SetSessionInfo(campaignId, uuid, "Reason", "ards_failed")
		SetSessionInfo(campaignId, uuid, "DialerStatus", "failed")
		go UploadSessionInfo(campaignId, uuid)
		fmt.Println(err.Error())
	}

	if resp != "" {
		var ardsRes = ArdsResult{}
		json.Unmarshal([]byte(resp), &ardsRes)
		if ardsRes.IsSuccess == false {
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
	pushD.From = "Campaign" + campaignId
	pushD.To = resourceInfo.ResourceInfo.ResourceName
	pushD.Direction = "STATEFUL"
	pushD.Message = reqOData.StrData
	pushD.CallbackURL = fmt.Sprintf("http://%s/DVP/DialerAPI/PreviewCallBack", CreateHost(lbIpAddress, lbPort))
	pushD.Ref = refDataStr

	jsonData, _ := json.Marshal(pushD)

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%s:%s", resourceInfo.Tenant, resourceInfo.Company)
	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/NotificationService/Notification/initiate", CreateHost(notificationServiceHost, notificationServicePort))
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	req.Header.Set("eventname", "preview_dialer_message")
	fmt.Println("request:", serviceurl)
	fmt.Println(string(jsonData))

	client := &http.Client{}
	redGreen := color.New(color.FgRed).Add(color.BgGreen)
	redGreen.Println("=========== Preview Data Set To UI ==========")
	resp, err := client.Do(req)
	if err != nil {
		redGreen.Println("=========== ERROR SENDING Preview Data Set To UI ==========")
		fmt.Println(err.Error())
	} else {
		t := strconv.FormatInt(time.Now().Unix(), 10)
		RedisHashSetField("CALLBACK_TIMEOUTS", campaignId+":"+resourceInfo.SessionID, t)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	//if done update Session ExpireTime
	SetSessionInfo(campaignId, resourceInfo.SessionID, "ResourceId", resourceInfo.ResourceInfo.ResourceId)
	SetSessionInfo(campaignId, resourceInfo.SessionID, "ArdsCategory", resourceInfo.RequestType)
}
