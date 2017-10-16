/**
	Initialize the fixtures tests
	this packege helpers for tests

	use parametr -p=1 to run all functional tests in one proccess
	example:

		go test ./... -p=1 -v

 */
package test

import (
	"os"
	"skylib/app"
	"skylib/components/redisCache"
	"skylib/components/tarantoolQ"
	"testing"
)

var isInit = false

type ConfigTests struct {
	IsRunFixtures    bool
	IsInitSql        bool
	IsCacheClear     bool
	IsDisableCache	 bool
	IsTarantoolQueue bool
	RootProjectDir   string
	LogFileName      string
	ExitAfterInit    bool
	ReinitConfig     bool
}

/**
	New init tests with configurable
 */
var InitConfigTests = func(confiForTest ConfigTests, t *testing.T) {
	if confiForTest.ReinitConfig {
		isInit = false
	}
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

		tarantoolQ.NewConnect()

		if confiForTest.IsDisableCache == false {
			redisCache.Init()
		}


		//tarantoolQ.InitFromConfigTarantoolQueue()

		if confiForTest.IsCacheClear {
			if redisCache.IsCacheEnable() {
				keys := redisCache.Client.Keys("*")
				for _, v := range keys.Val() {
					redisCache.Client.Del(v)
				}
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
