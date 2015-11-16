package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//func Post(serviceurl string, postData string) bool {
//	fmt.Println("URL:>", serviceurl)

//	postData = postData
//	fmt.Println("PostData:>", postData)

//	var jsonData = []byte(postData)
//	req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(jsonData))
//	req.Header.Set("Content-Type", "application/json")

//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		//panic(err)
//		return false
//	}
//	defer resp.Body.Close()

//	fmt.Println("response Status:", resp.Status)
//	fmt.Println("response Headers:", resp.Header)
//	body, _ := ioutil.ReadAll(resp.Body)
//	result := string(body)
//	fmt.Println("response Body:", result)
//	if result != "No matching resources at the moment" {
//		fmt.Println("Return true")
//		return true
//	}

//	fmt.Println("Return false")
//	return false
//}

func Get(serviceurl, path, param string) string {
	request := fmt.Sprintf("http://%s", serviceurl)

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
		return tmx
	}

	return ""
}
