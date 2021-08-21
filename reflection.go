package commons

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ErrNotBool = errors.New("not a boolean")

func copyStringField(inputFieldName string, inputField, outputField reflect.Value) error {
	// If input field is a string and output field is a number, parse input
	// field and assign.
	outputKind := outputField.Kind()
	if outputKind == reflect.Float32 || outputKind == reflect.Float64 {
		x, err := strconv.ParseFloat(inputField.String(), 64) //nolint:gomnd
		if err != nil {
			return err
		}

		outputField.SetFloat(x)
	} else if outputKind == reflect.Bool {
		value := strings.ToLower(inputField.String())
		switch value {
		case "false":
			outputField.SetBool(false)
		case "true":
			outputField.SetBool(true)
		default:
			return fmt.Errorf("field %s, value %s: %w", inputFieldName, value, ErrNotBool)
		}
	}

	return nil
}

func copyField(inputFieldName string, inputField, outputField reflect.Value) error {
	inputKind := inputField.Kind()
	outputKind := outputField.Kind()

	//nolint:exhaustive
	switch inputKind {
	case outputKind:
		outputField.Set(inputField)
	case reflect.Int64:
		// If input field is int64 and output is time.Time: convert time from
		// milliseconds and assign.
		if outputField.Type() == reflect.TypeOf(time.Time{}) {
			x := TimeFromMs(inputField.Int())
			outputField.Set(reflect.ValueOf(x))
		}
	case reflect.String:
		if err := copyStringField(inputFieldName, inputField, outputField); err != nil {
			return err
		}
	}

	return nil
}

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

		inputField := input.FieldByIndex(inputStructField.Index)
		outputField := output.FieldByName(inputStructField.Name)

		if err := copyField(inputStructField.Name, inputField, outputField); err != nil {
			return err
		}
	}

	return nil
}
