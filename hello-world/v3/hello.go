package main

import "fmt"

// Hello враћа персонализовани поздрав.
func Hello(name string) string {
	return "Hello, " + name
}

func main() {
	fmt.Println(Hello("world"))
}
