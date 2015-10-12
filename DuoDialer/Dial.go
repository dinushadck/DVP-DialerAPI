package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetUuid() string {
	uuidService := fmt.Sprintf("http://%s:%s/api/create_uuid", callServerHost, callServerPort)
	resp, err := http.Get(uuidService)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	} else {
		defer resp.Body.Close()
		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			return ""
		} else {
			tmx := string(response[:])
			fmt.Println(tmx)
			return tmx
		}
	}
}

func GetTrunkCode(authToken, ani, dnis string) (trunkCode, rAni, rDnis string) {
	fmt.Println("Start GetTrunkCode: ", authToken, ": ", ani, ": ", dnis)
	client := &http.Client{}

	request := fmt.Sprintf("http://%s:%s/DVP/API/6.0/CallRuleApi/CallRule/Outbound/ANI/%s/DNIS/%s", callRuleServiceHost, callRuleServicePort, ani, dnis)
	//request := fmt.Sprintf("%s?ANI=%s&DNIS=%s", callRuleService, ani, dnis)
	fmt.Println("Start GetTrunkCode request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", ""
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var apiResult CallRuleApiResult
	json.Unmarshal(response, &apiResult)
	if apiResult.IsSuccess == true {
		fmt.Println("callRule: ", apiResult.Result.GatewayCode, "ANI: ", apiResult.Result.ANI, "DNIS: ", apiResult.Result.DNIS)
		return apiResult.Result.GatewayCode, apiResult.Result.ANI, apiResult.Result.DNIS
	} else {
		return "", "", ""
	}
}

func Dial(server, params, furl, data string) (*http.Response, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Dial", r)
		}
	}()
	request := fmt.Sprintf("http://%s", server)
	path := fmt.Sprintf("api/originate?")
	param := fmt.Sprintf(" %s%s %s", params, furl, data)

	u, _ := url.Parse(request)
	u.Path += path
	u.Path += param

	fmt.Println(u.String())
	resp, err := http.Get(u.String())
	defer resp.Body.Close()
	return resp, err
}

func HandleDialResponse(resp *http.Response, err error, server CallServerInfo, campaignId, sessionId string) {
	if err != nil {
		DecrConcurrentChannelCount(server.CallServerId, campaignId)
		SetSessionInfo(campaignId, sessionId, "Reason", "dial_failed")
		SetSessionInfo(campaignId, sessionId, "DialerStatus", "failed")
		go UploadSessionInfo(campaignId, sessionId)
		fmt.Println(err.Error())
	}

	if resp != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		tmx := string(response[:])
		fmt.Println(tmx)
		resultInfo := strings.Split(tmx, " ")
		if len(resultInfo) > 0 {
			if resultInfo[0] == "-ERR" {
				DecrConcurrentChannelCount(server.CallServerId, campaignId)

				if len(resultInfo) > 1 {
					reason := resultInfo[1]
					if reason == "" {
						SetSessionInfo(campaignId, sessionId, "Reason", "not_specified")
					} else {
						SetSessionInfo(campaignId, sessionId, "Reason", reason)
					}
				} else {
					SetSessionInfo(campaignId, sessionId, "Reason", "not_specified")
				}
				SetSessionInfo(campaignId, sessionId, "DialerStatus", "not_connected")
				go UploadSessionInfo(campaignId, sessionId)
			} else {
				SetSessionInfo(campaignId, sessionId, "Reason", "dial_success")
				SetSessionInfo(campaignId, sessionId, "DialerStatus", "connected")
			}
		}
	}
}
