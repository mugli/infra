package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ValidateStruct validates that the values in the struct are valid according
// to the passed validation rules.
// If validation fails, the error will be of type Error.
// TODO: maybe un-export
func ValidateStruct(structPtr interface{}, rules ...ValidationRule) error {
	value := reflect.ValueOf(structPtr)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("struct argument must be pointer to struct, not %T", structPtr)
	}
	value = value.Elem()

	fieldErrors := make(map[string][]string)

	for _, rule := range rules {
		fv := rule.fieldValue()
		if fv.Kind() != reflect.Ptr {
			return fmt.Errorf("field value must be a pointer, not %v", fv.Kind())
		}
		ft := findStructField(value, fv)
		jsonName := fieldNameFromStructField(ft)
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
func fieldNameFromStructField(f *reflect.StructField) string {
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
