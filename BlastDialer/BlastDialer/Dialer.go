package main

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	//"github.com/jmcvetta/napping"
	"github.com/jmcvetta/restclient"
	//"strconv"
	//	"bufio"
	//	"io/ioutil"
	//"strconv"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var uuidService = "http://localhost:8080/api/create_uuid"
var maxLimit = 100
var callServer = "localhost:8080"
var extention = "1234"
var trunkCode = "SLTOB"
var fromNumber = "0117491700"

func Dial(CampName string) {

	GetPhonesFromList(CampName)
}

/*
func GetCampaign() []string {

	campz := []string{}
	fmt.Println("Camp array created")
	url := fmt.Sprintf("http://localhost:8083/DVP/API/1.0/DialerApi/GetCampaign")
	fmt.Println("URL hit")
	r := restclient.RequestResponse{
		Url:    url,
		Method: "GET",
		Result: &campz,
	}
	_, err := restclient.Do(&r)
	if err != nil {
		//panic(err)
		fmt.Println("Err", err)
	}

	//fmt.Println("element : ", campz[0])
	//camp[0] = "pp"
	//camp[1] = "cc"
	fmt.Println(campz[0])
	fmt.Println("Going to return %s", r.RawText)
	return campz

}
*/
func GetCampaignCount() int {

	//camp := make([]string, 0)
	var count = 0
	url := fmt.Sprintf("http://127.0.0.1:8083/DVP/API/1.0/DialerApi/GetCampaignCount")

	r := restclient.RequestResponse{
		Url:    url,
		Method: "GET",
		Result: &count,
	}
	_, err := restclient.Do(&r)
	if err != nil {
		//panic(err)
		//fmt.Println(status)
	}
	return count
}

func GetNumbers(CampName string, Max int, st int) {

	//Nums := []string{}
	for {
		cnt := GetPhoneCount(CampName)
		if cnt > 0 {
			fmt.Println("HIT GET NUMS...................")
			var p ResPhone

			if Max == 0 {
				SetCampaignStatus(CampName, "0")
				return
			}

			url := fmt.Sprintf("http://localhost:8083/DVP/API/1.0/DialerApi/FillCampaignPhones/%s/%d", CampName, Max)
			fmt.Println("URL hit ", url)
			r := restclient.RequestResponse{
				Url:    url,
				Method: "GET",
				Result: &p,
			}
			_, err := restclient.Do(&r)
			if err != nil {

				fmt.Println("Err", err)

			}
			if p.IsSuccess != true {
				SetCampaignStatus(CampName, "0")
				fmt.Printf("Error returns from service ", p.CustomMessage)
				return
			}

			c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
			if err == nil {

			} else {
				fmt.Println("GetFrom List Error ", err.Error())
			}
			for _, val := range p.Result {

				fmt.Println("Campname :", CampName)
				fmt.Println("Number :", val.Phone)
				Phn := c.Cmd("LPUSH", CampName, val.Phone)
				fmt.Println(Phn)

			}

			if st != 0 {
				fmt.Println("Not initial")
				GetPhonesFromList(CampName)
			} else {
				return
				fmt.Println("Initial")
			}

		} else {
			return
		}

		time.Sleep(1000 * time.Millisecond)
	}

}

func GetPhonesFromList(CampName string) {

	//CampSt := GetCampaignStatus(CampName)
	//need while loop to check status of campaign

	fmt.Println("Hit list", CampName)
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}
	CmpMn := fmt.Sprintf("%s_Min", CampName)
	CmpMx := fmt.Sprintf("%s_Max", CampName)
	LPhns := c.Cmd("LPOP", CampName).String()
	uuid := GetUuid()
	go DialServer(LPhns, uuid)

	fmt.Println("Poped ", LPhns)
	LenPhns, _ := c.Cmd("LLEN", CampName).Int()

	MinPhns, _ := c.Cmd("GET", CmpMn).Int()
	MaxPhns, _ := c.Cmd("GET", CmpMx).Int()

	if LenPhns <= MinPhns {
		NewFill := MaxPhns - LenPhns
		fmt.Println("Max %d- Length %d of %s", MaxPhns, LenPhns, CampName)
		fmt.Println("NewFill", NewFill)
		if GetPhoneCount(CampName) > 0 {
			GetNumbers(CampName, NewFill, 1)
		} else {
			if LenPhns > 0 {
				GetPhonesFromList(CampName)
			} else {
				return
			}

		}

		//GetPhonesFromList(CampName)
	} else {
		GetPhonesFromList(CampName)
	}

}

/*
func GetUuid() string {
	var uuid string
	r := restclient.RequestResponse{
		Url:    uuidService,
		Method: "GET",
		Result: &uuid,
	}
	_, err := restclient.Do(&r)
	if err != nil {
		//panic(err)
		fmt.Println("Err", err)
	}

	response, _ := ioutil.ReadAll(r.RawText)
	tmx := string(response[:])
	fmt.Println(tmx)
	return tmx
}
*/
/*
func DialServer(phoneNumber string, uuid string, numberListKey string) {
	request := fmt.Sprintf("http://%s", callServer)
	path := fmt.Sprintf("api/originate?")
	param := fmt.Sprintf(" {return_ring_ready=true,origination_uuid=%s,origination_caller_id_number=%s}sofia/gateway/%s/%s %s", uuid, fromNumber, trunkCode, phoneNumber, extention)

	u, _ := url.Parse(request)
	u.Path += path
	u.Path += param

	fmt.Println(u.String())

	resp, _ := http.Get(u.String())
	defer resp.Body.Close()

	if resp != nil {

		response, _ := ioutil.ReadAll(resp.Body)
		tmx := string(response[:])
		fmt.Println(tmx)
		//go AddPhoneNumberToCouch(numberListKey, phoneNumber)
	}
}
*/
func SetMaxMin(Campaign string, Min int, Max int) {
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	MnName := fmt.Sprintf("%s_Min", Campaign)
	MxName := fmt.Sprintf("%s_Max", Campaign)

	fmt.Println("Min ", MnName, "Mx ", MxName)
	SetMin := c.Cmd("SET", MnName, Min)
	SetMax := c.Cmd("SET", MxName, Max)
	fmt.Println("Campaign Max ", SetMax)
	fmt.Println("Campaign Min ", SetMin)

}
func SetCampaignStatus(Campaign string, St string) {

	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	CampSt := fmt.Sprintf("%s_St", Campaign)
	s := c.Cmd("SET", CampSt, St)
	fmt.Println(s)
}
func GetCampaignStatus(Campaign string) string {

	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	CampSt := fmt.Sprintf("%s_St", Campaign)
	st := c.Cmd("GET", CampSt).String()
	return st
}
func GetPhoneCount(CampName string) int {
	var pCount ResultPCount
	url := fmt.Sprintf("http://localhost:8083/DVP/API/1.0/DialerApi/PhoneCount/%s", CampName)
	fmt.Println("URL hit ", url)
	r := restclient.RequestResponse{
		Url:    url,
		Method: "GET",
		Result: &pCount,
	}
	_, err := restclient.Do(&r)
	if err != nil {
		//panic(err)
		fmt.Println("Err", err)
		//fmt.Println("Raw ", r.RawText)
	}
	//fmt.Println("Count ", pCount.Result)
	fmt.Println("Raw ", r.RawText)
	//I, _ := strconv.Atoi(pCount.Result)
	return pCount.Result

}
func GetUuid() string {
	resp, _ := http.Get(uuidService)
	defer resp.Body.Close()
	response, _ := ioutil.ReadAll(resp.Body)
	tmx := string(response[:])
	fmt.Println(tmx)
	return tmx
}
func DialServer(phoneNumber string, uuid string) {
	request := fmt.Sprintf("http://%s", callServer)
	path := fmt.Sprintf("api/originate?")
	param := fmt.Sprintf(" {return_ring_ready=true,origination_uuid=%s,origination_caller_id_number=%s}sofia/gateway/%s/%s %s", uuid, fromNumber, trunkCode, phoneNumber, extention)

	u, _ := url.Parse(request)
	u.Path += path
	u.Path += param

	fmt.Println(u.String())

	resp, _ := http.Get(u.String())
	defer resp.Body.Close()

	if resp != nil {

		response, _ := ioutil.ReadAll(resp.Body)
		tmx := string(response[:])
		fmt.Println(tmx)
		//go AddPhoneNumberToCouch(numberListKey, phoneNumber)
	}
}
