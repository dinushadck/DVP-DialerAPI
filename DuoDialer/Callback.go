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

	"github.com/fatih/color"
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
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaign/Configuration/%d/all", CreateHost(campaignServiceHost, campaignServicePort), configureId)
	DialerLog(fmt.Sprintf("Start RequestCampaignCallbackConfig request: %s", request))
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
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

	fmt.Println("Start UploadCallbackInfo")
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

	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/Callback/AddCallback", CreateHost(callbackServerHost, callbackServerPort))
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))

	fmt.Println("Add callback data:: ", serviceurl, "::", string(jsonData))

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
	fmt.Println("Key:GetAllowCallback:: ", hashKey)
	callbackIntervalStr := RedisGet(hashKey)
	fmt.Println("Result:GetAllowCallback:: ", callbackIntervalStr)
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

func ValidateDisconnectReason(disconnectReason string) (keyExist bool, value string) {
	confKey := fmt.Sprintf("CallbackReason:%s", disconnectReason)
	value = RedisGet(confKey)

	if value == "" {
		keyExist = false
	} else {
		keyExist = true
	}
	return
}

func AddContactToCallback(sessionInfo map[string]string) {
	DialerLog("start AddPhoneNumberToCallback")

	contactsList := []Contact{}

	_ = json.Unmarshal([]byte(sessionInfo["Contacts"]), &contactsList)

	if contactsList != nil && len(contactsList) > 0 {
		company, _ := strconv.Atoi(sessionInfo["CompanyId"])
		tenant, _ := strconv.Atoi(sessionInfo["TenantId"])
		campId, _ := strconv.Atoi(sessionInfo["CampaignId"])
		scheduleId, _ := strconv.Atoi(sessionInfo["ScheduleId"])
		isdisconnectReasonAllowed, hangupGruop := ValidateDisconnectReason(sessionInfo["Reason"])
		DialerLog(fmt.Sprintf("isdisconnectReasonAllowed:: %t", isdisconnectReasonAllowed))
		DialerLog(fmt.Sprintf("hangupGruop:: %s", hangupGruop))
		if isdisconnectReasonAllowed {

			campaignInfo, isCamExists := GetCampaign(company, tenant, campId)
			DialerLog(fmt.Sprintf("isCamExists:: %t", isCamExists))

			if isCamExists {

				scheduleInfo := CampaignShedule{}
				defaultScheduleInfo := CampaignShedule{}
				for _, schedule := range campaignInfo.CampScheduleInfo {
					if schedule.ScheduleId == scheduleId {
						scheduleInfo = schedule
						break
					}
				}

				if scheduleInfo != defaultScheduleInfo {

					location, _ := time.LoadLocation(scheduleInfo.TimeZone)

					tmNow := time.Now().In(location)
					callbackTime := time.Date(tmNow.Year(), tmNow.Month(), tmNow.Day(), tmNow.Hour(), tmNow.Minute(), tmNow.Second(), 0, location)
					fmt.Println("callbackTime:: ", callbackTime)

					//tempCampaignEndDate, _ := time.Parse(layout2, campaignInfo.CampConfigurations.EndDate)
					//campaignEndDate := time.Date(tempCampaignEndDate.Year(), tempCampaignEndDate.Month(), tempCampaignEndDate.Day(), tempCampaignEndDate.Hour(), tempCampaignEndDate.Minute(), tempCampaignEndDate.Second(), 0, location)

					scheduleEndDate := scheduleInfo.EndDate

					fmt.Println("Callback:scheduleEndDate:: ", scheduleEndDate)
					if scheduleEndDate.After(callbackTime) {
						fmt.Println("Start to build CallbackInfo")
						//scheduleIdStr := strconv.Itoa(scheduleInfo.ScheduleId)
						validateAppoinment := CheckAppoinmentForCallback(company, tenant, sessionInfo["ScheduleId"], callbackTime, scheduleInfo.TimeZone)
						fmt.Println("validateAppoinmentFor Callback:: ", validateAppoinment)
						if validateAppoinment {
							//HERE YOU HAVE TO ADD CALLBACK NUMBER TO THE FRONT OF THE NUMBER QUEUE AND RE ADJUST THE CONTACTS ARRAY
							nextContactNum := contactsList[0]
							r := append(contactsList[:0], contactsList[1:]...)
							fmt.Println(r)
							contactDet := ContactsDetails{phone: nextContactNum.contact, contacts: contactsList}
							AddContactToFront(company, tenant, sessionInfo["CampaignId"], contactDet)
							//go UploadCallbackInfo(_company, _tenant, callbackTime, campaignId, "DIALER", "CALLBACK", "INTERNAL", cbUrl, string(jsonData))
						}
					}
				} else {
					fmt.Println("Add Callback Failed, No Valied Schedule Found")
				}

			} else {
				fmt.Println("Add Callback Failed, No Existing Campaign Found")
			}
		}
	} else {
		color.Magenta("NO CONTACTS FOUND FOR RE DIALING")

	}

}

func AddPhoneNumberToCallback(company, tenant, tryCount, campaignId, scheduleId, phoneNumber, disConnectkReason string) {
	fmt.Println("start AddPhoneNumberToCallback")
	_company, _ := strconv.Atoi(company)
	_tenant, _ := strconv.Atoi(tenant)
	_tryCount, _ := strconv.Atoi(tryCount)
	isAllowCallback := GetAllowCallback(_company, _tenant, campaignId)
	fmt.Println("isAllowCallback:: ", isAllowCallback)
	isdisconnectReasonAllowed, hangupGruop := ValidateDisconnectReason(disConnectkReason)
	fmt.Println("isdisconnectReasonAllowed:: ", isdisconnectReasonAllowed)
	fmt.Println("hangupGruop:: ", hangupGruop)
	if isAllowCallback && isdisconnectReasonAllowed {
		maxCallbackCount, callbackInterval, isReasonExists := GetCallbackDetails(_company, _tenant, campaignId, hangupGruop)
		fmt.Println("maxCallbackCount:: ", maxCallbackCount)
		fmt.Println("callbackInterval:: ", callbackInterval)
		fmt.Println("isReasonExists:: ", isReasonExists)
		if isReasonExists {
			if maxCallbackCount > 0 && phoneNumber != "" && _tryCount > 0 && _tryCount <= maxCallbackCount {
				camIdInt, _ := strconv.Atoi(campaignId)
				scheduleIdInt, _ := strconv.Atoi(scheduleId)

				campaignInfo, isCamExists := GetCampaign(_company, _tenant, camIdInt)
				fmt.Println("isCamExists:: ", isCamExists)

				if isCamExists {

					scheduleInfo := CampaignShedule{}
					defaultScheduleInfo := CampaignShedule{}
					for _, schedule := range campaignInfo.CampScheduleInfo {
						if schedule.ScheduleId == scheduleIdInt {
							scheduleInfo = schedule
							break
						}
					}

					if scheduleInfo != defaultScheduleInfo {

						location, _ := time.LoadLocation(scheduleInfo.TimeZone)

						tmNow := time.Now().In(location)
						secCount := tmNow.Second() + callbackInterval
						callbackTime := time.Date(tmNow.Year(), tmNow.Month(), tmNow.Day(), tmNow.Hour(), tmNow.Minute(), secCount, 0, location)
						fmt.Println("callbackTime:: ", callbackTime)

						//tempCampaignEndDate, _ := time.Parse(layout2, campaignInfo.CampConfigurations.EndDate)
						//campaignEndDate := time.Date(tempCampaignEndDate.Year(), tempCampaignEndDate.Month(), tempCampaignEndDate.Day(), tempCampaignEndDate.Hour(), tempCampaignEndDate.Minute(), tempCampaignEndDate.Second(), 0, location)

						scheduleEndDate := scheduleInfo.EndDate

						fmt.Println("Callback:scheduleEndDate:: ", scheduleEndDate)
						if scheduleEndDate.After(callbackTime) {
							fmt.Println("Start to build CallbackInfo")
							//scheduleIdStr := strconv.Itoa(scheduleInfo.ScheduleId)
							validateAppoinment := CheckAppoinmentForCallback(_company, _tenant, scheduleId, callbackTime, scheduleInfo.TimeZone)
							fmt.Println("validateAppoinmentFor Callback:: ", validateAppoinment)
							if validateAppoinment {
								callbackObj := CampaignCallbackObj{}
								callbackObj.CampaignId = campaignId
								callbackObj.CallbackClass = "DIALER"
								callbackObj.CallbackType = "CALLBACK"
								callbackObj.CallbackCategory = "INTERNAL"
								callbackObj.CallBackCount = tryCount
								callbackObj.ContactId = phoneNumber
								callbackObj.DialoutTime = callbackTime

								dialerAPIUrl := fmt.Sprintf("http://%s", CreateHost(lbIpAddress, lbPort))
								path := fmt.Sprintf("DVP/DialerAPI/ResumeCallback")

								u, _ := url.Parse(dialerAPIUrl)
								u.Path += path

								fmt.Println(u.String())
								cbUrl := u.String()

								jsonData, _ := json.Marshal(callbackObj)
								go UploadCallbackInfo(_company, _tenant, callbackTime, campaignId, "DIALER", "CALLBACK", "INTERNAL", cbUrl, string(jsonData))
							}
						}
					} else {
						fmt.Println("Add Callback Failed, No Valied Schedule Found")
					}

				} else {
					fmt.Println("Add Callback Failed, No Existing Campaign Found")
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
