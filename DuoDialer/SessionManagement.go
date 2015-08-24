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

func InitiateSessionInfo(company, tenant int, tryCount, campaignId, sessionId, number, reason, dialerStatus, dialTime, serverId string) {
	companyStr := strconv.Itoa(company)
	tenantStr := strconv.Itoa(tenant)

	data := make(map[string]string)
	data["CompanyId"] = companyStr
	data["TenantId"] = tenantStr
	data["SessionId"] = sessionId
	data["Number"] = number
	data["DialerId"] = dialerId
	data["CampaignId"] = campaignId
	data["Dialtime"] = dialTime
	data["ChannelCreatetime"] = dialTime
	data["ChannelAnswertime"] = dialTime
	data["ServerId"] = serverId
	data["Reason"] = reason
	data["DialerStatus"] = dialerStatus
	data["TryCount"] = tryCount
	hashKey := fmt.Sprintf("sessionInfo:%s:%s", dialerId, sessionId)
	RedisHashSetMultipleField(hashKey, data)
}

func SetSessionInfo(sessionId, filed, value string) {
	hashKey := fmt.Sprintf("sessionInfo:%s:%s", dialerId, sessionId)
	RedisHashSetField(hashKey, filed, value)
}

func UploadSessionInfo(sessionId string) {
	hashKey := fmt.Sprintf("sessionInfo:%s:%s", dialerId, sessionId)
	sessionInfo := RedisHashGetAll(hashKey)
	RedisRemove(hashKey)
	sessionb, err := json.Marshal(sessionInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	text := string(sessionb)
	fmt.Println(text)

	//upload to campaign service
	serviceurl := fmt.Sprintf("%s/CampaignManager/Campaign/Session", campaignService)
	authToken := fmt.Sprintf("%s#%s", sessionInfo["TenantId"], sessionInfo["CompanyId"])

	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(sessionb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
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

func ClearTimeoutChannels() {
	sHashKey := fmt.Sprintf("sessionInfo:%s:*", dialerId)
	ongoingSessions := RedisSearchKeys(sHashKey)
	tn := time.Now()
	for _, session := range ongoingSessions {
		sessionInfo := RedisHashGetAll(session)
		dtime := sessionInfo["Dialtime"]
		ctime := sessionInfo["ChannelCreatetime"]
		atime := sessionInfo["ChannelAnswertime"]
		sid := sessionInfo["ServerId"]
		cid := sessionInfo["CampaignId"]
		sessionid := sessionInfo["SessionId"]

		dtt, _ := time.Parse(layout4, dtime)
		ctt, _ := time.Parse(layout4, ctime)
		if ctime == "" && tn.Sub(dtt).Seconds() > 240 {
			DecrConcurrentChannelCount(sid, cid)
			SetSessionInfo(sessionid, "reason", "ChannelCreate timeout")
			go UploadSessionInfo(sessionid)
		} else if atime == "" && ctime != "" && tn.Sub(ctt).Seconds() > 240 {
			DecrConcurrentChannelCount(sid, cid)
			SetSessionInfo(sessionid, "reason", "ChannelAnswer timeout")
			go UploadSessionInfo(sessionid)
		}
	}
}

func GetSpecificSessionFiled(sessionId, field string) string {
	hashKey := fmt.Sprintf("sessionInfo:%s:%s", dialerId, sessionId)
	return RedisHashGetField(hashKey, field)
}

func GetPhoneNumberAndTryCount(sessionId string) (string, int) {
	hashKey := fmt.Sprintf("sessionInfo:%s:%s", dialerId, sessionId)
	sessionInfo := RedisHashGetAll(hashKey)
	number := sessionInfo["Number"]
	tryCount, _ := strconv.Atoi(sessionInfo["TryCount"])
	return number, tryCount
}
