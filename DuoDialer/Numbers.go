package main

import (
	"fmt"
)

func GetNumbersFromNumberBase(company, tenant int, campaignId, scheduleId string) []string {
	numbers := make([]string, 0)
	// Get phone number from campign service and append
	return numbers
}
func LoadNumbers(company, tenant int, campaignId, scheduleId string) {
	listId := fmt.Sprintf("CampaignNumbers:%d:%d:%s:%s", company, tenant, campaignId, scheduleId)
	numbers := GetNumbersFromNumberBase(company, tenant, campaignId, scheduleId)
	for _, number := range numbers {
		RedisListRpush(listId, number)
	}
}
