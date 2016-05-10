package main

import (
	"time"
)

const layout1 = "2006-01-02T15:04:05Z07:00"
const layout2 = "2006-01-02"
const layout4 = "2006-01-02T15:04:05.999999-07:00"

//--------------------Dialer Configurations--------------------
type Configuration struct {
	SecurityIp               string
	SecurityPort             string
	RedisIp                  string
	RedisPort                string
	RedisPassword            string
	RedisDb                  int
	DialerId                 string
	CampaignLimit            int
	LbIpAddress              string
	LbPort                   string
	Port                     string
	CampaignRequestFrequency time.Duration
	CampaignServiceHost      string
	CampaignServicePort      string
	CallServerHost           string
	CallServerPort           string
	CallRuleServiceHost      string
	CallRuleServicePort      string
	ScheduleServiceHost      string
	ScheduleServicePort      string
	CallbackServerHost       string
	CallbackServerPort       string
	ArdsServiceHost          string
	ArdsServicePort          string
	NotificationServiceHost  string
	NotificationServicePort  string
	ClusterConfigServiceHost string
	ClusterConfigServicePort string
	CasServerHost            string
	V5_1SecurityToken        string
	AccessToken              string
}

type EnvConfiguration struct {
	SecurityIp               string
	SecurityPort             string
	RedisIp                  string
	RedisPort                string
	RedisPassword            string
	RedisDb                  string
	DialerId                 string
	CampaignLimit            string
	LbIpAddress              string
	LbPort                   string
	Port                     string
	CampaignRequestFrequency string
	CampaignServiceHost      string
	CampaignServicePort      string
	CallServerHost           string
	CallServerPort           string
	CallRuleServiceHost      string
	CallRuleServicePort      string
	ScheduleServiceHost      string
	ScheduleServicePort      string
	CallbackServerHost       string
	CallbackServerPort       string
	ArdsServiceHost          string
	ArdsServicePort          string
	NotificationServiceHost  string
	NotificationServicePort  string
	ClusterConfigServiceHost string
	ClusterConfigServicePort string
	CasServerHost            string
	V5_1SecurityToken        string
	AccessToken              string
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
	CampaignId  string
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

type EmailAdditionalData struct {
	FromAddresss string
	Cc           string
	Body         string
	Subject      string
}

type CampaignAdditionalData struct {
	Class          string
	Type           string
	Category       string
	AdditionalData string
}

type CampaignAdditionalDataResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        CampaignAdditionalData
}

//--------------------Dialer--------------------
type DialerInfo struct {
	DialerId      string
	CampaignLimit int
	HostIpAddress string
}

//--------------------Call Server--------------------
type ResourceServerInfo struct {
	ResourceServerId string
	Url              string
	MaxChannelCount  int
}

//--------------------ClusterConfig API--------------------
type CallServerResult struct {
	Id             int
	Activate       bool
	Class          string
	Type           string
	Category       string
	InternalMainIP string
	CompanyId      int
}
type ClusterConfigApiResult struct {
	IsSuccess bool
	Result    []CallServerResult
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
	ServerType  string
	RequestType string
	CallbackUrl string
	ServerID    string
}

type RequestOtherData struct {
	CampaignId string
	StrData    string
	DialoutMec string
}

type Request struct {
	ServerType      string
	RequestType     string
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
	ContactName string
	Domain      string
	ContactType string
	ResourceId  string
}

type ArdsCallbackInfo struct {
	Company      string
	Tenant       string
	ServerType   string
	RequestType  string
	SessionID    string
	OtherInfo    string
	ResourceInfo ResourceDetails
}

type ArdsResSlot struct {
	Company     string
	Tenant      string
	ReqCategory string
	State       string
	OtherInfo   string
}

//--------------------Notification--------------------
type PushData struct {
	From      string
	To        string
	Direction string
	Message   string
	Callback  string
	Ref       string
}

type ReplyData struct {
	Tkey    string
	Message string
}

type ReceiveData struct {
	Reply ReplyData
	Ref   string
}

//--------------------Response-------------------------
type Result struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        string
}

//-------------------CAS SMS---------------------------
type Sms struct {
	ShortMessageInfo SmsInfo
	SecurityToken    string
}

type SmsInfo struct {
	Attachments       SmsAttachments
	Date              string
	DeliveryRefID     string
	DeliveryStatus    string
	FromPhoneNumber   string
	GUReferenceID     int
	GUTranID          int
	GUVersionID       int
	GatewayName       int
	ID                int
	MessageContent    string
	MessageRefID      string
	OperationalStatus string
	PhoneNumbers      []string
}

type SmsAttachments struct {
	FileAttachments      []string
	ReferenceAttachments []string
}

//-------------------CAS Email---------------------------
type Email struct {
	EmailInformation EmailInformation
	SecurityToken    string
}

type EmailInformation struct {
	Attachments       EmailAttachments
	CcEmailAddresses  []string
	CompanyID         int
	Content           string
	Date              string
	FileAttachments   []string
	GUTranID          int
	ID                int
	MessageRefID      string
	OperationalStatus string
	Subject           string
	ToEmailAddresses  []string
	URLAttachments    []string
}

type EmailAttachments struct {
	FileAttachments      []string
	ReferenceAttachments []string
}
