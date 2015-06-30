package main

type Configuration struct {
	RedisIp                  string
	RedisDb                  int
	DialerId                 string
	CampaignLimit            int
	HostIpAddress            string
	CampaignRequestFrequency int
	CampaignService          string
}

type DialerInfo struct {
	DialerId      string
	CampaignLimit int
	HostIpAddress string
}

type Campaign struct {
	CampaignId string
	Company    int
	Tenant     int
	Calss      string
	Type       string
	Category   string
}
