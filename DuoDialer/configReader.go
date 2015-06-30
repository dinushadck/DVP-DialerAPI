package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var redisIp string
var redisDb int
var dialerId string
var campaignLimit int
var hostIpAddress string
var campaignService string
var campaignRequestFrequency int

func LoadConfiguration() {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Load Config error:", err)
		redisIp = "127.0.0.1:6379"
		redisDb = 5
		dialerId = "1"
		campaignLimit = 1
		hostIpAddress = "127.0.0.1"
		campaignRequestFrequency = 300
		campaignService = "http://localhost:7777/campaign"
	} else {
		redisIp = configuration.RedisIp
		redisDb = configuration.RedisDb
		dialerId = configuration.DialerId
		campaignLimit = configuration.CampaignLimit
		hostIpAddress = configuration.HostIpAddress
		campaignRequestFrequency = configuration.CampaignRequestFrequency
		campaignService = configuration.CampaignService
	}

	fmt.Println("redisIp:", redisIp)
	fmt.Println("redisDb:", redisDb)
	fmt.Println("dialerId:", dialerId)
	fmt.Println("campaignLimit:", campaignLimit)
}
