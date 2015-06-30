// BlastDialer project main.go
package main

import (
	"fmt"
	"time"
)

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	InitiateDuoDialer()
	for {
		onGoingCampaignCount := GetOnGoingCampaignCount()
		if onGoingCampaignCount < campaignLimit {
			campaign := RequestCampaign()
			if campaign != Campaign{}{
				//Start New Campaign
			}
		}
		time.Sleep(campaignRequestFrequency * time.Second)
	}
}
