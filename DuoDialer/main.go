// BlastDialer project main.go
package main

import (
	"code.google.com/p/gorest"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	InitiateDuoDialer()
	go InitiateService()
	for {
		onGoingCampaignCount := GetOnGoingCampaignCount()
		if onGoingCampaignCount < campaignLimit {
			campaigns := RequestCampaign(campaignLimit - onGoingCampaignCount)
			for _, campaign := range campaigns {
				AddCampaignToDialer(campaign)
			}
		}

		if onGoingCampaignCount > 0 {
			tm := time.Now()
			runningCampaigns := GetAllRunningCampaign()
			for _, campaign := range runningCampaigns {
				campIdStr := strconv.Itoa(campaign.CampaignId)
				campStatus := GetCampaignStatus(campIdStr, campaign.CompanyId, campaign.TenantId)
				fmt.Println("campStatus: ", campStatus)
				UpdateCampaignStatus(campaign.CompanyId, campaign.TenantId, campIdStr)

				if campStatus == "Resume" || campStatus == "Start" || campStatus == "PauseByDialer" || campStatus == "Waiting for Appoinment" {
					tempCampaignStartDate, _ := time.Parse(layout1, campaign.CampConfigurations.StartDate)
					tempCampaignEndDate, _ := time.Parse(layout1, campaign.CampConfigurations.EndDate)

					campaignStartDate := time.Date(tempCampaignStartDate.Year(), tempCampaignStartDate.Month(), tempCampaignStartDate.Day(), tempCampaignStartDate.Hour(), tempCampaignStartDate.Minute(), tempCampaignStartDate.Second(), 0, time.Local)
					campaignEndDate := time.Date(tempCampaignEndDate.Year(), tempCampaignEndDate.Month(), tempCampaignEndDate.Day(), tempCampaignEndDate.Hour(), tempCampaignEndDate.Minute(), tempCampaignEndDate.Second(), 0, time.Local)
					fmt.Println("Check Campaign: ", campIdStr)
					fmt.Println("campaignStartDate: ", campaignStartDate.String())
					fmt.Println("campaignEndDate: ", campaignEndDate.String())

					if campaignEndDate.Before(tm) {
						fmt.Println("campaignEndDate before: ", tm.String())
						RemoveCampaignFromDialer(campIdStr, campaign.CompanyId, campaign.TenantId)
					} else if campaignStartDate.Before(tm) && campaignEndDate.After(tm) {
						fmt.Println("Continue campaign: ", campIdStr)
						//ch := make(chan string)
						fmt.Println("StartCampaign: ", campIdStr)
						if len(campaign.CampScheduleInfo) > 0 {
							scheduleId := strconv.Itoa(campaign.CampScheduleInfo[0].ScheduleId)
							camScheduleId := strconv.Itoa(campaign.CampScheduleInfo[0].CamScheduleId)
							fmt.Println("SetscheduleId Pre: ", campaign.CampScheduleInfo[0].ScheduleId, " Post: ", scheduleId)
							go StartCampaign(campIdStr, scheduleId, camScheduleId, "*", campaign.Extensions, campaign.CampConfigurations.Caller, campaign.CompanyId, campaign.TenantId)
						}
					}
				}
			}
		}

		time.Sleep(campaignRequestFrequency * time.Second)
	}
}

func InitiateService() {
	gorest.RegisterService(new(DialerSelfHost))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":2223", nil)
}
