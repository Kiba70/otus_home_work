package hw02unpackstring

import (
	"errors"
	"log/slog"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrStrangeError  = errors.New("error executing Atoi")
	ErrAppendRune    = errors.New("error append Rune")
)

func Unpack(original string) (string, error) {
	var (
		prevRune  rune
		prevDigit bool
		screen    bool
		c         int
		err       error
		result    strings.Builder
	)

	if len(original) == 0 { // Пустая строка - корректное значение
		return "", nil
	}

	if original[0] >= '0' && original[0] <= '9' { // Руны 0-9 имеют размерность 1 байт, так делать можно
		return "", ErrInvalidString // Первый символ - цифра. Ошибка.
	}

	for _, r := range original {
		switch { // Последовательность case важна!
		case screen:
			if _, err = result.WriteRune(prevRune); err != nil {
				return "", ErrAppendRune
			}
			prevDigit = false
			screen = false
			prevRune = r
		case r == '\\':
			screen = true
		case r >= '0' && r <= '9' && prevDigit:
			return "", ErrInvalidString
		case r == '0':
			prevRune = 0
			prevDigit = true
		case r == '1':
			prevDigit = true
		case r >= '2' && r <= '9':
			prevDigit = true
			if c, err = strconv.Atoi(string(r)); err != nil { // Выполняем рекомендацию использовать strconv.Atoi
				// Странная ошибка, которой быть не должно
				return "", ErrStrangeError
			}
			//nolint: lll
			if _, err = result.WriteString(strings.Repeat(string(prevRune), c-1)); err != nil { // Выполняем рекомендацию использовать strings.Repeat
				return "", ErrAppendRune
			}
		case prevRune != 0:
			if _, err = result.WriteRune(prevRune); err != nil {
				return "", ErrAppendRune
			}
			fallthrough // Обязательно заполняем prevRune
		default:
			prevDigit = false
			prevRune = r
		}
	}
	if prevRune != 0 { // Если последний символ не шлак
		if _, err = result.WriteRune(prevRune); err != nil { // Последний символ
			return "", ErrAppendRune
		}
	}

	return result.String(), nil
}

// Вариант функции от моего сына.
// Он ради интереса и саморазвития делает "домашку" самостоятельно.
func UnpackFromMySon(text string) (string, error) {
	textRunes := []rune(text)
	if len(textRunes) < 1 { // empty input check
		return "", nil
	}

	numbers := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	if slices.Contains(numbers, textRunes[0]) {
		return "", ErrInvalidString
	}
	strBuilder := strings.Builder{}
	var escape, prtFlag, numFlag bool
	var output rune

	if textRunes[0] == '\\' {
		escape = true
	} else {
		output = textRunes[0]
		prtFlag = true
	}

	for _, symbol := range textRunes[1:] {
		if escape {
			switch {
			case symbol != '\\' && !slices.Contains(numbers, symbol):
				return "", ErrInvalidString
			default:
				output = symbol
				escape = false
				continue
			}
		}

		switch {
		case symbol == '\\':
			numFlag = false
			escape = true
			if !prtFlag {
				prtFlag = true
				continue
			}
			strBuilder.WriteRune(output)
		case slices.Contains(numbers, symbol):
			if numFlag {
				return "", ErrInvalidString
			}
			repeats, err := strconv.Atoi(string(symbol))
			if err != nil {
				slog.Error(err.Error())
			}
			strBuilder.WriteString(strings.Repeat(string(output), repeats))
			prtFlag = false
			numFlag = true
		default:
			numFlag = false
			if !prtFlag {
				prtFlag = true
				output = symbol
				continue
			}
			strBuilder.WriteRune(output)
			output = symbol
		}
	}

	switch { // handling last symbol (print it out if it wasn't, throw error if escaping end of string)
	case escape:
		return "", ErrInvalidString
	case prtFlag:
		strBuilder.WriteRune(output)
	}

	return strBuilder.String(), nil
}
