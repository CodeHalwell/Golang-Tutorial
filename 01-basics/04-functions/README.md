# Chapter 4: Functions and Multiple Return Values

## Learning Objectives

1. Declare and call functions
2. Return multiple values (Go's unique feature)
3. Use named return values
4. Work with variadic parameters
5. Understand function scope
6. Pass functions as values (first-class functions)
7. Implement closures

---

## Function Declaration

### Basic Function

```go
func greet(name string) {
    fmt.Println("Hello,", name)
}

greet("Alice")
```

### Function with Return Type

```go
func add(a int, b int) int {
    return a + b
}

result := add(3, 4)  // result = 7
```

### Multiple Return Values (Go's Unique Feature!)

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

result, err := divide(10, 2)
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Result:", result)
}
```

---

## Named Return Values

Return values can have names:

```go
func swap(x, y string) (first, second string) {
    first = y
    second = x
    return  // Returns named values
}

a, b := swap("hello", "world")
// a = "world", b = "hello"
```

**When to use**: When the meaning isn't obvious from context.

**When not to use**: For simple, obvious cases.

---

## Variadic Parameters

Functions can accept variable number of arguments:

```go
// Takes zero or more ints
func sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

sum()           // 0
sum(1)          // 1
sum(1, 2, 3)    // 6
sum(1, 2, 3, 4, 5)  // 15
```

**Key Point**: Variadic parameter must be last.

### Unpacking Slices

```go
numbers := []int{1, 2, 3, 4}
total := sum(numbers...)  // Unpack slice
```

---

## Function Parameters

### Pass by Value (Default)

```go
func increment(x int) {
    x = x + 1
}

n := 5
increment(n)
fmt.Println(n)  // Still 5 (not modified)
```

### Pass by Reference (Pointers)

```go
func increment(x *int) {
    *x = *x + 1
}

n := 5
increment(&n)
fmt.Println(n)  // 6 (modified!)
```

### Slices and Maps (Reference-Like Behavior)

```go
func append(slice []int, val int) {
    // Slice header is passed by value, but it contains a pointer
    // to the underlying array. Modifications to existing elements
    // affect the original, but appending may create a new array.
    // Always return and reassign: slice = append(slice, val)
}

func set(m map[string]string, key string, val string) {
    m[key] = val  // Maps are reference types - always affects original
}
```

**Important:** While slices behave like references for existing elements, operations that change the slice's length or capacity (like `append`) may allocate a new array. Always return and use the result: `slice = append(slice, val)`.

---

## Function Variables (First-Class Functions)

Functions can be assigned to variables:

```go
// Assign function to variable
var operation func(int, int) int

operation = add
result := operation(3, 4)  // 7

operation = multiply
result = operation(3, 4)   // 12
```

### Higher-Order Functions

Functions that take or return functions:

```go
func applyOperation(a int, b int, op func(int, int) int) int {
    return op(a, b)
}

add := func(x, y int) int { return x + y }
subtract := func(x, y int) int { return x - y }

applyOperation(10, 5, add)       // 15
applyOperation(10, 5, subtract)  // 5
```

---

## Anonymous Functions and Closures

### Anonymous Function

```go
// Define and call immediately
result := func(x int) int {
    return x * 2
}(5)  // Immediately call with 5
// result = 10
```

### Closure (Function Capturing Variables)

```go
func makeCounter() func() int {
    count := 0

    return func() int {
        count++
        return count
    }
}

counter := makeCounter()
counter()  // 1
counter()  // 2
counter()  // 3

counter2 := makeCounter()
counter2()  // 1 (separate counter)
```

**Key Point**: The returned function "captures" the `count` variable from its enclosing scope.

---

## Defer Statement

Execute function at end of current function:

```go
func readFile(filename string) ([]byte, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()  // Will be called at function end

    // Read file...
    return data, nil
    // file.Close() called here automatically
}
```

**Order**: Multiple defers execute in LIFO (last in, first out):

```go
func demo() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    // Output:
    // 3
    // 2
    // 1
}
```

---

## Panic and Recover

### Panic (Runtime Error)

```go
func divide(a, b float64) float64 {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}
```

### Recover (Catch Panic)

```go
func safeDivide(a, b float64) (result float64, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()

    return a / b, nil  // If b == 0, panics and is recovered
}
```

**When to use**: Rarely. Prefer returning errors.

---

## Practical Examples

### Example 1: Error Handling Function

```go
func validateEmail(email string) error {
    if email == "" {
        return errors.New("email is empty")
    }
    if !strings.Contains(email, "@") {
        return errors.New("email must contain @")
    }
    return nil
}

if err := validateEmail("test@example.com"); err != nil {
    fmt.Println("Invalid:", err)
}
```

### Example 2: Processing with Callbacks

```go
func processData(data []int, callback func(int) bool) {
    for _, item := range data {
        if !callback(item) {
            break
        }
    }
}

numbers := []int{1, 2, 3, 4, 5}
processData(numbers, func(n int) bool {
    fmt.Println(n)
    return n < 4
})
```

### Example 3: Filtering Function

```go
func filter(numbers []int, predicate func(int) bool) []int {
    var result []int
    for _, n := range numbers {
        if predicate(n) {
            result = append(result, n)
        }
    }
    return result
}

numbers := []int{1, 2, 3, 4, 5}
evens := filter(numbers, func(n int) bool {
    return n%2 == 0
})
// evens = [2, 4]
```

---

## Best Practices

1. **Use multiple returns for errors**
   ```go
   func read() ([]byte, error)  // ✓
   func read() ([]byte)         // ✗ Bad, no error handling
   ```

2. **Named returns for clarity**
   ```go
   func divide(a, b float64) (quotient float64, err error)  // ✓
   ```

3. **Error as last return value**
   ```go
   func process() (result string, err error)  // ✓
   ```

4. **Defer for cleanup**
   ```go
   defer file.Close()
   defer mutex.Unlock()
   defer db.Close()
   ```

5. **Don't panic for recoverable errors**
   ```go
   // ✗ Bad
   if err != nil {
       panic(err)
   }

   // ✓ Good
   if err != nil {
       return fmt.Errorf("failed: %w", err)
   }
   ```

---

## Exercises

### Exercise 1: Argument Validation
Write a function that validates multiple arguments and returns errors.

### Exercise 2: Map Function
Implement a generic map function that applies a function to each element.

### Exercise 3: Retry Function
Write a retry function that calls another function with exponential backoff.

### Exercise 4: Resource Management
Create a function that properly cleans up resources using defer.

---

## Summary

- **Multiple returns**: Go's signature feature for error handling
- **Variadic parameters**: Accept variable numbers of arguments
- **First-class functions**: Assign functions to variables
- **Closures**: Functions that capture enclosing scope
- **Defer**: Cleanup code guaranteed to run
- **Panic/Recover**: Runtime error handling (use sparingly)

Next: Chapter 5 - Structs and Interfaces!
