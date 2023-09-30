package main

import (
	"fmt"

	"github.com/GoesToEleven/GolangTraining/02_package/stringutil"
)

func reverser(in string) string {
	return stringutil.Reverse(in)
}

func printer(in string) {
	fmt.Println(in)
}

func main() {
	printer(reverser("Hello, OTUS!"))
}
