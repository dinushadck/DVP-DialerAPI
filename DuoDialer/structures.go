package main

import (
	"time"
)

const layout1 = "2006-01-02T15:04:05Z07:00"
const layout2 = "2006-01-02"
const layout3 = "15:04"

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

type CallServerInfo struct {
	CallServerId    string
	Url             string
	MaxChannelCount int
}

type CampaignResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        []Campaign
}

type CampaignConfigInfo struct {
	ChannelConcurrency int
	AllowCallBack      bool
	MaxCallBackCount   int
	Caller             string
	StartDate          string
	EndDate            string
}

type CampaignShedule struct {
	ScheduleId    int
	CamScheduleId int
}

type Campaign struct {
	CampaignName       string
	CampaignMode       string
	CampaignChannel    string
	DialoutMechanism   string
	CampaignId         int
	CompanyId          int
	TenantId           int
	Calss              string
	Type               string
	Category           string
	Extensions         string
	OperationalStatus  string
	CampScheduleInfo   []CampaignShedule
	CampConfigurations CampaignConfigInfo
}

type ContactInfo struct {
	ContactId string
}

type CampaignContactInfo struct {
	CampContactInfo ContactInfo
}

type PhoneNumberResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        []CampaignContactInfo
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

type SubEvents struct {
	SwitchName       string
	CampaignId       string
	SessionId        string
	EventClass       string
	EventType        string
	EventTime        string
	EventName        string
	EventData        string
	AuthData         string
	EventCategory    string
	DisconnectReason string
}

type CampaignStart struct {
	CampaignId int
	DialerId   string
}

type CampaignState struct {
	CampaignId    int
	CampaignState string
	DialerId      string
}
type CampaignStatusResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        CampaignState
}
