package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"github.com/jedi108/skylib/utils"
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


func SendStressRequestJsonStruct(uri string, structure interface{}) {
	jsonStruct, err := json.Marshal(structure)
	CheckErr(err)
	sendStressRequest(uri, utils.GetCrypted([]byte(jsonStruct)))
}

func SendRequestJsonString(uri string, stringJson string) sendReq {
	return SendRequest(uri, utils.GetCrypted([]byte(stringJson)))
}

func SendRequestJsonStruct(uri string, structure interface{}) sendReq {
	jsonStruct, err := json.Marshal(structure)
	CheckErr(err)
	return SendRequest(uri, utils.GetCrypted([]byte(jsonStruct)))
}

func SendRequest(uri string, cryptedJson []byte) sendReq {
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
	//T.Log(sendResponse.Body)
	return sendResponse
}


func PerfectSendRequestJsonStruct(uri string, structure interface{}) (sendReq, error) {
	jsonStruct, err := json.Marshal(structure)
	CheckErr(err)
	req, err := PerfectSendRequest(uri, utils.GetCrypted([]byte(jsonStruct)))
	if err != nil {
		return sendReq{}, err
	}
	return req, nil
}

func PerfectSendRequest(uri string, cryptedJson []byte) (sendReq, error) {
	ioEncryptJson := bytes.NewReader(cryptedJson)
	//T.Log(url+uri)
	req, err := http.NewRequest("POST", url+uri, ioEncryptJson)
	if err != nil {
		return sendReq{}, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("accept", "application/json, text/plain, */*")
	if tokenClient != "" {
		req.Header.Set("Token", tokenClient)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return sendReq{}, err
	}

	//CheckErr(err)
	var sendResponse = sendReq{}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		T.Log("Error response: ", err.Error())
	}
	sendResponse.Response = resp

	sendResponse.Body = string(buf.Bytes())
	//T.Log(sendResponse.Body)
	return sendResponse, nil
}

func sendStressRequest(uri string, cryptedJson []byte) {
	ioEncryptJson := bytes.NewReader(cryptedJson)

	for i:=1; i<100000; i++ {
		request(ioEncryptJson, uri)
	}

}

func request(ioEncryptJson *bytes.Reader, uri string)  {
	req, err := http.NewRequest("POST", url+uri, ioEncryptJson)

	CheckErr(err)
	if tokenClient != "" {
		req.Header.Set("Token", tokenClient)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("accept", "application/json, text/plain, */*")
	client := &http.Client{}
	client.Do(req)
	//defer resp.Body.Close()
}