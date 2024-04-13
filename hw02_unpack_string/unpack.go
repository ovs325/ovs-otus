package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	// На входе пустая строка
	if str == "" {
		return "", nil
	}
	// последовательность рун
	runes := []rune(str)
	// первая руна - число
	if unicode.IsNumber(runes[0]) {
		return "", ErrInvalidString
	}
	builder := strings.Builder{}
	// перебор последовательности рун парами
	for i := 0; i < len(runes)-1; i++ {
		// Пара число-число
		if unicode.IsNumber(runes[i]) && unicode.IsNumber(runes[i+1]) {
			return "", ErrInvalidString
		}
		// Пара чило-любая руна
		if unicode.IsNumber(runes[i]) {
			continue
		}
		// Пара символ-число
		if !unicode.IsNumber(runes[i]) && unicode.IsNumber(runes[i+1]) {
			num, err := strconv.Atoi(string(runes[i+1]))
			if err != nil {
				return "", ErrInvalidString
			}
			builder.WriteString(strings.Repeat(string(runes[i]), num))
		}
		// Пара симво-символ
		if !unicode.IsNumber(runes[i]) && !unicode.IsNumber(runes[i+1]) {
			builder.WriteString(string(runes[i]))
		}
	}
	// Последняя руна - символ
	if !unicode.IsNumber(runes[len(runes)-1]) {
		builder.WriteString(string(runes[len(runes)-1]))
	}

	return builder.String(), nil
}
