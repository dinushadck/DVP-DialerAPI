package main

import (
	"fmt"
	"strconv"
	"time"
)

func DirectDialCampaign(company, tenant, campaignId int, number string) bool {
	fmt.Println("start DialNumber")
	campaignIdStr := strconv.Itoa(campaignId)
	campaignInfo, isCamExists := GetCampaign(company, tenant, campaignId)
	if isCamExists {
		tmNowUTC := time.Now().UTC()

		tempCampaignEndDate, _ := time.Parse(layout1, campaignInfo.CampConfigurations.EndDate)
		campaignEndDate := time.Date(tempCampaignEndDate.Year(), tempCampaignEndDate.Month(), tempCampaignEndDate.Day(), tempCampaignEndDate.Hour(), tempCampaignEndDate.Minute(), tempCampaignEndDate.Second(), 0, time.UTC)

		if campaignEndDate.After(tmNowUTC) {
			scheduleIdStr := strconv.Itoa(campaignInfo.CampScheduleInfo[0].ScheduleId)
			camScheduleStr := strconv.Itoa(campaignInfo.CampScheduleInfo[0].CamScheduleId)
			validateAppoinment := CheckAppoinmentForCallback(company, tenant, scheduleIdStr, tmNowUTC)
			if validateAppoinment {
				numberWithTryCount := fmt.Sprintf("%s:%d", number, 1)
				return AddNumberToFront(company, tenant, campaignIdStr, camScheduleStr, numberWithTryCount)
			}
		}
	}
	return false
}

func DirectDial(company, tenant int, fromNumber, phoneNumber, extention, callServerId string) bool {

	authToken := fmt.Sprintf("%d#%d", tenant, company)
	trunkCode, ani, dnis := GetTrunkCode(authToken, fromNumber, phoneNumber)
	uuid := GetUuid()
	if trunkCode != "" && uuid != "" {
		fmt.Println("Start AddDirectDialRequest: ", uuid, ": ", ani, ": ", trunkCode, ": ", dnis, ": ", extention)
		campaignId := "DirectDial"
		callServer := GetCallServerInfo(company, tenant, callServerId)

		IncrConcurrentChannelCount(callServer.CallServerId, campaignId)
		IncrCampaignDialCount(company, tenant, campaignId)
		InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "DirectDial", "1", campaignId, uuid, dnis, "direct dial", "start", time.Now().UTC().Format(layout4), callServerId)
		SetSessionInfo(campaignId, uuid, "FromNumber", ani)
		SetSessionInfo(campaignId, uuid, "TrunkCode", trunkCode)
		SetSessionInfo(campaignId, uuid, "Extention", extention)

		fmt.Println("Start DialDirectNumber: ", uuid, ": ", ani, ": ", trunkCode, ": ", dnis, ": ", extention)
		customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)
		param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, uuid, ani)
		furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, dnis, extention)
		data := " xml dialer"

		SetSessionInfo(campaignId, uuid, "Reason", "Dial Number")

		resp, err := Dial(callServer.Url, param, furl, data)
		HandleDialResponse(resp, err, callServer, campaignId, uuid)
		return true
		//}
	}
	return false
}
