package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"skylib/utils"
)

var url = "http://127.0.0.1:9002/"

var tokenClient = ""

type sendReq struct {
	Response *http.Response
	Body     string
}

func SetUrl(apiUrl string)  {
	url = apiUrl
}

func SetToken(token string) {
	tokenClient = token
}

func SendRequestJsonString(uri string, stringJson string) sendReq {
	return sendRequest(uri, utils.GetCrypted([]byte(stringJson)))
}

func SendRequestJsonStruct(uri string, structure interface{}) sendReq {
	jsonStruct, err := json.Marshal(structure)
	CheckErr(err)
	return sendRequest(uri, utils.GetCrypted([]byte(jsonStruct)))
}

func sendRequest(uri string, cryptedJson []byte) sendReq {
	ioEncryptJson := bytes.NewReader(cryptedJson)

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
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		T.Log("Error response: ", err.Error())
	}
	sendResponse.Response = resp

	sendResponse.Body = string(buf.Bytes())
	T.Log(sendResponse.Body)
	return sendResponse
}
