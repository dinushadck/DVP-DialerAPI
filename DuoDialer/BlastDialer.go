package main

import (
	"fmt"
	"time"
)

func DialNumber(company, tenant int, resourceServer ResourceServerInfo, campaignId, scheduleId, campaignName, uuid, fromNumber, trunkCode, phoneNumber, xGateway, tryCount, extention string) {
	fmt.Println("Start DialNumber: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention, ": ", xGateway)
	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)

	var param string
	if xGateway != "" {
		param = fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CampaignName='%s',tenantid=%d,companyid=%d,CustomCompanyStr=%s,OperationType=Dialer,DVP_ACTION_CAT=DIALER,DVP_OPERATION_CAT=CUSTOMER,DVP_ADVANCED_OP_ACTION=BLAST,return_ring_ready=false,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30,sip_h_X-Gateway=%s}", subChannelName, campaignId, campaignName, tenant, company, customCompanyStr, uuid, fromNumber, xGateway)
	} else {
		param = fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CampaignName='%s',tenantid=%d,companyid=%d,CustomCompanyStr=%s,OperationType=Dialer,DVP_ACTION_CAT=DIALER,DVP_OPERATION_CAT=CUSTOMER,DVP_ADVANCED_OP_ACTION=BLAST,return_ring_ready=false,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, campaignName, tenant, company, customCompanyStr, uuid, fromNumber)
	}
	furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, phoneNumber, extention)
	data := " xml dialer"

	IncrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)
	IncrCampaignDialCount(company, tenant, campaignId)
	InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "BlastDial", tryCount, campaignId, scheduleId, campaignName, uuid, phoneNumber, "start", "dial_start", time.Now().UTC().Format(layout4), resourceServer.ResourceServerId)

	resp, err := Dial(resourceServer.Url, param, furl, data)
	HandleDialResponse(resp, err, resourceServer, campaignId, uuid)
}
