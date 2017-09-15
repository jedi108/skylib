package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

var config map[string]*json.RawMessage

var ThisDir = ""
var ConfigDevFileJson = "config.dev.json"
var ConfigTestsFileJson = "config.tests.json"
var ConfigProdFileJson = "config.json"

func GetConfig(section string) map[string]interface{} {
	err := InitAppConfig()
	if err != nil {
		panic(err)
	}

	conf, err := GetConfigSelection(section)
	if err != nil {
		panic(err)
	}

	return conf
}

func InitAppConfig() error {
	if (len(config) == 0) {
		bs, err := ioutil.ReadFile(ThisDir + ConfigDevFileJson)

		if err != nil {
			bs, err = ioutil.ReadFile(ThisDir + ConfigProdFileJson)
		}

		if err != nil {
			fmt.Println("config not open", err.Error())
			return err
		}
		str := string(bs)
		config_raw := []byte(str)
		err = json.Unmarshal(config_raw, &config)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetMapOfConfig() map[string]*json.RawMessage {
	return config
}

func GetConfigSelection(section string) (map[string]interface{}, error) {
	var config_section map[string]interface{}
	_, ok := config[section]
	if ok == false {
		return make(map[string]interface{}), errors.New("map not selection")
	}
	err := json.Unmarshal(*config[section], &config_section)
	if err != nil {
		return make(map[string]interface{}), err
	}
	return config_section, nil
}
