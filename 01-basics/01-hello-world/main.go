package main

import (
	"fmt"
	"time"
)

func main() {
	// Basic greeting
	fmt.Println("Hello, World!")
	fmt.Println("Welcome to Go!")

	// Print with formatting
	fmt.Printf("Current time: %v\n", time.Now().Format("2006-01-02 15:04:05"))

	// Print without newline
	fmt.Print("Go ")
	fmt.Print("is ")
	fmt.Println("awesome!")

	// Multiple values
	fmt.Println("Numbers:", 1, 2, 3)

	// Demonstrating format specifiers
	name := "Gopher"
	version := 1.21
	fmt.Printf("Welcome to Go %v, %s!\n", version, name)

	// Using Sprintf to create strings
	greeting := fmt.Sprintf("Hello from %s version %v", name, version)
	fmt.Println(greeting)
}
