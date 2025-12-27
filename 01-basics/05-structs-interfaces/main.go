package main

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

// ===== EXAMPLE 1: Basic Struct Declaration and Initialization =====

type Person struct {
	Name  string
	Age   int
	Email string
}

func example1() {
	fmt.Println("\n=== EXAMPLE 1: Struct Declaration & Initialization ===")

	// Method 1: Zero values
	var p1 Person
	fmt.Printf("Zero values: %+v\n", p1) // {Name: Age:0 Email:}

	// Method 2: Named fields
	p2 := Person{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}
	fmt.Printf("Named fields: %+v\n", p2)

	// Method 3: Positional (order matters!)
	p3 := Person{"Bob", 25, "bob@example.com"}
	fmt.Printf("Positional: %+v\n", p3)

	// Method 4: Partial initialization
	p4 := Person{Name: "Charlie"} // Age and Email are zero values
	fmt.Printf("Partial: %+v\n", p4)

	// Accessing fields
	fmt.Printf("Name: %s, Age: %d\n", p2.Name, p2.Age)

	// Modifying fields
	p2.Age = 31
	fmt.Printf("After birthday: %+v\n", p2)
}

// ===== EXAMPLE 2: Methods with Value Receiver =====

func (p Person) Describe() string {
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

func (p Person) IsAdult() bool {
	return p.Age >= 18
}

func example2() {
	fmt.Println("\n=== EXAMPLE 2: Methods with Value Receiver ===")

	p := Person{Name: "Diana", Age: 22}

	// Calling methods on value
	fmt.Println(p.Describe())
	fmt.Printf("Is adult: %v\n", p.IsAdult())

	// Value receiver gets a COPY - modifications don't affect original
	fmt.Printf("Before: %+v\n", p)
}

// ===== EXAMPLE 3: Methods with Pointer Receiver =====

func (p *Person) Birthday() {
	p.Age++ // Modifies the original struct
}

func (p *Person) SetEmail(email string) {
	if !strings.Contains(email, "@") {
		fmt.Println("Invalid email")
		return
	}
	p.Email = email
}

func example3() {
	fmt.Println("\n=== EXAMPLE 3: Methods with Pointer Receiver ===")

	p := Person{Name: "Eve", Age: 20}
	fmt.Printf("Before birthday: Age = %d\n", p.Age)

	p.Birthday()
	fmt.Printf("After birthday: Age = %d\n", p.Age)

	// Go automatically dereferences for method calls
	p.SetEmail("eve@example.com")
	fmt.Printf("Email set: %s\n", p.Email)

	// Pointer variable also works
	pPtr := &Person{Name: "Frank", Age: 30}
	pPtr.Birthday()
	fmt.Printf("Pointer variable: %+v\n", pPtr)
}

// ===== EXAMPLE 4: Struct Embedding (Composition) =====

type Animal struct {
	Name string
	Age  int
}

type Dog struct {
	Animal // Embedded struct
	Breed  string
}

type Cat struct {
	Animal
	Color string
}

func (a *Animal) Speak() {
	fmt.Printf("%s makes a sound\n", a.Name)
}

func (d *Dog) Speak() {
	fmt.Printf("%s barks\n", d.Name)
}

func example4() {
	fmt.Println("\n=== EXAMPLE 4: Struct Embedding (Composition) ===")

	// Create dog
	dog := Dog{
		Animal: Animal{Name: "Rex", Age: 5},
		Breed:  "Golden Retriever",
	}

	// Access promoted fields directly
	fmt.Printf("Dog name: %s (promoted field)\n", dog.Name)
	fmt.Printf("Dog breed: %s\n", dog.Breed)

	// Call promoted method
	dog.Speak()

	// Create cat
	cat := Cat{
		Animal: Animal{Name: "Whiskers", Age: 3},
		Color:  "Orange",
	}
	fmt.Printf("Cat: %s, Color: %s\n", cat.Name, cat.Color)
	cat.Speak() // Uses Animal's Speak, not overridden like Dog
}

// ===== EXAMPLE 5: Anonymous Structs =====

func example5() {
	fmt.Println("\n=== EXAMPLE 5: Anonymous Structs ===")

	// Create anonymous struct
	person := struct {
		Name string
		Age  int
	}{
		Name: "Grace",
		Age:  28,
	}

	fmt.Printf("Anonymous struct: %+v\n", person)
	fmt.Printf("Name: %s, Age: %d\n", person.Name, person.Age)

	// Anonymous structs in function returns
	result := func() struct {
		Success bool
		Data    string
		Error   string
	} {
		return struct {
			Success bool
			Data    string
			Error   string
		}{
			Success: true,
			Data:    "processed",
			Error:   "",
		}
	}()

	fmt.Printf("Result: Success=%v, Data=%s\n", result.Success, result.Data)
}

// ===== EXAMPLE 6: Struct Tags (JSON) =====

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email,omitempty"`
}

func example6() {
	fmt.Println("\n=== EXAMPLE 6: Struct Tags (JSON) ===")

	user := User{Name: "Henry", Age: 35, Email: "henry@example.com"}

	// Marshal to JSON
	jsonData, _ := json.Marshal(user)
	fmt.Printf("JSON: %s\n", jsonData)

	// Unmarshal from JSON
	jsonStr := `{"name":"Ivy","age":27}`
	var u User
	json.Unmarshal([]byte(jsonStr), &u)
	fmt.Printf("Unmarshaled: %+v\n", u)

	// Empty Email is omitted due to omitempty tag
	user2 := User{Name: "Jack", Age: 40}
	jsonData2, _ := json.Marshal(user2)
	fmt.Printf("Without email: %s\n", jsonData2)
}

// ===== EXAMPLE 7: Custom Struct Tags with Reflection =====

type Config struct {
	Host     string `config:"host" required:"true"`
	Port     int    `config:"port" min:"1" max:"65535"`
	Database string `config:"database" required:"true"`
}

func PrintTags(v interface{}) {
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		configTag := field.Tag.Get("config")
		requiredTag := field.Tag.Get("required")
		fmt.Printf("Field: %s, config tag: %s, required: %s\n",
			field.Name, configTag, requiredTag)
	}
}

func example7() {
	fmt.Println("\n=== EXAMPLE 7: Custom Struct Tags with Reflection ===")

	cfg := Config{Host: "localhost", Port: 8080, Database: "mydb"}
	PrintTags(cfg)
}

// ===== EXAMPLE 8: Struct Best Practices - Bank Account =====

type Account struct {
	AccountNumber string
	Owner         string
	Balance       float64
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive, got %f", amount)
	}
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive, got %f", amount)
	}
	if amount > a.Balance {
		return fmt.Errorf("insufficient funds: need %f, have %f", amount, a.Balance)
	}
	a.Balance -= amount
	return nil
}

func (a Account) String() string {
	return fmt.Sprintf("Account(%s): $%.2f", a.Owner, a.Balance)
}

func example8() {
	fmt.Println("\n=== EXAMPLE 8: Best Practice - Bank Account ===")

	account := Account{
		AccountNumber: "12345",
		Owner:         "Karen",
		Balance:       1000,
	}

	fmt.Println("Initial:", account)

	// Deposit
	if err := account.Deposit(500); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("After deposit:", account)
	}

	// Withdraw
	if err := account.Withdraw(300); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("After withdrawal:", account)
	}

	// Try invalid operations
	if err := account.Deposit(-100); err != nil {
		fmt.Println("Error (expected):", err)
	}

	if err := account.Withdraw(2000); err != nil {
		fmt.Println("Error (expected):", err)
	}
}

// ===== EXAMPLE 9: Struct Best Practices - Constructor Functions =====

type ValidatedUser struct {
	name  string // Private fields
	email string
	age   int
}

// Constructor with validation
func NewValidatedUser(name, email string, age int) (*ValidatedUser, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("invalid email format")
	}
	if age < 0 || age > 150 {
		return nil, fmt.Errorf("age must be between 0 and 150, got %d", age)
	}
	return &ValidatedUser{name: name, email: email, age: age}, nil
}

// Getter methods for private fields
func (u *ValidatedUser) Name() string  { return u.name }
func (u *ValidatedUser) Email() string { return u.email }
func (u *ValidatedUser) Age() int      { return u.age }

func (u *ValidatedUser) String() string {
	return fmt.Sprintf("User(%s, %s, %d)", u.name, u.email, u.age)
}

func example9() {
	fmt.Println("\n=== EXAMPLE 9: Best Practice - Constructor Functions ===")

	// Valid user
	user, err := NewValidatedUser("Leo", "leo@example.com", 28)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Created:", user)
	}

	// Invalid users
	_, err = NewValidatedUser("", "test@example.com", 30)
	fmt.Println("Error (empty name):", err)

	_, err = NewValidatedUser("Mia", "invalid-email", 25)
	fmt.Println("Error (invalid email):", err)

	_, err = NewValidatedUser("Noah", "noah@example.com", 200)
	fmt.Println("Error (invalid age):", err)
}

// ===== EXAMPLE 10: Struct Best Practices - Embedding for Code Reuse =====

type Named struct {
	Name string
}

func (n *Named) SetName(name string) {
	n.Name = name
}

func (n *Named) GetName() string {
	return n.Name
}

type Employee struct {
	Named  // Embedded - SetName and GetName are promoted
	Title  string
	Salary float64
}

type Manager struct {
	Employee // Can embed Employee which embeds Named
	Department string
}

func example10() {
	fmt.Println("\n=== EXAMPLE 10: Embedding for Code Reuse ===")

	emp := Employee{
		Named:  Named{Name: "Oliver"},
		Title:  "Engineer",
		Salary: 100000,
	}

	fmt.Printf("Employee: %s, Title: %s\n", emp.GetName(), emp.Title)

	// Promoted method works
	emp.SetName("Peter")
	fmt.Printf("After rename: %s\n", emp.GetName())

	// Manager embeds Employee
	mgr := Manager{
		Employee: Employee{
			Named:  Named{Name: "Quinn"},
			Title:  "Senior Engineer",
			Salary: 150000,
		},
		Department: "Engineering",
	}

	fmt.Printf("Manager: %s, Title: %s, Dept: %s\n",
		mgr.GetName(), mgr.Title, mgr.Department)
}

// ===== EXAMPLE 11: Struct with Multiple Tag Types =====

type Product struct {
	ID        int       `json:"id" db:"product_id" xml:"id"`
	Name      string    `json:"name" db:"product_name" xml:"name"`
	Price     float64   `json:"price" db:"price" xml:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at" xml:"created_at"`
}

func example11() {
	fmt.Println("\n=== EXAMPLE 11: Multiple Tag Types ===")

	product := Product{
		ID:        101,
		Name:      "Laptop",
		Price:     999.99,
		CreatedAt: time.Now(),
	}

	// JSON marshal
	jsonData, _ := json.Marshal(product)
	fmt.Printf("JSON: %s\n", jsonData)

	// Inspect tags
	t := reflect.TypeOf(product)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		dbTag := field.Tag.Get("db")
		fmt.Printf("%s: json=%s, db=%s\n", field.Name, jsonTag, dbTag)
	}
}

// ===== EXAMPLE 12: Struct Comparison and Equality =====

type Point struct {
	X float64
	Y float64
}

func example12() {
	fmt.Println("\n=== EXAMPLE 12: Struct Comparison ===")

	p1 := Point{X: 1.0, Y: 2.0}
	p2 := Point{X: 1.0, Y: 2.0}
	p3 := Point{X: 2.0, Y: 3.0}

	// Structs are comparable (if all fields are)
	fmt.Printf("p1 == p2: %v (same values)\n", p1 == p2)
	fmt.Printf("p1 == p3: %v (different values)\n", p1 == p3)

	// Different addresses but same value
	pPtr1 := &Point{X: 5.0, Y: 5.0}
	pPtr2 := &Point{X: 5.0, Y: 5.0}
	fmt.Printf("*pPtr1 == *pPtr2: %v (pointers point to same value)\n", *pPtr1 == *pPtr2)
	fmt.Printf("pPtr1 == pPtr2: %v (different pointers)\n", pPtr1 == pPtr2)
}

// ===== EXAMPLE 13: Method Sets and Receiver Types =====

type Counter struct {
	value int
}

// Value receiver
func (c Counter) Get() int {
	return c.value
}

// Pointer receiver
func (c *Counter) Increment() {
	c.value++
}

func example13() {
	fmt.Println("\n=== EXAMPLE 13: Method Sets and Receiver Types ===")

	c := Counter{value: 0}

	// Both work on value (Go auto-references for pointer methods)
	fmt.Printf("Initial value: %d\n", c.Get())

	c.Increment()
	fmt.Printf("After increment: %d\n", c.Get())

	// Also works with pointer
	cPtr := &Counter{value: 10}
	fmt.Printf("Pointer initial: %d\n", cPtr.Get())
	cPtr.Increment()
	fmt.Printf("Pointer after increment: %d\n", cPtr.Get())
}

// ===== EXAMPLE 14: Struct Interface Implementation =====

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func PrintShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func example14() {
	fmt.Println("\n=== EXAMPLE 14: Struct Interface Implementation ===")

	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 3}

	fmt.Println("Rectangle:")
	PrintShapeInfo(rect)

	fmt.Println("Circle:")
	PrintShapeInfo(circle)
}

// ===== EXAMPLE 15: Nil Receiver in Methods =====

type StringSet struct {
	values map[string]bool
}

func (s *StringSet) Add(value string) {
	if s == nil { // Safe nil check
		fmt.Println("Cannot add to nil StringSet")
		return
	}
	if s.values == nil {
		s.values = make(map[string]bool)
	}
	s.values[value] = true
}

func example15() {
	fmt.Println("\n=== EXAMPLE 15: Nil Receiver Handling ===")

	set := &StringSet{}
	set.Add("apple")
	set.Add("banana")
	fmt.Printf("Set: %v\n", set.values)

	// Nil receiver
	var nilSet *StringSet
	nilSet.Add("orange") // Safely handled
}

// ===== EXAMPLE 16: Struct with Interface Fields =====

type Logger interface {
	Log(message string)
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Log(message string) {
	fmt.Printf("[LOG] %s\n", message)
}

type Service struct {
	name   string
	logger Logger
}

func (s *Service) DoWork() {
	s.logger.Log(fmt.Sprintf("Service %s is working", s.name))
}

func example16() {
	fmt.Println("\n=== EXAMPLE 16: Struct with Interface Fields ===")

	service := &Service{
		name:   "PaymentService",
		logger: ConsoleLogger{},
	}

	service.DoWork()
}

// ===== EXAMPLE 17: Large Struct Best Practices =====

type UserProfile struct {
	ID           int64
	Username     string
	Email        string
	FullName     string
	Age          int
	Location     string
	Bio          string
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsActive     bool
	IsPremium    bool
	LastLoginAt  time.Time
	Preferences  map[string]string
}

// Constructor for large struct
func NewUserProfile(id int64, username, email string) *UserProfile {
	return &UserProfile{
		ID:          id,
		Username:    username,
		Email:       email,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
		Preferences: make(map[string]string),
	}
}

func example17() {
	fmt.Println("\n=== EXAMPLE 17: Large Struct Best Practices ===")

	user := NewUserProfile(1, "romeo", "romeo@example.com")
	user.FullName = "Romeo Smith"
	user.Age = 30

	fmt.Printf("User: %s (%s) - Created: %v\n", user.FullName, user.Email, user.CreatedAt)
}

// ===== MAIN FUNCTION =====

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Chapter 5: Structs and Methods - Comprehensive Examples   ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	example1()
	example2()
	example3()
	example4()
	example5()
	example6()
	example7()
	example8()
	example9()
	example10()
	example11()
	example12()
	example13()
	example14()
	example15()
	example16()
	example17()

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Key Takeaways:                                           ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Println("║ 1. Value receivers for immutable operations                 ║")
	fmt.Println("║ 2. Pointer receivers for mutations and large structs        ║")
	fmt.Println("║ 3. Embedding enables composition and code reuse             ║")
	fmt.Println("║ 4. Struct tags provide metadata for marshaling              ║")
	fmt.Println("║ 5. Constructor functions ensure valid initialization        ║")
	fmt.Println("║ 6. Go uses composition over inheritance                     ║")
	fmt.Println("║ 7. Methods make types expressive and powerful               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
