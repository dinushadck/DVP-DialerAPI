package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
)

func GetUuid(callServerHost string) string {
	uuidService := fmt.Sprintf("http://%s/api/create_uuid", callServerHost)
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

func GetTrunkCode(internalAuthToken, ani, dnis, context, businessUnit string) (trunkCode, rAni, rDnis, xGateway string) {
	fmt.Println("Start GetTrunkCode: ", internalAuthToken, ": ", ani, ": ", dnis)
	client := &http.Client{}

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CallRuleApi/CallRule/Outbound/ANI/%s/DNIS/%s/ByBU?Context=%s&BusinessUnit=%s", CreateHost(callRuleServiceHost, callRuleServicePort), ani, dnis, context, businessUnit)
	//request := fmt.Sprintf("%s?ANI=%s&DNIS=%s", callRuleService, ani, dnis)
	fmt.Println("Start GetTrunkCode request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", "", ""
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var apiResult CallRuleApiResult
	json.Unmarshal(response, &apiResult)
	if apiResult.IsSuccess == true {
		fmt.Println("callRule: ", apiResult.Result.GatewayCode, "ANI: ", apiResult.Result.ANI, "DNIS: ", apiResult.Result.DNIS)
		return apiResult.Result.GatewayCode, apiResult.Result.ANI, apiResult.Result.DNIS, apiResult.Result.IpUrl
	} else {
		return "", "", "", ""
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
	//defer resp.Body.Close()
	return resp, err
}

func DialNew(server, params, furl, data string) (*http.Response, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Dial", r)
		}
	}()
	request := fmt.Sprintf("http://%s", server)
	path := fmt.Sprintf("api/originate?")
	param := fmt.Sprintf(" %s %s%s", data, params, furl)

	u, _ := url.Parse(request)
	u.Path += path
	u.Path += param

	fmt.Println(u.String())
	resp, err := http.Get(u.String())
	//defer resp.Body.Close()
	return resp, err
}

func SendCustomerIntegrationData(campaignId, sessionId string) {
	hKey := fmt.Sprintf("sessionInfo:%s:%s", campaignId, sessionId)
	sessionInfo := RedisHashGetAll(hKey)
	color.Magenta(fmt.Sprintf(sessionInfo["IntegrationData"]))

	if sessionInfo != nil && sessionInfo["IntegrationData"] != "" {
		sessionInfo["EventType"] = "CUSTOMER_DISCONNECT"
		go ManageIntegrationData(sessionInfo, "CUSTOMER")
	} else {
		color.Magenta("NO INTEGRATION DATA")
	}
}

func SendAgentIntegrationData(sessionInfo map[string]string) {

	if sessionInfo != nil && sessionInfo["IntegrationData"] != "" {
		sessionInfo["EventType"] = "AGENT_DISCONNECTED"
		go ManageIntegrationData(sessionInfo, "AGENT")
	} else {
		color.Magenta("NO INTEGRATION DATA")
	}
}

func HandleDialResponse(resp *http.Response, err error, server ResourceServerInfo, campaignId, sessionId string) string {
	if err != nil {
		color.Red("=============HANDLE DIAL RESPONSE RETURNED ERROR=============")
		SetSessionInfo(campaignId, sessionId, "Reason", "dial_failed")
		SetAgentSessionInfo(campaignId, sessionId, "AgentReason", "dial_failed")
		SetSessionInfo(campaignId, sessionId, "DialerStatus", "dial_failed")

		SendCustomerIntegrationData(campaignId, sessionId)

		go UploadSessionInfo(campaignId, sessionId)
		fmt.Println(err.Error())

		return err.Error()
	}

	if resp != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		tmx := string(response[:])
		fmt.Println(tmx)
		resultInfo := strings.Split(tmx, " ")
		if len(resultInfo) > 0 {
			if resultInfo[0] == "-ERR" {

				if len(resultInfo) > 1 {
					reason := resultInfo[1]
					if reason == "" {
						SetAgentSessionInfo(campaignId, sessionId, "AgentReason", "not_specified")
						SetSessionInfo(campaignId, sessionId, "Reason", "not_specified")
					} else {
						reason = strings.TrimSuffix(reason, "\n")
						color.Red(reason)
						SetSessionInfo(campaignId, sessionId, "Reason", reason)
						SetAgentSessionInfo(campaignId, sessionId, "AgentReason", reason)
					}
				} else {
					SetSessionInfo(campaignId, sessionId, "Reason", "not_specified")
					SetAgentSessionInfo(campaignId, sessionId, "AgentReason", "not_specified")
				}
				SetSessionInfo(campaignId, sessionId, "DialerStatus", "dial_failed")
				SendCustomerIntegrationData(campaignId, sessionId)
				go UploadSessionInfo(campaignId, sessionId)
			} else {
				SetSessionInfo(campaignId, sessionId, "Reason", "dial_success")
				SetAgentSessionInfo(campaignId, sessionId, "AgentReason", "dial_success")
				SetSessionInfo(campaignId, sessionId, "DialerStatus", "dial_success")
			}
		}

		return tmx
	}
	return "SUCCESS"
}
