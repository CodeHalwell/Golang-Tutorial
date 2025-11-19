package main

import (
	"cmp"
	"fmt"
	"slices"
)

// === GO 1.18+ GENERICS ===
//
// Generics allow you to write functions and types that work with any type,
// while still being type-safe at compile time.
//
// Syntax: func Name[TypeParam Constraint](param TypeParam) TypeParam

// Example 1: Generic function that works with any type
// Comparable is a constraint that requires the type to support == and !=
func indexOf[T comparable](slice []T, target T) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

// Example 2: Generic function with multiple type parameters
func mapSlice[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Example 3: Generic type (data structure)
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if len(s.items) == 0 {
		return zero, fmt.Errorf("stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Example 4: Generic with ordered constraint (for sorting)
func findMax[T cmp.Ordered](slice []T) (T, error) {
	var zero T
	if len(slice) == 0 {
		return zero, fmt.Errorf("empty slice")
	}

	max := slice[0]
	for _, v := range slice[1:] {
		if v > max {
			max = v
		}
	}
	return max, nil
}

// Example 5: Custom constraint
type Number interface {
	int | int64 | float64
}

func sum[T Number](slice []T) T {
	var total T
	for _, v := range slice {
		total += v
	}
	return total
}

// Example 6: Generic with multiple constraints
type Stringable interface {
	String() string
}

type Comparable interface {
	Stringable
	comparable
}

// Example 7: Node for a generic linked list
type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type LinkedList[T any] struct {
	head *Node[T]
}

func (ll *LinkedList[T]) Add(value T) {
	newNode := &Node[T]{Value: value}
	if ll.head == nil {
		ll.head = newNode
		return
	}

	current := ll.head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
}

func (ll *LinkedList[T]) Print() {
	current := ll.head
	for current != nil {
		fmt.Printf("%v -> ", current.Value)
		current = current.Next
	}
	fmt.Println("nil")
}

// Example 8: Generic function with type constraints using union
func processNumbers[T int | float64](value T) T {
	return value * 2
}

// Example 9: Generic container
type Container[T any] struct {
	values []T
}

func (c *Container[T]) Add(value T) {
	c.values = append(c.values, value)
}

func (c *Container[T]) GetAll() []T {
	return c.values
}

func (c *Container[T]) Filter(predicate func(T) bool) []T {
	var result []T
	for _, v := range c.values {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Example 10: Generic pair
type Pair[K any, V any] struct {
	Key   K
	Value V
}

func (p *Pair[K, V]) String() string {
	return fmt.Sprintf("%v: %v", p.Key, p.Value)
}

func main() {
	fmt.Println("=== GO GENERICS EXAMPLES ===\n")

	// Example 1: Generic indexOf
	fmt.Println("1. Generic indexOf:")
	ints := []int{1, 2, 3, 4, 5}
	idx := indexOf(ints, 3)
	fmt.Printf("Index of 3 in %v: %d\n", ints, idx)

	strings := []string{"apple", "banana", "cherry"}
	idx = indexOf(strings, "banana")
	fmt.Printf("Index of 'banana' in %v: %d\n", strings, idx)

	// Example 2: Generic map
	fmt.Println("\n2. Generic mapSlice:")
	numbers := []int{1, 2, 3, 4}
	squared := mapSlice(numbers, func(n int) int {
		return n * n
	})
	fmt.Printf("Squared: %v\n", squared)

	stringified := mapSlice(numbers, func(n int) string {
		return fmt.Sprintf("num_%d", n)
	})
	fmt.Printf("Stringified: %v\n", stringified)

	// Example 3: Generic Stack
	fmt.Println("\n3. Generic Stack[int]:")
	intStack := &Stack[int]{}
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)

	for !intStack.IsEmpty() {
		val, _ := intStack.Pop()
		fmt.Printf("Popped: %d\n", val)
	}

	fmt.Println("\n   Generic Stack[string]:")
	stringStack := &Stack[string]{}
	stringStack.Push("hello")
	stringStack.Push("world")

	for !stringStack.IsEmpty() {
		val, _ := stringStack.Pop()
		fmt.Printf("Popped: %s\n", val)
	}

	// Example 4: Find max with ordered constraint
	fmt.Println("\n4. Find max (Ordered constraint):")
	intNumbers := []int{5, 2, 8, 1, 9}
	max, _ := findMax(intNumbers)
	fmt.Printf("Max int: %d\n", max)

	floatNumbers := []float64{5.5, 2.2, 8.8, 1.1, 9.9}
	maxFloat, _ := findMax(floatNumbers)
	fmt.Printf("Max float: %f\n", maxFloat)

	// Example 5: Sum with custom Number constraint
	fmt.Println("\n5. Sum (custom Number constraint):")
	sumInt := sum([]int{1, 2, 3, 4, 5})
	fmt.Printf("Sum of ints: %d\n", sumInt)

	sumFloat := sum([]float64{1.5, 2.5, 3.5})
	fmt.Printf("Sum of floats: %f\n", sumFloat)

	// Example 6: Generic Linked List
	fmt.Println("\n6. Generic LinkedList[int]:")
	intList := &LinkedList[int]{}
	intList.Add(1)
	intList.Add(2)
	intList.Add(3)
	fmt.Print("List: ")
	intList.Print()

	// Example 7: Process numbers
	fmt.Println("\n7. Process numbers (int | float64):")
	intResult := processNumbers(5)
	fmt.Printf("5 * 2 = %d\n", intResult)

	floatResult := processNumbers(2.5)
	fmt.Printf("2.5 * 2 = %f\n", floatResult)

	// Example 8: Generic Container with Filter
	fmt.Println("\n8. Generic Container with filter:")
	container := &Container[int]{}
	container.Add(1)
	container.Add(2)
	container.Add(3)
	container.Add(4)
	container.Add(5)

	evens := container.Filter(func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("Even numbers: %v\n", evens)

	// Example 9: Generic Pair
	fmt.Println("\n9. Generic Pair:")
	stringPair := &Pair[string, int]{Key: "age", Value: 30}
	fmt.Println("Pair:", stringPair)

	floatPair := &Pair[string, float64]{Key: "height", Value: 5.9}
	fmt.Println("Pair:", floatPair)

	// Example 10: Built-in generics from slices package
	fmt.Println("\n10. Built-in generics (slices package):")
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Printf("Original: %v\n", nums)

	// Sort (requires cmp.Ordered)
	slices.Sort(nums)
	fmt.Printf("Sorted: %v\n", nums)

	// SortFunc (with custom comparison)
	slices.SortFunc(nums, func(a, b int) int {
		if a < b {
			return 1  // Descending order
		}
		if a > b {
			return -1
		}
		return 0
	})
	fmt.Printf("Sorted descending: %v\n", nums)

	// Contains
	contains := slices.Contains(nums, 5)
	fmt.Printf("Contains 5: %v\n", contains)

	// Index
	index, found := slices.BinarySearch([]int{1, 1, 2, 3, 4, 5, 6, 9}, 4)
	fmt.Printf("BinarySearch for 4: index=%d, found=%v\n", index, found)

	fmt.Println("\n=== GENERICS BEST PRACTICES ===")
	fmt.Println("1. Use generics for reusable data structures (Stack, Queue, etc)")
	fmt.Println("2. Use constraints to limit type parameters (comparable, Ordered)")
	fmt.Println("3. Don't overuse generics - use only when beneficial")
	fmt.Println("4. Generics can make code clearer, but can also add complexity")
	fmt.Println("5. Consider readability - sometimes concrete types are better")
	fmt.Println("6. Use standard constraints from comparable and cmp packages")
	fmt.Println("7. Define custom constraints for domain-specific types")
	fmt.Println("8. Generics work with types, not values (compile-time only)")
}
