# Chapter 2: Variables, Types, and Constants

## Learning Objectives

1. Declare and initialize variables using different syntaxes
2. Understand Go's type system (static, strongly typed)
3. Work with primitive types (integers, floats, booleans, strings)
4. Use composite types (arrays, slices, maps)
5. Understand zero values
6. Create and use constants effectively
7. Type conversion and casting

---

## Go's Type System

Go is **statically typed** and **strongly typed**:

- **Statically typed**: Types are checked at compile time, not runtime
- **Strongly typed**: No implicit type conversions (you must be explicit)

```go
// ✓ Valid: Explicit type declaration
var x int = 10

// ✓ Valid: Type inference
var x = 10  // Go infers int

// ✗ Invalid: No implicit conversion
x = 10.5  // Compile error! Can't assign float64 to int
```

---

## Variable Declaration

### Syntax 1: `var` with Type

```go
var name string = "John"
var age int = 25
var score float64 = 95.5
```

### Syntax 2: `var` with Type Inference

```go
var name = "John"      // Type inferred as string
var age = 25           // Type inferred as int
var score = 95.5       // Type inferred as float64
```

### Syntax 3: Short Variable Declaration (Most Idiomatic)

```go
name := "John"
age := 25
score := 95.5
```

**Rules:**
- Only works inside functions
- Can't be used at package level
- Most used in Go code

### Syntax 4: Multiple Declarations

```go
var (
    name string = "John"
    age int = 25
    score float64 = 95.5
)

// Or with short syntax
name, age, score := "John", 25, 95.5
```

### Syntax 5: Blank Identifier for Unused Values

```go
// If we don't care about the second return value
result, _ := someFunction()
```

---

## Primitive Data Types

### Integers

```go
var signed8 int8 = 127           // Range: -128 to 127
var signed16 int16 = 32767       // Range: -32,768 to 32,767
var signed32 int32 = 2147483647  // Range: -2^31 to 2^31-1
var signed64 int64 = 9223372036854775807

var unsigned8 uint8 = 255        // Range: 0 to 255 (also called byte)
var unsigned16 uint16 = 65535
var unsigned32 uint32 = 4294967295
var unsigned64 uint64 = 18446744073709551615

// Platform-dependent (32-bit on 32-bit systems, 64-bit on 64-bit systems)
var platformInt int = 100
var platformUint uint = 100

// Aliases
var b byte = 255                 // Same as uint8
var r rune = 'A'                 // Same as int32, represents Unicode code point
```

**Size Chart:**
| Type | Size | Range |
|------|------|-------|
| int8 | 1 byte | -128 to 127 |
| int16 | 2 bytes | -32,768 to 32,767 |
| int32 | 4 bytes | ~2 billion |
| int64 | 8 bytes | ~9 quintillion |
| uint8 | 1 byte | 0 to 255 |
| uint16 | 2 bytes | 0 to 65,535 |
| uint32 | 4 bytes | ~4 billion |
| uint64 | 8 bytes | ~18 quintillion |

### Floating Point Numbers

```go
var f32 float32 = 3.14          // 32-bit float
var f64 float64 = 3.14159265359 // 64-bit float (default for literals)

// Scientific notation
var million float64 = 1e6  // 1,000,000
var tiny float64 = 1e-10   // 0.0000000001
```

### Booleans

```go
var active bool = true
var inactive bool = false

// Default value (zero value)
var uninitialized bool  // false
```

### Strings

```go
var message string = "Hello, World!"
var empty string        // "" (empty string)

// String literals
raw := `Line 1
Line 2
\n is literal`  // Raw string, backslashes not interpreted

interpreted := "Line 1\nLine 2"  // Interpreted string, \n is newline

// String concatenation
greeting := "Hello" + ", " + "World"
```

**String vs byte vs rune:**
```go
var str string = "Hello"
var b byte = str[0]      // 'H' (72) - First byte
var r rune = 'H'         // 'H' (72) - Unicode character

// Iterating over string
for i, r := range str {  // i is index, r is rune
    fmt.Printf("Index %d: %c (code %d)\n", i, r, r)
}
```

---

## Zero Values

Every type has a **zero value** (default value):

```go
var i int          // 0
var f float64      // 0.0
var s string       // ""
var b bool         // false
var p *int         // nil
var sl []int       // nil
var m map[string]int  // nil
```

This is why Go doesn't require explicit initialization:
```go
var age int        // Automatically 0, not undefined
```

---

## Type Conversion

Go doesn't have implicit type conversion, so you must be explicit:

```go
var f float64 = 3.14
var i int = int(f)        // Explicit conversion: i = 3 (loses decimal)

var u uint = uint(i)      // int to uint

// String conversions
var num int = 42
var str string = fmt.Sprintf("%d", num)  // int to string

// Using strconv package for more conversions
age, _ := strconv.Atoi("25")  // string to int
str := strconv.Itoa(25)       // int to string
```

---

## Constants

Constants are values that **cannot change** after initialization:

```go
const pi float64 = 3.14159
const greeting string = "Hello, World!"
const year int = 2024

// Multiple constants
const (
    maxAge = 120
    minAge = 0
    defaultName = "Guest"
)
```

### Typed vs Untyped Constants

```go
// Typed constant (must match type exactly)
const typedPI float64 = 3.14159

// Untyped constant (more flexible)
const untypedPI = 3.14159  // Can be assigned to float32 or float64
```

### Enumerated Constants with `iota`

```go
const (
    Sunday = iota   // 0
    Monday          // 1
    Tuesday         // 2
    Wednesday       // 3
    Thursday        // 4
    Friday          // 5
    Saturday        // 6
)

// Practical example: File permissions
const (
    Read = 1 << iota   // 1 (binary: 001)
    Write              // 2 (binary: 010)
    Execute            // 4 (binary: 100)
)
```

---

## Composite Types

### Arrays (Fixed Size)

```go
// Declare array of 5 integers
var numbers [5]int
numbers[0] = 10
numbers[1] = 20

// Array literal
scores := [3]float64{90.5, 85.0, 92.5}

// Let compiler count size
data := [...]int{1, 2, 3, 4, 5}  // Length is 5

// Accessing elements
fmt.Println(scores[0])  // 90.5
fmt.Println(len(scores))  // 3
```

### Slices (Dynamic Size)

```go
// Slice declaration (no fixed size)
var numbers []int
numbers = append(numbers, 1)
numbers = append(numbers, 2, 3)

// Slice literal
primes := []int{2, 3, 5, 7, 11}

// Slicing an array
arr := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
slice := arr[2:5]  // Elements at indices 2, 3, 4 → [2, 3, 4]

// Make slice with capacity
buffer := make([]int, 0, 10)  // len=0, capacity=10

// Common slice operations
slice = append(slice, 100)     // Append element
slice = append(slice, 101, 102)  // Append multiple
len(slice)                     // Get length
cap(slice)                     // Get capacity
```

### Maps (Key-Value Pairs)

```go
// Map declaration
var ages map[string]int
ages = make(map[string]int)
ages["Alice"] = 30
ages["Bob"] = 25

// Map literal
person := map[string]string{
    "name": "John",
    "city": "New York",
    "job": "Engineer",
}

// Accessing values
fmt.Println(person["name"])  // "John"

// Checking if key exists
job, ok := person["job"]  // ok = true
missing, ok := person["email"]  // ok = false

// Iterating over map
for key, value := range person {
    fmt.Printf("%s: %s\n", key, value)
}

// Deleting from map
delete(person, "job")
```

---

## Type Declarations

You can create new types based on existing types:

```go
// Create a new type based on int
type Age int

// Create a new type based on string
type UserID string

// Usage
var userAge Age = 30
var userID UserID = "user123"

// Note: Can't directly compare Age with int
// var i int = userAge  // Compile error!
var i int = int(userAge)  // Must convert explicitly
```

---

## Practical Examples

### Example 1: User Data Structure

```go
package main

import "fmt"

func main() {
    // User information
    firstName := "John"
    lastName := "Doe"
    age := 30
    email := "john@example.com"
    isActive := true

    fmt.Printf("Name: %s %s\n", firstName, lastName)
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Email: %s\n", email)
    fmt.Printf("Active: %v\n", isActive)
}
```

### Example 2: Collection Processing

```go
package main

import "fmt"

func main() {
    // Slice of numbers
    numbers := []int{1, 2, 3, 4, 5}

    // Sum all numbers
    sum := 0
    for _, num := range numbers {
        sum += num
    }
    fmt.Printf("Sum: %d\n", sum)
    fmt.Printf("Average: %f\n", float64(sum)/float64(len(numbers)))
}
```

### Example 3: Type Conversions

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // String to int
    strAge := "25"
    age, _ := strconv.Atoi(strAge)
    fmt.Printf("Age (int): %d\n", age)

    // Int to string
    num := 42
    str := strconv.Itoa(num)
    fmt.Printf("Number (string): %s\n", str)

    // Float to int (loses precision)
    f := 3.14159
    i := int(f)
    fmt.Printf("Float to int: %d\n", i)
}
```

---

## Common Mistakes

### Mistake 1: Implicit Type Conversion

```go
// ✗ Won't compile
var i int = 10
var f float64 = i  // Compile error!

// ✓ Correct
var f float64 = float64(i)
```

### Mistake 2: Mixing int and float

```go
// ✗ Won't compile
var result = 10 + 3.14  // Can't add int and float64

// ✓ Correct
var result = float64(10) + 3.14
```

### Mistake 3: Modifying Constants

```go
// ✗ Won't compile
const pi = 3.14159
pi = 3.14  // Compile error!

// ✓ Use var instead
var pi = 3.14159
pi = 3.14  // OK
```

---

## Best Practices

1. **Use short variable declaration** (`:=`) inside functions
2. **Use meaningful variable names** - `age` not `a`
3. **Prefer `const`** for values that don't change
4. **Use slice over array** unless size is truly fixed
5. **Explicit type conversions** show intent clearly
6. **Zero values are your friend** - they provide sensible defaults

---

## Exercises

### Exercise 1: Temperature Conversion
Write a program that converts Fahrenheit to Celsius:
- Formula: (F - 32) × 5/9
- Create variables for Fahrenheit and Celsius
- Display both values

### Exercise 2: Student Grades
Create variables for:
- Student name (string)
- Test scores (slice of ints)
- Calculate and display average grade

### Exercise 3: User Profile Map
Create a map with:
- name, email, age, city
- Display in a readable format
- Add error handling for missing fields

---

## Summary

- Go is statically and strongly typed
- Use `:=` for variable declaration inside functions
- Every type has a zero value
- Type conversion must be explicit
- Constants cannot change after initialization
- Slices are more flexible than arrays for most cases
- Maps provide key-value storage
- Zero values provide sensible defaults

---

## Next: Chapter 3 - Control Flow

We'll explore conditionals (if/else/switch) and loops (for/while).
