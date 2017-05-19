package main

import (
	"fmt"
	"time"
)

func DialNumber(company, tenant int, resourceServer ResourceServerInfo, campaignId, scheduleId, campaignName, uuid, fromNumber, trunkCode, phoneNumber, tryCount, extention string) {
	fmt.Println("Start DialNumber: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)
	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)
	param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,DVP_OPERATION_CAT=DIALER,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, uuid, fromNumber)
	furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, phoneNumber, extention)
	data := " xml dialer"

	IncrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "BlastDial", tryCount, campaignId, scheduleId, campaignName, uuid, phoneNumber, "start", "dial_start", time.Now().UTC().Format(layout4), resourceServer.ResourceServerId)

	resp, err := Dial(resourceServer.Url, param, furl, data)
	HandleDialResponse(resp, err, resourceServer, campaignId, uuid)
}
