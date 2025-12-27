# Chapter 1: Hello World and Go Fundamentals

## Learning Objectives

After completing this chapter, you will understand:

1. The structure of a Go program (packages, imports, main function)
2. How to compile and run Go code
3. The Go compilation model and why Go compiles to native binaries
4. The `fmt` package for formatted output
5. Comments and documentation
6. How Go differs from interpreted languages like Python
7. The entry point of Go programs

---

## Key Concepts

### Go Program Structure

Every Go program consists of:

```
1. Package declaration
2. Import statements
3. Function declarations
4. Code execution
```

### The `main` Package

- Go programs require a `main` package
- The `main` function is the entry point
- A program can only have ONE `main` function
- The `main` function takes NO arguments and returns NOTHING

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Compilation Model

Go is a **compiled language**, not interpreted:

```
Source Code (.go) → Go Compiler → Binary Executable
    hello.go     →  go build  →  hello (or hello.exe on Windows)
```

Benefits:
- Single binary with no runtime dependencies
- Fast execution
- Easy distribution (copy the binary)
- Early error detection

### Packages in Go

A package is a namespace for organizing code:

- All code belongs to a package
- `package main` is special - it's executable
- `package mypackage` creates a reusable library
- Import paths are derived from directory structure

---

## Detailed Examples

### Example 1: Simplest Program

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

**Execution:**
```bash
$ go run hello.go
Hello, World!
```

**What happens:**
1. `package main` declares this is an executable program
2. `import "fmt"` imports the formatting package from stdlib
3. `func main()` is the entry point
4. `fmt.Println()` prints text with a newline

### Example 2: Multiple Print Statements

```go
package main

import "fmt"

func main() {
    // Print with newline
    fmt.Println("Line 1")
    fmt.Println("Line 2")

    // Print without newline
    fmt.Print("No newline")
    fmt.Print(" here\n")

    // Formatted printing
    fmt.Printf("Number: %d\n", 42)
}
```

**Output:**
```
Line 1
Line 2
No newline here
Number: 42
```

### Example 3: Using Multiple Packages

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("Hello from Go!")
    fmt.Printf("Current time: %v\n", time.Now())
}
```

### Example 4: Comments and Documentation

```go
// This is a single-line comment

/*
This is a multi-line comment.
It can span multiple lines
and is useful for longer explanations.
*/

package main

import "fmt"

func main() {
    // Comments are ignored by the compiler
    fmt.Println("Comments are useful!")
}
```

### Example 5: Using Variables in Output

```go
package main

import "fmt"

func main() {
    name := "Gopher"
    version := 1.21

    fmt.Println("Welcome to Go", version)
    fmt.Printf("Hello, %s!\n", name)
}
```

---

## How to Run Go Programs

### Method 1: `go run` (Direct Execution)

```bash
$ go run main.go
Hello, World!
```

Best for: Development, quick testing, scripts

### Method 2: `go build` (Compile to Binary)

```bash
$ go build -o hello main.go
$ ./hello
Hello, World!
```

Best for: Distribution, running multiple times, production

### Method 3: `go install` (Compile and Install)

```bash
$ go install ./hello
$ $GOPATH/bin/hello
Hello, World!
```

Best for: Installing tools globally

---

## Behind the Scenes: The Compilation Process

### Step 1: Lexical Analysis
The source code is tokenized:
```
"fmt.Println" → [PACKAGE_NAME] [DOT] [FUNCTION_NAME]
```

### Step 2: Syntax Analysis
The tokens are parsed into an Abstract Syntax Tree (AST):
```
Import fmt
FunctionCall(Println, "Hello, World!")
```

### Step 3: Semantic Analysis
Type checking and validation occur:
- Does `fmt.Println` exist?
- Are the arguments correct types?

### Step 4: Code Generation
The AST is converted to intermediate code:
- Syntax is valid
- Types are correct
- Ready for machine code generation

### Step 5: Linking
Libraries are linked into a single binary:
- All imports are resolved
- External dependencies are included

### Step 6: Executable Creation
The final binary is created:
- Can be executed directly
- No runtime needed
- Portable within same OS/architecture

---

## Why Go's Compilation Model Matters

### Go vs Python
```python
# Python (interpreted)
python script.py  # Requires Python installed
```

```go
// Go (compiled)
go build -o script script.go
./script  // No Go installation needed!
```

### Go vs C
```go
// Go compilation is simple and automatic
go build main.go

// C requires manual compilation with flags
gcc -o main main.c -Wall -O2  # Much more complex
```

---

## Go vs Other Languages Comparison

| Language | Type | Distribution | Startup | Binary Size |
|----------|------|--------------|---------|-------------|
| Go | Compiled | Single binary | Instant | ~5-10 MB |
| Python | Interpreted | Needs runtime | Slow | N/A |
| Node.js | Interpreted | Needs runtime | Slow | N/A |
| Rust | Compiled | Single binary | Instant | ~3-5 MB |
| Java | Compiled to bytecode | Needs JVM | Slow (JVM startup) | ~100 MB+ |

---

## Essential Go Vocabulary

- **Package**: Namespace for organizing code
- **Import**: Including another package
- **Function**: Reusable block of code
- **Binary**: Compiled executable program
- **Goroutine**: Lightweight concurrency primitive (more later)
- **Channel**: Mechanism for goroutines to communicate (more later)

---

## Common Beginner Questions

### Q: Why do I need to write `package main`?
**A:** Go requires it to identify which code is executable vs. library code.

### Q: Can I name my file anything?
**A:** Yes, but `main.go` is conventional. Use `go run .` to run all files in a directory.

### Q: How do I pass command-line arguments?
**A:** Via the `os.Args` slice (covered in next chapters).

### Q: Why is the binary so large?
**A:** Go includes the entire runtime (memory management, garbage collector, etc.). For most programs, 5-10 MB is normal.

### Q: Can I compile for a different OS?
**A:** Yes! Use environment variables:
```bash
GOOS=linux GOARCH=amd64 go build -o hello main.go
```

---

## Production Considerations

### Error Handling (Even in Hello World)
```go
package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    // Even simple operations can fail
    _, err := io.WriteString(os.Stdout, "Hello, World!\n")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### Logging from the Start
```go
package main

import (
    "log"
)

func main() {
    log.Println("Application started")
    log.Println("Hello, World!")
    log.Println("Application finished")
}
```

---

## Exercises

### Exercise 1: Hello World Variations
Create a program that prints:
```
Hello, World!
Hello, Go!
Hello, Gopher!
```

**Solution:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    fmt.Println("Hello, Go!")
    fmt.Println("Hello, Gopher!")
}
```

### Exercise 2: Formatted Output
Write a program that prints a formatted message:
```
Name: John Doe
Age: 25
Language: Go
```

**Solution:**
```go
package main

import "fmt"

func main() {
    name := "John Doe"
    age := 25
    language := "Go"

    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Language: %s\n", language)
}
```

### Exercise 3: Multiple Files in Same Package
Create two files in the same package and run them together:

**hello.go:**
```go
package main

import "fmt"

func main() {
    sayHello()
}

func sayHello() {
    fmt.Println("Hello from hello.go")
}
```

**greeter.go:**
```go
package main

import "fmt"

func sayGoodbye() {
    fmt.Println("Goodbye from greeter.go")
}
```

**Run with:**
```bash
# Run all .go files in the current directory
go run .

# Or explicitly list the files
go run hello.go greeter.go
```

Both files are compiled together as part of the same `main` package. You can call functions from either file within the same package.

---

## Summary

- Go programs start with `package main` and a `main()` function
- Go is compiled to a native binary
- The compilation process ensures type safety before execution
- `go run` for development, `go build` for distribution
- Comments explain code; `fmt.Println` prints output
- No external runtime needed to run compiled Go binaries

---

## Next Steps

In Chapter 2, we'll explore **Variables and Types**, learning how to:
- Declare and initialize variables
- Understand Go's type system
- Work with different data types
- Use constants effectively

You now have a complete foundation for Go development!
