package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/fatih/color"
)

func PublishCampaignCallCounts(sessionId, category, comapnyId, tenantId, campaignId string) {
	tenant, _ := strconv.Atoi(tenantId)
	company, _ := strconv.Atoi(comapnyId)

	pubEventData := DashboardEvents{}

	pubEventData.SessionID = sessionId
	pubEventData.Tenant = tenant
	pubEventData.Company = company
	pubEventData.EventClass = "DIALER"
	pubEventData.EventType = "CAMPAIGN"
	pubEventData.EventCategory = category
	pubEventData.TimeStamp = time.Now().Format(layout1)
	pubEventData.Parameter1 = campaignId
	pubEventData.Parameter2 = ""
	pubEventData.BusinessUnit = "default"

	jvalue, _ := json.Marshal(pubEventData)
	jvalueStr := string(jvalue)
	
	color.Magenta(fmt.Sprintf("!!!!!!!!! DASHBOARD PUBLISH - %s - !!!!!!!!!!", jvalueStr))
	RabbitMQPublish("CampaignEvents", jvalue)

}