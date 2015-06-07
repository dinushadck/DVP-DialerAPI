// BlastDialer project main.go
package main

import (
	//"encoding/json"
	"fmt"
	"github.com/jmcvetta/restclient"
)

type Campaign struct {
	id           int
	CampaignName string
	Min          int
	Max          int
	StartTime    string
	EndTime      string
}

type Result struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        []Campaign
}

func main() {

	//var camp string //= make([]string, cnt)

	//camp := GetCampaign()
	GetCampaign()
	/*
		fmt.Println(camp)
		//camp[0] = "PP"

		for _, val := range camp {

			//fmt.Println("Num :", val)
			go GetNumbers(val, -1)
		}

		for _, CSVal := range camp {
			go GetPhonesFromList(CSVal)
		}
	*/

}

//func GetCampaign() []string {
func GetCampaign() {
	//campz := []string{}
	var s Result

	fmt.Println("Camp array created")
	url := fmt.Sprintf("http://localhost:8083/DVP/API/1.0/DialerApi/GetCampaign")
	fmt.Println("URL hit")
	r := restclient.RequestResponse{
		Url:    url,
		Method: "GET",
		Result: &s,
	}
	_, err := restclient.Do(&r)
	if err != nil {
		//panic(err)
		fmt.Println("Err", err)
	}

	//fmt.Println("element : ", campz[0])
	//camp[0] = "pp"
	//camp[1] = "cc"
	//fmt.Println(campz[0])
	fmt.Println("Going to return %s", s)
	//return campz

}
