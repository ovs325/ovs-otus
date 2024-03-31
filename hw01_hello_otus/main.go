package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	str := "Hello, OTUS!"
	reversed := stringutil.Reverse(str)
	fmt.Println(reversed)
}
