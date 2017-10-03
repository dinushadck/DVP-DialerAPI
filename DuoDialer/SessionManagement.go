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

//Initiate dial session for a number
func InitiateSessionInfo(company, tenant, sessionExprTime int, sclass, stype, scategory, tryCount, campaignId, scheduleId, campaignName, sessionId, number, reason, dialerStatus, dialTime, serverId string) {
	companyStr := strconv.Itoa(company)
	tenantStr := strconv.Itoa(tenant)
	sessionExprTimeStr := strconv.Itoa(sessionExprTime)

	data := make(map[string]string)
	data["Class"] = sclass
	data["Type"] = stype
	data["Category"] = scategory
	data["CompanyId"] = companyStr
	data["TenantId"] = tenantStr
	data["SessionId"] = sessionId
	data["Number"] = number
	data["DialerId"] = dialerId
	data["CampaignId"] = campaignId
	data["ScheduleId"] = scheduleId
	data["CampaignName"] = campaignName
	data["Dialtime"] = dialTime
	data["ChannelCreatetime"] = ""
	data["ChannelAnswertime"] = ""
	data["ServerId"] = serverId
	data["Reason"] = reason
	data["DialerStatus"] = dialerStatus
	data["TryCount"] = tryCount
	data["ExpireTime"] = sessionExprTimeStr
	hashKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
	RedisHashSetMultipleField(hashKey, data)
	PublishEvent(campaignId, sessionId)
	//RedisHashSetNxField(hashKey, "TryCount", tryCount)
}

func SetSessionInfo(campaignId, sessionId, filed, value string) {
	hashKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
	RedisHashSetField(hashKey, filed, value)
	PublishEvent(campaignId, sessionId)
}

func UploadSessionInfo(campaignId, sessionId string) {
	hashKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
	sessionInfo := RedisHashGetAll(hashKey)
	RedisRemove(hashKey)
	AddPhoneNumberToCallback(sessionInfo["CompanyId"], sessionInfo["TenantId"], sessionInfo["TryCount"], sessionInfo["CampaignId"], sessionInfo["ScheduleId"], sessionInfo["Number"], sessionInfo["Reason"])
	PublishEvent(campaignId, sessionId)
	UploadSessionInfoToCampaignManager(sessionInfo)
}

//Clear timed out sessions from campaign..
func ClearTimeoutChannels(campaignId string) {
	sHashKey := fmt.Sprintf("sessionInfo:%s:*", campaignId)
	ongoingSessions := RedisSearchKeys(sHashKey)
	tn := time.Now()
	for _, session := range ongoingSessions {
		sessionInfo := RedisHashGetAll(session)
		company := sessionInfo["CompanyId"]
		tenant := sessionInfo["TenantId"]
		dtime := sessionInfo["Dialtime"]
		ctime := sessionInfo["ChannelCreatetime"]
		atime := sessionInfo["ChannelAnswertime"]
		sid := sessionInfo["ServerId"]
		cid := sessionInfo["CampaignId"]
		sessionid := sessionInfo["SessionId"]
		category := sessionInfo["Category"]
		resourceId := sessionInfo["Resource"]
		ardsCategory := sessionInfo["ArdsCategory"]
		expierTimeStr := sessionInfo["ExpireTime"]

		expierTime, _ := strconv.ParseFloat(expierTimeStr, 64)

		dtt, _ := time.Parse(layout4, dtime)
		ctt, _ := time.Parse(layout4, ctime)

		if expierTime > 1 {
			if ctime == "" && tn.Sub(dtt).Seconds() > expierTime {
				if category == "PreviewDial" && resourceId != "" && ardsCategory != "" {
					go ClearResourceSlotWhenReject(company, tenant, ardsCategory, resourceId, sessionid)
				}
				DecrConcurrentChannelCount(sid, cid)
				SetSessionInfo(cid, sessionid, "reason", "ChannelCreate timeout")
				go UploadSessionInfo(cid, sessionid)
			} else if atime == "" && ctime != "" && tn.Sub(ctt).Seconds() > expierTime {
				if category == "PreviewDial" && resourceId != "" && ardsCategory != "" {
					go ClearResourceSlotWhenReject(company, tenant, ardsCategory, resourceId, sessionid)
				}
				DecrConcurrentChannelCount(sid, cid)
				SetSessionInfo(cid, sessionid, "reason", "ChannelAnswer timeout")
				go UploadSessionInfo(cid, sessionid)
			}
		}
	}
}

//func GetSpecificSessionFiled(campaignId, sessionId, field string) string {
//	hashKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
//	return RedisHashGetField(hashKey, field)
//}

//func GetPhoneNumberAndTryCount(campaignId, sessionId string) (string, int) {
//	hashKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
//	sessionInfo := RedisHashGetAll(hashKey)
//	number := sessionInfo["Number"]
//	tryCount, _ := strconv.Atoi(sessionInfo["TryCount"])
//	return number, tryCount
//}

//----------------Campaign Manager Service------------------------
func UploadSessionInfoToCampaignManager(sessionInfo map[string]string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in UploadSessionInfoToCampaignManager", r)
		}
	}()
	sessionb, err := json.Marshal(sessionInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	text := string(sessionb)
	fmt.Println(text)
	//upload to campaign service
	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaign/Session", CreateHost(campaignServiceHost, campaignServicePort))
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%s:%s", sessionInfo["TenantId"], sessionInfo["CompanyId"])

	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(sessionb))
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

	body, errb := ioutil.ReadAll(resp.Body)
	//if success remove hashInfo
	if errb != nil {
		fmt.Println(err.Error())
	} else {
		result := string(body)
		fmt.Println("response Body:", result)
	}
}
