package main

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	//"github.com/jmcvetta/napping"
	"github.com/jmcvetta/restclient"
	//"strconv"
	//	"bufio"
	//	"io/ioutil"
	"time"
)

var uuidService = "http://localhost:8080/api/create_uuid"
var maxLimit = 100
var callServer = "localhost:8080"
var extention = "1234"
var trunkCode = "SLTOB"
var fromNumber = "0117491700"

func Dial() {

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

func GetNumbers(CampName string, Max int) {

	Nums := []string{}
	fmt.Println("Camp array created")
	url := fmt.Sprintf("http://localhost:8083/DVP/API/1.0/DialerApi/FillCampaignPhones/%s/%d", CampName, Max)
	fmt.Println("URL hit")
	r := restclient.RequestResponse{
		Url:    url,
		Method: "GET",
		Result: &Nums,
	}
	_, err := restclient.Do(&r)
	if err != nil {
		//panic(err)
		fmt.Println("Err", err)
	}
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}
	for _, val := range Nums {

		//fmt.Println("Num :", val)
		Phn := c.Cmd("LPUSH", CampName, val)
		fmt.Println(Phn)

	}
	//fmt.Println("element : ", campz[0])
	//camp[0] = "pp"
	//camp[1] = "cc"
	//fmt.Println(campz[0])
	//fmt.Println("Going to return %s", r.RawText)
	//return campz

}

func GetPhonesFromList(CampName string) {

	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}
	CmpMn := fmt.Sprintf("Min_%s", CampName)
	CmpMx := fmt.Sprintf("Max_%s", CampName)
	//	LPhns := c.Cmd("LPOP", CampName)
	LenPhns, _ := c.Cmd("LLEN", CampName).Int()

	MinPhns, _ := c.Cmd("GET", CmpMn).Int()
	MaxPhns, _ := c.Cmd("GET", CmpMx).Int()

	if LenPhns <= MinPhns {
		NewFill := MaxPhns - LenPhns
		go GetNumbers(CampName, NewFill)
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
