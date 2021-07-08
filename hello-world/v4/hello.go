package main

import "fmt"

const englishHelloPrefix = "Hello, "

// Hello враћа персонализовани поздрав.
func Hello(name string) string {
	return englishHelloPrefix + name
}

func main() {
	fmt.Println(Hello("world"))
}
