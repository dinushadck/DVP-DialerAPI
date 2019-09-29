package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

func EnableConsoleInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		matched, _ := regexp.MatchString(`^(addreasons|)`, scanner.Text())

		if scanner.Text() == "logon" {
			enableLog = true
			fmt.Println("LOG ENABLED")
		} else if scanner.Text() == "logoff" {
			enableLog = false
			fmt.Println("LOG DISABLED")
		} else if scanner.Text() == "reloadreasons" {
			GetDisconnectReasons()
			fmt.Println("DISCONNECTION REASONS LOADED")
		} else if matched == true {
			fmt.Println("ADD REASONS")
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
