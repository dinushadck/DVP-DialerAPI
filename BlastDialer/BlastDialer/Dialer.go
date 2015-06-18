package main

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	//"github.com/jmcvetta/napping"
	"github.com/jmcvetta/restclient"
	//"strconv"
	//	"bufio"
	//	"io/ioutil"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var uuidService = "http://localhost:8080/api/create_uuid"
var maxLimit = 100
var callServer = "localhost:8080"
var extention = "1234"
var trunkCode = "SLTOB"
var fromNumber = "0117491700"

func Dial(CampName string) {

	//GetPhonesFromList(CampName)
	for GetCampaignStatus(CampName) == "1" {
		fmt.Println()
		fmt.Println("Campaign ", CampName, " Status 1")
		fmt.Println()
		if GetArrayLength(CampName) <= GetArrayMin(CampName) {
			if GetPhoneCount(CampName) == 0 {
				fmt.Println("Db is empty of ", CampName)

			} else {
				fmt.Println("ArrayLength < ArrayMin of ", CampName, " ", GetArrayLength(CampName), " < ", GetArrayMin(CampName))
				fmt.Println()
				require := GetArrayMax(CampName) - GetArrayLength(CampName)
				fmt.Println("Require ", require, " from ", CampName)
				fmt.Println()
				GetNumbers(CampName, require)
			}
		}
		if GetArrayLength(CampName) > 0 {
			fmt.Println("ArrayLength of ", CampName, "> 0  Length : ", GetArrayLength(CampName))
			fmt.Println()
			if GetConLimit(CampName) > 0 {

				fmt.Println("ConLimit of ", CampName, "> 0  Limit : ", GetConLimit(CampName))
				fmt.Println()

				if GetConLimit(CampName) < GetArrayLength(CampName) {

					fmt.Println("ConLimit of ", CampName, "<  ArrayLength ", GetConLimit(CampName), " < ", GetArrayLength(CampName))
					fmt.Println()
					var Phones = GetListPhoneSet(CampName, GetConLimit(CampName))
					if Phones != nil {

						for _, value := range Phones {
							//uuid := GetUuid()
							//go DialServer(value, uuid, CampName)
							fmt.Println("Poped ", value)
							fmt.Println()
						}

						//
					}
					//make([]string, GetConLimit(CampName))
					/*for i := 0; i < GetConLimit(CampName); i++ {
						//GetListPhones(CampName)
						Phones[i] = GetListPhones(CampName)
					}
					for i := 0; i < GetConLimit(CampName); i++ {
						//GetListPhones(CampName)
						Phones[i] = GetListPhones(CampName)
					}
					*/
				} else {
					fmt.Println("ConLimit of ", CampName, "<  ArrayLength ", GetConLimit(CampName), " < ", GetArrayLength(CampName))
					fmt.Println()
					var Phones = GetListPhoneSet(CampName, GetArrayLength(CampName))
					if Phones != nil {

						for _, value := range Phones {
							//uuid := GetUuid()
							//go DialServer(value, uuid, CampName)
							fmt.Println("Poped ", value)
							fmt.Println()
						}

						//
					}
				}

			} else {
				fmt.Println("Sleeps Con limit < = 0 of ", CampName)
				fmt.Println()
				//time.Sleep(1000 * time.Millisecond)
				return
			}
		} else {
			fmt.Println("Sleeps ArrayLength < = 0 of ", CampName)
			fmt.Println()
			//time.Sleep(1000 * time.Millisecond)
			return
		}
	}

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

	fmt.Println("Get numbers hit of campaign ", CampName)
	//Nums := []string{}
	fmt.Println()
	//for {
	cnt := GetPhoneCount(CampName)
	fmt.Println("Phone count of ", CampName, " ", cnt)
	if cnt > 0 {
		fmt.Println("COUNT > 0 ", CampName)
		fmt.Println()
		var p ResPhone

		if Max == 0 {
			SetCampaignStatus(CampName, "0")
			return
		}

		url := fmt.Sprintf("http://localhost:8083/DVP/API/1.0/DialerApi/FillCampaignPhones/%s/%d", CampName, Max)
		fmt.Println("URL hit ", url)
		fmt.Println()
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
			fmt.Println()
			fmt.Println("Number :", val.Phone)
			fmt.Println()
			Phn := c.Cmd("LPUSH", CampName, val.Phone)
			fmt.Println()
			fmt.Println(Phn, " PUSHED to ", CampName)
			fmt.Println()

		}

		/*if st != 0 {
			fmt.Println("Not initial")
			GetPhonesFromList(CampName)
		} else {
			fmt.Println("Initial")
			return

		}*/
		fmt.Println(CampName, " Returned ")
		return
	} else {
		fmt.Println("Phone count of ", CampName, " ", cnt, "and returned")
		return
	}

	//time.Sleep(1000 * time.Millisecond)
	//}

}

func GetPhonesFromList(CampName string) {

	//CampSt := GetCampaignStatus(CampName)
	//need while loop to check status of campaign

	fmt.Println("Hit list of ", CampName)
	fmt.Println()
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
		fmt.Println()
	}
	CmpMn := fmt.Sprintf("%s_Min", CampName)
	CmpMx := fmt.Sprintf("%s_Max", CampName)
	MinPhns, _ := c.Cmd("GET", CmpMn).Int()
	MaxPhns, _ := c.Cmd("GET", CmpMx).Int()
	LenPhns, _ := c.Cmd("LLEN", CampName).Int()
	fmt.Println(CampName, " Hit ", "Length ", LenPhns)
	if LenPhns == 0 {

		fmt.Println("Campaign Stopped ", CampName)
		return

	}
	LPhns := c.Cmd("LPOP", CampName).String()

	//uuid := GetUuid()
	//go DialServer(LPhns, uuid)

	fmt.Println("Poped ", LPhns, " of ", CampName)
	fmt.Println()

	newLen := LenPhns - 1

	if newLen <= MinPhns {

		NewFill := MaxPhns - newLen
		fmt.Println("Max %d- Length %d of %s", MaxPhns, newLen, CampName)
		fmt.Println()
		fmt.Println("NewFill", NewFill)
		fmt.Println()
		fmt.Println(CampName, " remaining in DB is ", GetPhoneCount(CampName))
		if GetPhoneCount(CampName) > 0 {
			fmt.Println("Load numbers from Campaign ", CampName)
			fmt.Println()
			GetNumbers(CampName, NewFill)
		} else {
			if newLen > 0 {
				fmt.Println()
				fmt.Println("Length of ", CampName, " is ", newLen)
				fmt.Println()
				fmt.Println("List is not empty of ", CampName)
				GetPhonesFromList(CampName)
			} else {
				fmt.Println()
				fmt.Println(CampName, " Stopped")
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

	SetMin := c.Cmd("SET", MnName, Min)
	SetMax := c.Cmd("SET", MxName, Max)
	fmt.Println("Min ", MnName, " ", Min, "Mx ", MxName, " ", Max)
	fmt.Println()
	fmt.Println("Campaign Max ", SetMax)
	fmt.Println()
	fmt.Println("Campaign Min ", SetMin)
	fmt.Println()

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
	fmt.Println("Campaign  ", Campaign, " Status 1 ", s)
	fmt.Println()
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
	fmt.Println("Getting phone count ", CampName)
	fmt.Println()
	var pCount ResultPCount
	url := fmt.Sprintf("http://localhost:8083/DVP/API/1.0/DialerApi/PhoneCount/%s", CampName)
	fmt.Println("URL hit ", url)
	fmt.Println()
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

	fmt.Println("Dial Server ", phoneNumber)
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
func SetCampaignConLimit(Campaign string, Limit int) {

	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	CampCon := fmt.Sprintf("%s_Con", Campaign)
	s := c.Cmd("SET", CampCon, Limit)
	fmt.Println("Campaign ", Campaign, "  ", CampCon, " ", Limit, " ", s)
	fmt.Println()
}
func GetArrayLength(Campaign string) int {

	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	//CampSt := fmt.Sprintf("%s_St", Campaign)
	st, _ := c.Cmd("LLEN", Campaign).Int()
	return st
}
func GetArrayMax(Campaign string) int {

	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	CampMax := fmt.Sprintf("%s_Max", Campaign)
	st, _ := c.Cmd("GET", CampMax).Int()
	return st
}
func GetArrayMin(Campaign string) int {

	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	CampMin := fmt.Sprintf("%s_Min", Campaign)
	st, _ := c.Cmd("GET", CampMin).Int()
	return st
}
func GetConLimit(Campaign string) int {

	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	CampCon := fmt.Sprintf("%s_Con", Campaign)
	st, _ := c.Cmd("GET", CampCon).Int()
	return st
}
func GetListPhones(CampName string) string {
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	st := c.Cmd("LPOP", CampName).String()
	return st
}
func GetListPhoneSet(CampName string, Size int) []string {
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err == nil {
		//fmt.Println("Redis client connected for GetFrom List. ", ListName)

	} else {
		fmt.Println("GetFrom List Error ", err.Error())
	}

	S := strconv.Itoa(Size)
	L := (Size - 1)
	E := strconv.Itoa(GetArrayLength(CampName))
	//fmt.Println(S)

	c.Append("LRANGE", CampName, 0, L)
	c.Append("LTRIM", CampName, S, E)

	listPop, errPop := c.GetReply().List()
	if errPop == nil {

		fmt.Println("Poped Numbers ", listPop)
		fmt.Println()
		_, errRem := c.GetReply().Str()
		if errRem == nil {
			return listPop
		} else {
			return nil
		}
	} else {
		return nil
	}

}
