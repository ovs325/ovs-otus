package hw09structvalidator

import (
	"fmt"
	rf "reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errors []string
	for _, err := range v {
		errors = append(errors, fmt.Sprintf("%s: %s", err.Field, err.Err.Error()))
	}
	return strings.Join(errors, ", ")
}

func Validate(v interface{}) error {
	val := rf.ValueOf(v)
	if val.Kind() != rf.Struct {
		return fmt.Errorf("ожидалась структура, получена %s", val.Kind())
	}
	if val.Kind() == rf.Ptr {
		val = val.Elem()
	}

	var errors ValidationErrors
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		var text string
		if text = field.Tag.Get("validate"); text == "" {
			continue
		}

		validatorsList := strings.Split(text, "|")
		for _, item := range validatorsList {
			parts := strings.SplitN(item, ":", 2)
			if len(parts) != 2 {
				appErr(&errors, field.Name, "плохой валидатор", item)
				continue
			}
			validType, validValue := parts[0], parts[1]
			switch validType {
			case "len":
				len_(validValue, value, field, &errors)
			case "min":
				min_(validValue, value, field, &errors)
			case "max":
				max_(validValue, value, field, &errors)
			case "regexp":
				regexp_(validValue, value, field, &errors)
			case "in":
				in_(validValue, value, field, &errors)
			default:
				appErr(&errors, field.Name, "неизвестный валидатор", validType)
			}
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func len_(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) {
	length, _ := strconv.Atoi(validValue)
	switch value.Kind() {
	case rf.String:
		if value.Len() != length {
			appErr(errors, field.Name, "длинна должна быть", length)
		}
	case rf.Slice:
		if value.Len() != length {
			appErr(errors, field.Name, "длинна должна быть", length)
		}
	default:
		appErr(errors, field.Name, "неправильный тип 'len' валидатора", value.Kind())
	}
}

func min_(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) {
	min, _ := strconv.Atoi(validValue)
	switch value.Kind() {
	case rf.Int:
		if value.Int() < int64(min) {
			appErr(errors, field.Name, "должно быть хотя бы", min)
		}
	default:
		appErr(errors, field.Name, "неправильный тип 'min' валидатора", value.Kind())
	}
}

func max_(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) {
	max, _ := strconv.Atoi(validValue)
	switch value.Kind() {
	case rf.Int:
		if value.Int() > int64(max) {
			appErr(errors, field.Name, "должно быть, не более", max)
		}
	default:
		appErr(errors, field.Name, "неправильный тип 'max' валидатора", value.Kind())
	}
}

func regexp_(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) {
	regex, err := regexp.Compile(validValue)
	if err != nil {
		appErr(errors, field.Name, "недопустимое регулярное выражение", err.Error())
		return
	}
	switch value.Kind() {
	case rf.String:
		if !regex.MatchString(value.String()) {
			appErr(errors, field.Name, "должно соответствовать регулярному выражению", validValue)
		}
	default:
		appErr(errors, field.Name, "неправильный тип 'regexp' валидатора", value.Kind())
	}
}

func in_(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) {
	options := strings.Split(validValue, ",")
	switch value.Kind() {
	case rf.String:
		found := false
		for _, option := range options {
			if value.String() == option {
				found = true
				break
			}
		}
		if !found {
			appErr(errors, field.Name, "должен быть одним из", validValue)
		}
	case rf.Int:
		found := false
		for _, option := range options {
			if value.Int() == int64(mustAtoi(option)) {
				found = true
				break
			}
		}
		if !found {
			appErr(errors, field.Name, "должен быть одним из", validValue)
		}
	default:
		appErr(errors, field.Name, "неправильный тип 'in' валидатора", value.Kind())
	}
}

func appErr(errors *ValidationErrors, name, msg string, val any) {
	*errors = append(*errors, ValidationError{
		Field: name,
		Err:   fmt.Errorf("%s: %v", msg, val),
	})
}
