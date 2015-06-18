// BlastDialer project main.go
package main

import (
	//"encoding/json"
	"fmt"
	"github.com/jmcvetta/restclient"
	//"strconv"
	//	"time"
)

type Campaign struct {
	id              int
	CampaignName    string
	Min             int
	Max             int
	StartTime       string
	EndTime         string
	LastUpdate      string
	ConcurrentLimit int
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
	//	var x = 0

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
		fmt.Println("Selected Campaign ", CampName)
		fmt.Println()
		SetMaxMin(CampName, Camps.Result[i].Min, Camps.Result[i].Max)
		SetCampaignStatus(CampName, "1")
		SetCampaignConLimit(CampName, Camps.Result[i].ConcurrentLimit)
		GetNumbers(CampName, value.Max)
		go Dial(CampName)
		//time.Sleep(1000 * time.Millisecond)
	}

	/*
		for _, val := range Camps.Result {

			//cnt := GetPhoneCount(val.CampaignName)
			fmt.Println("New Campaign ", val.CampaignName)
			//if cnt > 0 {
			fmt.Println()
			fmt.Println("Getting phones from redis ", val.CampaignName)
			fmt.Println()
			x = x + 1
			fmt.Println("Range is ", x)
			go GetPhonesFromList(val.CampaignName)

			//} else {
			//return
			//}

		}
	*/
	//var x = GetListPhoneSet("a")
	//fmt.Println(x)
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
	fmt.Println("Getting campagns Done")
	return s
	//return campz

}
