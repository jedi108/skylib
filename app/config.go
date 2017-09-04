package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var config map[string]*json.RawMessage

const PrivKey = "DM0HozajCtlOryLKnVnhS226Nq3Gsm7AGLeShIL7WBg="
const PubKey = "NvalbqygT8G2jp4IXJCW1OHia3LnDCqaqqV0i6w5Mys="

var ThisDir = ""
var ConfigDevFileJson = "config.dev.json"
var ConfigTestsFileJson = "config.tests.json"
var ConfigProdFileJson = "config.json"

func GetConfig(section string) map[string]interface{} {
	if (len(config) == 0) {

		bs, err := ioutil.ReadFile(ThisDir + ConfigDevFileJson)

		if err != nil {
			bs, err = ioutil.ReadFile(ThisDir + ConfigProdFileJson)
		}

		if err != nil {
			fmt.Println("config not open")
			panic(err)
		}
		str := string(bs)
		config_raw := []byte(str)
		err = json.Unmarshal(config_raw, &config)
		if err != nil {
			panic(err)
		}
	}
	var config_section map[string]interface{}
	err := json.Unmarshal(*config[section], &config_section)
	if err != nil {
		panic(err)
	}
	return config_section
}
