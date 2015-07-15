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

	campaignDetail1 := Campaign{}
	campaignDetail1.Calss = "Dialer"
	campaignDetail1.CampaignId = "1"
	campaignDetail1.Category = "Call"
	campaignDetail1.Company = 4
	campaignDetail1.EndDate = "2015-08-10T15:10:00.000Z"
	campaignDetail1.Extention = "1001"
	campaignDetail1.ScheduleId = "1"
	campaignDetail1.StartDate = "2015-07-10T11:11:00.000Z"
	campaignDetail1.Tenant = 2
	campaignDetail1.Type = "Callserver"
	campaignDetail1.CallServerId = "2"
	campaignDetail1.MaxCmpaignChannels = 0
	campaignDetail1.DefaultANI = "0888888881"

	campaignDetail2 := Campaign{}
	campaignDetail2.Calss = "Dialer"
	campaignDetail2.CampaignId = "2"
	campaignDetail2.Category = "Call"
	campaignDetail2.Company = 4
	campaignDetail2.EndDate = "2015-08-10T15:10:00.000Z"
	campaignDetail2.Extention = "1002"
	campaignDetail2.ScheduleId = "2"
	campaignDetail2.StartDate = "2015-07-10T11:11:00.000Z"
	campaignDetail2.Tenant = 2
	campaignDetail2.Type = "Callserver"
	campaignDetail2.CallServerId = "2"
	campaignDetail2.MaxCmpaignChannels = 0
	campaignDetail2.DefaultANI = "0888888882"

	campaignDetail3 := Campaign{}
	campaignDetail3.Calss = "Dialer"
	campaignDetail3.CampaignId = "3"
	campaignDetail3.Category = "Call"
	campaignDetail3.Company = 4
	campaignDetail3.EndDate = "2015-08-10T15:10:00.000Z"
	campaignDetail3.Extention = "1003"
	campaignDetail3.ScheduleId = "3"
	campaignDetail3.StartDate = "2015-07-10T11:11:00.000Z"
	campaignDetail3.Tenant = 2
	campaignDetail3.Type = "Callserver"
	campaignDetail3.CallServerId = "2"
	campaignDetail3.MaxCmpaignChannels = 0
	campaignDetail3.DefaultANI = "0888888883"
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
	campaignDetails = append(campaignDetails, campaignDetail1)
	campaignDetails = append(campaignDetails, campaignDetail2)
	campaignDetails = append(campaignDetails, campaignDetail3)
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

func StartCampaign(campaignId, scheduleId, callServerId, extention, defaultAni string, company, tenant int) {
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
								SetCampaignStatus(campaignId, "End", company, tenant)
								RemoveCampaignFromDialer(campaignId, company, tenant)
								return
								//ch <- "End"
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
					SetCampaignStatus(campaignId, "PauseByDialer", company, tenant)
					return
					//ch <- "PauseByDialer"
				}
			} else {
				switch campStatus {
				case "Stop":
					SetCampaignStatus(campaignId, "Stop", company, tenant)
					RemoveCampaignFromDialer(campaignId, company, tenant)
					return
					//ch <- "Stop"
				case "Pause":
					SetCampaignStatus(campaignId, "Pause", company, tenant)
					return
					//ch <- "Pause"
				case "End":
					SetCampaignStatus(campaignId, "End", company, tenant)
					RemoveCampaignFromDialer(campaignId, company, tenant)
					return
					//ch <- "End"
				case "PauseByDialer":
					SetCampaignStatus(campaignId, "PauseByDialer", company, tenant)
					return
					//ch <- "PauseByDialer"
				default:
					SetCampaignStatus(campaignId, "ForceFullyStop", company, tenant)
					return
					//ch <- "ForceFullyStop"
				}
			}
		}
	} else {
		SetCampaignStatus(campaignId, "Waiting for Appoinment", company, tenant)
		return
		//ch <- "Waiting for Appoinment"
	}
}
