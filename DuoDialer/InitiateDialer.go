package main

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
)

func InitiateDuoDialer() {
	//Get Environment variables and assign to variables
	LoadConfiguration()
	//Set up redis client
	InitiateRedis()
	//Get callback configurations from campaignmanager service
	LoadCallbackConfiguration()
	//Get disconnect reason
	GetDisconnectReasons()
	//Add dialer record to redis with dialer name set on env vars
	AddDialerInfoToRedis()
	//Create new go routing to listen to eventmonitor events EG:- CHANNEL_CREATE, CHANNEL_ANSWER
	go PubSub()
	go PubSubAgentChan()
}

func AddDialerInfoToRedis() {
	dialerInfo := DialerInfo{}
	dialerInfo.DialerId = dialerId
	dialerInfo.HostIpAddress = lbIpAddress
	dialerInfo.CampaignLimit = campaignLimit

	dialerKey := fmt.Sprintf("DialerInfo:%s:%s", lbIpAddress, dialerId)
	dialerInfoJson, _ := json.Marshal(dialerInfo)
	result := RedisAdd(dialerKey, string(dialerInfoJson))
	color.Green(fmt.Sprintf("Add DialerInfo to Redis: %s", result))
	if result == "OK" {
		AddOnGoingCampaignCount()
	}
}
