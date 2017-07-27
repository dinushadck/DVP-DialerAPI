package main

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
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
