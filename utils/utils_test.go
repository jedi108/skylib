package utils

import (
	"io"
	"testing"
)

var rowCrypted []byte
var str = `{"User":"whislai","98765431:"test","IPInfo":{"ip":"192.228.148.61","hostname":"broadband.time.net.my","city":"Shah Alam","region":"Selangor","country":"MY","loc":"3.0833,101.5333","org":"AS9930 TIMENET IP Hostmasters","postal":"40100"}}`

func TestEncrytpDecrypt(t *testing.T) {
	rowCrypted = GetCrypted([]byte(str))
	t.Log("rowCrypted:", rowCrypted)
	str2 := string(GetDecrypted(rowCrypted).Content)
	t.Log("GetDecrypted: ", str2)
	if str != str2 {
		t.Error("error decrypt")
	}
}

func TestReqEncrytpDecrypt(t *testing.T) {
	var vz io.Reader
	vz = ToBinJson(rowCrypted)

	t.Log("buffer:", vz)

	bodyValues := GetValuesFromArray(StreamToByte(vz))

	bw := StreamToByte(bodyValues)
	t.Log("StreamToByte bodyValues:", StreamToByte(bodyValues))

	t.Log("len1: ", len(string(rowCrypted)))
	t.Log("len2: ", len(string(bw)))

	if len(string(rowCrypted)) != len(bw) {
		t.Error("Not equal len")
	}

	bodyRaw := string(GetDecrypted(bw).Content)
	t.Log("str3:", bodyRaw)

	if bodyRaw != str {
		t.Error("Not equal json before and after cryptdecrypt")
	}
}
