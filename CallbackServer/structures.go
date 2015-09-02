package main

import (
	"time"
)

type Configuration struct {
	RedisIp                          string
	RedisPort                        string
	RedisDb                          int
	CallbackServerId                 string
	HostIpAddress                    string
	Port                             string
	ExternalCallbackRequestFrequency time.Duration
	CampaignService                  string
}

type EnvConfiguration struct {
	RedisIp                          string
	RedisPort                        string
	RedisDb                          string
	CallbackServerId                 string
	HostIpAddress                    string
	Port                             string
	ExternalCallbackRequestFrequency string
	CampaignService                  string
}

type CampaignCallback struct {
	Company     int
	Tenant      int
	Class       string
	Type        string
	Category    string
	DialoutTime time.Time
	CallbackUrl string
	CallbackObj string
}
