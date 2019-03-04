package main

import (
	"fmt"
	"time"
)

func DialNumberFIFO(company, tenant int, resourceServer ResourceServerInfo, campaignId, scheduleId, campaignName, uuid, fromNumber, trunkCode, phoneNumber, xGateway, extention string, integrationData *IntegrationConfig, contacts *[]Contact) {
	fmt.Println("Start DialNumber: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention, ": ", xGateway)
	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)

	var param string

	if xGateway != "" {
		param = fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CampaignName='%s',CustomCompanyStr=%s,OperationType=Dialer,DVP_ACTION_CAT=DIALER,DVP_OPERATION_CAT=CUSTOMER,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30,sip_h_X-Gateway=%s}", subChannelName, campaignId, campaignName, customCompanyStr, uuid, fromNumber, xGateway)
	} else {
		param = fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CampaignName='%s',CustomCompanyStr=%s,OperationType=Dialer,DVP_ACTION_CAT=DIALER,DVP_OPERATION_CAT=CUSTOMER,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, campaignName, customCompanyStr, uuid, fromNumber)
	}
	furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, phoneNumber, extention)
	data := " xml dialer"

	IncrConcurrentChannelCount(resourceServer.ResourceServerId, campaignId)
	InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "FIFODial", "1", campaignId, scheduleId, campaignName, uuid, phoneNumber, "start", "dial_start", time.Now().Format(layout4), resourceServer.ResourceServerId, integrationData, contacts, "", "")
	IncrCampaignDialCount(company, tenant, campaignId)

	resp, err := Dial(resourceServer.Url, param, furl, data)
	HandleDialResponse(resp, err, resourceServer, campaignId, uuid)
}
