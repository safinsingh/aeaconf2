package compat

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/pkg/errors"
	"github.com/safinsingh/aeaconf2"
)

func SeparateConfig(data []byte) ([]byte, []byte, error) {
	re := regexp.MustCompile(`(?m)^[[:space:]]*---[[:space:]]*$`)
	loc := re.FindIndex(data)
	if loc == nil {
		return nil, nil, errors.New("could not split configuration at separator, expected: '---' between header and checks")
	}

	return bytes.TrimSpace(data[:loc[0]]), bytes.TrimSpace(data[loc[1]:]), nil
}

func SeparateConfigFrom(fileName string) ([]byte, []byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("failed to open config file for reading: '%s'", fileName))
	}

	return SeparateConfig(data)
}

func ModifyConditionStrings(cond aeaconf2.Condition, fun func(string) string) {
	val := reflect.ValueOf(cond).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.String {
			field.SetString(fun(field.String()))
		} else if field.Kind() == reflect.Interface || field.Kind() == reflect.Ptr {
			fieldInterface := field.Interface()
			if nestedCond, ok := fieldInterface.(aeaconf2.Condition); ok {
				ModifyConditionStrings(nestedCond, fun)
			}
		} else if field.Kind() == reflect.Struct {
			// Handle BaseCondition
			for j := 0; j < field.NumField(); j++ {
				nestedField := field.Field(j)
				if nestedField.Kind() == reflect.String {
					nestedField.SetString(fun(nestedField.String()))
				}
			}
		}
	}
}
