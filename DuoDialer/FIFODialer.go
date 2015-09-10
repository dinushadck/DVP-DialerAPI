package main

import (
	"fmt"
	"time"
)

func DialNumberFIFO(company, tenant int, callServer CallServerInfo, campaignId, uuid, fromNumber, trunkCode, phoneNumber, extention string) {
	fmt.Println("Start DialNumber: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)
	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)

	param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, uuid, fromNumber)
	furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, phoneNumber, extention)
	data := " xml dialer"

	IncrConcurrentChannelCount(callServer.CallServerId, campaignId)
	InitiateSessionInfo(company, tenant, 240, "1", campaignId, uuid, phoneNumber, "start", "start", time.Now().Format(layout4), callServer.CallServerId)
	IncrCampaignDialCount(company, tenant, campaignId)

	resp, err := Dial(callServer.Url, param, furl, data)
	HandleDialResponse(resp, err, callServer, campaignId, uuid)
}
