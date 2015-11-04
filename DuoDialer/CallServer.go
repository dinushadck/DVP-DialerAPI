package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
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

	var clusterConfigApiResult ClusterConfigApiResult
	json.Unmarshal(response, &clusterConfigApiResult)
	if clusterConfigApiResult.IsSuccess == true {
		for _, callSvr := range clusterConfigApiResult.Result {
			activeCallServers = append(activeCallServers, callSvr)
			if len(activeCallServers) == 1 {
				return activeCallServers[0]
			} else if len(activeCallServers) > 1 {
				return activeCallServers[rand.Intn(len(activeCallServers))]
			} else {
				return CallServerResult{}
			}
		}
	}
	return CallServerResult{}
}

func RegisterCallServer(company, tenant int) CallServerInfo {
	defCallServerInfo := CallServerInfo{}

	//Get CallServer info
	pickedCallServer := GetCallserverInfo(company, tenant)
	if pickedCallServer.InternalMainIP != "" {
		callServerIdStr := strconv.Itoa(pickedCallServer.id)
		cs := CallServerInfo{}
		cs.CallServerId = callServerIdStr
		cs.MaxChannelCount = 50
		cs.Url = fmt.Sprintf("%s", CreateHost(pickedCallServer.InternalMainIP, callServerPort))

		callServerKey := fmt.Sprintf("CallServer:%s", cs.CallServerId)
		callServerjson, _ := json.Marshal(cs)
		addResult := RedisSet(callServerKey, string(callServerjson))

		if addResult == "OK" {
			csck := fmt.Sprintf("CallServerConcurrentCalls:%s", cs.CallServerId)
			csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", cs.CallServerId)
			countStr := strconv.Itoa(cs.MaxChannelCount)
			RedisSet(csck, "0")
			RedisSet(csmcl, countStr)
			return cs
		}
	}

	return defCallServerInfo
}

func GetCallServerInfo(company, tenant int, serverId string) CallServerInfo {
	callServerKey := fmt.Sprintf("CallServer:%s", serverId)
	csString := RedisGet(callServerKey)
	if csString != "" {
		var callServerInfo CallServerInfo
		json.Unmarshal([]byte(csString), &callServerInfo)
		return callServerInfo
	} else {
		return RegisterCallServer(company, tenant)
	}
}

func GetConcurrentChannelCount(serverId, campaignId string) (concurrentOnServer, concurrentOnCamp int) {
	csckC := fmt.Sprintf("CallServerConcurrentCalls:%s:%s", serverId, campaignId)
	csck := fmt.Sprintf("CallServerConcurrentCalls:%s", serverId)
	channelCountC := RedisGet(csckC)
	fmt.Println("RedisGet channelCountC: ", channelCountC)

	if channelCountC == "" {
		RedisSet(csckC, "0")
		channelCountC = "0"
	}

	channelCountS := RedisGet(csck)
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
	csckC := fmt.Sprintf("CallServerConcurrentCalls:%s:%s", serverId, campaignId)
	csck := fmt.Sprintf("CallServerConcurrentCalls:%s", serverId)
	RedisIncr(csckC)
	RedisIncr(csck)
}

func DecrConcurrentChannelCount(serverId, campaignId string) {
	csckC := fmt.Sprintf("CallServerConcurrentCalls:%s:%s", serverId, campaignId)
	csck := fmt.Sprintf("CallServerConcurrentCalls:%s", serverId)
	csckCExists := RedisCheckKeyExist(csckC)
	csckExists := RedisCheckKeyExist(csck)

	if csckCExists == true {
		RedisIncrBy(csckC, -1)
	}
	if csckExists == true {
		RedisIncrBy(csck, -1)
	}
}

func IncrMaxLimit(company, tenant int, serverId string) {
	callServerKey := fmt.Sprintf("CallServer:%s", serverId)
	csString := RedisGet(callServerKey)
	if csString == "" {
		RegisterCallServer(company, tenant)
	}

	csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", serverId)
	RedisIncr(csmcl)
}

func DecrMaxLimit(serverId string) {
	csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", serverId)
	decValue := RedisIncrBy(csmcl, -1)
	if decValue < 0 {
		RedisSet(csmcl, "0")
	}
}

func GetMaxChannelLimit(serverId string) int {
	csmcl := fmt.Sprintf("CallServerMaxCallLimit:%s", serverId)
	maxChannelCount := RedisGet(csmcl)
	value, err := strconv.Atoi(maxChannelCount)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	} else {
		return value
	}
}

func RemoveCampConcurrentChannelCount(campaignId string) {
	SKcsckC := fmt.Sprintf("CallServerConcurrentCalls:*:%s", campaignId)
	sResult := RedisSearchKeys(SKcsckC)
	if len(sResult) > 0 {
		RedisRemove(sResult[0])
	}
}
