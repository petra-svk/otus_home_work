package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	var res []rune
	prevDigit := false
	backslash := false
	for i, r := range str {
		switch {
		case unicode.IsLetter(r):
			if backslash {
				return "", ErrInvalidString
			}
			prevDigit = false
			res = append(res, r)

		case unicode.IsDigit(r):
			if i == 0 {
				return "", ErrInvalidString
			}

			if prevDigit {
				return "", ErrInvalidString
			}

			if backslash {
				backslash = false
				res = append(res, r)
				continue
			}

			prevDigit = true
			x, _ := strconv.Atoi(string(r))
			if x == 0 {
				res = res[:len(res)-1]
			} else {
				x--
				lastChar := res[len(res)-1]
				for j := 0; j < x; j++ {
					res = append(res, lastChar)
				}
			}

		case r == 92: // '\' = 92
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
