package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetUuid() string {
	resp, _ := http.Get(uuidService)
	defer resp.Body.Close()
	response, _ := ioutil.ReadAll(resp.Body)
	tmx := string(response[:])
	fmt.Println(tmx)
	return tmx
}

func GetTrunkCode(authToken, ani, dnis string) (trunkCode, rAni, rDnis string) {
	client := &http.Client{}

	request := fmt.Sprintf("%s/ANI/%s/DNIS/%s", callRuleService, ani, dnis)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var apiResult CallRuleApiResult
	json.Unmarshal(response, &apiResult)

	fmt.Println("callRule: ", apiResult.Result.GatewayCode, "ANI: ", apiResult.Result.ANI, "DNIS: ", apiResult.Result.DNIS)
	return apiResult.Result.GatewayCode, apiResult.Result.ANI, apiResult.Result.DNIS
}

func DialNumber(uuid, fromNumber, trunkCode, phoneNumber, extention string) {
	request := fmt.Sprintf("http://%s", callServer)
	path := fmt.Sprintf("api/originate?")
	param := fmt.Sprintf(" {return_ring_ready=true,origination_uuid=%s,origination_caller_id_number=%s}sofia/gateway/%s/%s %s", uuid, fromNumber, trunkCode, phoneNumber, extention)

	u, _ := url.Parse(request)
	u.Path += path
	u.Path += param

	fmt.Println(u.String())

	resp, _ := http.Get(u.String())
	defer resp.Body.Close()

	if resp != nil {
		response, _ := ioutil.ReadAll(resp.Body)
		tmx := string(response[:])
		fmt.Println(tmx)
	}
}
