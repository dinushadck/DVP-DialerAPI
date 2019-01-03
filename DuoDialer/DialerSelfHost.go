package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/DuoSoftware/gorest"
	"github.com/fatih/color"
)

type DVP struct {
	gorest.RestService     `root:"/DVP/" consumes:"application/json" produces:"application/json"`
	incrMaxChannelLimit    gorest.EndPoint `method:"POST" path:"/DialerAPI/IncrMaxChannelLimit/" postdata:"string"`
	decrMaxChannelLimit    gorest.EndPoint `method:"POST" path:"/DialerAPI/DecrMaxChannelLimit/" postdata:"string"`
	setMaxChannelLimit     gorest.EndPoint `method:"POST" path:"/DialerAPI/SetMaxChannelLimit/" postdata:"string"`
	previewCallBack        gorest.EndPoint `method:"POST" path:"/DialerAPI/PreviewCallBack/" postdata:"ReceiveData"`
	resumeCallback         gorest.EndPoint `method:"POST" path:"/DialerAPI/ResumeCallback/" postdata:"CallbackInfo"`
	getTotalDialCount      gorest.EndPoint `method:"GET" path:"/DialerAPI/GetTotalDialCount/{CompanyId:int}/{TenantId:int}/{CampaignId:string}" output:"int"`
	dialCall               gorest.EndPoint `method:"GET" path:"/DialerAPI/DialCall/{CampaignId:string}/{DialNumber:string}/{Agent:string}/{Domain:string}" output:"string"`
	getTotalConnectedCount gorest.EndPoint `method:"GET" path:"/DialerAPI/GetTotalConnectedCount/{CompanyId:int}/{TenantId:int}/{CampaignId:string}" output:"int"`
	dial                   gorest.EndPoint `method:"GET" path:"/DialerAPI/Dial/{AniNumber:string}/{DnisNumber:string}/{Extention:string}/{CallserverId:string}" output:"bool"`
	dialCampaign           gorest.EndPoint `method:"GET" path:"/DialerAPI/DialCampaign/{CampaignId:int}/{ScheduleId:int}/{ContactNumber:string}" output:"bool"`
	ardsCallback           gorest.EndPoint `method:"GET" path:"/DialerAPI/ArdsCallback/" output:"string"`
	sendSms                gorest.EndPoint `method:"GET" path:"/DialerAPI/SendSms/{DnisNumber:string}/{Message:string}" output:"bool"`
	clickToCall            gorest.EndPoint `method:"GET" path:"/DialerAPI/ClickToCall/{DnisNumber:string}/{Extention:string}" output:"bool"`
}

func (dvp DVP) IncrMaxChannelLimit(campaignId string) {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	if company != 0 && tenant != 0 {
		fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
		go IncrCampChannelMaxLimit(campaignId)
	} else {
		dvp.RB().SetResponseCode(403)
	}
}

func (dvp DVP) DecrMaxChannelLimit(campaignId string) {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	if company != 0 && tenant != 0 {
		fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
		go DecrCampChannelMaxLimit(campaignId)
	} else {
		dvp.RB().SetResponseCode(403)
	}
}

func (dvp DVP) SetMaxChannelLimit(campaignId string) {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	if company != 0 && tenant != 0 {
		fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
		go SetCampChannelMaxLimit(campaignId)
	} else {
		dvp.RB().SetResponseCode(403)
	}
}

func (dvp DVP) GetTotalDialCount(companyId, tenantId int, campaignId string) int {
	fmt.Println("Start GetTotalDialCount")
	company, tenant, _, msg := decodeJwtDialer(dvp, "dialer", "read")
	fmt.Println(company, tenant, msg)
	if company != 0 && tenant != 0 {
		fmt.Println("Start GetTotalDialCount CampaignId: ", campaignId)
		count := 0

		count = GetCampaignDialCount(company, tenant, campaignId)
		return count
	} else {
		dvp.RB().SetResponseCode(403)
		return 0
	}
}

func (dvp DVP) DialCall(campaignId string, dialNumber string, agent string, domain string) string {

	color.Green(fmt.Sprintf("Dial Number : %s, Campaign ID : %s, Agent: %s, Domain: %s", dialNumber, campaignId, agent, domain))

	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")

	campId, _ := strconv.Atoi(campaignId)

	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)

	fmt.Println(internalAuthToken)

	campaigninfo, result := GetCampaign(company, tenant, campId)

	if result {
		campStatus := GetCampaignStatus(campaignId, company, tenant)

		fmt.Println(campStatus)

		if (campStatus == "Running" || campStatus == "PauseByDialer") && campaigninfo.CampaignChannel == "CALL" {

			resourceServerInfos := GetResourceServerInfo(company, tenant, "*", campaigninfo.CampaignChannel)

			fmt.Println(resourceServerInfos)

			trunkCode, ani, dnis, xGateway := GetTrunkCode(internalAuthToken, campaigninfo.CampConfigurations.Caller, dialNumber)
			uuid := GetUuid(resourceServerInfos.Url)

			fmt.Println("UUID : " + uuid)

			scheduleId := fmt.Sprintf("%d", campaigninfo.CampScheduleInfo[0].ScheduleId)
			InitiateSessionInfo(company, tenant, 240, "Campaign", "Dialer", "API", "1", campaignId, scheduleId, campaigninfo.CampaignName, uuid, dnis, "api called", "dial_start", time.Now().UTC().Format(layout4), resourceServerInfos.ResourceServerId, &campaigninfo.CampConfigurations.IntegrationData, nil, "")

			SetSessionInfo(campaignId, uuid, "FromNumber", ani)
			SetSessionInfo(campaignId, uuid, "TrunkCode", trunkCode)
			SetSessionInfo(campaignId, uuid, "Extention", campaigninfo.Extensions)
			SetSessionInfo(campaignId, uuid, "XGateway", xGateway)

			customCompanyStr := fmt.Sprintf("%d_%d", company, tenant)

			var param string
			var furl string
			var data string
			var dial bool
			if agent != "" {
				dial = true
				param = fmt.Sprintf(" {leg_timeout=10,sip_h_DVP-DESTINATION-TYPE=PRIVATE_USER,DVP_CALL_DIRECTION=outbound,nolocal:DVP_CUSTOM_PUBID=%s,CustomCompanyStr=%s,CampaignId=%s,CampaignName=%s,tenantid=%d,companyid=%d,ards_client_uuid=%s,origination_uuid=%s,ards_servertype=DIALER,ards_requesttype=CALL,DVP_ACTION_CAT=DIALER,DVP_OPERATION_CAT=AGENT,return_ring_ready=true,ignore_early_media=true,origination_caller_id_number=%s}", subChannelName, customCompanyStr, campaignId, campaigninfo.CampaignName, tenant, company, uuid, uuid, dnis)
				furl = fmt.Sprintf("user/%s@%s", agent, domain)
			} else {
				dial = false
				fmt.Println("Invalid Operation")
			}

			data = fmt.Sprintf(" %s xml dialer", dialNumber)

			if dial == true {
				SetSessionInfo(campaignId, uuid, "Reason", "Dial Number")

				redwhite := color.New(color.FgRed).Add(color.BgWhite)
				redwhite.Println(fmt.Sprintf("DIALING OUT CALL - API CALL: %s | NUMBER : %s", campaigninfo.CampaignName, dialNumber))

				resp, err := Dial(resourceServerInfos.Url, param, furl, data)
				r := HandleDialResponse(resp, err, resourceServerInfos, campaignId, uuid)

				r = strings.TrimSuffix(r, "\n")

				if err != nil {
					w, _ := json.Marshal(DialResult{IsSuccess: false, Message: r})
					return string(w)
				} else {
					w, _ := json.Marshal(DialResult{IsSuccess: true, Message: r})
					return string(w)
				}
			} else {
				SetSessionInfo(campaignId, uuid, "Reason", "Invalied Operation")
				w, _ := json.Marshal(DialResult{IsSuccess: false, Message: "AGENT NOT FOUND"})
				return string(w)

			}

		} else {
			color.Red("CAMPAIGN IS NOT IN RUNNING STATE")
			w, _ := json.Marshal(DialResult{IsSuccess: false, Message: "CAMPAIGN IS NOT IN RUNNING STATE"})
			return string(w)
		}
	} else {
		color.Red("CAMPAIGN NOT FOUND API CALL")
		w, _ := json.Marshal(DialResult{IsSuccess: false, Message: "CAMPAIGN NOT FOUND"})
		return string(w)
	}
}

func (dvp DVP) GetTotalConnectedCount(companyId, tenantId int, campaignId string) int {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "read")
	if company != 0 && tenant != 0 {
		fmt.Println("Start GetTotalConnectedCount CampaignId: ", campaignId)
		count := 0

		count = GetCampaignConnectedCount(company, tenant, campaignId)
		return count
	} else {
		dvp.RB().SetResponseCode(403)
		return 0
	}
}

func (dvp DVP) ResumeCallback(callbackInfo CallbackInfo) {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	if company != 0 && tenant != 0 {

		fmt.Println("Company: ", company)
		fmt.Println("Tenant: ", tenant)

		callbackData, _ := json.Marshal(callbackInfo)
		fmt.Println("Start ResumeCallback :: ", string(callbackData))

		if strings.ToLower(callbackInfo["CallbackType"].(string)) == "callback" && strings.ToLower(callbackInfo["CallbackCategory"].(string)) == "internal" {

			callbackCount, _ := strconv.Atoi(callbackInfo["CallBackCount"].(string))
			campaignId, _ := strconv.Atoi(callbackInfo["CampaignId"].(string))

			ResumeCampaignCallback(company, tenant, callbackCount, campaignId, callbackInfo["ContactId"].(string))

		} else if strings.ToLower(callbackInfo["CallbackType"].(string)) == "schedulecallback" && strings.ToLower(callbackInfo["CallbackCategory"].(string)) == "agent" {

			attributeInterface := callbackInfo["AttributeInfo"].([]interface{})

			callbackDataAttributeInfo, _ := json.Marshal(callbackInfo["AttributeInfo"].([]interface{}))
			fmt.Println("callbackDataAttributeInfo :: ", string(callbackDataAttributeInfo))

			attributeDetails := make([]string, len(attributeInterface))
			fmt.Println("callbackDataAttributeInfo length :: ", len(attributeInterface))

			for i := range attributeInterface {
				attributeDetails[i] = attributeInterface[i].(string)
			}
			//for _, i := range attributeInterface {
			//	attributeInfoI, _ := json.Marshal(i)
			//	fmt.Println("i ::", attributeInfoI)
			//	attributeDetails = append(attributeDetails, i.(string))
			//}

			fmt.Println("attributeDetails ::", attributeDetails)

			SchedulePreviewCallback(company, tenant, callbackInfo["SessionId"].(string), callbackInfo["PhoneNumber"].(string), callbackInfo["PriviewData"].(string), "", attributeDetails)

		} else if strings.ToLower(callbackInfo["CallbackType"].(string)) == "schedulecallback" && strings.ToLower(callbackInfo["CallbackCategory"].(string)) == "ivr" {

			ScheduleIvrCallback(company, tenant, callbackInfo["SessionId"].(string), callbackInfo["PhoneNumber"].(string), callbackInfo["Extention"].(string))

		} else {

			fmt.Println("Invalied Callback")

		}

	} else {
		dvp.RB().SetResponseCode(403)
	}
}

func (dvp DVP) DialCampaign(campaignId, ScheduleId int, contactNumber string) bool {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	if company != 0 && tenant != 0 {
		log := fmt.Sprintf("Start Direct DialCampaign CampaignId:%d # DNIS:%s ", campaignId, contactNumber)
		fmt.Println(log)
		authHeaderStr := dvp.Context.Request().Header.Get("Authorization")
		fmt.Println(authHeaderStr)

		authHeaderInfo := strings.Split(authHeaderStr, "#")
		if len(authHeaderInfo) == 2 {
			tenant, _ := strconv.Atoi(authHeaderInfo[0])
			company, _ := strconv.Atoi(authHeaderInfo[1])
			fmt.Println("Company: ", company)
			fmt.Println("Tenant: ", tenant)
			return DirectDialCampaign(company, tenant, campaignId, ScheduleId, contactNumber)
		}
		return false
	} else {
		dvp.RB().SetResponseCode(403)
		return false
	}
}

func (dvp DVP) Dial(AniNumber, DnisNumber, Extention, CallserverId string) bool {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	if company != 0 && tenant != 0 {
		log := fmt.Sprintf("Start Direct Dial ANI:%s # DNIS:%s ", AniNumber, DnisNumber)
		fmt.Println(log)
		authHeaderStr := dvp.Context.Request().Header.Get("Authorization")
		fmt.Println(authHeaderStr)

		authHeaderInfo := strings.Split(authHeaderStr, "#")
		if len(authHeaderInfo) == 2 {
			tenant, _ := strconv.Atoi(authHeaderInfo[0])
			company, _ := strconv.Atoi(authHeaderInfo[1])
			fmt.Println("Company: ", company)
			fmt.Println("Tenant: ", tenant)
			return DirectDial(company, tenant, AniNumber, DnisNumber, Extention, CallserverId)
		}
		return false
	} else {
		dvp.RB().SetResponseCode(403)
		return false
	}
}

func (dvp DVP) ArdsCallback() string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in ArdsCallback", r)
		}
	}()
	//company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	//if company != 0 && tenant != 0 {
	redyellow := color.New(color.FgRed).Add(color.BgYellow)
	redyellow.Println("=========== ARDS CALL BACK RECEIVED ==========")
	jResult, _ := url.QueryUnescape(dvp.Context.Request().URL.RawQuery)
	log := fmt.Sprintf("Start ArdsCallback :%s ", jResult)
	redyellow.Println(log)

	var ardsCallbackInfo ArdsCallbackInfo
	var reqOData RequestOtherData
	json.Unmarshal([]byte(jResult), &ardsCallbackInfo)
	json.Unmarshal([]byte(ardsCallbackInfo.OtherInfo), &reqOData)

	//Send Agent Reserved Notification If Integration Data Exist
	SetSessionInfo(reqOData.CampaignId, ardsCallbackInfo.SessionID, "Agent", ardsCallbackInfo.ResourceInfo.ResourceName)
	SetSessionInfo(reqOData.CampaignId, ardsCallbackInfo.SessionID, "ResourceId", ardsCallbackInfo.ResourceInfo.ResourceId)

	hKey := fmt.Sprintf("sessionInfo:%s:%s", reqOData.CampaignId, ardsCallbackInfo.SessionID)
	sessionInfo := RedisHashGetAll(hKey)

	if sessionInfo != nil && sessionInfo["IntegrationData"] != "" {
		go ManageIntegrationData(sessionInfo, "AGENT")
	} else {
		color.Magenta("NO INTEGRATION DATA")
	}

	go RemoveRequest(ardsCallbackInfo.Company, ardsCallbackInfo.Tenant, ardsCallbackInfo.SessionID)

	SetSessionInfo(reqOData.CampaignId, ardsCallbackInfo.SessionID, "ArdsQueueName", ardsCallbackInfo.Skills)

	switch reqOData.DialoutMec {
	case "PREVIEW":
		SendPreviewDataToAgent(ardsCallbackInfo, reqOData)
		break
	case "ScheduledPreviewCallback":
		SendPreviewDataToAgent(ardsCallbackInfo, reqOData)
		break
	case "AGENT":
		log3 := fmt.Sprintf("Data:: ContactName: %s :: Domain: %s :: ContactType: %s ::ResourceId: %s  :: Company: %s :: Tenant: %s :: CampaignId: %s :: Class: %s :: Type: %s :: Category: %s :: SessionId: %s", ardsCallbackInfo.ResourceInfo.Extention, ardsCallbackInfo.ResourceInfo.Domain, ardsCallbackInfo.ResourceInfo.ContactType, ardsCallbackInfo.ResourceInfo.ResourceId, ardsCallbackInfo.Company, ardsCallbackInfo.Tenant, reqOData.CampaignId, ardsCallbackInfo.ServerType, ardsCallbackInfo.RequestType, ardsCallbackInfo.SessionID)
		fmt.Println(log3)
		DialAgent(ardsCallbackInfo.ResourceInfo.Extention, ardsCallbackInfo.ResourceInfo.Domain, ardsCallbackInfo.ResourceInfo.ContactType, ardsCallbackInfo.ResourceInfo.ResourceId, ardsCallbackInfo.Company, ardsCallbackInfo.Tenant, reqOData.CampaignId, ardsCallbackInfo.ServerType, ardsCallbackInfo.RequestType, ardsCallbackInfo.SessionID)
		break
	}

	//} else {
	//	dvp.RB().SetResponseCode(403)
	//}

	return ""
}

func (dvp DVP) PreviewCallBack(rdata ReceiveData) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PreviewCallBack", r)
		}
	}()
	//company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	//if company != 0 && tenant != 0 {
	log := fmt.Sprintf("Start PreviewCallBack Ref:%s ", rdata.Ref)
	log1 := fmt.Sprintf("Start PreviewCallBack TKey:%s ", rdata.Reply.Tkey)
	log2 := fmt.Sprintf("Start PreviewCallBack Message:%s ", rdata.Reply.Message)
	fmt.Println(log)
	fmt.Println(log1)
	fmt.Println(log2)

	var refData ArdsCallbackInfo
	json.Unmarshal([]byte(rdata.Ref), &refData)

	var reqOData RequestOtherData
	json.Unmarshal([]byte(refData.OtherInfo), &reqOData)

	if rdata.Reply.Message == "ACCEPTED" {
		fmt.Println("Start Dial Priview Number")
		log3 := fmt.Sprintf("Data:: ContactName: %s :: Domain: %s :: ContactType: %s ::ResourceId: %s  :: Company: %s :: Tenant: %s :: CampaignId: %s :: ServerType: %s :: RequestType: %s :: SessionId: %s", refData.ResourceInfo.ContactName, refData.ResourceInfo.Domain, refData.ResourceInfo.ContactType, refData.ResourceInfo.ResourceId, refData.Company, refData.Tenant, reqOData.CampaignId, refData.ServerType, refData.RequestType, refData.SessionID)
		fmt.Println(log3)
		DialAgent(refData.ResourceInfo.ContactName, refData.ResourceInfo.Domain, refData.ResourceInfo.ContactType, refData.ResourceInfo.ResourceId, refData.Company, refData.Tenant, reqOData.CampaignId, refData.ServerType, refData.RequestType, refData.SessionID)
	} else {
		fmt.Println("Start Reject Priview Number")
		AgentReject(refData.Company, refData.Tenant, reqOData.CampaignId, refData.SessionID, refData.RequestType, refData.ResourceInfo.ResourceId, "AgentRejected")
	}

	return
	//} else {
	//	dvp.RB().SetResponseCode(403)
	//	return
	//}
}

func (dvp DVP) SendSms(DnisNumber, Message string) bool {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	if company != 0 && tenant != 0 {
		log := fmt.Sprintf("Start Send SMS DNIS:%s # Message:%s ", DnisNumber, Message)
		fmt.Println(log)

		//SendSmsDirect(company, tenant, Message, DnisNumber)
		return true
	} else {
		dvp.RB().SetResponseCode(403)
		return false
	}
}

func (dvp DVP) ClickToCall(DnisNumber, Extention string) bool {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	if company != 0 && tenant != 0 {
		log := fmt.Sprintf("Start ClickToCall Dial DNIS:%s # Extension:%s ", DnisNumber, Extention)
		fmt.Println(log)

		fmt.Println("Company: ", company)
		fmt.Println("Tenant: ", tenant)
		return ClickToCall(company, tenant, DnisNumber, Extention, "1")
	} else {
		dvp.RB().SetResponseCode(403)
		return false
	}
}
