package main

import (
	"errors"
	"fmt"
)

var ErrStringInvalid = errors.New("string is invalid")
var ErrStringEmpty = errors.New("string is empty")

func Unpack(input string) (string, error) {
	// Place your code here.
	if len(input) == 0 {
		return "", ErrStringEmpty
	}
	return "", nil
}

func main() {
	_, e := Unpack("")
	fmt.Println(e.Error())

}
