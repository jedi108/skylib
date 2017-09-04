package test

import (
	"skylib/app"
	"runtime"
	"path/filepath"
)

var isInit = false

func InitPathConfigTest()  {
	_, b, _, _  := runtime.Caller(0)
	app.ThisDir = filepath.Dir( b) + "/../"
}

func InitTestDatabase()  {
	if isInit {
		return
	}
	isInit = true
	app.ConfigDevFileJson = app.ConfigTestsFileJson
}