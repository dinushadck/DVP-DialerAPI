package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var dirPath string
var redisIp string
var redisPort string
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
var callbackServerSelfHost string
var ardsService string
var port string

func GetDirPath() string {
	envPath := os.Getenv("GO_CONFIG_DIR")
	if envPath == "" {
		envPath = "./"
	}
	fmt.Println(envPath)
	return envPath
}

func GetDefaultConfig() Configuration {
	confPath := filepath.Join(dirPath, "conf.json")
	fmt.Println("GetDefaultConfig config path: ", confPath)
	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		defconfiguration.RedisIp = "127.0.0.1"
		defconfiguration.RedisPort = "6379"
		defconfiguration.RedisDb = 6
		defconfiguration.DialerId = "1"
		defconfiguration.CampaignLimit = 30
		defconfiguration.HostIpAddress = "127.0.0.1"
		defconfiguration.Port = "2226"
		defconfiguration.CampaignRequestFrequency = 300
		defconfiguration.CampaignService = "http://127.0.0.1:2222/DVP/API/6.0"
		defconfiguration.UuidService = "http://127.0.0.1:8080/api/create_uuid"
		defconfiguration.CallServer = "127.0.0.1:8080"
		defconfiguration.CallRuleService = "http://127.0.0.1/CallRuleRestApi/api/CallRuleOutbound"
		defconfiguration.ScheduleService = "http://127.0.0.1:2224/DVP/API/6.0/LimitAPI"
		defconfiguration.CallbackServerSelfHost = "http://127.0.0.1:2227/CallbackServerSelfHost"
		defconfiguration.ArdsService = "http://192.168.0.15:2225/DVP/API/1.0.0.0/ARDS"
	}

	return defconfiguration
}

func LoadDefaultConfig() {
	confPath := filepath.Join(dirPath, "conf.json")
	fmt.Println("LoadDefaultConfig config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		redisIp = "127.0.0.1:6379"
		redisPort = "6379"
		redisDb = 6
		dialerId = "1"
		campaignLimit = 30
		hostIpAddress = "127.0.0.1"
		port = "2226"
		campaignRequestFrequency = 300
		campaignService = "http://127.0.0.1:2222/DVP/API/6.0"
		uuidService = "http://127.0.0.1:8080/api/create_uuid"
		callServer = "127.0.0.1:8080"
		callRuleService = "http://127.0.0.1/CallRuleRestApi/api/CallRuleOutbound"
		scheduleService = "http://127.0.0.1:2224/DVP/API/6.0/LimitAPI"
		callbackServerSelfHost = "http://127.0.0.1:2227/CallbackServerSelfHost"
		ardsService = "http://192.168.0.15:2225/DVP/API/1.0.0.0/ARDS"
	} else {
		redisIp = fmt.Sprintf("%s:%s", defconfiguration.RedisIp, defconfiguration.RedisPort)
		redisPort = defconfiguration.RedisPort
		redisDb = defconfiguration.RedisDb
		dialerId = defconfiguration.DialerId
		campaignLimit = defconfiguration.CampaignLimit
		hostIpAddress = defconfiguration.HostIpAddress
		port = defconfiguration.Port
		campaignRequestFrequency = defconfiguration.CampaignRequestFrequency
		campaignService = defconfiguration.CampaignService
		uuidService = defconfiguration.UuidService
		callServer = defconfiguration.CallServer
		callRuleService = defconfiguration.CallRuleService
		scheduleService = defconfiguration.ScheduleService
		callbackServerSelfHost = defconfiguration.CallbackServerSelfHost
		ardsService = defconfiguration.ArdsService
	}
}

func LoadConfiguration() {
	dirPath = GetDirPath()
	confPath := filepath.Join(dirPath, "custom-environment-variables.json")
	fmt.Println("InitiateRedis config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	envconfiguration := EnvConfiguration{}
	enverr := json.Unmarshal(content, &envconfiguration)
	if enverr != nil {
		fmt.Println("error:", enverr)
		LoadDefaultConfig()
	} else {
		var converr error
		defConfig := GetDefaultConfig()
		redisIp = os.Getenv(envconfiguration.RedisIp)
		redisPort = os.Getenv(envconfiguration.RedisPort)
		redisDb, converr = strconv.Atoi(os.Getenv(envconfiguration.RedisDb))
		dialerId = os.Getenv(envconfiguration.DialerId)
		campaignLimit, converr = strconv.Atoi(os.Getenv(envconfiguration.CampaignLimit))
		hostIpAddress = os.Getenv(envconfiguration.HostIpAddress)
		port = os.Getenv(envconfiguration.Port)
		campaignRequestFrequencytemp := os.Getenv(envconfiguration.CampaignRequestFrequency)
		campaignService = os.Getenv(envconfiguration.CampaignService)
		uuidService = os.Getenv(envconfiguration.UuidService)
		callServer = os.Getenv(envconfiguration.CallServer)
		callRuleService = os.Getenv(envconfiguration.CallRuleService)
		scheduleService = os.Getenv(envconfiguration.ScheduleService)
		callbackServerSelfHost = os.Getenv(envconfiguration.CallbackServerSelfHost)
		ardsService = os.Getenv(envconfiguration.ArdsService)

		if redisIp == "" {
			redisIp = defConfig.RedisIp
		}
		if redisPort == "" {
			redisPort = defConfig.RedisPort
		}
		if redisDb == 0 || converr != nil {
			redisDb = defConfig.RedisDb
		}
		if dialerId == "" {
			dialerId = defConfig.DialerId
		}
		if campaignLimit == 0 || converr != nil {
			campaignLimit = defConfig.CampaignLimit
		}
		if hostIpAddress == "" {
			hostIpAddress = defConfig.HostIpAddress
		}
		if port == "" {
			port = defConfig.Port
		}
		if campaignRequestFrequencytemp == "" {
			campaignRequestFrequency = defConfig.CampaignRequestFrequency
		} else {
			campaignRequestFrequency, _ = time.ParseDuration(campaignRequestFrequencytemp)
		}
		if campaignService == "" {
			campaignService = defConfig.CampaignService
		}
		if uuidService == "" {
			uuidService = defConfig.UuidService
		}
		if callServer == "" {
			callServer = defConfig.CallServer
		}
		if callRuleService == "" {
			callRuleService = defConfig.CallRuleService
		}
		if scheduleService == "" {
			scheduleService = defConfig.ScheduleService
		}
		if callbackServerSelfHost == "" {
			callbackServerSelfHost = defConfig.CallbackServerSelfHost
		}
		if ardsService == "" {
			ardsService = defConfig.ArdsService
		}

		redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
	}

	fmt.Println("redisIp:", redisIp)
	fmt.Println("redisDb:", redisDb)
	fmt.Println("dialerId:", dialerId)
	fmt.Println("campaignLimit:", campaignLimit)
}

func LoadCallbackConfiguration() {
	dirPath = GetDirPath()
	confPath := filepath.Join(dirPath, "callbackConf.json")
	fmt.Println("InitiateCallback config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	callbackConf := CallbackConfiguration{}
	err := json.Unmarshal(content, &callbackConf)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		for _, conf := range callbackConf.DisconnectReasons {
			for _, reason := range conf.Values {
				confKey := fmt.Sprintf("CallbackReason:%s", reason)
				RedisSetNx(confKey, conf.Reason)
			}
		}
	}
}
