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
	LastUpdate   string
}
type Phones struct {
	Phone      string
	CampaignId string
}

type Result struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        []Campaign
}
type ResPhone struct {
	CustomMessage string
	IsSuccess     bool
	Result        []Phones
}

type ResultPCount struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        int
}

func main() {

	//go ptr("PAwan")
	//var camp string //= make([]string, cnt)
	var Camps Result
	//camp := GetCampaign()
	Camps = GetCampaign()

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

	for i, value := range Camps.Result {

		CampName := value.CampaignName
		fmt.Println("Campaign ", CampName)
		SetMaxMin(CampName, Camps.Result[i].Min, Camps.Result[i].Max)
		SetCampaignStatus(CampName, "1")
		GetNumbers(CampName, value.Max, 0)

	}
	for _, val := range Camps.Result {

		cnt := GetPhoneCount(val.CampaignName)
		if cnt > 0 {
			fmt.Println("Getting phones from redis ", val)
			go GetPhonesFromList(val.CampaignName)
		} else {
			return
		}
	}

	fmt.Scanln()

}

//func GetCampaign() []string {

func GetCampaign() Result {
	//campz := []string{}
	var s Result

	url := fmt.Sprintf("http://localhost:8083/DVP/API/1.0/DialerApi/GetCampaign")
	fmt.Println("URL ", url)
	r := restclient.RequestResponse{
		Url:    url,
		Method: "GET",
		Result: &s,
	}
	_, err := restclient.Do(&r)
	if err != nil {
		//panic(err)
		fmt.Println("Errz", err)
		fmt.Println(r)
	}

	//fmt.Println("element : ", campz[0])
	//camp[0] = "pp"
	//camp[1] = "cc"
	//fmt.Println(campz[0])
	//fmt.Println("Going to return Exception", s.Result[1].Min)
	return s
	//return campz

}
