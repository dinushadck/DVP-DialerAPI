package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
)


func AddCampaignDataRealtime(campaignData Campaign) {
	var campInfoRealTime map[string]string

	mapstructure.Decode(campaignData, &campInfoRealTime)

	color.Cyan("==============START=============");
	color.Cyan(fmt.Sprintf("%v", campInfoRealTime));
	color.Cyan("==============DONE=============");
	
}