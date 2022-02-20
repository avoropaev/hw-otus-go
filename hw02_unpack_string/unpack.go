package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(param string) (string, error) {
	stringBuilder := strings.Builder{}

	var prev *rune
	var prevAfterBackslash bool
	lastIndex := utf8.RuneCountInString(param) - 1

	for i, current := range param {
		current := current

		if err := validate(prev, &current, prevAfterBackslash, i == lastIndex); err != nil {
			return "", err
		}

		if prev == nil {
			prev = &current
			continue
		}

		if *prev == '\\' && !prevAfterBackslash {
			prevAfterBackslash = true
			prev = &current
			continue
		}

		if !unicode.IsDigit(current) {
			if !unicode.IsDigit(*prev) || prevAfterBackslash {
				if err := write(&stringBuilder, *prev, 1, &prevAfterBackslash); err != nil {
					return "", err
				}
			}

			prev = &current
			continue
		}

		currentAsInt, _ := strconv.Atoi(string(current))
		if err := write(&stringBuilder, *prev, currentAsInt, &prevAfterBackslash); err != nil {
			return "", err
		}

		prev = &current
		continue
	}

	if prev != nil && (!unicode.IsDigit(*prev) || prevAfterBackslash) {
		if err := write(&stringBuilder, *prev, 1, &prevAfterBackslash); err != nil {
			return "", err
		}
	}

	return stringBuilder.String(), nil
}

func validate(prev, current *rune, prevAfterBackslash bool, last bool) error {
	currentIsDigit := unicode.IsDigit(*current)

	// если первый символ - цифра
	if currentIsDigit && prev == nil {
		return ErrInvalidString
	}

	// если стоят 2 цифры подряд и перед первой не было \
	if !prevAfterBackslash && currentIsDigit && unicode.IsDigit(*prev) {
		return ErrInvalidString
	}

	// если перед предыдущим символом не было \, предыдущий символ = \ и текущий символ != цифре или \
	if !prevAfterBackslash && prev != nil && *prev == '\\' && !currentIsDigit && *current != '\\' {
		return ErrInvalidString
	}

	// последний элемент = \ и перед ним нет \ или перед ним есть \, но он уже "использован"
	if last && *current == '\\' && (prev == nil || *prev != '\\' || prevAfterBackslash) {
		return ErrInvalidString
	}

	return nil
}

func write(builder *strings.Builder, symbol rune, count int, prevAfterBackslash *bool) error {
	for i := 0; i < count; i++ {
		if _, err := builder.Write([]byte(string(symbol))); err != nil {
			return err
		}

		*prevAfterBackslash = false
	}

	return nil
}
