package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//----------Campaign Manager Service-----------------------
func RequestCampaignCallbackConfig(tenant, company, configureId int) ([]CampaignCallbackInfo, bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RequestCampaignCallbackConfig", r)
		}
	}()
	//Request campaign from Campaign Manager service
	campaignCallbackInfo := make([]CampaignCallbackInfo, 0)
	authToken := fmt.Sprintf("%d#%d", tenant, company)
	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaign/Configuration/%d/all", CreateHost(campaignServiceHost, campaignServicePort), configureId)
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

//----------CallbackServer Self Host-----------------------
func UploadCallbackInfo(company, tenant int, callbackTime time.Time, campaignId, cbClass, cbType, cbCategory, cbUrl, cbObj string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in UploadCallbackInfo", r)
		}
	}()
	callback := CampaignCallback{}

	callback.Company = company
	callback.Tenant = tenant
	callback.Class = cbClass
	callback.Type = cbType
	callback.Category = cbCategory
	callback.DialoutTime = callbackTime
	callback.CallbackUrl = cbUrl
	callback.CallbackObj = cbObj
	callback.CampaignId = campaignId

	jsonData, _ := json.Marshal(callback)

	serviceurl := fmt.Sprintf("http://%s/CallbackServerSelfHost/Callback/AddCallback", CreateHost(callbackServerHost, callbackServerPort))
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

//func GetMaxCallbackCount(company, tenant int, campaignId, reason string) int {
//	hashKey := fmt.Sprintf("CampaignCallbackConfig:%d:%d:%s:%s", company, tenant, campaignId, reason)
//	callbackCountStr := RedisHashGetField(hashKey, "MaxCallbackCount")
//	maxCallbackCount, _ := strconv.Atoi(callbackCountStr)
//	return maxCallbackCount
//}

//func GetCallbackInterval(company, tenant int, campaignId, reason string) int {
//	hashKey := fmt.Sprintf("CampaignCallbackConfig:%d:%d:%s:%s", company, tenant, campaignId, reason)
//	callbackIntervalStr := RedisHashGetField(hashKey, "CallbackInterval")
//	callbackInterval, _ := strconv.Atoi(callbackIntervalStr)
//	return callbackInterval
//}

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

func ValidateDisconnectReason(disconnectReason string) bool {
	confKey := fmt.Sprintf("CallbackReason:%s", disconnectReason)
	return RedisCheckKeyExist(confKey)
}

func AddPhoneNumberToCallback(company, tenant, tryCount, campaignId, phoneNumber, disConnectkReason string) {
	fmt.Println("start AddPhoneNumberToCallback")
	_company, _ := strconv.Atoi(company)
	_tenant, _ := strconv.Atoi(tenant)
	_tryCount, _ := strconv.Atoi(tryCount)
	isAllowCallback := GetAllowCallback(_company, _tenant, campaignId)
	//isAllowCallback = true
	isdisconnectReasonAllowed := ValidateDisconnectReason(disConnectkReason)
	//isdisconnectReasonAllowed = true
	if isAllowCallback && isdisconnectReasonAllowed {
		maxCallbackCount, callbackInterval, isReasonExists := GetCallbackDetails(_company, _tenant, campaignId, disConnectkReason)
		//maxCallbackCount, callbackInterval, isReasonExists = 3, 50, true
		if isReasonExists {
			if maxCallbackCount > 0 && phoneNumber != "" && _tryCount > 0 && _tryCount < maxCallbackCount {
				camIdInt, _ := strconv.Atoi(campaignId)

				campaignInfo, isCamExists := GetCampaign(_company, _tenant, camIdInt)
				if isCamExists {
					tmNow := time.Now()
					tmNowUTC := tmNow.UTC()
					secCount := tmNow.Second() + callbackInterval
					secCountUTC := tmNowUTC.Second() + callbackInterval
					callbackTime := time.Date(tmNow.Year(), tmNow.Month(), tmNow.Day(), tmNow.Hour(), tmNow.Minute(), secCount, 0, time.Local)
					callbackTimeUTC := time.Date(tmNowUTC.Year(), tmNowUTC.Month(), tmNowUTC.Day(), tmNowUTC.Hour(), tmNowUTC.Minute(), secCountUTC, 0, time.UTC)

					tempCampaignEndDate, _ := time.Parse(layout1, campaignInfo.CampConfigurations.EndDate)
					campaignEndDate := time.Date(tempCampaignEndDate.Year(), tempCampaignEndDate.Month(), tempCampaignEndDate.Day(), tempCampaignEndDate.Hour(), tempCampaignEndDate.Minute(), tempCampaignEndDate.Second(), 0, time.UTC)

					if campaignEndDate.After(callbackTimeUTC) {
						scheduleIdStr := strconv.Itoa(campaignInfo.CampScheduleInfo[0].ScheduleId)
						validateAppoinment := CheckAppoinmentForCallback(_company, _tenant, scheduleIdStr, callbackTimeUTC)
						if validateAppoinment {
							callbackObj := CampaignCallbackObj{}
							callbackObj.CampaignId = camIdInt
							callbackObj.CallbackClass = "DIALER"
							callbackObj.CallbackType = "CALLBACK"
							callbackObj.CallbackCategory = "INTERNAL"
							callbackObj.CallBackCount = _tryCount
							callbackObj.ContactId = phoneNumber
							callbackObj.DialoutTime = callbackTime

							dialerAPIUrl := fmt.Sprintf("http://%s", CreateHost(lbIpAddress, lbPort))
							path := fmt.Sprintf("DVP/DialerAPI/ResumeCallback")

							u, _ := url.Parse(dialerAPIUrl)
							u.Path += path

							fmt.Println(u.String())
							cbUrl := u.String()

							jsonData, _ := json.Marshal(callbackObj)
							go UploadCallbackInfo(_company, _tenant, callbackTimeUTC, campaignId, "DIALER", "CALLBACK", "INTERNAL", cbUrl, string(jsonData))
						}
					}
				}
			}
		}
	}
}

func ResumeCampaignCallback(company, tenant, callbackCount, campaignId int, number string) {
	fmt.Println("Start ResumeCampaignCallback")
	campaignIdStr := strconv.Itoa(campaignId)
	_tryCount := callbackCount + 1
	campaign, isCampaignExists := GetCampaign(company, tenant, campaignId)
	if isCampaignExists {
		camScheduleStr := strconv.Itoa(campaign.CampScheduleInfo[0].CamScheduleId)
		numberWithTryCount := fmt.Sprintf("%s:%d", number, _tryCount)
		AddNumberToFront(company, tenant, campaignIdStr, camScheduleStr, numberWithTryCount)
	}
}
