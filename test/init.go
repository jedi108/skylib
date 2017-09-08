package test

import (
	"os"
	"skylib/app"
	"testing"
)

var isInit = false

/** init for "go run project" */
var InitTest = func() {
	if isInit == false {
		startConfig()
	}
}

/** init for "go test -v ..." */
var InitForTest =  func(t *testing.T) {
	if isInit == false {
		T = t
		startConfig()
	}
}

func startConfig()  {
	app.ThisDir = DefaultDir
	app.ConfigDevFileJson = app.ConfigTestsFileJson
	isInit = true
	app.InitLog()
	app.GetConnection()
	RunFixtures()
}

func Args() []string {
	return os.Args[1:]
}