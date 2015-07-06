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
			defCampaign := Campaign{}
			if campaign != defCampaign {
				//Add Campaign
			}
		}
		time.Sleep(campaignRequestFrequency * time.Second)
	}
}
