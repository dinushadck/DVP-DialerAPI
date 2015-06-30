package main

import (
	"encoding/json"
	"fmt"
)

func InitiateDuoDialer() {
	LoadConfiguration()
	AddDialerInfoToRedis()
}

func AddDialerInfoToRedis() {
	dialerInfo := DialerInfo{}
	dialerInfo.DialerId = dialerId
	dialerInfo.HostIpAddress = hostIpAddress
	dialerInfo.CampaignLimit = campaignLimit

	dialerKey := fmt.Sprintf("DialerInfo:%s:%s", hostIpAddress, dialerId)
	dialerInfoJson, _ := json.Marshal(dialerInfo)
	result := RedisAdd(dialerKey, string(dialerInfoJson))
	fmt.Println("Add DialerInfo to Redis: ", result)
	if result == "OK" {
		AddOnGoingCampaignCount()
	}
}
