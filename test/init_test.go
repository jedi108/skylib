package test

import (
	"skylib/app"
	"testing"
)

func TestInit(t *testing.T)  {

	if isInit == true {
		t.Error("Is init")
	}

	InitTestDatabase()

	if app.ConfigDevFileJson != app.ConfigTestsFileJson {
		t.Error("Is not init")
	}
}
