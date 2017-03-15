package main

import (
	"encoding/json"
	"fmt"
	"github.com/DuoSoftware/log4go"
	"strconv"
	"strings"
	"time"
)

var eventLog = log4go.NewLogger()

func OnEvent(eventInfo SubEvents) {
	defEvent := SubEvents{}
	if eventInfo != defEvent {
		fmt.Println("AuthData: ", eventInfo.AuthData)
		fmt.Println("SessionId: ", eventInfo.SessionId)
		fmt.Println("EventName: ", eventInfo.EventName)
		fmt.Println("EventCategory: ", eventInfo.EventCategory)
		fmt.Println("SwitchName: ", eventInfo.SwitchName)

		authInfoArr := strings.Split(eventInfo.AuthData, "_")
		if len(authInfoArr) == 2 {
			company, _ := strconv.Atoi(authInfoArr[0])
			tenant, _ := strconv.Atoi(authInfoArr[1])

			fmt.Println("company: ", company)
			fmt.Println("tenant: ", tenant)

			switch eventInfo.EventCategory {
			case "CHANNEL_BRIDGE":
				fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
				break
			case "CHANNEL_CREATE":
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "ChannelCreatetime", time.Now().Format(layout4))
				fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
				break
			case "CHANNEL_ANSWER":
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "DialerStatus", "channel_answered")
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "ChannelAnswertime", time.Now().Format(layout4))
				IncrCampaignConnectedCount(company, tenant, eventInfo.CampaignId)
				fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
				break
			case "CHANNEL_DESTROY":
				LogEvent(eventInfo)
				hashKey := fmt.Sprintf("sessionInfo:%s:%s", eventInfo.CampaignId, eventInfo.SessionId)
				session := RedisCheckKeyExist(hashKey)
				if session {
					DecrConcurrentChannelCount(eventInfo.SwitchName, eventInfo.CampaignId)
					SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "Reason", eventInfo.DisconnectReason)
					go UploadSessionInfo(eventInfo.CampaignId, eventInfo.SessionId)
					fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
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

		Publish("SYS:MONITORING:DVPEVENTS", jvalueStr)
	}
}
