package main

import (
	"code.google.com/p/gorest"
	"fmt"
	"strconv"
	"strings"
)

type DVP struct {
	gorest.RestService     `root:"/DVP/" consumes:"application/json" produces:"application/json"`
	incrMaxChannelLimit    gorest.EndPoint `method:"POST" path:"/DialerAPI/IncrMaxChannelLimit/" postdata:"string"`
	decrMaxChannelLimit    gorest.EndPoint `method:"POST" path:"/DialerAPI/DecrMaxChannelLimit/" postdata:"string"`
	setMaxChannelLimit     gorest.EndPoint `method:"POST" path:"/DialerAPI/SetMaxChannelLimit/" postdata:"string"`
	getTotalDialCount      gorest.EndPoint `method:"GET" path:"/DialerAPI/GetTotalDialCount/{CompanyId:int}/{TenantId:int}/{CampaignId:string}" output:"int"`
	getTotalConnectedCount gorest.EndPoint `method:"GET" path:"/DialerAPI/GetTotalConnectedCount/{CompanyId:int}/{TenantId:int}/{CampaignId:string}" output:"int"`
	resumeCallback         gorest.EndPoint `method:"POST" path:"/DialerAPI/ResumeCallback/" postdata:"CampaignCallbackObj"`
	dial                   gorest.EndPoint `method:"GET" path:"/DialerAPI/Dial/{CampaignId:int}/{ContactNumber:string}" output:"bool"`
}

func (dvp DVP) IncrMaxChannelLimit(campaignId string) {
	fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
	go IncrCampChannelMaxLimit(campaignId)
	return
}

func (dvp DVP) DecrMaxChannelLimit(campaignId string) {
	fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
	go DecrCampChannelMaxLimit(campaignId)
	return
}

func (dvp DVP) SetMaxChannelLimit(campaignId string) {
	fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
	go SetCampChannelMaxLimit(campaignId)
	return
}

func (dvp DVP) GetTotalDialCount(companyId, tenantId int, campaignId string) int {
	fmt.Println("Start GetTotalDialCount CampaignId: ", campaignId)
	count := 0
	authHeaderStr := dvp.Context.Request().Header.Get("Authorization")
	fmt.Println(authHeaderStr)

	authHeaderInfo := strings.Split(authHeaderStr, "#")
	if len(authHeaderInfo) == 2 {
		tenant, _ := strconv.Atoi(authHeaderInfo[0])
		company, _ := strconv.Atoi(authHeaderInfo[1])
		count = GetCampaignDialCount(company, tenant, campaignId)
	}
	return count
}

func (dvp DVP) GetTotalConnectedCount(companyId, tenantId int, campaignId string) int {
	fmt.Println("Start GetTotalConnectedCount CampaignId: ", campaignId)
	count := 0
	authHeaderStr := dvp.Context.Request().Header.Get("Authorization")
	fmt.Println(authHeaderStr)

	authHeaderInfo := strings.Split(authHeaderStr, "#")
	if len(authHeaderInfo) == 2 {
		tenant, _ := strconv.Atoi(authHeaderInfo[0])
		company, _ := strconv.Atoi(authHeaderInfo[1])
		count = GetCampaignConnectedCount(company, tenant, campaignId)
	}
	return count
}

func (dvp DVP) ResumeCallback(callbackInfo CampaignCallbackObj) {
	log := fmt.Sprintf("Start ResumeCallback CampaignId:%d # ContactId:%s ", callbackInfo.CampaignId, callbackInfo.ContactId)
	fmt.Println(log)
	authHeaderStr := dvp.Context.Request().Header.Get("Authorization")
	fmt.Println(authHeaderStr)

	authHeaderInfo := strings.Split(authHeaderStr, "#")
	if len(authHeaderInfo) == 2 {
		tenant, _ := strconv.Atoi(authHeaderInfo[0])
		company, _ := strconv.Atoi(authHeaderInfo[1])
		fmt.Println("Company: ", company)
		fmt.Println("Tenant: ", tenant)
		ResumeCampaignCallback(company, tenant, callbackInfo.CallBackCount, callbackInfo.CampaignId, callbackInfo.ContactId)
	}
	return
}

func (dvp DVP) Dial(campaignId int, contactNumber string) bool {
	log := fmt.Sprintf("Start Direct Dial CampaignId:%d # DNIS:%s ", campaignId, contactNumber)
	fmt.Println(log)
	authHeaderStr := dvp.Context.Request().Header.Get("Authorization")
	fmt.Println(authHeaderStr)

	authHeaderInfo := strings.Split(authHeaderStr, "#")
	if len(authHeaderInfo) == 2 {
		tenant, _ := strconv.Atoi(authHeaderInfo[0])
		company, _ := strconv.Atoi(authHeaderInfo[1])
		fmt.Println("Company: ", company)
		fmt.Println("Tenant: ", tenant)
		return DirectDialNumber(company, tenant, campaignId, contactNumber)
	}
	return false
}
