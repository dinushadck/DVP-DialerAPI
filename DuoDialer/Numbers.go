package main

import (
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	//"net/http"
	"bufio"
	"os"
)

func GetNumbersFromNumberBase(company, tenant, numberLimit int, campaignId, scheduleId string) []string {
	numbers := make([]string, 0)
	pageKey := fmt.Sprintf("PhoneNumberPage:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
	pageNumberToRequest := RedisIncr(pageKey)
	fmt.Println("pageNumber: ", pageNumberToRequest)

	if pageNumberToRequest == 1 {
		file, err := os.Open("D:\\Duo Projects\\Version 5.1\\Documents\\GolangProjects\\CampaignManager\\NumberList4.txt")
		if err != nil {
			//log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			RedisListLpush("CampaignNumbers:4:2:1:1", scanner.Text())
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			//log.Fatal(err)
		}
	}

	/*// Get phone number from campign service and append
	authToken := fmt.Sprintf("%d#%d", company, tenant)
	fmt.Println("Start GetPhoneNumbers Auth: ", authToken, " CampaignId: ", campaignId, " SchedulrId: ", scheduleId)
	client := &http.Client{}

	request := fmt.Sprintf("%s/CampaignNumberUpload/%s/%s/%d/%d", campaignService, campaignId, scheduleId,numberLimit,pageNumberToRequest)
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
	}*/
	return numbers
}

func LoadNumbers(company, tenant, numberLimit int, campaignId, scheduleId string) {
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
	numbers := GetNumbersFromNumberBase(company, tenant, numberLimit, campaignId, scheduleId)
	if len(numbers) == 0 {
		numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
		RedisSet(numLoadingStatusKey, "done")
	} else {
		numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
		RedisSet(numLoadingStatusKey, "waiting")
		for _, number := range numbers {
			RedisListRpush(listId, number)
		}
	}
}

func LoadInitialNumberSet(company, tenant int, campaignId, scheduleId string) {
	numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
	LoadNumbers(company, tenant, 1000, campaignId, scheduleId)
	RedisSet(numLoadingStatusKey, "waiting")
}

func GetNumberToDial(company, tenant int, campaignId, scheduleId string) string {
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
	numLoadingStatusKey := fmt.Sprintf("PhoneNumberLoading:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
	numberCount := RedisListLlen(listId)
	numLoadingStatus := RedisGet(numLoadingStatusKey)

	if numLoadingStatus == "waiting" {
		if numberCount < 500 {
			LoadNumbers(company, tenant, 500, campaignId, scheduleId)
		}
	} else if numLoadingStatus == "done" && numberCount == 0 {
		pageKey := fmt.Sprintf("PhoneNumberPage:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
		RedisRemove(numLoadingStatusKey)
		RedisRemove(pageKey)
	}
	return RedisListLpop(listId)
}

func GetNumberCount(company, tenant int, campaignId, scheduleId string) int {
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
	return RedisListLlen(listId)
}

func RemoveNumbers(company, tenant int, campaignId string) {
	searchKey := fmt.Sprintf("CampaignNumbers:%d:%d:%s:*", company, tenant, campaignId)
	relatedNumberList := RedisSearchKeys(searchKey)
	for _, key := range relatedNumberList {
		RedisRemove(key)
	}
}
