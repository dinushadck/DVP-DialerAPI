package main

import (
	"code.google.com/p/gorest"
	"fmt"
)

type DialerSelfHost struct {
	gorest.RestService     `root:"/DialerSelfHost/" consumes:"application/json" produces:"application/json"`
	incrMaxChannelLimit    gorest.EndPoint `method:"POST" path:"/Campaign/IncrMaxChannelLimit/" postdata:"string"`
	decrMaxChannelLimit    gorest.EndPoint `method:"POST" path:"/Campaign/DecrMaxChannelLimit/" postdata:"string"`
	setMaxChannelLimit     gorest.EndPoint `method:"POST" path:"/Campaign/SetMaxChannelLimit/" postdata:"string"`
	getTotalDialCount      gorest.EndPoint `method:"GET" path:"/Campaign/GetTotalDialCount/{CompanyId:int}/{TenantId:int}/{CampaignId:string}" output:"int"`
	getTotalConnectedCount gorest.EndPoint `method:"GET" path:"/Campaign/GetTotalConnectedCount/{CompanyId:int}/{TenantId:int}/{CampaignId:string}" output:"int"`
}

func (dialerSelfHost DialerSelfHost) IncrMaxChannelLimit(campaignId string) {
	fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
	go IncrCampChannelMaxLimit(campaignId)
	return
}

func (dialerSelfHost DialerSelfHost) DecrMaxChannelLimit(campaignId string) {
	fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
	go DecrCampChannelMaxLimit(campaignId)
	return
}

func (dialerSelfHost DialerSelfHost) SetMaxChannelLimit(campaignId string) {
	fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
	go SetCampChannelMaxLimit(campaignId)
	return
}

func (dialerSelfHost DialerSelfHost) GetTotalDialCount(companyId, tenantId int, campaignId string) int {
	fmt.Println("Start GetTotalDialCount CampaignId: ", campaignId)
	count := GetCampaignDialCount(companyId, tenantId, campaignId)
	return count
}

func (dialerSelfHost DialerSelfHost) GetTotalConnectedCount(companyId, tenantId int, campaignId string) int {
	fmt.Println("Start GetTotalConnectedCount CampaignId: ", campaignId)
	count := GetCampaignConnectedCount(companyId, tenantId, campaignId)
	return count
}
