package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func FillStructureInterface(s interface{}, m *map[string]interface{}) error {
	for k, v := range *m {
		err := SetFieldInterface(&s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetFieldInterface(modelStructure *interface{}, keyMap string, valueMap interface{}) error {
	structData := reflect.ValueOf(*modelStructure).Elem()
	structValue := structData.FieldByName(keyMap)
	if !structValue.IsValid() {
		return nil
	}

	if !structValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", keyMap)
	}

	switch valueOfType := valueMap.(type) {
	case float32:
		if structValue.Type().Kind() != reflect.Float32 {
			return errors.New("No valid type Float32 of field `" + keyMap + "`")
		}
		switch structValue.Type().Kind() {
		case reflect.Float32:
			structValue.SetFloat(float64(valueOfType))
		default:
			return errors.New("No float valid type")
		}
		return nil
	case float64:
		if (structValue.Type().Kind() != reflect.Float64) && (structValue.Type().Kind() != reflect.Float32) && (structValue.Type().Kind() != reflect.Int) {
			return errors.New("No valid type Float64 of field `" + keyMap + "`")
		}
		switch structValue.Type().Kind() {
		case reflect.Float32:
			structValue.SetFloat(float64(valueOfType))
		case reflect.Float64:
			structValue.SetFloat(float64(valueOfType))
		case reflect.Int:
			structValue.SetInt(int64(valueOfType))
		default:
			return errors.New("No valid float type")
		}
		return nil
	case int:
		if structValue.Type().Kind() != reflect.Int {
			return errors.New("No valid type int of field `" + keyMap + "`")
		}
		structValue.SetInt(int64(valueOfType))
	case string:
		if structValue.Type().Kind() != reflect.String {
			return errors.New("No valid type of field `" + keyMap + "`")
		}
		structValue.SetString(valueOfType)
	case bool:
		if structValue.Type().Kind() != reflect.Bool {
			return nil
			//return errors.New("No valid type of field `" + keyMap + "`")
		}
		structValue.SetBool(valueOfType)
	default:
		return errors.New("No valid type")
	}
	return nil
}
