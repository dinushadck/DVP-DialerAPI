package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
				fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
				break
			case "CHANNEL_ANSWER":
				fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
				break
			case "CHANNEL_DESTROY":
				DecrConcurrentChannelCount(eventInfo.SwitchName, eventInfo.CampaignId)
				fmt.Println("SessionId: ", eventInfo.SessionId, " EventName: ", eventInfo.EventName, " EventCat: ", eventInfo.EventCategory)
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
