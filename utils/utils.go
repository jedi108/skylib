package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mdp/sodiumbox"
	"io"
	"strconv"
	"skylib/app"
	b64 "encoding/base64"
	"reflect"
	"io/ioutil"
	"net/http"
)

type keyPair struct {
	secretKey [32]byte
	publicKey [32]byte
}

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func getKeyPair() *keyPair {
	publicKeySlice, _ := b64.StdEncoding.DecodeString(app.PubKey)
	publicKey := *new([32]byte)
	copy(publicKey[:], publicKeySlice[0:32])

	privateKeySlice, _ := b64.StdEncoding.DecodeString(app.PrivKey)
	privateKey := *new([32]byte)
	copy(privateKey[:], privateKeySlice[0:32])
	return &keyPair{
		secretKey: privateKey,
		publicKey: publicKey,
	}
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
	de := GetDecrypted([]byte(buff.String()))
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
	publicKeySlice, _ := b64.StdEncoding.DecodeString(app.PubKey)
	publicKey := *new([32]byte)
	copy(publicKey[:], publicKeySlice[0:32])
	sealedMsg, _ := sodiumbox.Seal(myvar, &publicKey)
	return sealedMsg.Box
}

//Декриптования
func GetDecrypted(sealedMsg []byte) *sodiumbox.Message {
	testKeyPair := getKeyPair()
	msg, e := sodiumbox.SealOpen(sealedMsg, &testKeyPair.publicKey, &testKeyPair.secretKey)
	if e != nil {
		panic(e)
	}
	return msg
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
	bodyRaw := GetDecrypted(StreamToByte(bodyValues)).Content
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
