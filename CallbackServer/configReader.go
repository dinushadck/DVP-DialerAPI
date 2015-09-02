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
var callbackServerId string
var hostIpAddress string
var port string
var externalCallbackRequestFrequency time.Duration
var campaignService string

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
		defconfiguration.CallbackServerId = "1"
		defconfiguration.HostIpAddress = "127.0.0.1"
		defconfiguration.Port = "2226"
		defconfiguration.ExternalCallbackRequestFrequency = 300
		defconfiguration.CampaignService = "http://127.0.0.1:2222/DVP/API/6.0"
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
		callbackServerId = "1"
		hostIpAddress = "127.0.0.1"
		port = "2226"
		externalCallbackRequestFrequency = 300
		campaignService = "http://127.0.0.1:2222/DVP/API/6.0"
	} else {
		redisIp = fmt.Sprintf("%s:%s", defconfiguration.RedisIp, defconfiguration.RedisPort)
		redisPort = defconfiguration.RedisPort
		redisDb = defconfiguration.RedisDb
		callbackServerId = defconfiguration.CallbackServerId
		hostIpAddress = defconfiguration.HostIpAddress
		port = defconfiguration.Port
		externalCallbackRequestFrequency = defconfiguration.ExternalCallbackRequestFrequency
		campaignService = defconfiguration.CampaignService
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
		callbackServerId = os.Getenv(envconfiguration.CallbackServerId)
		hostIpAddress = os.Getenv(envconfiguration.HostIpAddress)
		port = os.Getenv(envconfiguration.Port)
		externalCallbackRequestFrequencyTemp := os.Getenv(envconfiguration.ExternalCallbackRequestFrequency)
		campaignService = os.Getenv(envconfiguration.CampaignService)

		if redisIp == "" {
			redisIp = defConfig.RedisIp
		}
		if redisPort == "" {
			redisPort = defConfig.RedisPort
		}
		if redisDb == 0 || converr != nil {
			redisDb = defConfig.RedisDb
		}
		if callbackServerId == "" {
			callbackServerId = defConfig.CallbackServerId
		}
		if hostIpAddress == "" {
			hostIpAddress = defConfig.HostIpAddress
		}
		if port == "" {
			port = defConfig.Port
		}
		if externalCallbackRequestFrequencyTemp == "" {
			externalCallbackRequestFrequency = defConfig.ExternalCallbackRequestFrequency
		} else {
			externalCallbackRequestFrequency, _ = time.ParseDuration(externalCallbackRequestFrequencyTemp)
		}
		if campaignService == "" {
			campaignService = defConfig.CampaignService
		}

		redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
	}

	fmt.Println("redisIp:", redisIp)
	fmt.Println("redisDb:", redisDb)
}
