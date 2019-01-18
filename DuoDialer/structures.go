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
	ContactServiceHost       string
	ContactServicePort       string
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
	RabbitMQHost             string
	RabbitMQPort             string
	RabbitMQUser             string
	RabbitMQPassword         string
	FileServiceHost          string
	FileServicePort          string
	RedisMode                string
	RedisClusterName         string
	SentinelHosts            string
	SentinelPort             string
	DvpEventType             string
	UseAmqpAdapter           string
	AmqpAdapterPort          string
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
	ContactServiceHost       string
	ContactServicePort       string
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
	RabbitMQHost             string
	RabbitMQPort             string
	RabbitMQUser             string
	RabbitMQPassword         string
	FileServiceHost          string
	FileServicePort          string
	RedisMode                string
	RedisClusterName         string
	SentinelHosts            string
	SentinelPort             string
	DvpEventType             string
	UseAmqpAdapter           string
	AmqpAdapterPort          string
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
	CampAdditionalData CampaignAdditionalData
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
	ConfigureId         int
	MaxCallBackCount    int
	CallBackInterval    int
	ReasonId            int
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
	CampaignId       string
	CallbackClass    string
	CallbackType     string
	CallbackCategory string
	ContactId        string
	CallBackCount    string
	DialoutTime      time.Time
	OtherContacts    []Contact
	PreviewData      string
}

type CallbackInfo map[string]interface{}

type CampaignConfigInfo struct {
	ConfigureId         int
	ChannelConcurrency  int
	AllowCallBack       bool
	Caller              string
	StartDate           time.Time
	EndDate             time.Time
	StartTimeZone       string
	EndTimeZone         string
	IntegrationData     IntegrationConfig
	NumberLoadingMethod string
}

type IntegrationConfig struct {
	Agent    IntegrationInfo
	Customer IntegrationInfo
}

type IntegrationInfo struct {
	Url    string
	Method string
	Params []string
}

type CampaignShedule struct {
	ScheduleId    int
	CamScheduleId int
	TimeZone      string
	StartDate     time.Time
	EndDate       time.Time
}

type ContactInfo struct {
	ContactId string
}

type CampaignContactInfo struct {
	ExtraData       string
	CampContactInfo ContactInfo
}

type DialResult struct {
	IsSuccess bool
	Message   string
}

type PhoneNumberResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        []CampaignContactInfo
}

type Contact struct {
	Contact  string
	Display  string
	Verified bool
}

type ContactsDetails struct {
	Phone        string
	PreviewData  string
	Api_Contacts []Contact
}

type ContactsResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        []ContactsDetails
}

type DncNumberResult struct {
	CustomMessage string
	IsSuccess     bool
	Result        []string
}

type CallbackConfiguration struct {
	IsSuccess bool
	Result    []DisconnectReason
}

type DisconnectReason struct {
	Reason      string
	HangupCause []string
	Status      bool
	ReasonId    int
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
	Result        []CampaignAdditionalData
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
	MainIp         string
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
	IpUrl       string
}

//--------------------Limit API--------------------
type ScheduleDetails struct {
	CustomMessage string
	IsSuccess     bool
	Result        []Schedule
}

type Schedule struct {
	ScheduleName string
	TimeZone     string
	StartDate    string
	EndDate      string
	Appointment  []Appoinment
}

type Appoinment struct {
	id                int
	AppointmentName   string
	Action            string
	ExtraData         string
	StartDate         string
	EndDate           string
	StartTime         string
	EndTime           string
	RecurrencePattern string
	DaysOfWeek        string
	ObjClass          string
	ObjType           string
	ObjCategory       string
	CompanyId         int
	TenantId          int
	createdAt         string
	updatedAt         string
	ScheduleId        int
}

//--------------------Events--------------------
type SubEvents struct {
	TenantId         string
	CompanyId        string
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
	ServerType     string
	RequestType    string
	CallbackUrl    string
	ServerID       string
	CallbackOption string
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
	ContactName  string
	Domain       string
	ContactType  string
	ResourceId   string
	ResourceName string
	Extention    string
}

type ArdsCallbackInfo struct {
	Company      string
	Tenant       string
	ServerType   string
	RequestType  string
	SessionID    string
	OtherInfo    string
	Skills       string
	ResourceInfo ResourceDetails
}

type ArdsResSlot struct {
	Company     string
	Tenant      string
	ReqCategory string
	State       string
	OtherInfo   string
}

type ArdsResource struct {
	ServerType   string
	RequestType  string
	State        string
	OtherInfo    string
	Reason       string
	Company      string
	Tenant       string
	BusinessUnit string
}

//--------------------Notification--------------------
type PushData struct {
	From        string
	To          string
	Direction   string
	Message     string
	CallbackURL string
	Ref         string
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

//-------------------SmsAndEmail---------------------------
type SmsAndEmail struct {
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
