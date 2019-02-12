package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/fatih/color"
)

func AddRequestServer() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in AddRequestServer", r)
		}
	}()
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	dialerAPIUrl := fmt.Sprintf("http://%s", CreateHost(lbIpAddress, lbPort))
	path := fmt.Sprintf("DVP/DialerAPI/ArdsCallback")

	u, _ := url.Parse(dialerAPIUrl)
	u.Path += path

	fmt.Println(u.String())
	cbUrl := u.String()

	var reqServer = RequestServer{}
	reqServer.ServerID = dialerId
	reqServer.ServerType = "DIALER"
	reqServer.RequestType = "CALL"
	reqServer.CallbackUrl = cbUrl
	reqServer.CallbackOption = "GET"

	jsonData, _ := json.Marshal(reqServer)

	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/requestserver", CreateHost(ardsServiceHost, ardsServicePort))
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	fmt.Println("request:", serviceurl)
	fmt.Println(jsonData)

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

func AddRequest(company, tenant int, uuid, otherData string, attributes []string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in AddRequest", r)
		}
	}()

	var ardsReq = Request{}
	ardsReq.SessionId = uuid
	ardsReq.ServerType = "DIALER"
	ardsReq.RequestType = "CALL"
	ardsReq.Priority = "0"
	ardsReq.RequestServerId = dialerId
	ardsReq.Attributes = attributes
	ardsReq.OtherInfo = otherData

	jsonData, _ := json.Marshal(ardsReq)

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/request", CreateHost(ardsServiceHost, ardsServicePort))
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)

	whiteblue := color.New(color.FgWhite).Add(color.BgBlue)
	whiteblue.Println(fmt.Sprintf("ARDS REQUEST : %s | DATA : %s", serviceurl, string(jsonData)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	response, _ := ioutil.ReadAll(resp.Body)
	result := string(response)
	fmt.Println("response Body:", result)

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	return result, err
}

func RemoveRequest(company, tenant, sessionId string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in AddRequest", r)
		}
	}()
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%s:%s", tenant, company)
	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/request/%s/NONE", CreateHost(ardsServiceHost, ardsServicePort), sessionId)
	fmt.Println("Start RemoveRequest: ", request)
	req, _ := http.NewRequest("DELETE", request, nil)
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(response))
}

func RemoveRequestNoSession(company, tenant, sessionId string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in AddRequest", r)
		}
	}()
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%s:%s", tenant, company)
	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/request/%s/NoSession", CreateHost(ardsServiceHost, ardsServicePort), sessionId)
	fmt.Println("Start RemoveRequest: ", request)
	req, _ := http.NewRequest("DELETE", request, nil)
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(response))
}

func RejectRequest(company, tenant, sessionId string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RejectRequest", r)
		}
	}()
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%s:%s", tenant, company)
	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/request/%s/reject/AgentRejected", CreateHost(ardsServiceHost, ardsServicePort), sessionId)
	fmt.Println("Start RejectRequest: ", request)
	req, _ := http.NewRequest("DELETE", request, nil)
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(response))
}

func ClearResourceSlotWhenReject(company, tenant, reqCategory, resId, sessionId string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in ClearResourceSlotWhenReject", r)
		}
	}()

	var ardsResSlot = ArdsResSlot{}
	ardsResSlot.ReqCategory = reqCategory
	ardsResSlot.State = "Available"
	ardsResSlot.OtherInfo = "Reject"

	jsonData, _ := json.Marshal(ardsResSlot)

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%s:%s", tenant, company)
	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/resource/%s/concurrencyslot/session/%s", CreateHost(ardsServiceHost, ardsServicePort), resId, sessionId)
	req, err := http.NewRequest("PUT", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	fmt.Println("request:", serviceurl)
	fmt.Println(string(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	response, _ := ioutil.ReadAll(resp.Body)
	result := string(response)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Body:", result)
	defer resp.Body.Close()
}

func SetAgentStatusArds(company, tenant, reqCategory, resId, sessionId, state, ardsServerType, ardsRequestType string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in ClearResourceSlotWhenReject", r)
		}
	}()

	ardsResSlot := ArdsResource{}
	ardsResSlot.ServerType = ardsServerType
	ardsResSlot.RequestType = ardsRequestType
	ardsResSlot.State = state
	ardsResSlot.OtherInfo = ""
	ardsResSlot.Reason = ""
	ardsResSlot.Company = company
	ardsResSlot.Tenant = tenant
	ardsResSlot.BusinessUnit = "default"

	jsonData, _ := json.Marshal(ardsResSlot)

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%s:%s", tenant, company)
	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/ARDS/resource/%s/concurrencyslot/session/%s?direction=inbound", CreateHost(ardsServiceHost, ardsServicePort), resId, sessionId)
	req, err := http.NewRequest("PUT", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	fmt.Println("request:", serviceurl)
	fmt.Println(string(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	response, _ := ioutil.ReadAll(resp.Body)
	result := string(response)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Body:", result)
	defer resp.Body.Close()
}
