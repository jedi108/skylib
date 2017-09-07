package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdp/sodiumbox"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"skylib/app"
	"strconv"
	b64 "encoding/base64"
)

type keyPair struct {
	secretKey [32]byte
	publicKey [32]byte
}

var isInitKeysConfig = false;
var PublicKeydata = ""
var PublicKeyResponse = ""
var PrivateKeyData = ""

func InitKeys() {
	if isInitKeysConfig == false {
		ks := app.GetConfig("keys")
		PublicKeydata = ks["PublicKeydata"].(string)
		PublicKeyResponse = ks["PublicKeyResponse"].(string)
		PrivateKeyData = ks["PrivateKeyData"].(string)
		isInitKeysConfig = true
	}
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func getKeyPairInclude(strPub string, strPriv string) *keyPair {
	publicKeySlice, _ := b64.StdEncoding.DecodeString(strPub)
	publicKey := *new([32]byte)
	copy(publicKey[:], publicKeySlice[0:32])

	privateKeySlice, _ := b64.StdEncoding.DecodeString(strPriv)
	privateKey := *new([32]byte)
	copy(privateKey[:], privateKeySlice[0:32])
	return &keyPair{
		secretKey: privateKey,
		publicKey: publicKey,
	}
}

func getKeyPair() *keyPair {
	InitKeys()
	return getKeyPairInclude(PublicKeydata, PrivateKeyData)
}

func getKeyPairRespone() *keyPair {
	InitKeys()
	return getKeyPairInclude(PublicKeyResponse, PrivateKeyData)
}

func DecryptRequest(j []byte) []byte {
	var raw map[int]byte
	e := json.Unmarshal(j, &raw)
	if e != nil {
		panic(e)
	}
	var buffStr string
	buff := bytes.NewBufferString(buffStr)
	for i := 0; i < len(raw); i++ {
		buff.WriteByte(raw[i])
	}
	de, _ := GetDecrypted([]byte(buff.String()))
	return de.Content
}

func GetValuesFromArray(j []byte) (io.Reader) {
	var raw map[int]byte
	e := json.Unmarshal(j, &raw)
	if e != nil {
		errors.New("No unmarshal")
		panic(e)
	}
	var buffStr string
	buff := bytes.NewBufferString(buffStr)
	for i := 0; i < len(raw); i++ {
		buff.WriteByte(raw[i])
	}
	return buff
}

//Криптование
func GetCrypted(myvar []byte) []byte {
	InitKeys()
	return getEncrypt(myvar, PublicKeydata)
}

//Криптование
func GetCryptedResponse(myvar []byte) []byte {
	InitKeys()
	return getEncrypt(myvar, PublicKeyResponse)
}

func getEncrypt(myvar []byte, strPub string) []byte {
	publicKeySlice, _ := b64.StdEncoding.DecodeString(strPub)
	publicKey := *new([32]byte)
	copy(publicKey[:], publicKeySlice[0:32])
	sealedMsg, err := sodiumbox.Seal(myvar, &publicKey)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	return sealedMsg.Box
}

//Декриптования
func GetDecrypted(sealedMsg []byte) (*sodiumbox.Message, error) {
	testKeyPair := getKeyPair()
	return getDecrypt(sealedMsg, testKeyPair)
}

//Декриптования
func GetDecryptedResponse(sealedMsg []byte) (*sodiumbox.Message, error) {
	testKeyPair := getKeyPairRespone()
	return getDecrypt(sealedMsg, testKeyPair)
}

func getDecrypt(sealedMsg []byte, testKeyPair *keyPair) (*sodiumbox.Message, error) {
	if len(sealedMsg)<33 {
		return nil, errors.New("open message error")
	}
	msg, e := sodiumbox.SealOpen(sealedMsg, &testKeyPair.publicKey, &testKeyPair.secretKey)
	if e != nil {
		return nil, errors.New("encrypt error")
	}
	return msg, e
}

func ToBinJson(rc []byte) io.Reader {
	var buffStr string
	buff := bytes.NewBufferString(buffStr)
	buff.WriteString("{")
	for i := 0; i < len(rc); i++ {
		buff.Write([]byte("\""))
		buff.WriteString(strconv.Itoa(i))
		buff.Write([]byte("\""))
		buff.Write([]byte(": "))
		buff.WriteString(fmt.Sprintf("%v", rc[i]))

		if i < len(rc)-1 {
			buff.Write([]byte(", "))
		}
	}
	buff.WriteString("}")
	return buff
}

func FromEncryptRespToJsonString(req *http.Request) map[string]string {
	body, _ := ioutil.ReadAll(req.Body)
	bodyValues := GetValuesFromArray(body)
	bodyCrypted, _ := GetDecrypted(StreamToByte(bodyValues))
	bodyRaw := bodyCrypted.Content
	var bodyJson map[string]string
	json.Unmarshal(bodyRaw, &bodyJson)
	return bodyJson
}

func Empty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan:
		return v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	}
	return false
}
