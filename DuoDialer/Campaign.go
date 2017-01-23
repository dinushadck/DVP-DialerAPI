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

//----------Ongoing Campaign Count-----------------------
func AddOnGoingCampaignCount() {
	dialerCampaignCountKey := fmt.Sprintf("DialerOnGoingCampaignCount:%s:%s", lbIpAddress, dialerId)
	result := RedisAdd(dialerCampaignCountKey, "0")
	fmt.Println("Add DialerOnGoingCampaignCount to Redis: ", result)
}

func GetOnGoingCampaignCount() int {
	dialerCampaignCountKey := fmt.Sprintf("DialerOnGoingCampaignCount:%s:%s", lbIpAddress, dialerId)
	countStr := RedisGet(dialerCampaignCountKey)
	count, _ := strconv.Atoi(countStr)
	fmt.Println("OnGoingCampaignCount: ", countStr)
	return count
}

func IncrementOnGoingCampaignCount() {
	dialerCampaignCountKey := fmt.Sprintf("DialerOnGoingCampaignCount:%s:%s", lbIpAddress, dialerId)
	incValue := RedisIncr(dialerCampaignCountKey)
	fmt.Println("IncrementedOnGoingCampaignCount: ", incValue)
}

func DecrementOnGoingCampaignCount() {
	dialerCampaignCountKey := fmt.Sprintf("DialerOnGoingCampaignCount:%s:%s", lbIpAddress, dialerId)
	dncValue := RedisIncrBy(dialerCampaignCountKey, -1)
	fmt.Println("DecrementOnGoingCampaignCount: ", dncValue)
}

func DecrementOnGoingCampaignCountOther(oDialerId string) {
	dialerCampaignCountSKey := fmt.Sprintf("DialerOnGoingCampaignCount:*:%s", oDialerId)
	searchdialer := RedisSearchKeys(dialerCampaignCountSKey)
	if len(searchdialer) > 0 {
		dncValue := RedisIncrBy(searchdialer[0], -1)
		fmt.Println("DecrementOnGoingCampaignCountOther: ", dncValue)
	}
}

//----------Campaign Status-----------------------
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

func RemoveCampaignStatusOther(oDialerId, campaignId string, company, tenant int) {
	campaignStatusKey := fmt.Sprintf("CampaignStatus:%s:%d:%d:%s", oDialerId, company, tenant, campaignId)
	result := RedisRemove(campaignStatusKey)
	fmt.Println("RemoveCampaignStatusOther CampaignId: ", campaignStatusKey, " Result: ", result)
}

//----------Campaign-----------------------
func AddCampaignToDialer(campaignD Campaign) {
	campaignKey := fmt.Sprintf("Campaign:%s:%d:%d:%d", dialerId, campaignD.CompanyId, campaignD.TenantId, campaignD.CampaignId)
	searchCamp := fmt.Sprintf("Campaign:*:%d:%d:%d", campaignD.CompanyId, campaignD.TenantId, campaignD.CampaignId)
	existingKeys := RedisSearchKeys(searchCamp)

	companyToken := fmt.Sprintf("%d:%d", campaignD.TenantId, campaignD.CompanyId)
	scheduleId := strconv.Itoa(campaignD.CampScheduleInfo[0].ScheduleId)
	timeZone := GetTimeZoneFroSchedule(companyToken, scheduleId)

	fmt.Println("Add Time Zone::", timeZone)
	campaignD.TimeZone = timeZone

	if len(existingKeys) == 0 {
		campaignJson, _ := json.Marshal(campaignD)
		result := RedisAdd(campaignKey, string(campaignJson))
		fmt.Println("Add Campaign to Redis: ", campaignKey, " Result: ", result)
		if result == "OK" {
			campIdStr := strconv.Itoa(campaignD.CampaignId)
			channelCountStr := strconv.Itoa(campaignD.CampConfigurations.ChannelConcurrency)
			//SetCampaignTimeZone(campIdStr, campaignD.CompanyId, campaignD.TenantId, timeZone)
			IncrementOnGoingCampaignCount()
			SetCampChannelMaxLimitDirect(campIdStr, channelCountStr)
			AddCampaignCallbackConfigInfo(campaignD.CompanyId, campaignD.TenantId, campaignD.CampaignId, campaignD.CampConfigurations.ConfigureId)
			SetCampaignStatus(campIdStr, "Start", campaignD.CompanyId, campaignD.TenantId)
			UpdateCampaignStartStatus(campaignD.CompanyId, campaignD.TenantId, campIdStr)
		}
	} else {
		splitVals := strings.Split(existingKeys[0], ":")
		preDialerId := splitVals[1]
		campIdStr := strconv.Itoa(campaignD.CampaignId)
		RemoveCampaignFromOtherDialer(preDialerId, campIdStr, campaignD.CompanyId, campaignD.TenantId)

		campaignJson, _ := json.Marshal(campaignD)
		result := RedisAdd(campaignKey, string(campaignJson))
		fmt.Println("Add Campaign to Redis: ", campaignKey, " Result: ", result)
		if result == "OK" {
			//SetCampaignTimeZone(campIdStr, campaignD.CompanyId, campaignD.TenantId, timeZone)
			IncrementOnGoingCampaignCount()
			SetCampaignStatus(campIdStr, "Resume", campaignD.CompanyId, campaignD.TenantId)
			UpdateCampaignStartStatus(campaignD.CompanyId, campaignD.TenantId, campIdStr)
		}
	}
}

func GetCampaign(company, tenant, campaignId int) (Campaign, bool) {
	campaignKey := fmt.Sprintf("Campaign:%s:%d:%d:%d", dialerId, company, tenant, campaignId)
	isExists := RedisCheckKeyExist(campaignKey)
	if isExists {
		campJson := RedisGet(campaignKey)
		var camp Campaign
		json.Unmarshal([]byte(campJson), &camp)
		return camp, isExists
	} else {
		return Campaign{}, isExists
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
		RemoveCampChannelMaxLimit(campaignId)
		RemoveCampaignConnectedCount(company, tenant, campaignId)
		RemoveCampaignDialCount(company, tenant, campaignId)
		RemoveCampConcurrentChannelCount(campaignId)
		RemoveCampaignCallbackConfigInfo(company, tenant, campaignId)
	}
}

func RemoveCampaignFromOtherDialer(oDialerId, campaignId string, company, tenant int) {
	campaignKey := fmt.Sprintf("Campaign:%s:%d:%d:%s", oDialerId, company, tenant, campaignId)
	result := RedisRemove(campaignKey)
	fmt.Println("Remove Campaign from Redis: ", campaignKey, " Result: ", result)
	if result == true {
		DecrementOnGoingCampaignCountOther(oDialerId)
		RemoveCampaignStatusOther(oDialerId, campaignId, company, tenant)
	}
}

func SetCampaignTimeZone(campaignId string, company, tenant int, timeZone string) {
	campaignKey := fmt.Sprintf("Campaign:%s:%d:%d:%s", dialerId, company, tenant, campaignId)

	campJson := RedisGet(campaignKey)
	var campObj Campaign
	json.Unmarshal([]byte(campJson), &campObj)

	campObj.TimeZone = timeZone

	campaignJson, _ := json.Marshal(campObj)

	result := RedisAdd(campaignKey, string(campaignJson))
	fmt.Println("Add Campaign to Redis: ", campaignKey, " Result: ", result)
	if result == "OK" {
		fmt.Println("Update Campaign TimeZone success")
	}
}

//----------Campaign Manager Service-----------------------
func RequestCampaign(requestCount int) []Campaign {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RequestCampaign", r)
		}
	}()
	//Request campaign from Campaign Manager service
	campaignDetails := make([]Campaign, 0)

	client := &http.Client{}

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaigns/State/Pending/%d", CreateHost(campaignServiceHost, campaignServicePort), requestCount)
	fmt.Println("Start RequestCampaign request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", jwtToken)
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in UpdateCampaignStatus", r)
		}
	}()
	//Send CampaignStatus to Campaign Manager
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	fmt.Println("Start UpdateCampaignStatus Auth: ", internalAuthToken, " CampaignId: ", campaignId, " DialerId: ", dialerId)
	client := &http.Client{}

	currentState := GetCampaignStatus(campaignId, company, tenant)
	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaign/%s/Operations/State/%s/%s", CreateHost(campaignServiceHost, campaignServicePort), campaignId, dialerId, currentState)
	fmt.Println("Start UpdateCampaignStatus request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
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
					break
				case "Pause":
					SetCampaignStatus(campIdStr, "Pause", company, tenant)
					break
				case "Resume":
					SetCampaignStatus(campIdStr, "Resume", company, tenant)
					break
				case "End":
					SetCampaignStatus(campIdStr, "End", company, tenant)
					break
				default:
					break
				}
			}
		}
	}
}

func UpdateCampaignStartStatus(company, tenant int, campaignId string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in UpdateCampaignStartStatus", r)
		}
	}()
	//Send CampaignStatus to Campaign Manager
	state := CampaignStart{}
	camIdInt, _ := strconv.Atoi(campaignId)
	state.CampaignId = camIdInt
	state.DialerId = dialerId

	jsonData, _ := json.Marshal(state)

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaign/%s/Operations/%s", CreateHost(campaignServiceHost, campaignServicePort), campaignId, dialerId)
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
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

//----------Campaign Channel Max Limit-----------------------
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

func SetCampChannelMaxLimitDirect(campaignId, channelcount string) {
	cmcl := fmt.Sprintf("CampaignMaxCallLimit:%s", campaignId)
	RedisSet(cmcl, channelcount)
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

func RemoveCampChannelMaxLimit(campaignId string) {
	cmcl := fmt.Sprintf("CampaignMaxCallLimit:%s", campaignId)
	RedisRemove(cmcl)
}

//----------Campaign Dial Count-----------------------
func IncrCampaignDialCount(company, tenant int, campaignId string) {
	cmcl := fmt.Sprintf("CampaignDialCount:%d:%d:%s", company, tenant, campaignId)
	RedisIncr(cmcl)
}

func GetCampaignDialCount(company, tenant int, campaignId string) int {
	cmcl := fmt.Sprintf("CampaignDialCount:%d:%d:%s", company, tenant, campaignId)
	value := RedisGet(cmcl)
	count, _ := strconv.Atoi(value)
	return count
}

func RemoveCampaignDialCount(company, tenant int, campaignId string) {
	cmcl := fmt.Sprintf("CampaignDialCount:%d:%d:%s", company, tenant, campaignId)
	RedisRemove(cmcl)
}

//----------Campaign Connected Count-----------------------
func IncrCampaignConnectedCount(company, tenant int, campaignId string) {
	cmcl := fmt.Sprintf("CampaignConnectedCount:%d:%d:%s", company, tenant, campaignId)
	RedisIncr(cmcl)
}

func GetCampaignConnectedCount(company, tenant int, campaignId string) int {
	cmcl := fmt.Sprintf("CampaignConnectedCount:%d:%d:%s", company, tenant, campaignId)
	value := RedisGet(cmcl)
	count, _ := strconv.Atoi(value)
	return count
}

func RemoveCampaignConnectedCount(company, tenant int, campaignId string) {
	cmcl := fmt.Sprintf("CampaignConnectedCount:%d:%d:%s", company, tenant, campaignId)
	RedisRemove(cmcl)
}

//----------Run Campaign-----------------------
func StartCampaign(campaignId, dialoutMec, CampaignChannel, camClass, camType, camCategory, scheduleId, camScheduleId, resourceServerId, extention, defaultAni string, company, tenant, campaignMaxChannelCount int) {
	emtAppoinment := Appoinment{}
	defResourceServerInfo := ResourceServerInfo{}
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)

	appment, _, appmntEndTime := CheckAppoinmentForCampaign(internalAuthToken, scheduleId)

	resourceServerInfos := GetResourceServerInfo(company, tenant, resourceServerId, CampaignChannel)

	if appment != emtAppoinment && resourceServerInfos != defResourceServerInfo {
		campStatus := GetCampaignStatus(campaignId, company, tenant)
		if campStatus == "Start" {
			LoadInitialNumberSet(company, tenant, campaignId, camScheduleId)
		}

		SetCampaignStatus(campaignId, "Running", company, tenant)
		maxChannelLimitStr := strconv.Itoa(campaignMaxChannelCount)
		SetCampChannelMaxLimitDirect(campaignId, maxChannelLimitStr)

		//endTime, _ := time.Parse(layout1, appment.EndTime)
		//timeNow := time.Now().UTC()
		//appmntEndTime := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), 0, time.UTC)

		for {
			campStatus = GetCampaignStatus(campaignId, company, tenant)
			if campStatus == "Running" {
				tm := time.Now().UTC()
				fmt.Println("endTime: ", appmntEndTime.String())
				fmt.Println("timeNW: ", tm.String())
				if tm.Before(appmntEndTime) {
					cchannelCountS, cchannelCountC := GetConcurrentChannelCount(resourceServerInfos.ResourceServerId, campaignId)
					maxChannelLimit := GetMaxChannelLimit(resourceServerInfos.ResourceServerId)
					maxCampaignChannelLimit := GetCampMaxChannelLimit(campaignId)
					fmt.Println("resourceServerInfos.CallServerId: ", resourceServerInfos.ResourceServerId)
					fmt.Println("MaxCallServerChannelLimit: ", maxChannelLimit)
					fmt.Println("maxCampaignChannelLimit: ", maxCampaignChannelLimit)
					fmt.Println("ConcurrentResourceServerChannel: ", cchannelCountS)
					fmt.Println("ConcurrentCampaignChannel: ", cchannelCountC)

					if cchannelCountS < maxChannelLimit && cchannelCountC < maxCampaignChannelLimit {
						number, tryCount, numExtraData := GetNumberToDial(company, tenant, campaignId, camScheduleId)
						if number == "" {
							numberCount := GetNumberCount(company, tenant, campaignId, camScheduleId)
							if numberCount == 0 {
								SetCampaignStatus(campaignId, "End", company, tenant)
								RemoveCampaignFromDialer(campaignId, company, tenant)
								return
							}
						} else {
							switch CampaignChannel {
							case "CALL":
								//trunkCode, ani, dnis := "OutTrunk001", defaultAni, number
								trunkCode, ani, dnis := GetTrunkCode(internalAuthToken, defaultAni, number)
								uuid := GetUuid()
								if trunkCode != "" && uuid != "" {
									switch dialoutMec {
									case "BLAST":
										go DialNumber(company, tenant, resourceServerInfos, campaignId, uuid, ani, trunkCode, dnis, tryCount, extention)
										break
									case "FIFO":
										go DialNumberFIFO(company, tenant, resourceServerInfos, campaignId, uuid, ani, trunkCode, dnis, extention)
										break
									case "PREVIEW":
										go AddPreviewDialRequest(company, tenant, resourceServerInfos, campaignId, dialoutMec, uuid, ani, trunkCode, dnis, numExtraData, tryCount, extention)
										break
									case "AGENT":
										go AddAgentDialRequest(company, tenant, resourceServerInfos, campaignId, dialoutMec, uuid, ani, trunkCode, dnis, numExtraData, tryCount, extention)
										break
									}
								}
								break
							case "SMS":
								//message := RequestCampaignAddtionalData(company, tenant, campaignId, "SMS", "mode1", "BLAST")
								//if message != "" {
								//	go SendSms(company, tenant, resourceServerInfos, campaignId, camClass, camType, camCategory, defaultAni, message, number)
								//}
								break
							case "EMAIL":
								defEmailInfo := EmailAdditionalData{}
								emailInfo := RequestEmailInfo(company, tenant, campaignId)
								if emailInfo != defEmailInfo {
									go SendEmail(company, tenant, resourceServerInfos, campaignId, camClass, camType, camCategory, emailInfo.FromAddresss, emailInfo.Subject, emailInfo.Body, number)
								}
								break
							}
							time.Sleep(100 * time.Millisecond)
						}
					} else {
						fmt.Println("dialer waiting...")
						time.Sleep(500 * time.Millisecond)
					}
				} else {
					SetCampaignStatus(campaignId, "PauseByDialer", company, tenant)
					SetCampChannelMaxLimitDirect(campaignId, "0")
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
		SetCampaignStatus(campaignId, "Waiting for Appoinment", company, tenant)
		return
	}
}
