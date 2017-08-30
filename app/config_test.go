package app

import (
	"testing"
	"path/filepath"
	//"runtime"
	"runtime"
)

func TestGetConfig(t *testing.T) {
	_, b, _, _  := runtime.Caller(0)
	ThisDir = filepath.Dir( b) + "/../"
	db := GetConfig("database")
	if db["username"] == nil {
		t.Error("Error get config")
	} //else {
		//t.Error("Error get config")
		//t.Log(db["username"])
	//}
}

func TestGetConnection(t *testing.T) {
	_, b, _, _  := runtime.Caller(0)
	ThisDir = filepath.Dir( b) + "/../"

	db := GetConfig("database")

	InitLog()

	t.Log(db["username"])

	GetConnection()
	if DB == nil {
		t.Error("Error connect")
	}
}