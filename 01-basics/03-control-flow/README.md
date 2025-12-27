# Chapter 3: Control Flow - if/else, switch, and loops

## Learning Objectives

1. Use `if/else` statements effectively
2. Understand the `switch` statement variants
3. Master `for` loops (the only loop in Go)
4. Use `break` and `continue` appropriately
5. Work with `range` for iterating collections
6. Handle early returns for code clarity
7. Write idiomatic Go control flow

---

## if/else Statements

### Basic if

```go
if condition {
    // execute if true
}
```

### if/else

```go
if x > 10 {
    fmt.Println("x is greater than 10")
} else {
    fmt.Println("x is 10 or less")
}
```

### if/else if/else chain

```go
if x < 0 {
    fmt.Println("Negative")
} else if x == 0 {
    fmt.Println("Zero")
} else {
    fmt.Println("Positive")
}
```

### Variable Declaration in if

Go allows variable declaration in the if statement:

```go
if x := getValue(); x > 10 {
    fmt.Println("x is large:", x)
} // x is not accessible here
```

**Key Point**: Variable `x` is scoped to the if block.

### No Parentheses Required

Go doesn't require parentheses around conditions:

```go
// ✓ Valid
if x > 10 {
    // ...
}

// ✗ Syntax error
if (x > 10) {
    // Parentheses are allowed but not idiomatic
}
```

---

## switch Statements

Go's `switch` statement is more powerful than C-style switches.

### Basic switch

```go
switch day {
case "Monday":
    fmt.Println("Start of week")
case "Friday":
    fmt.Println("Almost weekend")
case "Saturday", "Sunday":
    fmt.Println("Weekend")
default:
    fmt.Println("Midweek")
}
```

**Important**: Go's switch doesn't fall through by default. No `break` needed.

### switch with no expression

Use switch like if/else:

```go
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
```

### Type switch

Switch on type instead of value:

```go
func describe(x interface{}) {
    switch v := x.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case bool:
        fmt.Printf("Boolean: %v\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}

describe(42)      // Integer: 42
describe("hello") // String: hello
describe(true)    // Boolean: true
```

### Fallthrough (rare)

To fall through to next case:

```go
switch x {
case 1:
    fmt.Println("One")
    fallthrough  // Falls through to next case
case 2:
    fmt.Println("Two")
}

// Output if x==1:
// One
// Two
```

---

## for Loops (The Only Loop)

Go has one loop construct: `for`. It's flexible and covers all looping needs.

### Classic C-style for

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)  // Prints 0, 1, 2, 3, 4
}
```

**Parts**:
- `i := 0` - initialization
- `i < 5` - condition (checked before each iteration)
- `i++` - increment (after each iteration)

### Omit parts

```go
// While loop (omit init and increment)
i := 0
for i < 5 {
    fmt.Println(i)
    i++
}

// Infinite loop (omit all parts)
for {
    fmt.Println("Infinite")
    break  // Must have break or return
}
```

### range for (iterating collections)

```go
// Slice
nums := []int{10, 20, 30}
for i, val := range nums {
    fmt.Printf("Index %d: %d\n", i, val)
}

// Map
person := map[string]string{
    "name": "John",
    "city": "NYC",
}
for key, value := range person {
    fmt.Printf("%s: %s\n", key, value)
}

// String (iterates as runes)
for i, ch := range "hello" {
    fmt.Printf("%d: %c\n", i, ch)
}
```

### Skip values in range

```go
// Skip index
for _, val := range nums {
    fmt.Println(val)
}

// Skip value
for i := range nums {
    fmt.Println(i)
}
```

---

## break and continue

### break

Exit the loop immediately:

```go
for i := 0; i < 10; i++ {
    if i == 5 {
        break  // Exit loop
    }
    fmt.Println(i)  // Prints 0, 1, 2, 3, 4
}
```

### continue

Skip to next iteration:

```go
for i := 0; i < 5; i++ {
    if i == 2 {
        continue  // Skip i==2
    }
    fmt.Println(i)  // Prints 0, 1, 3, 4
}
```

### Labeled break/continue

Break/continue specific loop in nested loops:

```go
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i == 1 && j == 1 {
            break outer  // Break outer loop
        }
        fmt.Printf("(%d, %d) ", i, j)
    }
}
```

---

## Practical Examples

### Example 1: Fizz Buzz

```go
func fizzbuzz(n int) {
    for i := 1; i <= n; i++ {
        if i%15 == 0 {
            fmt.Println("FizzBuzz")
        } else if i%3 == 0 {
            fmt.Println("Fizz")
        } else if i%5 == 0 {
            fmt.Println("Buzz")
        } else {
            fmt.Println(i)
        }
    }
}
```

### Example 2: Find element in slice

```go
func findIndex(slice []int, target int) int {
    for i, val := range slice {
        if val == target {
            return i  // Early return
        }
    }
    return -1  // Not found
}
```

### Example 3: Count occurrences

```go
func countChar(s string, char rune) int {
    count := 0
    for _, c := range s {
        if c == char {
            count++
        }
    }
    return count
}
```

### Example 4: Nested loops with break

```go
func findPair(matrix [][]int, target int) (int, int, bool) {
    for i := 0; i < len(matrix); i++ {
        for j := 0; j < len(matrix[i]); j++ {
            if matrix[i][j] == target {
                return i, j, true
            }
        }
    }
    return -1, -1, false
}
```

---

## Early Returns (Idiomatic Go)

Instead of deeply nested if/else, use early returns:

```go
// ✗ Not idiomatic (pyramid of doom)
func processUser(id string) error {
    user, err := getUser(id)
    if err == nil {
        if user.Active {
            if user.Premium {
                // Process premium user (deeply nested)
                return nil
            }
        }
    }
    return fmt.Errorf("can't process")
}

// ✓ Idiomatic (early returns)
func processUser(id string) error {
    user, err := getUser(id)
    if err != nil {
        return err
    }

    if !user.Active {
        return fmt.Errorf("user not active")
    }

    if !user.Premium {
        return fmt.Errorf("user not premium")
    }

    // Process premium user (clear path)
    return nil
}
```

---

## Guard Clauses

Check preconditions at start:

```go
func calculateDiscount(age int, income int) float64 {
    // Guard clauses
    if age < 0 || income < 0 {
        return 0
    }

    if age >= 65 {
        return 0.2  // Senior discount
    }

    if income < 30000 {
        return 0.1  // Low income discount
    }

    return 0  // No discount
}
```

---

## Common Patterns

### Iterate until condition

```go
for {
    input, _ := reader.ReadString('\n')
    if input == "quit" {
        break
    }
    process(input)
}
```

### Iterate with step

Go doesn't have step-based ranges, so use classic for:

```go
// Skip every 2 elements
for i := 0; i < len(items); i += 2 {
    fmt.Println(items[i])
}
```

### Iterate map in order

Maps don't have guaranteed order. For ordered iteration:

```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)

for _, k := range keys {
    fmt.Println(k, m[k])
}
```

---

## Exercises

### Exercise 1: Grade Calculator
Write a function that prints grades based on score:
- 90-100: A
- 80-89: B
- 70-79: C
- 60-69: D
- Below 60: F

### Exercise 2: Prime Number Checker
Write a function that checks if a number is prime.

### Exercise 3: Pattern Printing
Print a pyramid pattern:
```
*
**
***
****
*****
```

### Exercise 4: Find duplicates
Find duplicate elements in a slice.

---

## Best Practices

1. **Use early returns** to avoid nesting
2. **Prefer switch** over long if/else chains
3. **Use range** for simplicity when possible
4. **Break/continue** for complex loops
5. **Guard clauses** at function start
6. **Type switch** for interface{} handling

---

## Summary

- Go has **one loop construct**: `for`
- **switch** doesn't fall through by default
- **range** is for iterating collections
- **Early returns** improve readability
- **Guard clauses** prevent nesting

Next chapter: Functions and multiple return values!
