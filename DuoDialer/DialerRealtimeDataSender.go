package main

import (
	"strconv"
	"fmt"
	"github.com/fatih/color"
)


func AddCampaignDataRealtime(campaignData Campaign) {
	color.Cyan(fmt.Sprintf("Adding Campaign Realtime Data"))
	campInfoRealTime := make(map[string]string)

	campInfoRealTime["CampaignId"] = strconv.Itoa(campaignData.CampaignId)
	campInfoRealTime["CampaignName"] = campaignData.CampaignName
	campInfoRealTime["StartTime"] = campaignData.CampConfigurations.StartDate.Format("02 Jan 06 15:04 -0700")
	campInfoRealTime["EndTime"] = campaignData.CampConfigurations.EndDate.Format("02 Jan 06 15:04 -0700")
	campInfoRealTime["CampaignMode"] = campaignData.CampaignMode
	campInfoRealTime["CampaignChannel"] = campaignData.CampaignChannel
	campInfoRealTime["DialoutMechanism"] = campaignData.DialoutMechanism
	campInfoRealTime["Extension"] = campaignData.Extensions
	campInfoRealTime["OperationalStatus"] = "WAITING"

	key := fmt.Sprintf("RealTimeCampaign:%d:%d:%d", campaignData.TenantId, campaignData.CompanyId, campaignData.CampaignId)

	RedisHMSet(key, campInfoRealTime)
	
}

func AddCampaignCallsRealtime(PhoneNumber, TryCount, DialState, TenantId, CompanyId, CampaignId, ScheduleId string) {
	color.Cyan(fmt.Sprintf("Adding Campaign Call Realtime Data"))
	campCallRealTime := make(map[string]string)

	campCallRealTime["PhoneNumber"] = PhoneNumber
	campCallRealTime["TryCount"] = TryCount
	campCallRealTime["DialState"] = DialState
	campCallRealTime["TenantId"] = TenantId
	campCallRealTime["CompanyId"] = CompanyId
	campCallRealTime["CampaignId"] = CampaignId
	campCallRealTime["ScheduleId"] = ScheduleId

	key := fmt.Sprintf("RealTimeCampaignCalls:%s:%s:%s:%s", TenantId, CompanyId, CampaignId, ScheduleId)

	RedisHMSet(key, campCallRealTime)
	
}

func UpdateCampaignRealtimeField(fieldName, val string, tenantId, companyId, campaignId int) {
	color.Cyan(fmt.Sprintf("Updating Campaign Realtime Field"))

	key := fmt.Sprintf("RealTimeCampaign:%d:%d:%d", tenantId, companyId, campaignId)

	RedisHashSetField(key, fieldName, val)
	
}

func UpdateCampaignCallRealtimeField(fieldName, val, tenantId, companyId, campaignId, scheduleId string) {
	color.Cyan(fmt.Sprintf("Updating Campaign Realtime Field"))

	key := fmt.Sprintf("RealTimeCampaignCalls:%s:%s:%s:%s", tenantId, companyId, campaignId, scheduleId)

	RedisHashSetField(key, fieldName, val)
	
}

func RemoveCampaignRealtime(tenantId, companyId, campaignId int) {
	color.Cyan(fmt.Sprintf("Removing Campaign Realtime"))

	key := fmt.Sprintf("RealTimeCampaign:%d:%d:%d", tenantId, companyId, campaignId)

	RedisRemove(key)
	
}

func RemoveCampaignCallRealtime(tenantId, companyId, campaignId, scheduleId string) {
	color.Cyan(fmt.Sprintf("Removing Campaign Realtime"))

	key := fmt.Sprintf("RealTimeCampaignCalls:%s:%s:%s:%s", tenantId, companyId, campaignId, scheduleId)

	RedisRemove(key)
	
}