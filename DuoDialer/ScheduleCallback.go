package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func SchedulePreviewCallback(company, tenant int, sessionId, phoneNumber, previewData, extention string, attributeInfo []string) {

	campaignId := "ScheduleCallbak"
	campaignName := "ScheduleCallbak"
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)

	resourceServerInfos := RegisterCallServer(company, tenant)
	trunkCode, ani, dnis, xGateway := GetTrunkCode(internalAuthToken, "", phoneNumber)

	InitiateSessionInfo(company, tenant, 240, "Campaign", "ScheduleCallbak", "PreviewDial", "1", campaignId, "", campaignName, sessionId, dnis, "ards added", "dial_start", time.Now().UTC().Format(layout4), resourceServerInfos.ResourceServerId)
	SetSessionInfo(campaignId, sessionId, "FromNumber", ani)
	SetSessionInfo(campaignId, sessionId, "TrunkCode", trunkCode)
	SetSessionInfo(campaignId, sessionId, "Extention", extention)
	SetSessionInfo(campaignId, sessionId, "XGateway", xGateway)

	reqOtherData := RequestOtherData{}
	reqOtherData.CampaignId = campaignId
	reqOtherData.StrData = previewData
	reqOtherData.DialoutMec = "ScheduledPreviewCallback"
	tmpReqOtherData, _ := json.Marshal(reqOtherData)

	resp, err := AddRequest(company, tenant, sessionId, string(tmpReqOtherData), attributeInfo)
	if err != nil {
		SetSessionInfo(campaignId, sessionId, "Reason", "ards_failed")
		SetSessionInfo(campaignId, sessionId, "DialerStatus", "failed")
		go UploadSessionInfo(campaignId, sessionId)
		fmt.Println(err.Error())
	}

	if resp != "" {
		var ardsRes = ArdsResult{}
		json.Unmarshal([]byte(resp), &ardsRes)
		if ardsRes.IsSuccess == false {
			SetSessionInfo(campaignId, sessionId, "Reason", ardsRes.CustomMessage)
			SetSessionInfo(campaignId, sessionId, "DialerStatus", "failed")
			go UploadSessionInfo(campaignId, sessionId)
		}
	}

}

func ScheduleIvrCallback(company, tenant int, sessionId, phoneNumber, extention string) {

	campaignId := "ScheduleCallbak"
	campaignName := "ScheduleCallbak"
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)

	resourceServerInfos := RegisterCallServer(company, tenant)
	trunkCode, ani, dnis, xGateway := GetTrunkCode(internalAuthToken, "", phoneNumber)

	InitiateSessionInfo(company, tenant, 240, "Campaign", "ScheduleCallbak", "IVR", "1", campaignId, "", campaignName, sessionId, dnis, "start", "dial_start", time.Now().UTC().Format(layout4), resourceServerInfos.ResourceServerId)

	fmt.Println("Start DialNumber: ", sessionId, ": ", ani, ": ", trunkCode, ": ", dnis, ": ", extention)

	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)

	var param string
	if xGateway != "" {
		param = fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CampaignName=%s,CustomCompanyStr=%s,OperationType=Dialer,DVP_ACTION_CAT=DIALER,DVP_OPERATION_CAT=CUSTOMER,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30, sip_h_X-Gateway=%s}", subChannelName, campaignId, campaignName, customCompanyStr, sessionId, ani, xGateway)
	} else {
		param = fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CampaignName=%s,CustomCompanyStr=%s,OperationType=Dialer,DVP_ACTION_CAT=DIALER,DVP_OPERATION_CAT=CUSTOMER,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, campaignName, customCompanyStr, sessionId, ani)
	}
	furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, phoneNumber, extention)
	data := " xml dialer"

	resp, err := Dial(resourceServerInfos.Url, param, furl, data)

	if err != nil {
		SetSessionInfo(campaignId, sessionId, "Reason", "dial_failed")
		SetSessionInfo(campaignId, sessionId, "DialerStatus", "dial_failed")
		go UploadSessionInfo(campaignId, sessionId)
		fmt.Println(err.Error())
	}

	if resp != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		tmx := string(response[:])
		fmt.Println(tmx)
		resultInfo := strings.Split(tmx, " ")
		if len(resultInfo) > 0 {
			if resultInfo[0] == "-ERR" {

				if len(resultInfo) > 1 {
					reason := resultInfo[1]
					if reason == "" {
						SetSessionInfo(campaignId, sessionId, "Reason", "not_specified")
					} else {
						SetSessionInfo(campaignId, sessionId, "Reason", reason)
					}
				} else {
					SetSessionInfo(campaignId, sessionId, "Reason", "not_specified")
				}
				SetSessionInfo(campaignId, sessionId, "DialerStatus", "dial_failed")
				go UploadSessionInfo(campaignId, sessionId)
			} else {
				SetSessionInfo(campaignId, sessionId, "Reason", "dial_success")
				SetSessionInfo(campaignId, sessionId, "DialerStatus", "dial_success")
			}
		}
	}

}
