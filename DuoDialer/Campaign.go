package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
	campaignKey := fmt.Sprintf("Campaign:%s:%d:%d:%d", dialerId, campaignD.CompanyId, campaignD.TenantId, campaignD.CampaignId)
	DefCampConfig := CampaignConfigInfo{}
	if campaignD.CampConfigurations != DefCampConfig {
		campaignJson, _ := json.Marshal(campaignD)
		result := RedisAdd(campaignKey, string(campaignJson))
		fmt.Println("Add Campaign to Redis: ", campaignKey, " Result: ", result)
		if result == "OK" {
			campId64 := int64(campaignD.CampaignId)
			campIdStr := strconv.FormatInt(campId64, 32)
			IncrementOnGoingCampaignCount()
			SetCampaignStatus(campIdStr, "Start", campaignD.CompanyId, campaignD.TenantId)
			UpdateCampaignStartStatus(campaignD.CompanyId, campaignD.TenantId, campIdStr)
		}
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

	client := &http.Client{}

	request := fmt.Sprintf("%s/CampaignManager/Handler/pending/%d", campaignService, requestCount)
	fmt.Println("Start RequestCampaign request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", "")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return campaignDetails
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var campaignResult CampaignResult
	json.Unmarshal(response, &campaignResult)
	if campaignResult.IsSuccess == true {
		for _, camRes := range campaignResult.Result {
			campaignDetails = append(campaignDetails, camRes)
		}
	}
	return campaignDetails
}

func UpdateCampaignStatus(company, tenant int, campaignId string) {
	//Send CampaignStatus to Campaign Manager
	authToken := fmt.Sprintf("%d#%d", company, tenant)
	fmt.Println("Start UpdateCampaignStatus Auth: ", authToken, " CampaignId: ", campaignId, " DialerId: ", dialerId)
	client := &http.Client{}

	currentState := GetCampaignStatus(campaignId, company, tenant)
	request := fmt.Sprintf("%s/CampaignManager/Operations/OperationState/%s/%s/%s", campaignService, campaignId, dialerId, currentState)
	fmt.Println("Start UpdateCampaignStatus request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	var campaignStatusResult CampaignStatusResult
	json.Unmarshal(response, &campaignStatusResult)
	if campaignStatusResult.IsSuccess == true {
		emtResult := CampaignState{}
		if campaignStatusResult.Result != emtResult {
			state := campaignStatusResult.Result.CampaignState
			camId := campaignStatusResult.Result.CampaignId
			dId := campaignStatusResult.Result.DialerId

			campIdStr := strconv.Itoa(camId)

			if campaignId == campIdStr && dialerId == dId && currentState != state {
				switch state {
				case "Stop":
					SetCampaignStatus(campIdStr, "Stop", company, tenant)
					RemoveCampaignFromDialer(campIdStr, company, tenant)
					break
				case "Pause":
					SetCampaignStatus(campIdStr, "Pause", company, tenant)
					break
				case "End":
					SetCampaignStatus(campIdStr, "End", company, tenant)
					RemoveCampaignFromDialer(campIdStr, company, tenant)
					break
				default:
					break
				}
			}
		}
	}
}

func UpdateCampaignStartStatus(company, tenant int, campaignId string) {
	//Send CampaignStatus to Campaign Manager
	state := CampaignStart{}
	camIdInt, _ := strconv.Atoi(campaignId)
	state.CampaignId = camIdInt
	state.DialerId = dialerId

	jsonData, _ := json.Marshal(state)

	serviceurl := fmt.Sprintf("%s/CampaignManager/Operations", campaignService)
	authToken := fmt.Sprintf("%d#%d", company, tenant)
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	fmt.Println("request:", serviceurl)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, errb := ioutil.ReadAll(resp.Body)
	if errb != nil {
		fmt.Println(err.Error())
	} else {
		result := string(body)
		fmt.Println("response Body:", result)
	}
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

func SetCampChannelMaxLimit(campaignId string) {
	ids := strings.Split(campaignId, "_")
	if len(ids) == 2 {
		cmcl := fmt.Sprintf("CampaignMaxCallLimit:%s", ids[0])
		RedisSet(cmcl, ids[1])
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

func StartCampaign(campaignId, scheduleId, camScheduleId, callServerId, extention, defaultAni string, company, tenant int) {
	emtAppoinment := Appoinment{}
	defCallServerInfo := CallServerInfo{}
	authToken := fmt.Sprintf("%d#%d", company, tenant)
	appment := CheckAppoinmentForCampaign(authToken, scheduleId)
	callServerInfos := GetCallServerInfo(callServerId)
	if appment != emtAppoinment && callServerInfos != defCallServerInfo {
		campStatus := GetCampaignStatus(campaignId, company, tenant)
		if campStatus == "Start" {
			LoadInitialNumberSet(company, tenant, campaignId, camScheduleId)
		}

		SetCampaignStatus(campaignId, "Running", company, tenant)

		appmntEndTime, _ := time.Parse(layout1, appment.EndTime)

		for {
			campStatus = GetCampaignStatus(campaignId, company, tenant)
			if campStatus == "Running" {
				tm := time.Now()
				if appmntEndTime.Before(tm) {
					cchannelCountS, cchannelCountC := GetConcurrentChannelCount(callServerInfos.CallServerId, campaignId)
					maxChannelLimit := GetMaxChannelLimit(callServerInfos.CallServerId)
					maxCampaignChannelLimit := GetCampMaxChannelLimit(campaignId)
					fmt.Println("callServerInfos.CallServerId: ", callServerInfos.CallServerId)
					fmt.Println("MaxCallServerChannelLimit: ", maxChannelLimit)
					fmt.Println("maxCampaignChannelLimit: ", maxCampaignChannelLimit)
					fmt.Println("ConcurrentCallServerChannel: ", cchannelCountS)
					fmt.Println("ConcurrentCampaignChannel: ", cchannelCountC)

					if cchannelCountS < maxChannelLimit && cchannelCountC < maxCampaignChannelLimit {
						number := GetNumberToDial(company, tenant, campaignId, camScheduleId)
						if number == "" {
							numberCount := GetNumberCount(company, tenant, campaignId, camScheduleId)
							if numberCount == 0 {
								SetCampaignStatus(campaignId, "End", company, tenant)
								RemoveCampaignFromDialer(campaignId, company, tenant)
								return
							}
						} else {
							//trunkCode, ani, dnis := "OutTrunk001", defaultAni, number
							trunkCode, ani, dnis := GetTrunkCode(authToken, defaultAni, number)
							uuid := GetUuid()
							if trunkCode != "" && uuid != "" {
								go DialNumber(company, tenant, callServerInfos, campaignId, uuid, ani, trunkCode, dnis, extention)
								//go DialNumberFIFO(company, tenant, callServerInfos, campaignId, uuid, ani, trunkCode, dnis, extention)
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
				}
			} else {
				switch campStatus {
				case "Stop":
					SetCampaignStatus(campaignId, "Stop", company, tenant)
					RemoveCampaignFromDialer(campaignId, company, tenant)
					return
				case "Pause":
					SetCampaignStatus(campaignId, "Pause", company, tenant)
					return
				case "End":
					SetCampaignStatus(campaignId, "End", company, tenant)
					RemoveCampaignFromDialer(campaignId, company, tenant)
					return
				case "PauseByDialer":
					SetCampaignStatus(campaignId, "PauseByDialer", company, tenant)
					return
				default:
					SetCampaignStatus(campaignId, "ForceFullyStop", company, tenant)
					return
				}
			}
		}
	} else {
		campStatus := GetCampaignStatus(campaignId, company, tenant)
		switch campStatus {
		case "Stop":
			SetCampaignStatus(campaignId, "Stop", company, tenant)
			RemoveCampaignFromDialer(campaignId, company, tenant)
			return
		case "Pause":
			SetCampaignStatus(campaignId, "Pause", company, tenant)
			return
		case "End":
			SetCampaignStatus(campaignId, "End", company, tenant)
			RemoveCampaignFromDialer(campaignId, company, tenant)
			return
		case "PauseByDialer":
			SetCampaignStatus(campaignId, "PauseByDialer", company, tenant)
			return
		default:
			SetCampaignStatus(campaignId, "Waiting for Appoinment", company, tenant)
			return
		}
	}
}
