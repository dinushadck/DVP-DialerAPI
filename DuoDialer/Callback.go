package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//----------Campaign Manager Service-----------------------
func RequestCampaignCallbackConfig(tenant, company, configureId int) ([]CampaignCallbackInfo, bool) {
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
		return campaignCallbackInfo, false
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var campaignCallbackResult CampaignCallbackInfoResult
	json.Unmarshal(response, &campaignCallbackResult)
	if campaignCallbackResult.IsSuccess == true {
		if len(campaignCallbackResult.Result) > 0 {
			camRes := campaignCallbackResult.Result[0]
			for _, config := range camRes.CampCallbackConfigurations {
				campaignCallbackInfo = append(campaignCallbackInfo, config)
			}
			return campaignCallbackInfo, camRes.AllowCallBack
		}
	}
	return campaignCallbackInfo, false
}

func UploadCallbackInfo(company, tenant, callBackCount, callbackInterval int, campaignId, contactId string) {
	callback := CampaignCallback{}
	camIdInt, _ := strconv.Atoi(campaignId)
	tmNow := time.Now()
	secCount := tmNow.Second() + callbackInterval

	callback.CampaignId = camIdInt
	callback.CallBackCount = callBackCount
	callback.ContactId = contactId
	callback.DialoutTime = time.Date(tmNow.Year(), tmNow.Month(), tmNow.Day(), tmNow.Hour(), tmNow.Minute(), secCount, 0, tmNow.Location())

	jsonData, _ := json.Marshal(callback)

	serviceurl := fmt.Sprintf("%s/CampaignManager/Callback", campaignService)
	authToken := fmt.Sprintf("%d#%d", tenant, company)
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

//----------Campaign Callback Info-----------------------
func AddCampaignCallbackConfigInfo(company, tenant, campaignId, configureId int) {
	callbackConfigList, AllowCallback := RequestCampaignCallbackConfig(tenant, company, configureId)
	allowCallbackKey := fmt.Sprintf("CampaignAllowCallback:%d:%d:%d", company, tenant, campaignId)
	isSet := RedisSet(allowCallbackKey, strconv.FormatBool(AllowCallback))
	if isSet == "OK" {
		for _, callbackConfig := range callbackConfigList {
			hashKey := fmt.Sprintf("CampaignCallbackConfig:%d:%d:%d:%s", company, tenant, campaignId, callbackConfig.CampCallBackReasons.Reason)
			data := make(map[string]string)
			data["MaxCallbackCount"] = strconv.Itoa(callbackConfig.MaxCallBackCount)
			data["CallbackInterval"] = strconv.Itoa(callbackConfig.CallBackInterval)
			RedisHashSetMultipleField(hashKey, data)
		}
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

func GetCallbackDetails(company, tenant int, campaignId, reason string) (maxCallbackCount, callbackInterval int, isReasonExists bool) {
	hashKey := fmt.Sprintf("CampaignCallbackConfig:%d:%d:%s:%s", company, tenant, campaignId, reason)
	isReasonExists = RedisCheckKeyExist(hashKey)
	if isReasonExists {
		callbackIntervalStr := RedisHashGetAll(hashKey)
		callbackInterval, _ := strconv.Atoi(callbackIntervalStr["CallbackInterval"])
		maxCallbackCount, _ := strconv.Atoi(callbackIntervalStr["MaxCallbackCount"])
		return maxCallbackCount, callbackInterval, isReasonExists
	}
	return 0, 0, isReasonExists
}

func GetAllowCallback(company, tenant int, campaignId string) bool {
	hashKey := fmt.Sprintf("CampaignAllowCallback:%d:%d:%s", company, tenant, campaignId)
	callbackIntervalStr := RedisGet(hashKey)
	allowCallback, _ := strconv.ParseBool(callbackIntervalStr)
	return allowCallback
}

func RemoveCampaignCallbackConfigInfo(company, tenant int, campaignId string) {
	searchKey := fmt.Sprintf("CampaignCallbackConfig:%d:%d:%s:*", company, tenant, campaignId)
	allowCallbackKey := fmt.Sprintf("CampaignAllowCallback:%d:%d:%s", company, tenant, campaignId)
	RedisRemove(allowCallbackKey)
	callbackKeys := RedisSearchKeys(searchKey)
	for _, key := range callbackKeys {
		RedisRemove(key)
	}
}

func AddPhoneNumberToCallback(company, tenant int, campaignId, sessionId, disConnectkReason string) {
	isAllowCallback := GetAllowCallback(company, tenant, campaignId)
	if isAllowCallback == true {
		maxCallbackCount, callbackInterval, isReasonExists := GetCallbackDetails(company, tenant, campaignId, disConnectkReason)
		if isReasonExists {
			number, tryCount := GetPhoneNumberAndTryCount(sessionId)
			if maxCallbackCount > 0 && number != "" && tryCount > 0 && tryCount < maxCallbackCount {
				go UploadCallbackInfo(company, tenant, tryCount, callbackInterval, campaignId, number)
			}
		}
	}
}
