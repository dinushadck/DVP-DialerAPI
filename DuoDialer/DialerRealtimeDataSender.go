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
	campInfoRealTime["OperationalStatus"] = campaignData.OperationalStatus

	key := fmt.Sprintf("RealTimeCampaign:%d:%d:%d", campaignData.TenantId, campaignData.CompanyId, campaignData.CampaignId)

	RedisHMSet(key, campInfoRealTime)
	
}

func UpdateCampaignRealtimeField(fieldName, val string, tenantId, companyId, campaignId int) {
	color.Cyan(fmt.Sprintf("Updating Campaign Realtime Field"))

	key := fmt.Sprintf("RealTimeCampaign:%d:%d:%d", tenantId, companyId, campaignId)

	RedisHashSetField(key, fieldName, val)
	
}