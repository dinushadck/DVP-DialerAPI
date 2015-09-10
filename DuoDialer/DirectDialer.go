package main

import (
	"fmt"
	"strconv"
	"time"
)

func DirectDialNumber(company, tenant, campaignId int, number string) bool {
	fmt.Println("start DialNumber")
	campaignIdStr := strconv.Itoa(campaignId)
	campaignInfo, isCamExists := GetCampaign(company, tenant, campaignId)
	if isCamExists {
		tmNowUTC := time.Now().UTC()

		tempCampaignEndDate, _ := time.Parse(layout1, campaignInfo.CampConfigurations.EndDate)
		campaignEndDate := time.Date(tempCampaignEndDate.Year(), tempCampaignEndDate.Month(), tempCampaignEndDate.Day(), tempCampaignEndDate.Hour(), tempCampaignEndDate.Minute(), tempCampaignEndDate.Second(), 0, time.UTC)

		if campaignEndDate.After(tmNowUTC) {
			scheduleIdStr := strconv.Itoa(campaignInfo.CampScheduleInfo[0].ScheduleId)
			camScheduleStr := strconv.Itoa(campaignInfo.CampScheduleInfo[0].CamScheduleId)
			validateAppoinment := CheckAppoinmentForCallback(company, tenant, scheduleIdStr, tmNowUTC)
			if validateAppoinment {
				numberWithTryCount := fmt.Sprintf("%s:%d", number, 1)
				return AddNumberToFront(company, tenant, campaignIdStr, camScheduleStr, numberWithTryCount)
			}
		}
	}
	return false
}
