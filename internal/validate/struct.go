package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Validate that the values in the struct are valid according to the validation rules.
// If validation fails the error will be of type Error.
func Validate(req Request) error {
	value := reflect.ValueOf(req)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("request argument must be pointer to struct, not %T", req)
	}
	value = value.Elem()

	fieldErrors := make(map[string][]string)

	for _, rule := range req.ValidationRules() {
		fv := rule.fieldValue()
		if fv.Kind() != reflect.Ptr {
			return fmt.Errorf("field value must be a pointer, not %v", fv.Kind())
		}
		ft := findStructField(value, fv)
		jsonName := jsonNameFromStructField(ft)
		if jsonName == "" {
			continue
		}
		if err := rule.validate(); err != nil {
			// TODO: test nested errors
			var nestedError Error
			if errors.As(err, &nestedError) {
				for k, v := range nestedError.FieldErrors {
					fieldErrors[jsonName+"."+k] = v
				}
				continue
			}
			fieldErrors[jsonName] = append(fieldErrors[jsonName], err.Error())
		}
	}

	if len(fieldErrors) > 0 {
		return Error{FieldErrors: fieldErrors}
	}
	return nil
}

// RulesToMap translates creates a map of validation rules, where the key is
// the struct field name (not the JSON field name), that the rules apply to.
// TODO: maybe accept the reflect.Value instead.
func RulesToMap(req Request) (map[string][]ValidationRule, error) {
	value := reflect.ValueOf(req)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("request argument must be pointer to struct, not %T", req)
	}
	value = value.Elem()

	result := make(map[string][]ValidationRule)
	for _, rule := range req.ValidationRules() {
		fv := rule.fieldValue()
		if fv.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("field value must be a pointer, not %v", fv.Kind())
		}
		ft := findStructField(value, fv)
		result[ft.Name] = append(result[ft.Name], rule)
	}
	return result, nil
}

// findStructField looks for a field by pointer address in the structValue.
// If found, the field info will be returned. Otherwise, nil will be returned.
// Copied from ozzo-validation.
// TODO: copy test cases
func findStructField(structValue reflect.Value, fieldValue reflect.Value) *reflect.StructField {
	ptr := fieldValue.Pointer()
	for i := structValue.NumField() - 1; i >= 0; i-- {
		sf := structValue.Type().Field(i)
		if ptr == structValue.Field(i).UnsafeAddr() {
			// do additional type comparison because it's possible that the address of
			// an embedded struct is the same as the first field of the embedded struct
			if sf.Type == fieldValue.Elem().Type() {
				return &sf
			}
		}
		if sf.Anonymous {
			// delve into anonymous struct to look for the field
			fi := structValue.Field(i)
			if sf.Type.Kind() == reflect.Ptr {
				fi = fi.Elem()
			}
			if fi.Kind() == reflect.Struct {
				if f := findStructField(fi, fieldValue); f != nil {
					return f
				}
			}
		}
	}
	return nil
}

// TODO: test cases (name in tag, empty name in tag, no tag, "-" tag)
func jsonNameFromStructField(f *reflect.StructField) string {
	tag := f.Tag.Get("json")
	name, _, _ := strings.Cut(tag, ",")
	switch name {
	case "":
		return f.Name
	case "-":
		return ""
	}
	return name
}
