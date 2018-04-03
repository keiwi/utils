package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func toUpper(name string) string {
	return strings.ToUpper(name[0:1]) + name[1:]
}

func fixName(name string) string {
	var n string
	for _, s := range strings.Split(name, "_") {
		n += toUpper(s)
	}
	return n
}

func fieldByName(v reflect.Value, name string) reflect.Value {
	max := v.NumField()
	for i := 0; i < max; i++ {
		field := v.Type().Field(i)
		json := field.Tag.Get("json")
		if json != "" && json != "-" {
			if json == name {
				return v.Field(i)
			} else {
				continue
			}
		}
		if field.Name == fixName(name) {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := fieldByName(structValue, name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		return nil
	}

	if structFieldType != val.Type() {
		return errors.New("provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func MapToStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
