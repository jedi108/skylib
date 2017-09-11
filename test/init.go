/**
	Initialize the fixtures tests
	this packege helpers for tests

	use parametr -p=1 to run all functional tests in one proccess
	example:

		go test ./... -p=1 -v

 */
package test

import (
	"skylib/app"
	"testing"
)

var isInit = false

/** init for "go test ./... -p=1" */
var InitForTest = func(t *testing.T) {
	if isInit == false {
		T = t
		app.ThisDir = DefaultDir
		app.ConfigDevFileJson = app.ConfigTestsFileJson
		isInit = true
		app.InitLog()
		app.GetConnection()
		RunFixtures()
	}
}

/** init for "go test ./... -p=1" */
var InitForStressTest = func(t *testing.T) {
	if isInit == false {
		T = t
		app.ThisDir = DefaultDir
	}
}
