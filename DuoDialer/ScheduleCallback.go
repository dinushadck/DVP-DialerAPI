package main

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"strings"
	"time"
)

type ScheduleCallback struct {
}

func (scheduleCallback ScheduleCallback) AddPreviewCallback(company, tenant int, phoneNumber, previewData, extention string, attributeInfo []string) {

	campaignId := "ScheduleCallbak"
	campaignName := "ScheduleCallbak"
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	sessionId := uuid.NewV4().String()

	resourceServerInfos := RegisterCallServer(company, tenant)
	trunkCode, ani, dnis := GetTrunkCode(internalAuthToken, "", phoneNumber)

	InitiateSessionInfo(company, tenant, 240, "Campaign", "ScheduleCallbak", "PreviewDial", "1", campaignId, "", campaignName, sessionId, dnis, "ards added", "dial_start", time.Now().UTC().Format(layout4), resourceServerInfos.ResourceServerId)
	SetSessionInfo(campaignId, sessionId, "FromNumber", ani)
	SetSessionInfo(campaignId, sessionId, "TrunkCode", trunkCode)
	SetSessionInfo(campaignId, sessionId, "Extention", extention)

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

func (scheduleCallback ScheduleCallback) DialIvrCallback(company, tenant int, phoneNumber, extention string) {

	sessionId := uuid.NewV4().String()

	campaignId := "ScheduleCallbak"
	campaignName := "ScheduleCallbak"
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)

	resourceServerInfos := RegisterCallServer(company, tenant)
	trunkCode, ani, dnis := GetTrunkCode(internalAuthToken, "", phoneNumber)

	fmt.Println("Start DialNumber: ", sessionId, ": ", ani, ": ", trunkCode, ": ", dnis, ": ", extention)

	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)
	param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,DVP_OPERATION_CAT=DIALER,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, sessionId, ani)
	furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, phoneNumber, extention)
	data := " xml dialer"

	InitiateSessionInfo(company, tenant, 240, "Campaign", "ScheduleCallbak", "IVR", "1", campaignId, "", campaignName, sessionId, dnis, "start", "dial_start", time.Now().UTC().Format(layout4), resourceServerInfos.ResourceServerId)

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
