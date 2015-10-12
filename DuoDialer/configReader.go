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
var port string
var campaignRequestFrequency time.Duration
var campaignServiceHost string
var campaignServicePort string
var callServerHost string
var callServerPort string
var callRuleServiceHost string
var callRuleServicePort string
var scheduleServiceHost string
var scheduleServicePort string
var callbackServerHost string
var callbackServerPort string
var ardsServiceHost string
var ardsServicePort string

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
		defconfiguration.RedisIp = "192.168.3.200"
		defconfiguration.RedisPort = "6379"
		defconfiguration.RedisDb = 5
		defconfiguration.DialerId = "Dialer2"
		defconfiguration.CampaignLimit = 30
		defconfiguration.HostIpAddress = "192.168.0.15"
		defconfiguration.Port = "2226"
		defconfiguration.CampaignRequestFrequency = 300
		defconfiguration.CampaignServiceHost = "192.168.0.143"
		defconfiguration.CampaignServicePort = "2222"
		defconfiguration.CallServerHost = "192.168.0.53"
		defconfiguration.CallServerPort = "8080"
		defconfiguration.CallRuleServiceHost = "192.168.0.89"
		defconfiguration.CallRuleServicePort = "2220"
		defconfiguration.ScheduleServiceHost = "192.168.3.200"
		defconfiguration.ScheduleServicePort = "2224"
		defconfiguration.CallbackServerHost = "192.168.0.15"
		defconfiguration.CallbackServerPort = "2227"
		defconfiguration.ArdsServiceHost = "192.168.0.15"
		defconfiguration.ArdsServicePort = "2225"
	}

	return defconfiguration
}

func LoadDefaultConfig() {

	defconfiguration := GetDefaultConfig()

	redisIp = defconfiguration.RedisIp
	redisPort = defconfiguration.RedisPort
	redisDb = defconfiguration.RedisDb
	dialerId = defconfiguration.DialerId
	campaignLimit = defconfiguration.CampaignLimit
	hostIpAddress = defconfiguration.HostIpAddress
	port = defconfiguration.Port
	campaignRequestFrequency = defconfiguration.CampaignRequestFrequency
	campaignServiceHost = defconfiguration.CampaignServiceHost
	campaignServicePort = defconfiguration.CampaignServicePort
	callServerHost = defconfiguration.CallServerHost
	callServerPort = defconfiguration.CallServerPort
	callRuleServiceHost = defconfiguration.CallRuleServiceHost
	callRuleServicePort = defconfiguration.CallRuleServicePort
	scheduleServiceHost = defconfiguration.ScheduleServiceHost
	scheduleServicePort = defconfiguration.ScheduleServicePort
	callbackServerHost = defconfiguration.CallbackServerHost
	callbackServerPort = defconfiguration.CallbackServerPort
	ardsServiceHost = defconfiguration.ArdsServiceHost
	ardsServicePort = defconfiguration.ArdsServicePort

	redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
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
		redisPort = os.Getenv(envconfiguration.RedisIp)
		redisDb, converr = strconv.Atoi(os.Getenv(envconfiguration.RedisDb))
		dialerId = os.Getenv(envconfiguration.DialerId)
		campaignLimit, converr = strconv.Atoi(os.Getenv(envconfiguration.CampaignLimit))
		hostIpAddress = os.Getenv(envconfiguration.HostIpAddress)
		port = os.Getenv(envconfiguration.Port)
		campaignRequestFrequencytemp := os.Getenv(envconfiguration.CampaignRequestFrequency)
		campaignServiceHost = os.Getenv(envconfiguration.CampaignServiceHost)
		campaignServicePort = os.Getenv(envconfiguration.CampaignServicePort)
		callServerHost = os.Getenv(envconfiguration.CallServerHost)
		callServerPort = os.Getenv(envconfiguration.CallServerPort)
		callRuleServiceHost = os.Getenv(envconfiguration.CallRuleServiceHost)
		callRuleServicePort = os.Getenv(envconfiguration.CallRuleServicePort)
		scheduleServiceHost = os.Getenv(envconfiguration.ScheduleServiceHost)
		scheduleServicePort = os.Getenv(envconfiguration.ScheduleServicePort)
		callbackServerHost = os.Getenv(envconfiguration.CallbackServerHost)
		callbackServerPort = os.Getenv(envconfiguration.CallbackServerPort)
		ardsServiceHost = os.Getenv(envconfiguration.ArdsServiceHost)
		ardsServicePort = os.Getenv(envconfiguration.ArdsServiceHost)

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
		if campaignServiceHost == "" {
			campaignServiceHost = defConfig.CampaignServiceHost
		}
		if campaignServicePort == "" {
			campaignServicePort = defConfig.CampaignServicePort
		}
		if callServerHost == "" {
			callServerHost = defConfig.CallServerHost
		}
		if callServerPort == "" {
			callServerPort = defConfig.CallServerPort
		}
		if callRuleServiceHost == "" {
			callRuleServiceHost = defConfig.CallRuleServiceHost
		}
		if callRuleServicePort == "" {
			callRuleServicePort = defConfig.CallRuleServicePort
		}
		if scheduleServiceHost == "" {
			scheduleServiceHost = defConfig.ScheduleServiceHost
		}
		if scheduleServicePort == "" {
			scheduleServicePort = defConfig.ScheduleServicePort
		}
		if callbackServerHost == "" {
			callbackServerHost = defConfig.CallbackServerHost
		}
		if callbackServerPort == "" {
			callbackServerPort = defConfig.CallbackServerPort
		}
		if ardsServiceHost == "" {
			ardsServiceHost = defConfig.ArdsServiceHost
		}
		if ardsServicePort == "" {
			ardsServicePort = defConfig.ArdsServicePort
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
