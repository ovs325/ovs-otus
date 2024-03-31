package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	str := "Hello, World!"
	reversed := stringutil.Reverse(str)
	fmt.Println(reversed)
}
