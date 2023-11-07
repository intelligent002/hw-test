package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidString          = errors.New("string is invalid")
	ErrUnsupportedCharacters  = errors.New("unsupported characters detected")
	ErrFailedToConvertInteger = errors.New("failed to convert rune to integer")
)

type RuneKind int

const (
	Unsupported RuneKind = iota
	Letter
	Digit
	Backslash
)

type character struct {
	letter  rune
	width   int
	escaped bool
	kind    RuneKind
}

func (char *character) detectKind() {
	str := string(char.letter)
	switch {
	case "0" <= str && str <= "9":
		{
			// arabic-ascii digits will go as digits, the rest (https://pkg.go.dev/unicode#Adlam) will not
			char.kind = Digit
			break
		}
	case str == "\\":
		{
			// backslashes are special animal
			char.kind = Backslash
			break
		}
	default:
		{
			// all other characters: include punctuation, emoticons etc
			char.kind = Letter
			break
		}
	}
}

func getChar(input string, offset int) character {
	var char character
	char.letter, char.width = utf8.DecodeRuneInString(input[offset:])
	char.escaped = false

	if string(char.letter) == "\\" {
		// it is an escaped char
		var widthNext int
		char.letter, widthNext = utf8.DecodeRuneInString(input[offset+char.width:])
		char.width += widthNext
		char.escaped = true
	}

	// detect kind only once
	char.detectKind()

	return char
}
func handleDigit(charCurr, charNext character, result string, index int) (string, int, error) {
	switch {
	case !charCurr.escaped:
		{
			// Curr is a non escaped digit, bad string
			return "", index, ErrInvalidString
		}
	case charNext.kind == Digit && !charNext.escaped:
		{
			// Curr is an escaped digit, next is a non escaped digit - multiply Curr by digit
			number, err := strconv.Atoi(string(charNext.letter))
			if err != nil {
				// something went wrong during ATOI
				return "", index, ErrFailedToConvertInteger
			}
			result += strings.Repeat(string(charCurr.letter), number)
			// and skip the upcoming digit
			index += charNext.width
		}
	default:
		{
			// in any other case - just add the digit (stripping the escaping)
			result += string(charCurr.letter)
		}
	}
	return result, index, nil
}
func handleLetter(charCurr, charNext character, result string, index int) (string, int, error) {
	switch {
	case charCurr.escaped:
		{
			// Curr is an escaped letter, bad string
			return "", index, ErrInvalidString
		}
	case charNext.kind == Digit && !charNext.escaped:
		{
			// Curr is a non escaped letter, next is a non escaped digit - multiply Curr by digit
			number, err := strconv.Atoi(string(charNext.letter))
			if err != nil {
				// Curr is a non escaped digit, bad string
				return "", index, ErrFailedToConvertInteger
			}
			result += strings.Repeat(string(charCurr.letter), number)
			// and skip the upcoming digit
			index += charNext.width
		}
	default:
		{
			// in any other case - just add the letter
			result += string(charCurr.letter)
		}
	}
	return result, index, nil
}
func handleBackslash(charCurr, charNext character, result string, index int) (string, int, error) {
	switch {
	case charNext.kind == Digit && !charNext.escaped:
		{
			// Curr is an escaped backslash, next is a non escaped digit - multiply Curr by digit
			number, err := strconv.Atoi(string(charNext.letter))
			if err != nil {
				// Curr is a non escaped digit, bad string
				return "", index, ErrFailedToConvertInteger
			}
			result += strings.Repeat(string(charCurr.letter), number)
			// and skip the upcoming digit
			index += charNext.width
		}
	default:
		{
			// in any other case - just add the backslash (stripping escaping)
			result += string(charCurr.letter)
		}
	}
	return result, index, nil
}
func Unpack(input string) (string, error) {
	// init result
	var result string
	// init error just for case :)
	var err error
	// init iteration runes
	var charCurr, charNext character

	// iterate over the string by runes
	for index := 0; index < len(input); index += charCurr.width {
		// get rune Curr
		charCurr = getChar(input, index)
		// get rune next
		charNext = getChar(input, index+charCurr.width)
		// decide on operation
		switch charCurr.kind {
		case Digit:
			{
				result, index, err = handleDigit(charCurr, charNext, result, index)
				if err != nil {
					return "", err
				}
			}
		case Letter:
			{
				result, index, err = handleLetter(charCurr, charNext, result, index)
				if err != nil {
					return "", err
				}
			}
		case Backslash:
			{
				result, index, err = handleBackslash(charCurr, charNext, result, index)
				if err != nil {
					return "", err
				}
			}
		case Unsupported:
			{
				// Curr is not a slash/digit/letter, bad string (unreachable)
				return "", ErrUnsupportedCharacters
			}
		}
	}
	return result, nil
}
