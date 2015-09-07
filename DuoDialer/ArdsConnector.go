package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func AddRequestServer() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in AddRequestServer", r)
		}
	}()
	dialerAPIUrl := fmt.Sprintf("http://%s:%s", hostIpAddress, port)
	path := fmt.Sprintf("DVP/DialerAPI/ArdsCallback")

	u, _ := url.Parse(dialerAPIUrl)
	u.Path += path

	fmt.Println(u.String())
	cbUrl := u.String()

	var reqServer = RequestServer{}
	reqServer.ServerID = dialerId
	reqServer.Class = "CALLSERVER"
	reqServer.Type = "ARDS"
	reqServer.Category = "CALL"
	reqServer.CallbackUrl = cbUrl

	jsonData, _ := json.Marshal(reqServer)

	serviceurl := fmt.Sprintf("%s/requestserver", ardsService)
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
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

func AddRequest(company, tenant int, uuid string, attributes []string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in AddRequest", r)
		}
	}()

	var ardsReq = Request{}
	ardsReq.SessionId = uuid
	ardsReq.Class = "CALLSERVER"
	ardsReq.Type = "ARDS"
	ardsReq.Category = "CALL"
	ardsReq.Priority = "L"
	ardsReq.RequestServerId = dialerId
	ardsReq.Attributes = attributes
	ardsReq.OtherInfo = ""

	jsonData, _ := json.Marshal(ardsReq)

	authToken := fmt.Sprintf("%d#%d", tenant, company)
	serviceurl := fmt.Sprintf("%s/request", ardsService)
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authToken)
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

func RemoveRequest(company, tenant int, sessionId string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in AddRequest", r)
		}
	}()
	authToken := fmt.Sprintf("%d#%d", tenant, company)
	client := &http.Client{}

	request := fmt.Sprintf("%s/request/%s", ardsService, sessionId)
	fmt.Println("Start RemoveRequest: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(response))
}
