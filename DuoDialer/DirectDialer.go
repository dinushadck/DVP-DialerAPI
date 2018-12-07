package main

import (
	"fmt"
	"strconv"
	"time"
)

func DirectDialCampaign(company, tenant, campaignId, ScheduleId int, number string) bool {
	fmt.Println("start DialNumber")
	campaignIdStr := strconv.Itoa(campaignId)
	campaignInfo, isCamExists := GetCampaign(company, tenant, campaignId)
	if isCamExists {

		scheduleInfo := CampaignShedule{}
		defaultScheduleInfo := CampaignShedule{}
		for _, schedule := range campaignInfo.CampScheduleInfo {
			if schedule.ScheduleId == ScheduleId {
				scheduleInfo = schedule
				break
			}
		}

		if scheduleInfo != defaultScheduleInfo {
			location, _ := time.LoadLocation(scheduleInfo.TimeZone)
			tmNow := time.Now().In(location)

			//tempCampaignEndDate, _ := time.Parse(layout1, campaignInfo.CampConfigurations.EndDate)
			scheduleEndDate := scheduleInfo.EndDate

			if scheduleEndDate.After(tmNow) {
				scheduleIdStr := strconv.Itoa(scheduleInfo.ScheduleId)
				camScheduleStr := strconv.Itoa(scheduleInfo.CamScheduleId)
				validateAppoinment := CheckAppoinmentForCallback(company, tenant, scheduleIdStr, tmNow, scheduleInfo.TimeZone)
				if validateAppoinment {
					numberWithTryCount := fmt.Sprintf("%s:%d", number, 1)
					return AddNumberToFront(company, tenant, campaignIdStr, camScheduleStr, numberWithTryCount)
				}
			}
		} else {
			return false
		}
	}
	return false
}

func DirectDial(company, tenant int, fromNumber, phoneNumber, extention, resourceServerId string) bool {

	internalAccessToken := fmt.Sprintf("%d:%d", tenant, company)
	trunkCode, ani, dnis, xGateway := GetTrunkCode(internalAccessToken, fromNumber, phoneNumber)

	if trunkCode != "" {

		resourceServer := GetResourceServerInfo(company, tenant, resourceServerId, "call")
		uuid := GetUuid(resourceServer.Url)

		if uuid != "" {
			fmt.Println("Start AddDirectDialRequest: ", uuid, ": ", ani, ": ", trunkCode, ": ", dnis, ": ", extention)
			campaignId := "DirectDial"

			IncrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)
			IncrCampaignDialCount(company, tenant, campaignId)
			InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "DirectDial", "1", campaignId, "0", campaignId, uuid, dnis, "direct dial", "dial_start", time.Now().UTC().Format(layout4), resourceServerId, nil, nil)
			SetSessionInfo(campaignId, uuid, "FromNumber", ani)
			SetSessionInfo(campaignId, uuid, "TrunkCode", trunkCode)
			SetSessionInfo(campaignId, uuid, "Extention", extention)

			fmt.Println("Start DialDirectNumber: ", uuid, ": ", ani, ": ", trunkCode, ": ", dnis, ": ", extention)
			customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)

			var param string
			if xGateway != "" {
				param = fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CampaignName='%s',CustomCompanyStr=%s,OperationType=Dialer,DVP_ACTION_CAT=DIALER,DVP_OPERATION_CAT=CUSTOMER,DVP_ADVANCED_OP_ACTION=DIRECT,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30,sip_h_X-Gateway=%s}", subChannelName, campaignId, campaignId, customCompanyStr, uuid, ani, xGateway)
			} else {
				param = fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CampaignName='%s',CustomCompanyStr=%s,OperationType=Dialer,DVP_ACTION_CAT=DIALER,DVP_OPERATION_CAT=CUSTOMER,DVP_ADVANCED_OP_ACTION=DIRECT,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, campaignId, customCompanyStr, uuid, ani)
			}
			furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, dnis, extention)
			data := " xml dialer"

			SetSessionInfo(campaignId, uuid, "Reason", "Dial Number")

			resp, err := Dial(resourceServer.Url, param, furl, data)
			HandleDialResponse(resp, err, resourceServer, campaignId, uuid)
			return true
		}
	}
	return false
}

func ClickToCall(company, tenant int, phoneNumber, extention, resourceServerId string) bool {

	//internalAccessToken := fmt.Sprintf("%d:%d", tenant, company)
	//trunkCode, ani, dnis := GetTrunkCode(internalAccessToken, "", phoneNumber)
	//uuid := GetUuid()
	//if trunkCode != "" && uuid != "" {
	//	fmt.Println("Start Add ClickToCall Request: ", uuid, ": ", ani, ": ", trunkCode, ": ", dnis, ": ", extention)
	//	campaignId := "ClickToCall"
	//	resourceServer := GetResourceServerInfo(company, tenant, resourceServerId, "call")

	//	IncrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)
	//	IncrCampaignDialCount(company, tenant, campaignId)
	//	InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "DirectDial", "1", campaignId, campaignId, uuid, dnis, "direct dial", "dial_start", time.Now().UTC().Format(layout4), resourceServerId)
	//	SetSessionInfo(campaignId, uuid, "FromNumber", ani)
	//	SetSessionInfo(campaignId, uuid, "TrunkCode", trunkCode)
	//	SetSessionInfo(campaignId, uuid, "Extention", extention)

	//	fmt.Println("Start DialDirectNumber: ", uuid, ": ", ani, ": ", trunkCode, ": ", dnis, ": ", extention)
	//	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)
	//	param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,companyid=%d,tenantid=%d,OperationType=Dialer,DVP_OPERATION_CAT=DIALER,dvp_app_type=HTTAPI,return_ring_ready=true,ignore_early_media=true,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, company, tenant, uuid, ani)
	//	furl := fmt.Sprintf("sofia/gateway/%s/%s", trunkCode, dnis)

	//	data := fmt.Sprintf(" &transfer(%s xml )", extention)

	//	SetSessionInfo(campaignId, uuid, "Reason", "Dial Number")

	//	resp, err := Dial(resourceServer.Url, param, furl, data)
	//	HandleDialResponse(resp, err, resourceServer, campaignId, uuid)
	//	return true
	//	//}
	//}
	//return false

	//internalAccessToken := fmt.Sprintf("%d:%d", tenant, company)
	//trunkCode, ani, dnis := GetTrunkCode(internalAccessToken, "", phoneNumber)

	resourceServer := GetResourceServerInfo(company, tenant, resourceServerId, "call")
	uuid := GetUuid(resourceServer.Url)
	if uuid != "" {
		fmt.Println("Start Add ClickToCall Request: ", uuid, ": ", ": ", phoneNumber, ": ", extention)
		campaignId := "ClickToCall"

		IncrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)
		IncrCampaignDialCount(company, tenant, campaignId)
		InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "DirectDial", "1", campaignId, "0", campaignId, uuid, phoneNumber, "direct dial", "dial_start", time.Now().UTC().Format(layout4), resourceServerId, nil, nil)

		SetSessionInfo(campaignId, uuid, "Extention", extention)

		fmt.Println("Start DialDirectNumber: ", uuid, ": ", phoneNumber, ": ", extention)
		//customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)
		//param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,companyid=%d,tenantid=%d,OperationType=Dialer,DVP_OPERATION_CAT=DIALER,dvp_app_type=HTTAPI,return_ring_ready=true,ignore_early_media=true,origination_caller_id_number=%s,origination_uuid=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, company, tenant, ani, uuid)
		param := fmt.Sprintf(" {companyid=%d,tenantid=%d,origination_caller_id_number=%s,DVP_CLICKTOCALL=C2C,originate_timeout=30,force_transfer_context=PBXFeatures|%d|%d}", company, tenant, phoneNumber, tenant, company)
		furl := fmt.Sprintf("user/%s", extention)

		data := fmt.Sprintf(" &transfer(%s)", phoneNumber)

		SetSessionInfo(campaignId, uuid, "Reason", "Dial Number")

		resp, err := Dial(resourceServer.Url, param, furl, data)
		HandleDialResponse(resp, err, resourceServer, campaignId, uuid)
		return true
		//}
	}
	return false

}
