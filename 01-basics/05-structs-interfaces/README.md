# Chapter 5: Structs and Methods

## Learning Objectives

1. Define structs with fields and types
2. Create and initialize struct instances
3. Implement methods on types
4. Understand receivers (value vs pointer)
5. Design structs effectively
6. Work with anonymous structs
7. Embed structs for composition

---

## Structs - Grouping Related Data

A struct is a collection of fields (data members):

```go
type Person struct {
    Name  string
    Age   int
    Email string
}
```

### Declaring Struct Variables

```go
// Method 1: Zero values
var p Person  // p = Person{"", 0, ""}

// Method 2: Field names
p := Person{
    Name:  "Alice",
    Age:   30,
    Email: "alice@example.com",
}

// Method 3: Positional (order matters!)
p := Person{"Alice", 30, "alice@example.com"}

// Method 4: Partial initialization
p := Person{Name: "Alice"}  // Age and Email are zero values
```

### Accessing Fields

```go
fmt.Println(p.Name)  // "Alice"
p.Age = 31           // Modify field
```

---

## Methods - Functions on Types

A method is a function with a **receiver**:

```go
func (p Person) Describe() string {
    return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

p := Person{Name: "Alice", Age: 30}
fmt.Println(p.Describe())  // "Alice is 30 years old"
```

### Receiver Types: Value vs Pointer

#### Value Receiver (Receiver gets a copy)

```go
func (p Person) Birthday() {
    p.Age++  // Only modifies the copy, not original
}

p := Person{Name: "Alice", Age: 30}
p.Birthday()
fmt.Println(p.Age)  // Still 30!
```

#### Pointer Receiver (Receiver modifies original)

```go
func (p *Person) Birthday() {
    p.Age++  // Modifies original
}

p := Person{Name: "Alice", Age: 30}
p.Birthday()
fmt.Println(p.Age)  // 31 (modified!)
```

### When to Use Each

**Value Receiver**:
- Small structs
- No modifications needed
- Type is immutable

**Pointer Receiver**:
- Large structs (avoid copying)
- Need to modify
- Implement interfaces consistently
- Reference semantics needed

### Pointer Receiver on Pointer Variables

```go
p := &Person{Name: "Alice", Age: 30}
p.Birthday()  // Go automatically dereferences for methods
```

Go handles `(*p).Birthday()` as `p.Birthday()` automatically.

---

## Struct Composition (Not Inheritance)

Go doesn't have inheritance. Instead, use **embedding** for composition:

```go
type Animal struct {
    Name string
    Age  int
}

type Dog struct {
    Animal  // Embedded struct
    Breed   string
}

// Create a Dog
dog := Dog{
    Animal: Animal{Name: "Rex", Age: 5},
    Breed:  "Golden Retriever",
}

// Access embedded fields directly
fmt.Println(dog.Name)   // "Rex" (promoted field)
fmt.Println(dog.Breed)  // "Golden Retriever"
```

### Promoted Methods

Methods from embedded structs are promoted:

```go
type Animal struct {
    Name string
}

func (a *Animal) Speak() {
    fmt.Println(a.Name, "makes a sound")
}

type Dog struct {
    Animal
    Breed string
}

dog := Dog{Animal: Animal{Name: "Rex"}, Breed: "Labrador"}
dog.Speak()  // "Rex makes a sound" (promoted method)
```

---

## Anonymous Structs

Structs without names, useful for one-off uses:

```go
// Create anonymous struct
person := struct {
    Name string
    Age  int
}{
    Name: "Alice",
    Age:  30,
}

fmt.Println(person.Name)  // "Alice"
```

**Use Cases**:
- JSON unmarshaling
- Temporary data grouping
- Function arguments
- Return values

```go
// Return multiple values with meaningful names
func getUserData() struct {
    Name  string
    Email string
} {
    return struct {
        Name  string
        Email string
    }{
        Name:  "Alice",
        Email: "alice@example.com",
    }
}
```

---

## Struct Tags

Metadata attached to struct fields:

```go
type User struct {
    Name  string `json:"name" validate:"required"`
    Age   int    `json:"age" validate:"min=0,max=150"`
    Email string `json:"email" validate:"email"`
}
```

### Common Tags

| Tag | Purpose | Example |
|-----|---------|---------|
| `json` | JSON marshaling | `json:"field_name"` |
| `xml` | XML marshaling | `xml:"field_name"` |
| `validate` | Validation | `validate:"required"` |
| `db` | Database mapping | `db:"user_id"` |
| `orm` | ORM mapping | `gorm:"primary_key"` |

### Using Tags

```go
import (
    "encoding/json"
)

user := User{Name: "Alice", Age: 30, Email: "alice@example.com"}

// Marshal to JSON
jsonData, _ := json.Marshal(user)
fmt.Println(string(jsonData))
// {"name":"Alice","age":30,"email":"alice@example.com"}

// Unmarshal from JSON
var u User
json.Unmarshal(jsonData, &u)
```

### Reflect to Read Tags

```go
import "reflect"

func PrintTags(v interface{}) {
    t := reflect.TypeOf(v)
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        tag := field.Tag.Get("json")
        fmt.Printf("%s -> %s\n", field.Name, tag)
    }
}
```

---

## Struct Best Practices

### 1. Small Structs with Single Responsibility

```go
// ✓ Good: Single responsibility
type User struct {
    Name  string
    Email string
}

type Address struct {
    Street string
    City   string
    ZIP    string
}

// ✗ Bad: Too many concerns
type Person struct {
    Name    string
    Email   string
    Street  string
    City    string
    ZIP     string
    Phone   string
    // ... 10 more fields
}
```

### 2. Use Value Receivers for Immutable Types

```go
type Point struct {
    X, Y float64
}

// ✓ Good: Value receiver for immutable operation
func (p Point) Distance() float64 {
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// ✓ Good: Pointer receiver for mutation
func (p *Point) Move(dx, dy float64) {
    p.X += dx
    p.Y += dy
}
```

### 3. Document Exported Types

```go
// ✓ Good: Public types are documented
// User represents a user in the system
type User struct {
    // Name is the user's full name
    Name string
    // Email is the user's email address
    Email string
}

// GetUser retrieves a user by ID
func GetUser(id string) (*User, error) {
    // ...
}
```

### 4. Constructor Functions for Complex Initialization

```go
// ✓ Good: Constructor ensures valid state
func NewUser(name, email string) (*User, error) {
    if name == "" {
        return nil, fmt.Errorf("name required")
    }
    if !strings.Contains(email, "@") {
        return nil, fmt.Errorf("invalid email")
    }
    return &User{Name: name, Email: email}, nil
}

// Usage
user, err := NewUser("Alice", "alice@example.com")
```

### 5. Use Embedding for Code Reuse

```go
// ✓ Good: Embedding promotes methods
type Named struct {
    Name string
}

func (n *Named) SetName(name string) {
    n.Name = name
}

type Employee struct {
    Named  // Embedded - SetName is promoted
    Title  string
    Salary float64
}

emp := Employee{Named: Named{Name: "Alice"}, Title: "Engineer"}
emp.SetName("Bob")  // Works due to promoted method
```

---

## Practical Examples

### Example 1: Bank Account

```go
type Account struct {
    AccountNumber string
    Owner         string
    Balance       float64
}

func (a *Account) Deposit(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("deposit amount must be positive")
    }
    a.Balance += amount
    return nil
}

func (a *Account) Withdraw(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("withdrawal amount must be positive")
    }
    if amount > a.Balance {
        return fmt.Errorf("insufficient funds")
    }
    a.Balance -= amount
    return nil
}

func (a Account) String() string {
    return fmt.Sprintf("%s: $%.2f", a.Owner, a.Balance)
}

// Usage
account := Account{
    AccountNumber: "12345",
    Owner:         "Alice",
    Balance:       1000,
}

account.Deposit(500)    // Balance: 1500
account.Withdraw(300)   // Balance: 1200
fmt.Println(account)    // "Alice: $1200.00"
```

### Example 2: Configuration with Validation

```go
type Config struct {
    Host     string `json:"host" validate:"required"`
    Port     int    `json:"port" validate:"min=1,max=65535"`
    Database string `json:"database" validate:"required"`
}

func (c *Config) Validate() error {
    if c.Host == "" {
        return fmt.Errorf("host is required")
    }
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("port must be 1-65535, got %d", c.Port)
    }
    if c.Database == "" {
        return fmt.Errorf("database is required")
    }
    return nil
}

func (c *Config) ConnectionString() string {
    return fmt.Sprintf("postgres://%s:%d/%s", c.Host, c.Port, c.Database)
}
```

### Example 3: Composition Over Inheritance

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Logger struct {
    name string
}

func (l *Logger) Log(msg string) {
    fmt.Printf("[%s] %s\n", l.name, msg)
}

type Service struct {
    Logger  // Embedded - Log method is promoted
    Reader  // Embedded interface
    Config  Config
}

service := Service{
    Logger: Logger{name: "myservice"},
    Config: Config{Host: "localhost", Port: 8080},
}

service.Log("Starting service")  // Works due to promoted method
```

---

## Exercises

### Exercise 1: Shape Types
Create Rectangle and Circle types with Area() and Perimeter() methods.

### Exercise 2: Constructor Functions
Build a User struct with email validation in constructor.

### Exercise 3: Struct Embedding
Create Manager and Developer types that embed Employee.

### Exercise 4: JSON Tags
Create a struct that marshals/unmarshals from JSON with custom field names.

---

## Summary

- **Structs**: Group related data
- **Methods**: Functions on types
- **Value receiver**: Immutable operations
- **Pointer receiver**: Mutations, large structs
- **Embedding**: Composition, code reuse
- **Tags**: Metadata for marshaling and validation
- **Constructors**: Ensure valid initialization

This completes the foundation! Next: Module 2 - Advanced Concepts.
