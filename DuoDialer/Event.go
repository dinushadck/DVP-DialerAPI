package main

import (
	"code.google.com/p/log4go"
	"encoding/json"
	"fmt"
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
				SetSessionInfo(eventInfo.CampaignId, eventInfo.SessionId, "ChannelAnswertime", time.Now().Format(layout4))
				IncrCampaignConnectedCount(company, tenant, eventInfo.CampaignId)
				fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
				break
			case "CHANNEL_DESTROY":
				LogEvent(eventInfo)
				hashKey := fmt.Sprintf("sessionInfo:%s:%s", dialerId, eventInfo.SessionId)
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
