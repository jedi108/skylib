package test

import (
	"path/filepath"
	"runtime"
	"skylib/app"
)

var isInit = false

var DefaultDir = "../"

var InitTest = func() {
	if isInit == false {
		app.ThisDir = DefaultDir
		isInit = true
		app.InitLog()
		app.GetConnection()
	}
}

func InitPathConfigTest() {
	_, b, _, _ := runtime.Caller(0)
	app.ThisDir = filepath.Dir(b) + "/../"
}

func InitTestDatabase() {
	if isInit {
		return
	}
	isInit = true
	app.ConfigDevFileJson = app.ConfigTestsFileJson
}
