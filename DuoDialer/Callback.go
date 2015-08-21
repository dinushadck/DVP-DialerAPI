package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

//----------Campaign Manager Service-----------------------
func RequestCampaignCallbackConfig(tenant, company, configureId int) []CampaignCallbackInfo {
	//Request campaign from Campaign Manager service
	campaignCallbackInfo := make([]CampaignCallbackInfo, 0)
	authToken := fmt.Sprintf("%d#%d", tenant, company)
	client := &http.Client{}

	request := fmt.Sprintf("%s/CampaignManager/Campaign/Configuration/%d/all", campaignService, configureId)
	fmt.Println("Start RequestCampaignCallbackConfig request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return campaignCallbackInfo
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var campaignCallbackResult CampaignCallbackInfoResult
	json.Unmarshal(response, &campaignCallbackResult)
	if campaignCallbackResult.IsSuccess == true {
		for _, camRes := range campaignCallbackResult.Result {
			for _, config := range camRes.CampCallbackConfigurations {
				campaignCallbackInfo = append(campaignCallbackInfo, config)
			}
		}
	}
	return campaignCallbackInfo
}

//----------Campaign Callback Info-----------------------
func AddCampaignCallbackConfigInfo(company, tenant, campaignId, configureId int) {
	callbackConfigList := RequestCampaignCallbackConfig(tenant, company, configureId)
	for _, callbackConfig := range callbackConfigList {
		hashKey := fmt.Sprintf("CampaignCallbackConfig:%d:%d:%d:%s", company, tenant, campaignId, callbackConfig.CampCallBackReasons.Reason)
		data := make(map[string]string)
		data["MaxCallbackCount"] = strconv.Itoa(callbackConfig.MaxCallBackCount)
		data["CallbackInterval"] = strconv.Itoa(callbackConfig.CallBackInterval)
		RedisHashSetMultipleField(hashKey, data)
	}
}

func GetMaxCallbackCount(company, tenant int, campaignId, reason string) int {
	hashKey := fmt.Sprintf("CampaignCallbackConfig:%d:%d:%s:%s", company, tenant, campaignId, reason)
	callbackCountStr := RedisHashGetField(hashKey, "MaxCallbackCount")
	maxCallbackCount, _ := strconv.Atoi(callbackCountStr)
	return maxCallbackCount
}

func GetCallbackInterval(company, tenant int, campaignId, reason string) int {
	hashKey := fmt.Sprintf("CampaignCallbackConfig:%d:%d:%s:%s", company, tenant, campaignId, reason)
	callbackIntervalStr := RedisHashGetField(hashKey, "CallbackInterval")
	callbackInterval, _ := strconv.Atoi(callbackIntervalStr)
	return callbackInterval
}

func RemoveCampaignCallbackConfigInfo(company, tenant int, campaignId string) {
	searchKey := fmt.Sprintf("CampaignCallbackConfig:%d:%d:%s:*", company, tenant, campaignId)
	callbackKeys := RedisSearchKeys(searchKey)
	for _, key := range callbackKeys {
		RedisRemove(key)
	}
}
