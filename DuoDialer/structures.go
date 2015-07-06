package main

import (
	"time"
)

type Configuration struct {
	RedisIp                  string
	RedisDb                  int
	DialerId                 string
	CampaignLimit            int
	HostIpAddress            string
	CampaignRequestFrequency time.Duration
	CampaignService          string
	UuidService              string
	CallServer               string
	CallRuleService          string
	ScheduleService          string
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
	Extention  string
	StartDate  string
	EndDate    string
}

type CallRuleApiResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        CallRule
}

type CallRule struct {
	GatewayCode string
	DNIS        string
	ANI         string
}

type ScheduleDetails struct {
	CustomMessage string
	IsSuccess     bool
	Result        []Appoinment
}

type Appoinment struct {
	id              int
	AppointmentName string
	Action          string
	ExtraData       string
	StartDate       string
	EndDate         string
	StartTime       string
	EndTime         string
	DaysOfWeek      string
	ObjClass        string
	ObjType         string
	ObjCategory     string
	CompanyId       int
	TenantId        int
	createdAt       string
	updatedAt       string
	ScheduleId      int
}
