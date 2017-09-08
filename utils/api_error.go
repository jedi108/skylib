package utils

import (
	"encoding/json"
	"net/http"
	"log"
)

type Error struct {
	Field   string
	Message string
	Code    int
}

type ApiErrors struct {
	Code	int
	Message string
	Errors []Error
}

func ResponseError(HttpError int, s string, resp http.ResponseWriter) {
	log.Println(s)
	resp.WriteHeader(HttpError)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(GetApiError(s))
}

func GetApiError(error string) []byte {

	var s []Error

	e := &Error{
		Message: error,
	}
	s = append(s, *e)
	c := &ApiErrors{
		Code: 400,
		Message: "test",
		Errors: s,
	}

	rr, _ := json.Marshal(c)
	return rr
}

func GetApiErrors(errors []string) []byte {

	var s []Error

	//x := new(ApiErrors)

	for _, v := range errors {

		e := &Error{
			Message: v,
		}

		s = append(s, *e)
	}

	c := &ApiErrors{
		Code: 400,
		Message: "Seriously exception",
		Errors: s,
	}

	rr, _ := json.Marshal(c)
	return rr
}

func FinishApiErrors(s []Error) []byte {
	c := &ApiErrors{
		Code: 400,
		Message: "Molodec! :-)",
		Errors: s,
	}

	rr, _ := json.Marshal(c)
	return rr

}

