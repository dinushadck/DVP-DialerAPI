package main

import (
	"encoding/json"
	"fmt"
)

func InitiateDuoDialer() {
	LoadConfiguration()
	LoadCallbackConfiguration()
	AddDialerInfoToRedis()
	go PubSub()
}

func AddDialerInfoToRedis() {
	dialerInfo := DialerInfo{}
	dialerInfo.DialerId = dialerId
	dialerInfo.HostIpAddress = lbIpAddress
	dialerInfo.CampaignLimit = campaignLimit

	dialerKey := fmt.Sprintf("DialerInfo:%s:%s", lbIpAddress, dialerId)
	dialerInfoJson, _ := json.Marshal(dialerInfo)
	result := RedisAdd(dialerKey, string(dialerInfoJson))
	fmt.Println("Add DialerInfo to Redis: ", result)
	if result == "OK" {
		AddOnGoingCampaignCount()
	}
}
