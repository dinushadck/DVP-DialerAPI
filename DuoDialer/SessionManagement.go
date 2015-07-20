package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func InitiateSessionInfo(company, tenant int, tryCount, campaignId, sessionId, reason, dialerStatus, dialTime string) {
	companyStr := strconv.Itoa(company)
	tenantStr := strconv.Itoa(tenant)

	data := make(map[string]string)
	data["companyId"] = companyStr
	data["tenantId"] = tenantStr
	data["sessionId"] = sessionId
	data["dialerId"] = dialerId
	data["campaignId"] = campaignId
	data["dialtime"] = dialTime
	data["reason"] = reason
	data["dialerStatus"] = dialerStatus
	data["tryCount"] = tryCount
	hashKey := fmt.Sprintf("sessionInfo:%s", sessionId)
	RedisHashSetMultipleField(hashKey, data)
}

func SetSessionInfo(sessionId, filed, value string) {
	hashKey := fmt.Sprintf("sessionInfo:%s", sessionId)
	RedisHashSetField(hashKey, filed, value)
}

func UploadSessionInfo(sessionId string) {
	hashKey := fmt.Sprintf("sessionInfo:%s", sessionId)
	sessionInfo := RedisHashGetAll(hashKey)
	bytes, err := json.Marshal(sessionInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	text := string(bytes)
	fmt.Println(text)

	//upload to campaign service
	//if success remove hashInfo
	//RedisRemove(hashKey)
}
