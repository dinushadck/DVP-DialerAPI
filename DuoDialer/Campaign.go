package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func AddOnGoingCampaignCount() {
	dialerCampaignCountKey := fmt.Sprintf("DialerOnGoingCampaignCount:%s:%s", hostIpAddress, dialerId)
	result := RedisAdd(dialerCampaignCountKey, "0")
	fmt.Println("Add DialerOnGoingCampaignCount to Redis: ", result)
}

func GetOnGoingCampaignCount() int {
	dialerCampaignCountKey := fmt.Sprintf("DialerOnGoingCampaignCount:%s:%s", hostIpAddress, dialerId)
	countStr := RedisGet(dialerCampaignCountKey)
	count, _ := strconv.Atoi(countStr)
	fmt.Println("OnGoingCampaignCount: ", countStr)
	return count
}

func IncrementOnGoingCampaignCount() {
	dialerCampaignCountKey := fmt.Sprintf("DialerOnGoingCampaignCount:%s:%s", hostIpAddress, dialerId)
	incValue := RedisIncr(dialerCampaignCountKey)
	fmt.Println("IncrementedOnGoingCampaignCount: ", incValue)
}

func DecrementOnGoingCampaignCount() {
	dialerCampaignCountKey := fmt.Sprintf("DialerOnGoingCampaignCount:%s:%s", hostIpAddress, dialerId)
	dncValue := RedisIncrBy(dialerCampaignCountKey, -1)
	fmt.Println("DecrementOnGoingCampaignCount: ", dncValue)
}

func SetCampaignStatus(campaignId, status string, company, tenant int) {
	campaignStatusKey := fmt.Sprintf("CampaignStatus:%s:%d:%d:%s", dialerId, company, tenant, campaignId)
	result := RedisSet(campaignStatusKey, status)
	fmt.Println("SetCampaignStatus CampaignId: ", campaignStatusKey, " Result: ", result)
}

func GetCampaignStatus(campaignId string, company, tenant int) string {
	campaignStatusKey := fmt.Sprintf("CampaignStatus:%s:%d:%d:%s", dialerId, company, tenant, campaignId)
	result := RedisGet(campaignStatusKey)
	fmt.Println("GetCampaignStatus CampaignId: ", campaignStatusKey, " Result: ", result)
	return result
}

func RemoveCampaignStatus(campaignId string, company, tenant int) {
	campaignStatusKey := fmt.Sprintf("CampaignStatus:%s:%d:%d:%s", dialerId, company, tenant, campaignId)
	result := RedisRemove(campaignStatusKey)
	fmt.Println("GetCampaignStatus CampaignId: ", campaignStatusKey, " Result: ", result)
}

func AddCampaignToDialer(campaignD Campaign) {
	campaignKey := fmt.Sprintf("Campaign:%s:%d:%d:%s", dialerId, campaignD.Company, campaignD.Tenant, campaignD.CampaignId)
	campaignJson, _ := json.Marshal(campaignD)
	result := RedisAdd(campaignKey, string(campaignJson))
	fmt.Println("Add Campaign to Redis: ", campaignKey, " Result: ", result)
	if result == "OK" {
		IncrementOnGoingCampaignCount()
		SetCampaignStatus(campaignD.CampaignId, "Start", campaignD.Company, campaignD.Tenant)
		UpdateCampaignStartStatus(campaignD.Company, campaignD.Tenant, campaignD.CampaignId)
	}
}

func GetAllRunningCampaign() []Campaign {
	searchKey := fmt.Sprintf("Campaign:%s:*", dialerId)
	allCampaignKeys := RedisSearchKeys(searchKey)

	allCampaigns := make([]Campaign, 0)
	for _, key := range allCampaignKeys {
		campJson := RedisGet(key)
		var camp Campaign
		json.Unmarshal([]byte(campJson), &camp)
		allCampaigns = append(allCampaigns, camp)
	}
	return allCampaigns
}

func RemoveCampaignFromDialer(campaignId string, company, tenant int) {
	campaignKey := fmt.Sprintf("Campaign:%s:%d:%d:%s", dialerId, company, tenant, campaignId)
	result := RedisRemove(campaignKey)
	fmt.Println("Remove Campaign from Redis: ", campaignKey, " Result: ", result)
	if result == true {
		DecrementOnGoingCampaignCount()
		RemoveCampaignStatus(campaignId, company, tenant)
		RemoveNumbers(company, tenant, campaignId)
	}
}

func RequestCampaign(requestCount int) []Campaign {
	//Request campaign from Campaign Manager service
	campaignDetails := make([]Campaign, 0)

	campaignDetail := Campaign{}
	campaignDetail.Calss = "Dialer"
	campaignDetail.CampaignId = "1"
	campaignDetail.Category = "Call"
	campaignDetail.Company = 4
	campaignDetail.EndDate = "2015-08-10T15:10:00.000Z"
	campaignDetail.Extention = "5000"
	campaignDetail.ScheduleId = "1"
	campaignDetail.StartDate = "2015-07-10T11:11:00.000Z"
	campaignDetail.Tenant = 2
	campaignDetail.Type = "Callserver"
	campaignDetail.CallServerId = "2"
	campaignDetail.MaxCmpaignChannels = 0
	campaignDetail.DefaultANI = "077555555"

	/*authToken := fmt.Sprintf("%d#%d", company, tenant)
	fmt.Println("Start GetPhoneNumbers Auth: ", authToken, " CampaignId: ", campaignId, " SchedulrId: ", scheduleId)
	client := &http.Client{}

	request := fmt.Sprintf("%s/CampaignManager/Handler/pending/%d", requestCount)
	fmt.Println("Start RequestCampaign request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return numbers
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var campaignResult CampaignResult
	json.Unmarshal(response, &campaignResult)
	if campaignResult.IsSuccess == true {
		for _, camRes := range campaignResult.Result {
			campaignDetails = append(campaignDetails, camRes)
		}
	}*/
	campaignDetails = append(campaignDetails, campaignDetail)
	return campaignDetails
}

func UpdateCampaignStatus(company, tenant int, campaignId string) {
	//Send CampaignStatus to Campaign Manager
}

func UpdateCampaignStartStatus(company, tenant int, campaignId string) {
	//Send CampaignStatus to Campaign Manager
}

func IncrCampChannelMaxLimit(campaignId string) {
	cmcl := fmt.Sprintf("CampaignMaxCallLimit:%s", campaignId)
	RedisIncr(cmcl)
}

func DecrCampChannelMaxLimit(campaignId string) {
	cmcl := fmt.Sprintf("CampaignMaxCallLimit:%s", campaignId)
	decValue := RedisIncrBy(cmcl, -1)
	if decValue < 0 {
		RedisSet(cmcl, "0")
	}
}

func GetCampMaxChannelLimit(campaignId string) int {
	cmcl := fmt.Sprintf("CampaignMaxCallLimit:%s", campaignId)
	maxChannelCount := RedisGet(cmcl)
	if maxChannelCount == "" {
		RedisSet(cmcl, "0")
	}
	value, err := strconv.Atoi(maxChannelCount)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	} else {
		return value
	}
}

func StartCampaign(campaignId, scheduleId, callServerId, extention, defaultAni string, company, tenant int, ch chan string) {
	emtAppoinment := Appoinment{}
	defCallServerInfo := CallServerInfo{}
	authToken := fmt.Sprintf("%d#%d", company, tenant)
	appment := CheckAppoinmentForCampaign(authToken, scheduleId)
	callServerInfo := GetCallServerInfo(callServerId)
	if appment != emtAppoinment && callServerInfo != defCallServerInfo {
		campStatus := GetCampaignStatus(campaignId, company, tenant)
		if campStatus == "Start" {
			LoadInitialNumberSet(company, tenant, campaignId, scheduleId)
		}

		SetCampaignStatus(campaignId, "Running", company, tenant)

		appmntEndTime, _ := time.Parse(layout1, appment.EndTime)

		for {
			campStatus = GetCampaignStatus(campaignId, company, tenant)
			if campStatus == "Running" {
				tm := time.Now()
				if appmntEndTime.Before(tm) {
					cchannelCountS, cchannelCountC := GetConcurrentChannelCount(callServerId, campaignId)
					maxChannelLimit := GetMaxChannelLimit(callServerId)
					maxCampaignChannelLimit := GetCampMaxChannelLimit(campaignId)

					fmt.Println("MaxCallServerChannelLimit: ", maxChannelLimit)
					fmt.Println("maxCampaignChannelLimit: ", maxCampaignChannelLimit)
					fmt.Println("ConcurrentCallServerChannel: ", cchannelCountS)
					fmt.Println("ConcurrentCampaignChannel: ", cchannelCountC)
					//Check channel count
					//fmt.Println("FreeChannels: ", maxChannelLimit-cchannelCount)
					if cchannelCountS < maxChannelLimit && cchannelCountC < maxCampaignChannelLimit {
						number := GetNumberToDial(company, tenant, campaignId, scheduleId)
						if number == "" {
							numberCount := GetNumberCount(company, tenant, campaignId, scheduleId)
							if numberCount == 0 {
								ch <- "End"
							}
						} else {
							trunkCode, ani, dnis := "OutTrunk001", defaultAni, number
							uuid := GetUuid()
							if trunkCode != "" && uuid != "" {
								//go DialNumber(company, tenant, callServerInfo.Url, campaignId, uuid, ani, trunkCode, dnis, extention)
								go DialNumberFIFO(company, tenant, callServerInfo.Url, campaignId, uuid, ani, trunkCode, dnis, extention)
								time.Sleep(100 * time.Millisecond)
							}
						}
					} else {
						fmt.Println("dialer waiting...")
						time.Sleep(500 * time.Millisecond)
					}
				} else {
					ch <- "PauseByDialer"
				}
			} else {
				switch campStatus {
				case "Stop":
					ch <- "Stop"
				case "Pause":
					ch <- "Pause"
				case "End":
					ch <- "End"
				case "PauseByDialer":
					ch <- "PauseByDialer"
				default:
					ch <- "ForceFullyStop"
				}
			}
		}
	} else {
		ch <- "Waiting for Appoinment"
	}
}
