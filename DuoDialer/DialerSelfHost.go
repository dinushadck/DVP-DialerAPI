package main

import (
	"encoding/json"
	"fmt"
	"github.com/DuoSoftware/gorest"
	"net/url"
	"strconv"
	"strings"
)

type DVP struct {
	gorest.RestService     `root:"/DVP/" consumes:"application/json" produces:"application/json"`
	incrMaxChannelLimit    gorest.EndPoint `method:"POST" path:"/DialerAPI/IncrMaxChannelLimit/" postdata:"string"`
	decrMaxChannelLimit    gorest.EndPoint `method:"POST" path:"/DialerAPI/DecrMaxChannelLimit/" postdata:"string"`
	setMaxChannelLimit     gorest.EndPoint `method:"POST" path:"/DialerAPI/SetMaxChannelLimit/" postdata:"string"`
	previewCallBack        gorest.EndPoint `method:"POST" path:"/DialerAPI/PreviewCallBack/" postdata:"ReceiveData"`
	resumeCallback         gorest.EndPoint `method:"POST" path:"/DialerAPI/ResumeCallback/" postdata:"CampaignCallbackObj"`
	getTotalDialCount      gorest.EndPoint `method:"GET" path:"/DialerAPI/GetTotalDialCount/{CompanyId:int}/{TenantId:int}/{CampaignId:string}" output:"int"`
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

func (dvp DVP) ResumeCallback(callbackInfo CampaignCallbackObj) {
	company, tenant, _, _ := decodeJwtDialer(dvp, "dialer", "write")
	if company != 0 && tenant != 0 {
		log := fmt.Sprintf("Start ResumeCallback CampaignId:%d # ContactId:%s ", callbackInfo.CampaignId, callbackInfo.ContactId)
		fmt.Println(log)

		fmt.Println("Company: ", company)
		fmt.Println("Tenant: ", tenant)
		ResumeCampaignCallback(company, tenant, callbackInfo.CallBackCount, callbackInfo.CampaignId, callbackInfo.ContactId)
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
	fmt.Println("---------------Start ArdsCallback---------")
	jResult, _ := url.QueryUnescape(dvp.Context.Request().URL.RawQuery)
	log := fmt.Sprintf("Start ArdsCallback :%s ", jResult)
	fmt.Println(log)

	var ardsCallbackInfo ArdsCallbackInfo
	var reqOData RequestOtherData
	json.Unmarshal([]byte(jResult), &ardsCallbackInfo)
	json.Unmarshal([]byte(ardsCallbackInfo.OtherInfo), &reqOData)

	go RemoveRequest(ardsCallbackInfo.Company, ardsCallbackInfo.Tenant, ardsCallbackInfo.SessionID)

	switch reqOData.DialoutMec {
	case "PREVIEW":
		SendPreviewDataToAgent(ardsCallbackInfo, reqOData)
		break
	case "AGENT":
		log3 := fmt.Sprintf("Data:: ContactName: %s :: Domain: %s :: ContactType: %s ::ResourceId: %s  :: Company: %s :: Tenant: %s :: CampaignId: %s :: Class: %s :: Type: %s :: Category: %s :: SessionId: %s", ardsCallbackInfo.ResourceInfo.ContactName, ardsCallbackInfo.ResourceInfo.Domain, ardsCallbackInfo.ResourceInfo.ContactType, ardsCallbackInfo.ResourceInfo.ResourceId, ardsCallbackInfo.Company, ardsCallbackInfo.Tenant, reqOData.CampaignId, ardsCallbackInfo.ServerType, ardsCallbackInfo.RequestType, ardsCallbackInfo.SessionID)
		fmt.Println(log3)
		DialAgent(ardsCallbackInfo.ResourceInfo.ContactName, ardsCallbackInfo.ResourceInfo.Domain, ardsCallbackInfo.ResourceInfo.ContactType, ardsCallbackInfo.ResourceInfo.ResourceId, ardsCallbackInfo.Company, ardsCallbackInfo.Tenant, reqOData.CampaignId, ardsCallbackInfo.ServerType, ardsCallbackInfo.RequestType, ardsCallbackInfo.SessionID)
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
		log3 := fmt.Sprintf("Data:: ContactName: %s :: Domain: %s :: ContactType: %s ::ResourceId: %s  :: Company: %s :: Tenant: %s :: CampaignId: %s :: Class: %s :: Type: %s :: Category: %s :: SessionId: %s", refData.ResourceInfo.ContactName, refData.ResourceInfo.Domain, refData.ResourceInfo.ContactType, refData.ResourceInfo.ResourceId, refData.Company, refData.Tenant, reqOData.CampaignId, refData.ServerType, refData.RequestType, refData.SessionID)
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

		SendSmsDirect(company, tenant, Message, DnisNumber)
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
