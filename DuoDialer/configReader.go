package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var redisIp string
var redisDb int
var dialerId string
var campaignLimit int
var hostIpAddress string
var campaignService string
var campaignRequestFrequency time.Duration
var uuidService string
var callServer string
var callRuleService string
var scheduleService string

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
		uuidService = "http://192.168.2.101:8080/api/create_uuid"
		callServer = "192.168.2.101:8080"
		callRuleService = "http://ip:port/DVP/API/:version/CallRule/Outbound/"
		scheduleService = "http://192.168.0.51:8083/DVP/API/6.0/LimitAPI"
	} else {
		redisIp = configuration.RedisIp
		redisDb = configuration.RedisDb
		dialerId = configuration.DialerId
		campaignLimit = configuration.CampaignLimit
		hostIpAddress = configuration.HostIpAddress
		campaignRequestFrequency = configuration.CampaignRequestFrequency
		campaignService = configuration.CampaignService
		uuidService = configuration.UuidService
		callServer = configuration.CallServer
		callRuleService = configuration.CallRuleService
		scheduleService = configuration.ScheduleService
	}

	fmt.Println("redisIp:", redisIp)
	fmt.Println("redisDb:", redisDb)
	fmt.Println("dialerId:", dialerId)
	fmt.Println("campaignLimit:", campaignLimit)
}
