package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("=== Control Flow Examples ===\n")

	// === if/else ===
	fmt.Println("1. if/else statements:")
	x := 15

	if x > 10 {
		fmt.Println("x is greater than 10")
	} else if x == 10 {
		fmt.Println("x is exactly 10")
	} else {
		fmt.Println("x is less than 10")
	}

	// if with variable declaration
	if y := getValue(); y > 20 {
		fmt.Printf("y is large: %d\n", y)
	}

	// === switch ===
	fmt.Println("\n2. switch statements:")
	day := "Friday"

	switch day {
	case "Monday":
		fmt.Println("Start of week")
	case "Friday":
		fmt.Println("Almost weekend")
	case "Saturday", "Sunday":
		fmt.Println("Weekend!")
	default:
		fmt.Println("Midweek")
	}

	// switch with no expression (like if/else)
	fmt.Println("\n3. switch as if/else:")
	score := 85

	switch {
	case score >= 90:
		fmt.Println("Grade: A")
	case score >= 80:
		fmt.Println("Grade: B")
	case score >= 70:
		fmt.Println("Grade: C")
	default:
		fmt.Println("Below C")
	}

	// Type switch
	fmt.Println("\n4. Type switch:")
	testTypeSwitch(42)
	testTypeSwitch("hello")
	testTypeSwitch(3.14)
	testTypeSwitch(true)

	// === for loops ===
	fmt.Println("\n5. for loop (classic C-style):")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// while loop style
	fmt.Println("\n6. while loop style:")
	i := 0
	for i < 3 {
		fmt.Printf("%d ", i)
		i++
	}
	fmt.Println()

	// === range for ===
	fmt.Println("\n7. range for slice:")
	numbers := []int{10, 20, 30, 40}
	for i, val := range numbers {
		fmt.Printf("Index %d: %d\n", i, val)
	}

	// range with map
	fmt.Println("\n8. range for map:")
	person := map[string]string{
		"name": "John",
		"city": "NYC",
		"job":  "Engineer",
	}

	for key, value := range person {
		fmt.Printf("%s: %s\n", key, value)
	}

	// range with string (runes)
	fmt.Println("\n9. range for string (runes):")
	for i, ch := range "Go" {
		fmt.Printf("Index %d: %c (code: %d)\n", i, ch, ch)
	}

	// Skip values in range
	fmt.Println("\n10. Skipping values in range:")
	fmt.Print("Only values: ")
	for _, val := range numbers {
		fmt.Printf("%d ", val)
	}
	fmt.Println()

	fmt.Print("Only indices: ")
	for i := range numbers {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// === break and continue ===
	fmt.Println("\n11. break statement:")
	for i := 0; i < 10; i++ {
		if i == 5 {
			fmt.Println("Breaking at i=5")
			break
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	fmt.Println("\n12. continue statement:")
	for i := 0; i < 5; i++ {
		if i == 2 {
			fmt.Println("Skipping i=2")
			continue
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// === Practical Examples ===
	fmt.Println("\n13. Fizz Buzz:")
	fizzBuzz(15)

	fmt.Println("\n14. Find element in slice:")
	items := []int{1, 2, 3, 4, 5}
	idx := findIndex(items, 3)
	fmt.Printf("Index of 3: %d\n", idx)
	idx = findIndex(items, 10)
	fmt.Printf("Index of 10: %d\n", idx)

	fmt.Println("\n15. Count character:")
	count := countChar("hello world", 'l')
	fmt.Printf("Count of 'l': %d\n", count)

	fmt.Println("\n16. Prime number checker:")
	for num := 1; num <= 20; num++ {
		if isPrime(num) {
			fmt.Printf("%d ", num)
		}
	}
	fmt.Println()

	fmt.Println("\n17. Guard clauses example:")
	result := processUser("user123", true, false)
	fmt.Println("Result:", result)
	result = processUser("user123", true, true)
	fmt.Println("Result:", result)

	fmt.Println("\n18. Nested loop with break:")
	found := false
	for i := 0; i < 3 && !found; i++ {
		for j := 0; j < 3; j++ {
			fmt.Printf("(%d,%d) ", i, j)
			if i == 1 && j == 1 {
				fmt.Println("\nFound target at (1,1)")
				found = true
				break
			}
		}
	}

	fmt.Println("\n19. Early returns (idiomatic):")
	errs := []error{
		validateEmail("invalid"),
		validateEmail("valid@example.com"),
	}
	for _, e := range errs {
		fmt.Println(e)
	}

	fmt.Println("\n20. Sort map keys and iterate:")
	scores := map[string]int{
		"Alice": 90,
		"Bob":   80,
		"Carol": 85,
	}
	keys := make([]string, 0, len(scores))
	for k := range scores {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println("Sorted by name:")
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, scores[k])
	}
}

// Helper functions

func getValue() int {
	return 25
}

func testTypeSwitch(v interface{}) {
	switch val := v.(type) {
	case int:
		fmt.Printf("Integer: %d\n", val)
	case string:
		fmt.Printf("String: %s\n", val)
	case float64:
		fmt.Printf("Float: %f\n", val)
	case bool:
		fmt.Printf("Boolean: %v\n", val)
	default:
		fmt.Printf("Unknown: %T\n", v)
	}
}

func fizzBuzz(n int) {
	for i := 1; i <= n; i++ {
		if i%15 == 0 {
			fmt.Print("FizzBuzz ")
		} else if i%3 == 0 {
			fmt.Print("Fizz ")
		} else if i%5 == 0 {
			fmt.Print("Buzz ")
		} else {
			fmt.Printf("%d ", i)
		}
	}
	fmt.Println()
}

func findIndex(slice []int, target int) int {
	for i, val := range slice {
		if val == target {
			return i
		}
	}
	return -1
}

func countChar(s string, char rune) int {
	count := 0
	for _, c := range s {
		if c == char {
			count++
		}
	}
	return count
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func processUser(id string, active bool, premium bool) string {
	// Guard clauses
	if id == "" {
		return "Invalid ID"
	}

	if !active {
		return "User not active"
	}

	if !premium {
		return "User not premium"
	}

	// Process premium user
	return "Premium user processed successfully"
}

func validateEmail(email string) error {
	// Early return pattern
	if email == "" {
		return fmt.Errorf("email is empty")
	}

	if !contains(email, "@") {
		return fmt.Errorf("email missing @")
	}

	return nil
}

func contains(s string, ch string) bool {
	for _, c := range s {
		if string(c) == ch {
			return true
		}
	}
	return false
}
