package main

import (
	"code.google.com/p/gorest"
	"fmt"
	"strconv"
	"strings"
)

type CallbackServerSelfHost struct {
	gorest.RestService `root:"/CallbackServerSelfHost/" consumes:"application/json" produces:"application/json"`
	addCallback        gorest.EndPoint `method:"POST" path:"/Callback/AddCallback/" postdata:"CampaignCallback"`
}

func (callbackServerSelfHost CallbackServerSelfHost) AddCallback(callbackInfo CampaignCallback) {
	authHeaderStr := callbackServerSelfHost.Context.Request().Header.Get("Authorization")
	fmt.Println("Start AddCallback: ", callbackInfo.CallbackUrl, "#", callbackInfo.DialoutTime.String())
	fmt.Println(authHeaderStr)

	authHeaderInfo := strings.Split(authHeaderStr, "#")
	if len(authHeaderInfo) == 2 {
		tenant, _ := strconv.Atoi(authHeaderInfo[0])
		company, _ := strconv.Atoi(authHeaderInfo[1])

		go AddCallbackInfoToRedis(company, tenant, callbackInfo)
		if callbackInfo.Class == "DIALER" && callbackInfo.Type == "CALLBACK" && callbackInfo.Category == "INTERNAL" {
			go UploadCampaignMgrCallbackInfo(company, tenant, callbackInfo.CallbackObj)
		}
	}
	return
}
