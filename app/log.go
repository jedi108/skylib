package app

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logday int
var config_log map[string]interface{}
var fileLogName string

func SetFileLogName(name string) {
	fileLogName = name
}

func InitLog() {
	config_log = GetConfig("log")
	createNewLogFile()
}

func createNewLogFile() {
	t := time.Now()
	logday = t.Day()
	filepath := fmt.Sprintf("%v/%v%v-%v-%v.log", config_log["path"], fileLogName, t.Year(), int(t.Month()), t.Day())
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	//defer f.Close()
	log.SetOutput(f)
}

func Log(s string) {
	t := time.Now()
	if (logday != t.Day()) {
		createNewLogFile()
	}
	log.Printf(s)
}
