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
	"skylib/components/redisCache"
	"testing"
	"os"
)

var isInit = false

type ConfigTests struct {
	IsRunFixtures    bool
	IsInitSql        bool
	IsCacheClear     bool
	IsTarantoolQueue bool
	RootProjectDir   string
	LogFileName      string
	ExitAfterInit    bool
}

/**
	New init tests with configurable
 */
var InitConfigTests = func(confiForTest ConfigTests, t *testing.T) {
	if isInit == false {
		T = t

		if confiForTest.RootProjectDir == "" {
			app.ThisDir = DefaultDir
		} else {
			app.ThisDir = confiForTest.RootProjectDir
			DirFixtures = app.ThisDir + "tests/fixtures/"
		}

		app.ConfigDevFileJson = app.ConfigTestsFileJson

		app.SetFileLogName(confiForTest.LogFileName)
		app.InitLog()

		if confiForTest.IsInitSql {
			app.GetConnection()
		}

		if confiForTest.IsRunFixtures {
			RunFixtures()
		}

		t.Log("cache+++++++++++++++")
		if confiForTest.IsCacheClear {
			redisCache.Init()
			t.Log("clear cache+++++++++++++++")
			keys := redisCache.Client.Keys("*")
			for _, v := range keys.Val() {
				redisCache.Client.Del(v)
			}
		}

		isInit = true
	}
	if confiForTest.ExitAfterInit == true {
		os.Exit(0)
	}
}

/**
	DEPRECATED
	@see: InitConfigTests
*/
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

/**
	DEPRECATED
	@see: InitConfigTests
*/
var InitForStressTest = func(t *testing.T) {
	if isInit == false {
		T = t
		app.ConfigDevFileJson = app.ConfigTestsFileJson
		app.ThisDir = DefaultDir
		RunFixtures()
	}
}
