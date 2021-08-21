package commons

import (
	"reflect"
	"strconv"
)

func SmartCopy(dst, src interface{}) error {
	input := reflect.ValueOf(src).Elem()
	output := reflect.ValueOf(dst).Elem()

	hasField := func(name string) bool {
		for _, field := range reflect.VisibleFields(input.Type()) {
			if field.Name == name {
				return true
			}
		}

		return false
	}

	for _, inputStructField := range reflect.VisibleFields(input.Type()) {
		if !hasField(inputStructField.Name) {
			continue
		}

		outputField := output.FieldByName(inputStructField.Name)
		inputField := input.FieldByIndex(inputStructField.Index)
		inputKind := inputField.Kind()

		if outputField.Kind() == inputKind {
			outputField.Set(inputField)
		} else if inputKind == reflect.String {
			kind := outputField.Kind()
			if kind == reflect.Float32 || kind == reflect.Float64 {
				x, err := strconv.ParseFloat(inputField.String(), 64) //nolint:gomnd
				if err != nil {
					return err
				}

				outputField.SetFloat(x)
			}
		}
	}

	return nil
}
