package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"skylib/utils"
)

const url = "http://127.0.0.1:9002/"

var tokenClient = ""

type sendReq struct {
	Response *http.Response
	Body     string
}

func SetToken(token string) {
	tokenClient = token
}

func SendRequest(uri string, structure interface{}) sendReq {
	jsonStruct, err := json.Marshal(structure)
	CheckErr(err)
	PrintLn(string(jsonStruct))
	encryptedJsonStruct := utils.GetCrypted([]byte(jsonStruct))
	ioEncryptJson := bytes.NewReader(encryptedJsonStruct)

	req, err := http.NewRequest("POST", url+uri, ioEncryptJson)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("accept", "application/json, text/plain, */*")
	if tokenClient != "" {
		req.Header.Set("Token", tokenClient)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	CheckErr(err)
	var sendResponse = sendReq{}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	sendResponse.Response = resp

	sendResponse.Body = string(buf.Bytes())
	T.Log(sendResponse.Body)
	return sendResponse
}
