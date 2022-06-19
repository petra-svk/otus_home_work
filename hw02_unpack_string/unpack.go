package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
	"unicode/utf8"
)

const BackslashSymbol rune = 92

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}
	if !utf8.ValidString(str) {
		return "", ErrInvalidString
	}
	var res []rune
	prevDigit := false
	backslash := false
	for i, r := range str {
		switch {
		case unicode.IsSpace(r):
			res = append(res, r)
		case unicode.IsLetter(r):
			if backslash {
				return "", ErrInvalidString
			}
			prevDigit = false
			res = append(res, r)

		case unicode.IsDigit(r):
			if i == 0 || prevDigit {
				return "", ErrInvalidString
			}

			if backslash {
				backslash = false
				res = append(res, r)
				continue
			}

			prevDigit = true
			x, err := strconv.Atoi(string(r))
			if err != nil {
				return "", ErrInvalidString
			}
			if x == 0 {
				res = res[:len(res)-1]
			} else {
				x--
				lastChar := res[len(res)-1]
				for j := 0; j < x; j++ {
					res = append(res, lastChar)
				}
			}
		case r == BackslashSymbol:
			prevDigit = false
			if backslash {
				backslash = false
				res = append(res, r)
			} else {
				backslash = true
			}
		default:
			return "", ErrInvalidString
		}
	}
	return string(res), nil
}
