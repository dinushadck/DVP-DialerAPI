package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"bufio"
	//"os"
)

func GetNumbersFromNumberBase(company, tenant, numberLimit int, campaignId, camScheduleId string) []string {
	numbers := make([]string, 0)
	pageKey := fmt.Sprintf("PhoneNumberPage:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	pageNumberToRequest := RedisIncr(pageKey)
	fmt.Println("pageNumber: ", pageNumberToRequest)

	// Get phone number from campign service and append
	authToken := fmt.Sprintf("%d#%d", tenant, company)
	fmt.Println("Start GetPhoneNumbers Auth: ", authToken, " CampaignId: ", campaignId, " camScheduleId: ", camScheduleId)
	client := &http.Client{}

	request := fmt.Sprintf("%s/CampaignManager/Campaign/%s/Numbers/%s/%d/%d", campaignService, campaignId, camScheduleId, numberLimit, pageNumberToRequest)
	fmt.Println("Start GetPhoneNumbers request: ", request)
	req, _ := http.NewRequest("GET", request, nil)
	req.Header.Add("Authorization", authToken)
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
			numbers = append(numbers, numRes.CampContactInfo.ContactId)
		}
	}
	return numbers
}

func LoadNumbers(company, tenant, numberLimit int, campaignId, camScheduleId string) {
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	numbers := GetNumbersFromNumberBase(company, tenant, numberLimit, campaignId, camScheduleId)
	if len(numbers) == 0 {
		numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
		RedisSet(numLoadingStatusKey, "done")
	} else {
		numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
		RedisSet(numLoadingStatusKey, "waiting")
		for _, number := range numbers {
			RedisListRpush(listId, number)
		}
	}
}

func LoadInitialNumberSet(company, tenant int, campaignId, camScheduleId string) {
	numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	LoadNumbers(company, tenant, 1000, campaignId, camScheduleId)
	RedisSet(numLoadingStatusKey, "waiting")
}

func GetNumberToDial(company, tenant int, campaignId, camScheduleId string) string {
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
	numberCount := RedisListLlen(listId)
	numLoadingStatus := RedisGet(numLoadingStatusKey)

	if numLoadingStatus == "waiting" {
		if numberCount < 500 {
			LoadNumbers(company, tenant, 500, campaignId, camScheduleId)
		}
	} else if numLoadingStatus == "done" && numberCount == 0 {
		pageKey := fmt.Sprintf("PhoneNumberPage:%d:%d:%s:%s", company, tenant, campaignId, camScheduleId)
		RedisRemove(numLoadingStatusKey)
		RedisRemove(pageKey)
	}
	return RedisListLpop(listId)
}

func GetNumberCount(company, tenant int, campaignId, camScheduleId string) int {
	fmt.Println("Start GetNumberCount")
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
