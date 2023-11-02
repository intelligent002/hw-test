package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidString = errors.New("string is invalid")
	ErrEmptyString   = errors.New("string is empty")
)

func getChar(input string, offset int) (rune, int, bool) {
	char, widthCurrent := utf8.DecodeRuneInString(input[offset:])

	if string(char) == "\\" {
		var widthNext int
		char, widthNext = utf8.DecodeRuneInString(input[offset+widthCurrent:])
		return char, widthCurrent + widthNext, true
	}

	return char, widthCurrent, false
}

func Unpack(input string) (string, error) {
	// init result
	var result string
	var runeCurrent, runeNext rune
	var index, widthCurrent, widthNext int
	var escapedCurrent, escapedNext bool

	// validate empty input
	if len(input) == 0 {
		return result, ErrEmptyString
	}

	// iterate over the string by runes
	for index = 0; index < len(input); index += widthCurrent {
		// get rune current
		runeCurrent, widthCurrent, escapedCurrent = getChar(input, index)

		// get rune next
		runeNext, widthNext, escapedNext = getChar(input, index+widthCurrent)

		// decide on operation
		switch {
		case unicode.IsDigit(runeCurrent):
			{
				if !escapedCurrent {
					// current is a non escaped digit, bad string
					return "", ErrInvalidString
				}
				if unicode.IsDigit(runeNext) && !escapedNext {
					// current is an escaped digit, next is a non escaped digit - multiply current by digit
					result += strings.Repeat(string(runeCurrent), int(runeNext-'0'))
					// and skip the upcoming digit
					index += widthNext
					continue
				}
				// in any other case - just add the digit (stripping the escaping)
				result += string(runeCurrent)
				continue
			}
		case unicode.IsLetter(runeCurrent):
			{
				if escapedCurrent {
					// current is an escaped letter, bad string
					return "", ErrInvalidString
				}
				if unicode.IsDigit(runeNext) && !escapedNext {
					// current is a non escaped letter, next is a non escaped digit - multiply current by digit
					result += strings.Repeat(string(runeCurrent), int(runeNext-'0'))
					// and skip the upcoming digit
					index += widthNext
					continue
				}
				// in any other case - just add the letter
				result += string(runeCurrent)
				continue
			}
		case string(runeCurrent) == "\\":
			{
				if unicode.IsDigit(runeNext) && !escapedNext {
					// current is an escaped backslash, next is a non escaped digit - multiply current by digit
					result += strings.Repeat(string(runeCurrent), int(runeNext-'0'))
					// and skip the upcoming digit
					index += widthNext
					continue
				}
				// in any other case - just add the backslash (stripping escaping)
				result += string(runeCurrent)
				continue
			}
		}
	}
	return result, nil
}

func main() {
	result, e := Unpack("a4")
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println("result is: ", result)
	}
}
