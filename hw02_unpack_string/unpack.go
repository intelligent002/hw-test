package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("string is invalid")
var ErrEmptyString = errors.New("string is empty")

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
	for index, widthCurrent = 0, 0; index < len(input); index += widthCurrent {
		// get rune current
		runeCurrent, widthCurrent, escapedCurrent = getChar(input, index)
		// get rune next
		runeNext, widthNext, escapedNext = getChar(input, index+widthCurrent)
		// decide on operation
		switch {
		case widthNext == 0:
			{
				// last rune - just add
				result += string(runeCurrent)
			}
		case unicode.IsDigit(runeCurrent):
			{
				if !escapedCurrent {
					// current is a non escaped digit, bad string
					return "", ErrInvalidString
				}
				if unicode.IsLetter(runeNext) {
					// current is an escaped digit, next is a letter - add current
					result += string(runeCurrent)
				}
				if unicode.IsDigit(runeNext) {
					if escapedNext {
						// current is an escaped digit, next is an escaped digit - add current
						result += string(runeCurrent)
					}
					if !escapedNext {
						// current is an escaped digit, next is a non escaped digit - multiply current by digit
						result += strings.Repeat(string(runeCurrent), int(runeNext-'0'))
						// and skip the upcoming digit
						index += widthNext
					}
				}
			}
		case unicode.IsLetter(runeCurrent):
			{
				switch {
				case escapedCurrent:
					{
						// current is an escaped letter, bad string
						return "", ErrInvalidString
					}
				case unicode.IsLetter(runeNext):
					{
						// current is a letter, next is letter - add current
						result += string(runeCurrent)
					}
				case string(runeNext) == "\\":
					{
						// current is a letter, next is an escaped backslash - add current
						result += string(runeCurrent)
					}
				case unicode.IsDigit(runeNext) && escapedNext:
					{
						// current is a letter, next is an escaped digit - add current
						result += string(runeCurrent)
					}
				case unicode.IsDigit(runeNext) && !escapedNext:
					{
						// current is a letter, next is a non escaped digit - multiply current by digit
						result += strings.Repeat(string(runeCurrent), int(runeNext-'0'))
						// and skip the upcoming digit
						index += widthNext
					}
				}
			}
		case string(runeCurrent) == "\\":
			{
				if unicode.IsLetter(runeNext) {
					// current is an escaped backslash, next is any letter - add current
					result += string(runeCurrent)
				}
				if unicode.IsDigit(runeNext) {
					if escapedNext {
						// current is an escaped backslash, next is an escaped digit - add current
						result += string(runeCurrent)
					}
					if !escapedNext {
						// current is an escaped backslash, next is a non escaped digit - multiply current by digit
						result += strings.Repeat(string(runeCurrent), int(runeNext-'0'))
						// and skip the upcoming digit
						index += widthNext
					}
				}
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
