package main

import (
	"code.google.com/p/gorest"
	"fmt"
)

type DialerSelfHost struct {
	gorest.RestService  `root:"/DialerSelfHost/" consumes:"application/json" produces:"application/json"`
	incrMaxChannelLimit gorest.EndPoint `method:"POST" path:"/Campaign/IncrMaxChannelLimit/" postdata:"string"`
	decrMaxChannelLimit gorest.EndPoint `method:"POST" path:"/Campaign/DecrMaxChannelLimit/" postdata:"string"`
}

func (dialerSelfHost DialerSelfHost) IncrMaxChannelLimit(campaignId string) {
	fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
	go IncrCampChannelMaxLimit(campaignId)
	return
}

func (dialerSelfHost DialerSelfHost) DecrMaxChannelLimit(campaignId string) {
	fmt.Println("Start IncrMaxChannelLimit ServerId: ", campaignId)
	go DecrCampChannelMaxLimit(campaignId)
	return
}
