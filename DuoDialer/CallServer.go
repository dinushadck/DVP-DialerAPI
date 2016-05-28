package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

//---------------------ClusterConfigService------------------------
func GetCallserverInfo(company, tenant int) CallServerResult {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetCallserverInfo", r)
		}
	}()
	//Request campaign from Campaign Manager service
	activeCallServers := make([]CallServerResult, 0)
	authToken := fmt.Sprintf("%d#%d", tenant, company)

	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CloudConfiguration/CallserversByCompany", CreateHost(clusterConfigServiceHost, clusterConfigServicePort))
	fmt.Println("Start CallserversByCompany request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return CallServerResult{}
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Result: ", string(response))

	var clusterConfigApiResult ClusterConfigApiResult
	json.Unmarshal(response, &clusterConfigApiResult)
	if clusterConfigApiResult.IsSuccess == true {
		for _, callSvr := range clusterConfigApiResult.Result {
			if callSvr.Activate == true {
				activeCallServers = append(activeCallServers, callSvr)
			}
		}
		if len(activeCallServers) == 1 {
			return activeCallServers[0]
		} else if len(activeCallServers) > 1 {
			return activeCallServers[rand.Intn(len(activeCallServers))]
		} else {
			return CallServerResult{}
		}
	}
	return CallServerResult{}
}

func GetSmsServerInfo(company, tenant int) ResourceServerInfo {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetSmsserverInfo", r)
		}
	}()

	var smsServerApiResult ResourceServerInfo
	smsServerApiResult.ResourceServerId = "SMS1"
	smsServerApiResult.Url = "159.203.109.43:1401"
	smsServerApiResult.MaxChannelCount = 5

	return smsServerApiResult
}

func GetEmailServerInfo(company, tenant int) ResourceServerInfo {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetSmsserverInfo", r)
		}
	}()
	authToken := fmt.Sprintf("%d:%d", tenant, company)

	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DuoMessageTemplate/EmailService.svc/Json/getEmailServeDetails/%d/%d", casServerHost, company, tenant)
	fmt.Println("Start GetSmsserverInfo request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return ResourceServerInfo{}
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Result: ", string(response))

	var smsServerApiResult ResourceServerInfo
	var resData string
	conErr1 := json.Unmarshal(response, &resData)
	if conErr1 != nil {
		fmt.Println(conErr1.Error())
	}
	fmt.Println("tempData: ", resData)
	conErr := json.Unmarshal([]byte(resData), &smsServerApiResult)
	if conErr != nil {
		fmt.Println(conErr.Error())
	}
	if smsServerApiResult.ResourceServerId != "" {
		return smsServerApiResult
	}
	return ResourceServerInfo{}
}

func RegisterCallServer(company, tenant int) ResourceServerInfo {
	defResourceServerInfo := ResourceServerInfo{}

	//Get CallServer info
	pickedCallServer := GetCallserverInfo(company, tenant)
	log := fmt.Sprintf("Callserver id: %d :: ip: %s :: CompanyId: %d", pickedCallServer.Id, pickedCallServer.InternalMainIP, pickedCallServer.CompanyId)
	fmt.Println(log)
	if pickedCallServer.InternalMainIP != "" {
		callServerIdStr := strconv.Itoa(pickedCallServer.Id)
		rs := ResourceServerInfo{}
		rs.ResourceServerId = callServerIdStr
		rs.MaxChannelCount = 50
		rs.Url = fmt.Sprintf("%s", CreateHost(pickedCallServer.InternalMainIP, callServerPort))

		resourceServerKey := fmt.Sprintf("ResourceServer:%s", rs.ResourceServerId)
		resourceServerjson, _ := json.Marshal(rs)
		addResult := RedisSet(resourceServerKey, string(resourceServerjson))

		if addResult == "OK" {
			rsck := fmt.Sprintf("ResourceServerConcurrentCalls:%s", rs.ResourceServerId)
			rsmcl := fmt.Sprintf("ResourceServerMaxCallLimit:%s", rs.ResourceServerId)
			countStr := strconv.Itoa(rs.MaxChannelCount)
			RedisSet(rsck, "0")
			RedisSet(rsmcl, countStr)
			return rs
		}
	}

	return defResourceServerInfo
}

func RegisterSmsServer(company, tenant int) ResourceServerInfo {
	defResourceServerInfo := ResourceServerInfo{}

	//Get CallServer info
	pickedSmsServer := GetSmsServerInfo(company, tenant)
	log := fmt.Sprintf("SmsSrver id: %d :: URL: %s", pickedSmsServer.ResourceServerId, pickedSmsServer.Url)
	fmt.Println(log)
	if pickedSmsServer.ResourceServerId != "" {
		resourceServerKey := fmt.Sprintf("ResourceServer:%s", pickedSmsServer.ResourceServerId)
		resourceServerjson, _ := json.Marshal(pickedSmsServer)
		addResult := RedisSet(resourceServerKey, string(resourceServerjson))

		if addResult == "OK" {
			rsck := fmt.Sprintf("ResourceServerConcurrentCalls:%s", pickedSmsServer.ResourceServerId)
			rsmcl := fmt.Sprintf("ResourceServerMaxCallLimit:%s", pickedSmsServer.ResourceServerId)
			countStr := strconv.Itoa(pickedSmsServer.MaxChannelCount)
			RedisSet(rsck, "0")
			RedisSet(rsmcl, countStr)
			return pickedSmsServer
		}
	}

	return defResourceServerInfo
}

func RegisterEmailServer(company, tenant int) ResourceServerInfo {
	defResourceServerInfo := ResourceServerInfo{}

	//Get CallServer info
	pickedEmailServer := GetEmailServerInfo(company, tenant)
	log := fmt.Sprintf("SmsSrver id: %d :: URL: %s", pickedEmailServer.ResourceServerId, pickedEmailServer.Url)
	fmt.Println(log)
	if pickedEmailServer.ResourceServerId != "" {
		resourceServerKey := fmt.Sprintf("ResourceServer:%s", pickedEmailServer.ResourceServerId)
		resourceServerjson, _ := json.Marshal(pickedEmailServer)
		addResult := RedisSet(resourceServerKey, string(resourceServerjson))

		if addResult == "OK" {
			rsck := fmt.Sprintf("ResourceServerConcurrentCalls:%s", pickedEmailServer.ResourceServerId)
			rsmcl := fmt.Sprintf("ResourceServerMaxCallLimit:%s", pickedEmailServer.ResourceServerId)
			countStr := strconv.Itoa(pickedEmailServer.MaxChannelCount)
			RedisSet(rsck, "0")
			RedisSet(rsmcl, countStr)
			return pickedEmailServer
		}
	}

	return defResourceServerInfo
}

func GetResourceServerInfo(company, tenant int, serverId, serverType string) ResourceServerInfo {
	resourceServerKey := fmt.Sprintf("ResourceServer:%s", serverId)
	rsString := RedisGet(resourceServerKey)
	if rsString != "" {
		var resourceServerInfo ResourceServerInfo
		json.Unmarshal([]byte(rsString), &resourceServerInfo)
		return resourceServerInfo
	} else {
		//add swith case to pick server for campanign type eg:- CAll, SMS, Email
		switch strings.ToLower(serverType) {
		case "call":
			return RegisterCallServer(company, tenant)
		case "sms":
			return RegisterSmsServer(company, tenant)
		case "email":
			return RegisterEmailServer(company, tenant)
		}

		return ResourceServerInfo{}
	}
}

func GetConcurrentChannelCount(serverId, campaignId string) (concurrentOnServer, concurrentOnCamp int) {
	rsckC := fmt.Sprintf("ResourceServerConcurrentCalls:%s:%s", serverId, campaignId)
	rsck := fmt.Sprintf("ResourceServerConcurrentCalls:%s", serverId)
	channelCountC := RedisGet(rsckC)
	fmt.Println("RedisGet channelCountC: ", channelCountC)

	if channelCountC == "" {
		RedisSet(rsckC, "0")
		channelCountC = "0"
	}

	channelCountS := RedisGet(rsck)
	fmt.Println("RedisGet channelCountS: ", channelCountS)
	valueC, err := strconv.Atoi(channelCountC)
	valueS, _ := strconv.Atoi(channelCountS)
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0
	} else {
		return valueS, valueC
	}
}

func IncrConcurrentChannelCount(serverId, campaignId string) {
	rsckC := fmt.Sprintf("ResourceServerConcurrentCalls:%s:%s", serverId, campaignId)
	rsck := fmt.Sprintf("ResourceServerConcurrentCalls:%s", serverId)
	RedisIncr(rsckC)
	RedisIncr(rsck)
}

func DecrConcurrentChannelCount(serverId, campaignId string) {
	rsckC := fmt.Sprintf("ResourceServerConcurrentCalls:%s:%s", serverId, campaignId)
	rsck := fmt.Sprintf("ResourceServerConcurrentCalls:%s", serverId)
	rsckCExists := RedisCheckKeyExist(rsckC)
	rsckExists := RedisCheckKeyExist(rsck)

	if rsckCExists == true {
		RedisIncrBy(rsckC, -1)
	}
	if rsckExists == true {
		RedisIncrBy(rsck, -1)
	}
}

//func IncrMaxLimit(company, tenant int, serverId string) {
//	callServerKey := fmt.Sprintf("CallServer:%s", serverId)
//	csString := RedisGet(callServerKey)
//	if csString == "" {
//		RegisterCallServer(company, tenant)
//	}

//	csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", serverId)
//	RedisIncr(csmcl)
//}

//func DecrMaxLimit(serverId string) {
//	csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", serverId)
//	decValue := RedisIncrBy(csmcl, -1)
//	if decValue < 0 {
//		RedisSet(csmcl, "0")
//	}
//}

func GetMaxChannelLimit(serverId string) int {
	rsmcl := fmt.Sprintf("ResourceServerMaxCallLimit:%s", serverId)
	maxChannelCount := RedisGet(rsmcl)
	value, err := strconv.Atoi(maxChannelCount)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	} else {
		return value
	}
}

func RemoveCampConcurrentChannelCount(campaignId string) {
	SKrsckC := fmt.Sprintf("ResourceServerConcurrentCalls:*:%s", campaignId)
	sResult := RedisSearchKeys(SKrsckC)
	if len(sResult) > 0 {
		RedisRemove(sResult[0])
	}
}
