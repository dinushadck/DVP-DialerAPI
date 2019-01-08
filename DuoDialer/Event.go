package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/DuoSoftware/log4go"
	"github.com/fatih/color"
)

var eventLog = log4go.NewLogger()

func OnEvent(eventInfo SubEvents) {
	defEvent := SubEvents{}
	if eventInfo != defEvent {
		fmt.Println("AuthData: ", eventInfo.AuthData)
		fmt.Println("SessionId: ", eventInfo.SessionId)
		fmt.Println("EventName: ", eventInfo.EventName)
		fmt.Println("EventCategory: ", eventInfo.EventCategory)
		fmt.Println("CampaignId: ", eventInfo.CampaignId)
		fmt.Println("Tenant: ", eventInfo.TenantId)
		fmt.Println("Company: ", eventInfo.CompanyId)

		//authInfoArr := strings.Split(eventInfo.AuthData, "_")
		if eventInfo.TenantId != "" && eventInfo.CompanyId != "" {
			company, _ := strconv.Atoi(eventInfo.CompanyId)
			tenant, _ := strconv.Atoi(eventInfo.TenantId)

			fmt.Println("company: ", company)
			fmt.Println("tenant: ", tenant)

			switch eventInfo.EventCategory {
			case "CHANNEL_BRIDGE":
				fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
				break
			case "CHANNEL_CREATE":
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "ChannelCreatetime", time.Now().Format(layout4))
				color.Magenta(fmt.Sprintf("EventName: %s, SessionId: %s, EventCat: %s", eventInfo.EventName, eventInfo.SessionId, eventInfo.EventCategory))
				break
			case "CHANNEL_ANSWER":
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "DialerStatus", "channel_answered")
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "ChannelAnswertime", time.Now().Format(layout4))
				IncrCampaignConnectedCount(company, tenant, eventInfo.CampaignId)
				color.Magenta(fmt.Sprintf("EventName: %s, SessionId: %s, EventCat: %s", eventInfo.EventName, eventInfo.SessionId, eventInfo.EventCategory))
				break
			case "CHANNEL_DESTROY":
				//LogEvent(eventInfo)
				color.Magenta(fmt.Sprintf("EventName: %s, SessionId: %s, EventCat: %s, DisconnectReason : %s", eventInfo.EventName, eventInfo.SessionId, eventInfo.EventCategory, eventInfo.DisconnectReason))
				hashKey := fmt.Sprintf("sessionInfo:%s:%s", eventInfo.CampaignId, eventInfo.SessionId)
				session := RedisCheckKeyExist(hashKey)
				if session {
					color.Magenta("==========Session Found============")
					DecrConcurrentChannelCount(eventInfo.SwitchName, eventInfo.CampaignId)
					SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "Reason", eventInfo.DisconnectReason)

					hKey := fmt.Sprintf("sessionInfo:%s:%s", eventInfo.CampaignId, eventInfo.SessionId)
					sessionInfo := RedisHashGetAll(hKey)

					color.Magenta("=============DISCONNECT=============")
					color.Magenta(fmt.Sprintf(sessionInfo["IntegrationData"]))

					if sessionInfo != nil && sessionInfo["IntegrationData"] != "" {
						sessionInfo["EventType"] = "CUST	OMER_DISCONNECT"
						go ManageIntegrationData(sessionInfo, "CUSTOMER")
					} else {
						color.Magenta("NO INTEGRATION DATA")
					}

					go UploadSessionInfo(eventInfo.CampaignId, eventInfo.SessionId)
					//fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
				} else {
					color.Magenta("==========Session Not Found : " + hashKey)
				}
				break
			default:
				break
			}
		} else {
			fmt.Println("Auth error")
			fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
		}
	} else {
		fmt.Println("Empty Event")
	}
}

func OnEventAgent(eventInfo SubEvents) {
	redGreen := color.New(color.FgRed).Add(color.BgGreen)
	defEventAgent := SubEvents{}
	if eventInfo != defEventAgent {
		if eventInfo.TenantId != "" && eventInfo.CompanyId != "" {
			switch eventInfo.EventName {
			case "CHANNEL_CREATE":
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "ChannelCreatetime", time.Now().Format(layout4))
				redGreen.Println(fmt.Sprintf("EventName: %s, SessionId: %s", eventInfo.EventName, eventInfo.SessionId))
				hKey := fmt.Sprintf("sessionInfo:%s:%s", eventInfo.CampaignId, eventInfo.SessionId)
				sessionInfo := RedisHashGetAll(hKey)

				if sessionInfo != nil && sessionInfo["IntegrationData"] != "" {
					sessionInfo["EventType"] = "AGENT_RINGING"
					go ManageIntegrationData(sessionInfo, "AGENT")
				}
				break
			case "CHANNEL_ANSWER":
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "DialerStatus", "channel_answered")
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "ChannelAnswertime", time.Now().Format(layout4))
				redGreen.Println(fmt.Sprintf("EventName: %s, SessionId: %s", eventInfo.EventName, eventInfo.SessionId))
				hKey := fmt.Sprintf("sessionInfo:%s:%s", eventInfo.CampaignId, eventInfo.SessionId)
				sessionInfo := RedisHashGetAll(hKey)

				if sessionInfo != nil && sessionInfo["IntegrationData"] != "" {
					sessionInfo["EventType"] = "AGENT_ANSWERED"
					go ManageIntegrationData(sessionInfo, "AGENT")
				}
				break
			case "CHANNEL_HANGUP":
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "Reason", eventInfo.DisconnectReason)
				redGreen.Println(fmt.Sprintf("EventName: %s, SessionId: %s", eventInfo.EventName, eventInfo.SessionId))
				hKey := fmt.Sprintf("sessionInfo:%s:%s", eventInfo.CampaignId, eventInfo.SessionId)
				sessionInfo := RedisHashGetAll(hKey)

				if sessionInfo != nil && sessionInfo["IntegrationData"] != "" {
					sessionInfo["EventType"] = "AGENT_DISCONNECTED"
					go ManageIntegrationData(sessionInfo, "AGENT")
				}

				break
			default:
				break
			}
		} else {
			fmt.Println("Auth error")
			fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
		}
	} else {
		fmt.Println("Empty Event")
	}
}

func LogEvent(eventInfo SubEvents) {
	eventLog.AddFilter("file", log4go.FINE, log4go.NewFileLogWriter("EventLog.txt", false))
	logData, _ := json.Marshal(eventInfo)

	eventLog.Info("------------------------------------------\n")
	eventLog.Info(string(logData), "\n")
	eventLog.Close()
}

func PublishEvent(campaignId, sessionId string) {
	sessionInfoKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
	if RedisCheckKeyExist(sessionInfoKey) {
		sessionInfo := RedisHashGetAll(sessionInfoKey)

		tenant, _ := strconv.Atoi(sessionInfo["TenantId"])
		company, _ := strconv.Atoi(sessionInfo["CompanyId"])

		pubEventData := PubEvents{}

		pubEventData.SessionId = sessionId
		pubEventData.TenantId = tenant
		pubEventData.CompanyId = company
		pubEventData.EventClass = sessionInfo["Class"]
		pubEventData.EventType = sessionInfo["Type"]
		pubEventData.EventCategory = sessionInfo["Category"]
		pubEventData.EventName = "DialInfo"
		pubEventData.EventData = sessionInfo["DialerStatus"]
		pubEventData.EventParams = sessionInfo["Reason"]
		pubEventData.EventTime = time.Now().Local().String()

		jvalue, _ := json.Marshal(pubEventData)
		jvalueStr := string(jvalue)
		fmt.Println("Event Pub value: ", jvalueStr)

		if dvpEventType == "amqp" {

			fmt.Println("Start Publish Event to rabbitMQ")
			RabbitMQPublish("DVPEVENTS", jvalue)

		} else {

			Publish("SYS:MONITORING:DVPEVENTS", jvalueStr)

		}

	}
}
