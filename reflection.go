package commons

import (
	"reflect"
	"strconv"
)

// ParseFields expects two pointers to structs.  It parses all string
// fields from src and assigns them to dst.  Fields are matched by
// name.  ParseFields returns the number of fields successfully read,
// parsed and assigned.
func ParseFields(dst, src interface{}) int {
	input := reflect.ValueOf(src).Elem()
	inputFields := reflect.VisibleFields(input.Type())

	output := reflect.ValueOf(dst).Elem()

	match := func(field reflect.StructField) bool {
		if field.Type.Kind() != reflect.String {
			return false
		}

		outputField := output.FieldByName(field.Name)

		if outputField.IsZero() {
			return false
		}

		return outputField.Type().Kind() == reflect.Float64
	}

	ans := 0

	for _, inputField := range inputFields {
		if match(inputField) {
			stringValue := input.FieldByName(inputField.Name).String()

			floatValue, err := strconv.ParseFloat(stringValue, 64) //nolint:gomnd
			if err != nil {
				continue
			}

			outputField := output.FieldByName(inputField.Name)
			if outputField.IsZero() {
				panic("DEBA")
			}

			outputField.SetFloat(floatValue)
			ans++
		}
	}

	return ans
}
