package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

/* var newLoger *logrus.Logger

func InitializeLogrusLogger() {

	mate, _ := logrus_mate.NewLogrusMate(logrus_mate.ConfigFile("mate.conf"))

	newLoger := logrus.New()

	mate.Hijack(newLoger, "mike")

	newLoger.Infoln("hello std logger is hijack by mike")
	newLoger.Warningln("hello std logger is hijack by mike")
	newLoger.Debugln("hello std logger is hijack by mike")
	newLoger.Errorln("hello std logger is hijack by mike")

} */

var enableLog bool = true

func isJSON(s string) (bool, map[string]string) {
	var dcReasonData map[string]string
	result := json.Unmarshal([]byte(s), &dcReasonData)
	return (result == nil), dcReasonData

}

func EnableConsoleInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		matched, _ := regexp.MatchString("^(addreasons)", scanner.Text())

		if scanner.Text() == "logon" {
			enableLog = true
			fmt.Println("LOG ENABLED")
		} else if scanner.Text() == "logoff" {
			enableLog = false
			fmt.Println("LOG DISABLED")
		} else if scanner.Text() == "reloadreasons" {
			GetDisconnectReasons()
			fmt.Println("DISCONNECTION REASONS LOADED")
		} else if scanner.Text() == "flushreasons" {
			RedisRemove("DisconnectReasonMap")
			fmt.Println("ALL REASONS FLUSHED SUCCESSFULLY")
		} else if matched == true {
			inputReasons := strings.Split(scanner.Text(), "|")
			if len(inputReasons) == 2 {
				isJSONStr, jsonData := isJSON(inputReasons[1])

				if isJSONStr {
					RedisHMSet("DisconnectReasonMap", jsonData)
					fmt.Println("ADD REASONS SUCCESS")
				} else {
					fmt.Println("INVALID FORMAT - DATA IS NOT A VALID JSON")
				}

			} else {
				fmt.Println("INVALID FORMAT - PLEASE USE | TO SEPARATE COMMAND AND JSON DATA")
			}
		} else {
			fmt.Println("")
		}
	}

}

func DialerLog(message string) {
	if enableLog {
		fmt.Println(message)
	}
}
