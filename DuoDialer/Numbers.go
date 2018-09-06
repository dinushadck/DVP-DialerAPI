package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

func GetNumbersFromNumberBase(company, tenant, numberLimit int, campaignId, camScheduleId string) []string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetNumbersFromNumberBase", r)
		}
	}()
	numbers := make([]string, 0)
	pageKey := fmt.Sprintf("PhoneNumberPage:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)

	numberOffsetToRequest := "0"

	if numberLimit == 500 {
		numberOffsetToRequest = RedisGet(pageKey)
	}
	DialerLog(fmt.Sprintf("numberOffsetToRequest: %s", numberOffsetToRequest))

	// Get phone number from campign service and append
	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	DialerLog(fmt.Sprintf("Start GetPhoneNumbers Auth: %s  CampaignId: %s  camScheduleId: %s", internalAuthToken, campaignId, camScheduleId))
	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaign/%s/NumbersByOffset/%s/%d/%s", CreateHost(campaignServiceHost, campaignServicePort), campaignId, camScheduleId, numberLimit, numberOffsetToRequest)
	DialerLog(fmt.Sprintf("Start GetPhoneNumbers request: ", request))
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return numbers
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var phoneNumberResult PhoneNumberResult
	json.Unmarshal(response, &phoneNumberResult)
	if phoneNumberResult.IsSuccess == true {
		for _, numRes := range phoneNumberResult.Result {
			if numRes.ExtraData != "" {
				numberWithExtraD := fmt.Sprintf("%s:%s:%s", numRes.CampContactInfo.ContactId, "1", numRes.ExtraData)
				numbers = append(numbers, numberWithExtraD)
			} else {
				numberWithData := strings.Split(numRes.CampContactInfo.ContactId, ":")
				if len(numberWithData) > 1 {
					exData := strings.Join(numberWithData[1:], ":")
					numberAndExtraD := fmt.Sprintf("%s:%s:%s", numberWithData[0], "1", exData)
					numbers = append(numbers, numberAndExtraD)
				} else {
					numberWithoutExtraData := fmt.Sprintf("%s:%s:", numRes.CampContactInfo.ContactId, "1")
					numbers = append(numbers, numberWithoutExtraData)
				}
			}
		}

		RedisIncrBy(pageKey, len(phoneNumberResult.Result))
	}
	return numbers
}

func SetDncNumbersFromNumberBase(company, tenant int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SetDncNumbersFromNumberBase", r)
		}
	}()

	jwtToken := fmt.Sprintf("Bearer %s", accessToken)
	internalAuthToken := fmt.Sprintf("%d:%d", tenant, company)
	DialerLog(fmt.Sprintf("Start SetDncNumbersFromNumberBase Auth: %s", internalAuthToken))
	client := &http.Client{}

	request := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Dnc", CreateHost(campaignServiceHost, campaignServicePort))
	DialerLog(fmt.Sprintf("Start GetDncPhoneNumbers request: %s", request))
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", jwtToken)
	req.Header.Set("companyinfo", internalAuthToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)

	var dncNumberResult DncNumberResult
	json.Unmarshal(response, &dncNumberResult)
	if dncNumberResult.IsSuccess == true {
		dncNumberKey := fmt.Sprintf("DncNumber:%d:%d", tenant, company)
		RedisSetAdd(dncNumberKey, dncNumberResult.Result)
	}
}

func LoadNumbers(company, tenant, numberLimit int, campaignId, camScheduleId string) {
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	numbers := GetNumbersFromNumberBase(company, tenant, numberLimit, campaignId, camScheduleId)

	DialerLog(fmt.Sprintf("Number count = %d", len(numbers)))
	if len(numbers) == 0 {
		numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
		RedisSet(numLoadingStatusKey, "waiting")
	} else {
		numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
		dncNumberKey := fmt.Sprintf("DncNumber:%d:%d", tenant, company)
		RedisSet(numLoadingStatusKey, "waiting")
		for _, number := range numbers {
			if !RedisSetIsMember(dncNumberKey, number) {
				fmt.Println("Adding number to campaign: ", number)
				RedisListRpush(listId, number)
			}
		}
	}
}

func LoadInitialNumberSet(company, tenant int, campaignId, camScheduleId string) {
	numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	LoadNumbers(company, tenant, 1000, campaignId, camScheduleId)
	RedisSet(numLoadingStatusKey, "waiting")
}

func GetNumberToDial(company, tenant int, campaignId, camScheduleId string) (string, string, string) {
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	numberCount := RedisListLlen(listId)
	numLoadingStatus := RedisGet(numLoadingStatusKey)

	if numLoadingStatus == "waiting" {
		if numberCount < 500 {
			LoadNumbers(company, tenant, 500, campaignId, camScheduleId)
		}
	}

	color.Red(fmt.Sprintf("======= NUMBER LOADING STATUS : %s", numLoadingStatus))
	numberWithTryCount := RedisListLpop(listId)
	numberInfos := strings.Split(numberWithTryCount, ":")
	if len(numberInfos) > 3 {
		return numberInfos[0], numberInfos[1], strings.Join(numberInfos[2:], ":")
	} else if len(numberInfos) == 3 {
		return numberInfos[0], numberInfos[1], numberInfos[2]
	} else if len(numberInfos) == 2 {
		return numberInfos[0], numberInfos[1], ""
	} else {
		return "", "", ""
	}
}

func GetNumberCount(company, tenant int, campaignId, camScheduleId string) int {
	DialerLog("Start GetNumberCount")
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	return RedisListLlen(listId)
}

func RemoveNumbers(company, tenant int, campaignId string) {
	searchKey := fmt.Sprintf("CampaignNumbers:%d:%d:%s:*", company, tenant, campaignId)
	relatedNumberList := RedisSearchKeys(searchKey)
	for _, key := range relatedNumberList {
		RedisRemove(key)
	}
}

func RemoveNumberStatusKey(company, tenant int, campaignId string) {

	searchKeyPhoneNumberPage := fmt.Sprintf("PhoneNumberPage:%d:%d:%s:*", company, tenant, campaignId)
	searchKeyPhoneNumberLoading := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:*", company, tenant, campaignId)
	relatedNumberStatusKey := RedisSearchKeys(searchKeyPhoneNumberPage)
	relatedNumberStatusKeyNumberLoading := RedisSearchKeys(searchKeyPhoneNumberLoading)

	for _, key := range relatedNumberStatusKey {
		RedisRemove(key)
	}

	for _, key := range relatedNumberStatusKeyNumberLoading {
		RedisRemove(key)
	}
}

func AddNumberToFront(company, tenant int, campaignId, camScheduleId, number string) bool {
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	return RedisListLpush(listId, number)
}
