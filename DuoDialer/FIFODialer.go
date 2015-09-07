package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func DialNumberFIFO(company, tenant int, callServer CallServerInfo, campaignId, uuid, fromNumber, trunkCode, phoneNumber, extention string) {
	fmt.Println("Start DialNumber: ", uuid, ": ", fromNumber, ": ", trunkCode, ": ", phoneNumber, ": ", extention)
	customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)

	param := fmt.Sprintf(" {DVP_CUSTOM_PUBID=%s,CampaignId=%s,CustomCompanyStr=%s,OperationType=Dialer,return_ring_ready=true,ignore_early_media=false,origination_uuid=%s,origination_caller_id_number=%s,originate_timeout=30}", subChannelName, campaignId, customCompanyStr, uuid, fromNumber)
	furl := fmt.Sprintf("sofia/gateway/%s/%s %s", trunkCode, phoneNumber, extention)
	data := " xml dialer"

	IncrConcurrentChannelCount(callServer.CallServerId, campaignId)
	InitiateSessionInfo(company, tenant, "1", campaignId, uuid, phoneNumber, "start", "start", time.Now().Format(layout4), callServer.CallServerId)
	IncrCampaignDialCount(company, tenant, campaignId)

	resp, err := Dial(callServer.Url, param, furl, data)
	if err != nil {
		DecrConcurrentChannelCount(callServer.CallServerId, campaignId)
		SetSessionInfo(campaignId, uuid, "Reason", "dial_failed")
		SetSessionInfo(campaignId, uuid, "DialerStatus", "failed")
		go UploadSessionInfo(campaignId, uuid)
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	if resp != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		tmx := string(response[:])
		fmt.Println(tmx)
		resultInfo := strings.Split(tmx, " ")
		if len(resultInfo) > 0 {
			if resultInfo[0] == "-ERR" {
				//DecrConcurrentChannelCount(callServer.CallServerId, campaignId)

				if len(resultInfo) > 1 {
					reason := resultInfo[1]
					if reason == "" {
						SetSessionInfo(campaignId, uuid, "Reason", "not_specified")
					} else {
						SetSessionInfo(campaignId, uuid, "Reason", reason)
					}
				} else {
					SetSessionInfo(campaignId, uuid, "Reason", "not_specified")
				}
				SetSessionInfo(campaignId, uuid, "DialerStatus", "not_connected")
				//go UploadSessionInfo(uuid)
			} else {
				SetSessionInfo(campaignId, uuid, "Reason", "dial_success")
				SetSessionInfo(campaignId, uuid, "DialerStatus", "connected")
			}
		}
	}
}
