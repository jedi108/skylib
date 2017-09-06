package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

var FillStructVerbose = false

var RandomText = func(length int) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func FillStructure(s interface{}, m *map[string]string) error {
	for k, v := range *m {
		err := setField(&s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func setField(obj *interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(*obj).Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
		return nil
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {

		var val64 float64
		var valInt int64
		var errType error

		switch structFieldType.Kind() {
		case reflect.Float64, reflect.Float32:
			val64, errType = strconv.ParseFloat(val.String(), 32)
			structFieldValue.SetFloat(val64)

			if FillStructVerbose {
				fmt.Println("val64:", val64)
				fmt.Println(" structFieldValue", structFieldValue.String())
				fmt.Println(" val", val)
			}

			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valInt, errType = strconv.ParseInt(val.String(), 64, 32)
			structFieldValue.SetInt(valInt)

			if FillStructVerbose {
				fmt.Println("valInt:", valInt)
				fmt.Println("val.String():", val.String())
			}

			return nil
		default:

			if FillStructVerbose {
				fmt.Println("-----------------------------")
				fmt.Println("name", name)
				fmt.Println("val", val)
				fmt.Println("structFieldType", structFieldType)
				fmt.Println("val.Type()", val.Type())
				fmt.Println("structFieldValue", structFieldValue)
				fmt.Println("structValue", structValue)
			}

			invalidTypeError := errors.New("Value type didn't match obj field type")
			return invalidTypeError
		}

		if errType != nil {

			if FillStructVerbose {
				fmt.Println("-----------------------------")
				fmt.Println("name", name)
				fmt.Println("structFieldType", structFieldType)
				fmt.Println("val.Type()", val.Type())
				fmt.Println("val", val)
				fmt.Println("structFieldValue", structFieldValue)
				fmt.Println("structValue", structValue)
			}

			invalidTypeError := errors.New("Provided value type didn't match obj field type")
			return invalidTypeError
		}

		return nil
	} else {

		if FillStructVerbose {
			fmt.Println("-------is ok")
			fmt.Println("name", name)
			fmt.Println("structFieldType", structFieldType)
			fmt.Println("val.Type()", val.Type())
			fmt.Println("val", val)
			fmt.Println("structFieldValue", structFieldValue)
			fmt.Println("structValue", structValue)
		}

	}

	structFieldValue.Set(val)
	return nil
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
