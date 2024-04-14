package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" { // На входе пустая строка
		return "", nil
	}
	runes := []rune(str)            // последовательность рун
	if unicode.IsNumber(runes[0]) { // первая руна - число
		return "", ErrInvalidString
	}
	builder := strings.Builder{}
	for i := 0; i < len(runes)-1; i++ { // перебор рун парами
		if unicode.IsNumber(runes[i]) && unicode.IsNumber(runes[i+1]) { // Пара число-число
			return "", ErrInvalidString
		}
		if unicode.IsNumber(runes[i]) { // Пара чило-любая руна
			continue
		}
		if !unicode.IsNumber(runes[i]) && unicode.IsNumber(runes[i+1]) { // Пара символ-число
			num, err := strconv.Atoi(string(runes[i+1]))
			if err != nil {
				return "", ErrInvalidString
			}
			builder.WriteString(strings.Repeat(string(runes[i]), num))
		}
		if !unicode.IsNumber(runes[i]) && !unicode.IsNumber(runes[i+1]) { // Пара симво-символ
			builder.WriteString(string(runes[i]))
		}
	}
	if !unicode.IsNumber(runes[len(runes)-1]) { // Последняя руна - символ
		builder.WriteString(string(runes[len(runes)-1]))
	}
	return builder.String(), nil
}
