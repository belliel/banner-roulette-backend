package utils

import (
	"reflect"
)

func ToMSI(s interface{}, tagName string, ignore []string) map[string]interface{} {
	ignoreListO1 := make(map[string]interface{})

	for _, ignoreItem := range ignore { ignoreListO1[ignoreItem] = "" }

	result := make(map[string]interface{})

	val := reflect.ValueOf(s).Elem()

	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get(tagName)
		if _, exist := ignoreListO1[tag]; !exist {
			result[tag] = val.Field(i)
		}
	}

	return result
}
