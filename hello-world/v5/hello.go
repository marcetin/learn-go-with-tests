package main

import "fmt"

const englishHelloPrefix = "Hello, "

// Hello враћа персонализовани поздрав, подразумевано је "Hello, world" ако се проследи празно име.
func Hello(name string) string {
	if name == "" {
		name = "World"
	}
	return englishHelloPrefix + name
}

func main() {
	fmt.Println(Hello("world"))
}
