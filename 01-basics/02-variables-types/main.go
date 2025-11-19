package main

import (
	"fmt"
	"strconv"
)

func main() {
	// === VARIABLE DECLARATION ===
	fmt.Println("=== Variable Declaration ===")

	// Short declaration (most idiomatic)
	name := "Alice"
	age := 28
	score := 95.5

	fmt.Printf("Name: %s, Age: %d, Score: %f\n", name, age, score)

	// === PRIMITIVE TYPES ===
	fmt.Println("\n=== Primitive Types ===")

	var (
		intVal    int    = 42
		floatVal  float64 = 3.14159
		boolVal   bool   = true
		stringVal string = "Hello, Go!"
		byteVal   byte   = 255
		runeVal   rune   = 'A' // Unicode character
	)

	fmt.Printf("int: %d\n", intVal)
	fmt.Printf("float64: %f\n", floatVal)
	fmt.Printf("bool: %v\n", boolVal)
	fmt.Printf("string: %s\n", stringVal)
	fmt.Printf("byte: %d\n", byteVal)
	fmt.Printf("rune: %c (code: %d)\n", runeVal, runeVal)

	// === ZERO VALUES ===
	fmt.Println("\n=== Zero Values ===")
	var zeroInt int
	var zeroFloat float64
	var zeroString string
	var zeroBool bool

	fmt.Printf("Zero int: %d\n", zeroInt)
	fmt.Printf("Zero float: %f\n", zeroFloat)
	fmt.Printf("Zero string: '%s'\n", zeroString)
	fmt.Printf("Zero bool: %v\n", zeroBool)

	// === ARRAYS ===
	fmt.Println("\n=== Arrays ===")

	var arr [5]int
	arr[0] = 10
	arr[1] = 20
	arr[2] = 30

	fmt.Printf("Array: %v\n", arr)
	fmt.Printf("Length: %d\n", len(arr))

	// Array literal
	primes := [...]int{2, 3, 5, 7, 11, 13}
	fmt.Printf("Primes: %v\n", primes)

	// === SLICES ===
	fmt.Println("\n=== Slices ===")

	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("Original slice: %v\n", numbers)
	fmt.Printf("Length: %d, Capacity: %d\n", len(numbers), cap(numbers))

	// Append
	numbers = append(numbers, 6, 7)
	fmt.Printf("After append: %v\n", numbers)

	// Slice from existing array
	arr2 := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	slice := arr2[2:5]
	fmt.Printf("Sliced [2:5]: %v\n", slice)

	// === MAPS ===
	fmt.Println("\n=== Maps ===")

	person := map[string]string{
		"name": "John",
		"city": "New York",
		"job":  "Engineer",
	}

	fmt.Printf("Person: %v\n", person)
	fmt.Printf("Name: %s\n", person["name"])

	// Check if key exists
	job, ok := person["job"]
	fmt.Printf("Job exists: %v, Value: %s\n", ok, job)

	missing, ok := person["email"]
	fmt.Printf("Email exists: %v, Value: %s\n", ok, missing)

	// Add to map
	person["email"] = "john@example.com"
	fmt.Printf("After adding email: %v\n", person)

	// === CONSTANTS ===
	fmt.Println("\n=== Constants ===")

	const pi = 3.14159265359
	const maxAge int = 150
	const minAge int = 0

	fmt.Printf("Pi: %f\n", pi)
	fmt.Printf("Max age: %d, Min age: %d\n", maxAge, minAge)

	// === TYPE CONVERSION ===
	fmt.Println("\n=== Type Conversion ===")

	f := 3.14
	i := int(f)
	fmt.Printf("Float to int: %f → %d\n", f, i)

	str := "25"
	convertedInt, _ := strconv.Atoi(str)
	fmt.Printf("String to int: %s → %d\n", str, convertedInt)

	num := 42
	strFromInt := strconv.Itoa(num)
	fmt.Printf("Int to string: %d → %s\n", num, strFromInt)

	// === ITERATING OVER COLLECTIONS ===
	fmt.Println("\n=== Iterating Over Collections ===")

	// Slice
	fruits := []string{"apple", "banana", "cherry"}
	for i, fruit := range fruits {
		fmt.Printf("Index %d: %s\n", i, fruit)
	}

	// Map
	fmt.Println("Person details:")
	for key, value := range person {
		fmt.Printf("  %s: %s\n", key, value)
	}

	// String (iterates as runes)
	fmt.Println("Characters in 'Go':")
	for i, r := range "Go" {
		fmt.Printf("  Index %d: %c (code: %d)\n", i, r, r)
	}

	// === MULTIPLE RETURN VALUES ===
	fmt.Println("\n=== Unpacking Values ===")
	result, remainder := divideWithRemainder(10, 3)
	fmt.Printf("10 / 3 = %d remainder %d\n", result, remainder)
}

func divideWithRemainder(a, b int) (int, int) {
	return a / b, a % b
}
