// BlastDialer project main.go
package main

import (
	"fmt"
	"github.com/DuoSoftware/gorest"
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
	go AddRequestServer()

	//AddPhoneNumberToCallback("1", "1", "1", "1", "0112546969", "USER_BUSY")
	for {
		onGoingCampaignCount := GetOnGoingCampaignCount()
		if onGoingCampaignCount < campaignLimit {
			//campaigns := RequestCampaign(campaignLimit - onGoingCampaignCount)
			//for _, campaign := range campaigns {
			//	AddCampaignToDialer(campaign)
			//}
		}

		if onGoingCampaignCount > 0 {

			tm := time.Now().UTC()
			runningCampaigns := GetAllRunningCampaign()
			for _, campaign := range runningCampaigns {
				campIdStr := strconv.Itoa(campaign.CampaignId)

				go ClearTimeoutChannels(campIdStr)

				campStatus := GetCampaignStatus(campIdStr, campaign.CompanyId, campaign.TenantId)
				fmt.Println("campStatus: ", campStatus)
				UpdateCampaignStatus(campaign.CompanyId, campaign.TenantId, campIdStr)

				if campStatus == "Resume" || campStatus == "Start" || campStatus == "PauseByDialer" || campStatus == "Waiting for Appoinment" {
					tempCampaignStartDate, _ := time.Parse(layout1, campaign.CampConfigurations.StartDate)
					tempCampaignEndDate, _ := time.Parse(layout1, campaign.CampConfigurations.EndDate)

					if campStatus == "Resume" {
						UpdateCampaignStartStatus(campaign.CompanyId, campaign.TenantId, campIdStr)
					}

					campaignStartDate := time.Date(tempCampaignStartDate.Year(), tempCampaignStartDate.Month(), tempCampaignStartDate.Day(), tempCampaignStartDate.Hour(), tempCampaignStartDate.Minute(), tempCampaignStartDate.Second(), 0, time.UTC)
					campaignEndDate := time.Date(tempCampaignEndDate.Year(), tempCampaignEndDate.Month(), tempCampaignEndDate.Day(), tempCampaignEndDate.Hour(), tempCampaignEndDate.Minute(), tempCampaignEndDate.Second(), 0, time.UTC)
					fmt.Println("Check Campaign: ", campIdStr)
					fmt.Println("campaignStartDate: ", campaignStartDate.String())
					fmt.Println("campaignEndDate: ", campaignEndDate.String())

					if campaignEndDate.Before(tm) {
						fmt.Println("campaignEndDate before: ", tm.String())
						RemoveCampaignFromDialer(campIdStr, campaign.CompanyId, campaign.TenantId)
					} else if campaignStartDate.Before(tm) && campaignEndDate.After(tm) {
						fmt.Println("Continue campaign: ", campIdStr)
						if len(campaign.CampScheduleInfo) > 0 {
							scheduleId := strconv.Itoa(campaign.CampScheduleInfo[0].ScheduleId)
							camScheduleId := strconv.Itoa(campaign.CampScheduleInfo[0].CamScheduleId)
							go StartCampaign(campIdStr, campaign.DialoutMechanism, campaign.CampaignChannel, campaign.Class, campaign.Type, campaign.Category, scheduleId, camScheduleId, "*", campaign.Extensions, campaign.CampConfigurations.Caller, campaign.CompanyId, campaign.TenantId, campaign.CampConfigurations.ChannelConcurrency)
						}
					}
				} else {
					switch campStatus {
					case "Stop":
						SetCampaignStatus(campIdStr, "Stop", campaign.CompanyId, campaign.TenantId)
						RemoveCampaignFromDialer(campIdStr, campaign.CompanyId, campaign.TenantId)
						break
					case "End":
						SetCampaignStatus(campIdStr, "End", campaign.CompanyId, campaign.TenantId)
						RemoveCampaignFromDialer(campIdStr, campaign.CompanyId, campaign.TenantId)
						break
					default:
						break
					}
				}
			}
		}

		time.Sleep(campaignRequestFrequency * time.Second)
	}
}

func InitiateService() {
	jwtMiddleware := loadJwtMiddleware()
	gorest.RegisterService(new(DVP))
	//http.Handle("/", gorest.Handle())
	app := jwtMiddleware.Handler(gorest.Handle())
	addr := fmt.Sprintf(":%s", port)
	fmt.Println(addr)
	http.ListenAndServe(addr, app)
}
