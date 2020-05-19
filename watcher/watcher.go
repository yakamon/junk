package main

import "fmt"

func main() {
	for {
		fmt.Printf("Please input a text: ")
		var input string
		fmt.Scanln(&input)
		fmt.Printf("You input a text: %s\n", input)
	}
}
