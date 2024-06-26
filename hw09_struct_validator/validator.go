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
	errs := make([]string, 0, len(v))
	for _, err := range v {
		errs = append(errs, fmt.Sprintf("%s: %s", err.Field, err.Err.Error()))
	}
	return strings.Join(errs, ", ")
}

func Validate(v interface{}) error {
	val := rf.ValueOf(v)
	if val.Kind() != rf.Struct {
		return fmt.Errorf("ожидалась структура, получена %v", val.Kind())
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
		if err := verification(text, value, field, &errors); err != nil {
			return fmt.Errorf("процесс проверки завершился ошибкой: %w", err)
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func verification(text string, value rf.Value, field rf.StructField, errors *ValidationErrors) error {
	validatorsList := strings.Split(text, "|")
	for _, item := range validatorsList {
		parts := strings.SplitN(item, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("инвалидный валидатор: %s", item)
		}
		validType, validValue := parts[0], parts[1]
		switch validType {
		case "len":
			if err := lenV(validValue, value, field, errors); err != nil {
				return err
			}
		case "min":
			if err := minV(validValue, value, field, errors); err != nil {
				return err
			}
		case "max":
			if err := maxV(validValue, value, field, errors); err != nil {
				return err
			}
		case "regexp":
			if err := regexpV(validValue, value, field, errors); err != nil {
				return err
			}
		case "in":
			if err := inV(validValue, value, field, errors); err != nil {
				return err
			}
		default:
			return fmt.Errorf("неизвестный валидатор: %v", validType)
		}
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

func lenV(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) error {
	length, _ := strconv.Atoi(validValue)
	vs := value.String()
	vl := value.Len()
	_ = vs
	switch value.Kind() { // nolint:exhaustive,nolintlint
	case rf.String:
		if vl != length {
			AppErr(errors, field.Name, "длинна должна быть = ", length)
		}
	case rf.Slice:
		for i := 0; i < value.Len(); i++ {
			if value.Index(i).Len() != length {
				msg := fmt.Sprintf("длина элемента %d должна быть = ", i)
				AppErr(errors, field.Name, msg, length)
			}
		}
	default:
		return fmt.Errorf("неправильный тип 'min' валидатора: %v", value.Kind())
	}
	return nil
}

func minV(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) error {
	min, _ := strconv.Atoi(validValue)
	switch value.Kind() { // nolint:exhaustive,nolintlint
	case rf.Int:
		if value.Int() < int64(min) {
			AppErr(errors, field.Name, "должно быть не менее ", min)
		}
	case rf.Slice:
		for i := 0; i < value.Len(); i++ {
			if value.Index(i).Int() < int64(min) {
				AppErr(errors, field.Name, fmt.Sprintf("элемент %d должен быть не менее ", i), min)
			}
		}
	default:
		return fmt.Errorf("неправильный тип 'min' валидатора: %v", value.Kind())
	}
	return nil
}

func maxV(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) error {
	max, _ := strconv.Atoi(validValue)
	switch value.Kind() { // nolint:exhaustive,nolintlint
	case rf.Int:
		if value.Int() > int64(max) {
			AppErr(errors, field.Name, "должно быть, не более ", max)
		}
	case rf.Slice:
		for i := 0; i < value.Len(); i++ {
			if value.Index(i).Int() > int64(max) {
				AppErr(errors, field.Name, fmt.Sprintf("элемент %d должен быть, не более ", i), max)
			}
		}
	default:
		return fmt.Errorf("неправильный тип 'max' валидатора: %v", value.Kind())
	}
	return nil
}

func regexpV(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) error {
	regex, err := regexp.Compile(validValue)
	if err != nil {
		return fmt.Errorf("недопустимое регулярное выражение: %w", err)
	}
	switch value.Kind() { // nolint:exhaustive,nolintlint
	case rf.String:
		if !regex.MatchString(value.String()) {
			AppErr(errors, field.Name, "не соответствует регулярному выражению ", validValue)
		}
	case rf.Slice:
		for i := 0; i < value.Len(); i++ {
			if !regex.MatchString(value.Index(i).String()) {
				AppErr(errors, field.Name, fmt.Sprintf("элемент %d не соответствует регулярному выражению ", i), validValue)
			}
		}
	default:
		return fmt.Errorf("неправильный тип 'regexp' валидатора: %v", value.Kind())
	}
	return nil
}

func inV(validValue string, value rf.Value, field rf.StructField, errors *ValidationErrors) error {
	options := strings.Split(validValue, ",")
	switch value.Kind() { // nolint:exhaustive,nolintlint
	case rf.String:
		found := false
		for _, option := range options {
			if value.String() == option {
				found = true
				break
			}
		}
		if !found {
			AppErr(errors, field.Name, "должен быть из множества ", validValue)
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
			AppErr(errors, field.Name, "должен быть из множества ", validValue)
		}
	case rf.Slice:
		if err := verifySlice(value, options, field, validValue, errors); err != nil {
			return err
		}
	default:
		return fmt.Errorf("неправильный тип 'in' валидатора: %v", value.Kind())
	}
	return nil
}

func verifySlice(
	value rf.Value,
	options []string,
	field rf.StructField,
	validValue string,
	errors *ValidationErrors,
) error {
	typ := value.Type().Elem().Kind()
	switch typ { // nolint:exhaustive,nolintlint
	case rf.String:
		for i := 0; i < value.Len(); i++ {
			found := false
			for _, option := range options {
				if value.Index(i).String() == option {
					found = true
					break
				}
			}
			if !found {
				AppErr(errors, field.Name, fmt.Sprintf("элемент %d должен быть из множества ", i), validValue)
			}
		}
	case rf.Int:
		for i := 0; i < value.Len(); i++ {
			found := false
			for _, option := range options {
				if value.Index(i).Int() == int64(mustAtoi(option)) {
					found = true
					break
				}
			}
			if !found {
				AppErr(errors, field.Name, fmt.Sprintf("элемент %d должен быть из множества ", i), validValue)
			}
		}
	default:
		return fmt.Errorf("неправильный тип элементов слайса для 'in' валидатора: %v", typ)
	}
	return nil
}

func AppErr(errors *ValidationErrors, name, msg string, val any) {
	*errors = append(*errors, ValidationError{
		Field: name,
		Err:   fmt.Errorf("%s %v", msg, val),
	})
}
