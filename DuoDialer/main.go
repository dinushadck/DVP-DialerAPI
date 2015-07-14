// BlastDialer project main.go
package main

import (
	"code.google.com/p/gorest"
	"fmt"
	"net/http"
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
				defCampaign := Campaign{}
				if campaign != defCampaign {
					AddCampaignToDialer(campaign)
				} else {
					fmt.Println("No campaign to add")
				}
			}
		}

		if onGoingCampaignCount > 0 {
			tm := time.Now()
			runningCampaigns := GetAllRunningCampaign()
			for _, campaign := range runningCampaigns {
				campStatus := GetCampaignStatus(campaign.CampaignId, campaign.Company, campaign.Tenant)
				fmt.Println("campStatus: ", campStatus)
				UpdateCampaignStatus(campaign.Company, campaign.Tenant, campaign.CampaignId)

				if campStatus == "Resume" || campStatus == "Start" || campStatus == "PauseByDialer" || campStatus == "Waiting for Appoinment" {
					tempCampaignStartDate, _ := time.Parse(layout1, campaign.StartDate)
					tempCampaignEndDate, _ := time.Parse(layout1, campaign.EndDate)

					campaignStartDate := time.Date(tempCampaignStartDate.Year(), tempCampaignStartDate.Month(), tempCampaignStartDate.Day(), tempCampaignStartDate.Hour(), tempCampaignStartDate.Minute(), tempCampaignStartDate.Second(), 0, time.Local)
					campaignEndDate := time.Date(tempCampaignEndDate.Year(), tempCampaignEndDate.Month(), tempCampaignEndDate.Day(), tempCampaignEndDate.Hour(), tempCampaignEndDate.Minute(), tempCampaignEndDate.Second(), 0, time.Local)
					fmt.Println("Check Campaign: ", campaign.CampaignId)
					fmt.Println("campaignStartDate: ", campaignStartDate.String())
					fmt.Println("campaignEndDate: ", campaignEndDate.String())

					if campaignEndDate.Before(tm) {
						fmt.Println("campaignEndDate before: ", tm.String())
						RemoveCampaignFromDialer(campaign.CampaignId, campaign.Company, campaign.Tenant)
					} else if campaignStartDate.Before(tm) && campaignEndDate.After(tm) {
						fmt.Println("Continue campaign: ", campaign.CampaignId)
						ch := make(chan string)
						fmt.Println("StartCampaign: ", campaign.CampaignId)
						go StartCampaign(campaign.CampaignId, campaign.ScheduleId, campaign.CallServerId, campaign.Extention, campaign.DefaultANI, campaign.Company, campaign.Tenant, ch)
						status := <-ch

						SetCampaignStatus(campaign.CampaignId, status, campaign.Company, campaign.Tenant)
						if status == "End" || status == "Stop" {
							//RemoveCampaignFromDialer(campaign.CampaignId, campaign.Company, campaign.Tenant)
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
	http.ListenAndServe(":6565", nil)
}
