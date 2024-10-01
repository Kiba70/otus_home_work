package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	// Используем функцию reverse.String согласно требованиям ДЗ
	fmt.Println(reverse.String("Hello, OTUS!"))
}
