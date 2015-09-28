package main

import (
	"time"
)

const layout1 = "2006-01-02T15:04:05Z07:00"
const layout2 = "2006-01-02"
const layout3 = "15:04"
const layout4 = "2006-01-02T15:04:05.999999-07:00"

//--------------------Dialer Configurations--------------------
type Configuration struct {
	RedisIp                  string
	RedisPort                string
	RedisDb                  int
	DialerId                 string
	CampaignLimit            int
	HostIpAddress            string
	Port                     string
	CampaignRequestFrequency time.Duration
	CampaignService          string
	UuidService              string
	CallServer               string
	CallRuleService          string
	ScheduleService          string
	CallbackServerSelfHost   string
	ArdsService              string
}

type EnvConfiguration struct {
	RedisIp                  string
	RedisPort                string
	RedisDb                  string
	DialerId                 string
	CampaignLimit            string
	HostIpAddress            string
	Port                     string
	CampaignRequestFrequency string
	CampaignService          string
	UuidService              string
	CallServer               string
	CallRuleService          string
	ScheduleService          string
	CallbackServerSelfHost   string
	ArdsService              string
}

//--------------------Campaign--------------------
type Campaign struct {
	CampaignName       string
	CampaignMode       string
	CampaignChannel    string
	DialoutMechanism   string
	CampaignId         int
	CompanyId          int
	TenantId           int
	Class              string
	Type               string
	Category           string
	Extensions         string
	OperationalStatus  string
	CampScheduleInfo   []CampaignShedule
	CampConfigurations CampaignConfigInfo
}

type CampaignResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        []Campaign
}

type CampaignCallbackInfoResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        []CampaignCallbackConfigInfo
}

type CampaignCallbackReason struct {
	Reason string
}

type CampaignCallbackInfo struct {
	CallBackConfId      int
	MaxCallBackCount    int
	CallBackInterval    int
	ReasonId            string
	CampCallBackReasons CampaignCallbackReason
}

type CampaignCallbackConfigInfo struct {
	AllowCallBack              bool
	CampCallbackConfigurations []CampaignCallbackInfo
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

type CampaignCallbackObj struct {
	CampaignId       int
	CallbackClass    string
	CallbackType     string
	CallbackCategory string
	ContactId        string
	CallBackCount    int
	DialoutTime      time.Time
}

type CampaignConfigInfo struct {
	ConfigureId        int
	ChannelConcurrency int
	AllowCallBack      bool
	Caller             string
	StartDate          string
	EndDate            string
}

type CampaignShedule struct {
	ScheduleId    int
	CamScheduleId int
}

type ContactInfo struct {
	ContactId string
}

type CampaignContactInfo struct {
	ExtraData       string
	CampContactInfo ContactInfo
}

type PhoneNumberResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        []CampaignContactInfo
}

type CallbackConfiguration struct {
	DisconnectReasons []DisconnectReason
}

type DisconnectReason struct {
	Reason string
	Values []string
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

//--------------------Dialer--------------------
type DialerInfo struct {
	DialerId      string
	CampaignLimit int
	HostIpAddress string
}

//--------------------Call Server--------------------
type CallServerInfo struct {
	CallServerId    string
	Url             string
	MaxChannelCount int
}

//--------------------Rule API--------------------
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

//--------------------Limit API--------------------
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

//--------------------Events--------------------
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

type PubEvents struct {
	EventClass    string
	EventType     string
	EventCategory string
	EventName     string
	EventData     string
	EventParams   string
	SessionId     string
	EventTime     string
	CompanyId     int
	TenantId      int
}

//--------------------Ards--------------------
type RequestServer struct {
	Class       string
	Type        string
	Category    string
	CallbackUrl string
	ServerID    string
}

type Request struct {
	Class           string
	Type            string
	Category        string
	SessionId       string
	Attributes      []string
	RequestServerId string
	Priority        string
	OtherInfo       string
}

type ArdsResult struct {
	CustomMessage string
	IsSuccess     bool
}

type ResourceDetails struct {
	Extention    int
	DialHostName string
}

type ArdsCallback struct {
	SessionID    string
	OtherInfo    string
	ResourceInfo ResourceDetails
}
