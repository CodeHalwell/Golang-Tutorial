package main

import (
	"encoding/json"
	"testing"
)

// ===== Tests for Struct Declaration and Initialization =====

func TestPersonZeroValues(t *testing.T) {
	var p Person
	if p.Name != "" || p.Age != 0 || p.Email != "" {
		t.Errorf("Expected zero values, got %+v", p)
	}
}

func TestPersonInitialization(t *testing.T) {
	tests := []struct {
		name     string
		person   Person
		expected Person
	}{
		{
			name:     "Named fields",
			person:   Person{Name: "Alice", Age: 30, Email: "alice@example.com"},
			expected: Person{Name: "Alice", Age: 30, Email: "alice@example.com"},
		},
		{
			name:     "Positional",
			person:   Person{"Bob", 25, "bob@example.com"},
			expected: Person{Name: "Bob", Age: 25, Email: "bob@example.com"},
		},
		{
			name:     "Partial",
			person:   Person{Name: "Charlie"},
			expected: Person{Name: "Charlie", Age: 0, Email: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.person != tt.expected {
				t.Errorf("Expected %+v, got %+v", tt.expected, tt.person)
			}
		})
	}
}

// ===== Tests for Value Receivers =====

func TestPersonDescribe(t *testing.T) {
	p := Person{Name: "Diana", Age: 22}
	expected := "Diana is 22 years old"
	result := p.Describe()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestPersonIsAdult(t *testing.T) {
	tests := []struct {
		name     string
		age      int
		expected bool
	}{
		{"Adult", 30, true},
		{"Just adult", 18, true},
		{"Not adult", 17, false},
		{"Child", 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Person{Age: tt.age}
			if p.IsAdult() != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, p.IsAdult())
			}
		})
	}
}

// ===== Tests for Pointer Receivers =====

func TestBirthday(t *testing.T) {
	p := Person{Age: 20}
	p.Birthday()

	if p.Age != 21 {
		t.Errorf("Expected age 21, got %d", p.Age)
	}
}

func TestSetEmail(t *testing.T) {
	p := Person{}

	p.SetEmail("test@example.com")
	if p.Email != "test@example.com" {
		t.Errorf("Expected email set, got %q", p.Email)
	}

	// Invalid email should not be set
	originalEmail := p.Email
	p.SetEmail("invalid-email")
	if p.Email != originalEmail {
		t.Errorf("Expected email unchanged, got %q", p.Email)
	}
}

// ===== Tests for Struct Embedding =====

func TestDogEmbedding(t *testing.T) {
	dog := Dog{
		Animal: Animal{Name: "Rex", Age: 5},
		Breed:  "Golden Retriever",
	}

	// Promoted field access
	if dog.Name != "Rex" {
		t.Errorf("Expected promoted Name field, got %q", dog.Name)
	}

	if dog.Breed != "Golden Retriever" {
		t.Errorf("Expected Breed, got %q", dog.Breed)
	}
}

func TestPromotedMethods(t *testing.T) {
	dog := Dog{
		Animal: Animal{Name: "Buddy", Age: 3},
		Breed:  "Labrador",
	}

	// Dog has Speak method (overridden)
	// Animal has Speak method too
	// This should call Dog's Speak
	dog.Speak() // Just verify it doesn't panic
}

// ===== Tests for Struct Tags and JSON =====

func TestJSONMarshal(t *testing.T) {
	user := User{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	var unmarshaled User
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if unmarshaled.Name != user.Name || unmarshaled.Age != user.Age {
		t.Errorf("Marshal/unmarshal failed: %+v != %+v", unmarshaled, user)
	}
}

func TestJSONOmitEmpty(t *testing.T) {
	user := User{Name: "Bob", Age: 25} // No email

	jsonData, _ := json.Marshal(user)
	jsonStr := string(jsonData)

	// "email" key should not be in JSON because it's empty and has omitempty
	if contains(jsonStr, `"email"`) {
		t.Errorf("Expected email field to be omitted, but found in: %s", jsonStr)
	}
}

// ===== Tests for Bank Account =====

func TestAccountDeposit(t *testing.T) {
	acc := Account{Balance: 100}

	err := acc.Deposit(50)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if acc.Balance != 150 {
		t.Errorf("Expected balance 150, got %f", acc.Balance)
	}
}

func TestAccountDepositInvalid(t *testing.T) {
	acc := Account{Balance: 100}

	tests := []float64{-50, 0, -0.01}

	for _, amount := range tests {
		err := acc.Deposit(amount)
		if err == nil {
			t.Errorf("Expected error for deposit %f, got nil", amount)
		}
	}
}

func TestAccountWithdraw(t *testing.T) {
	acc := Account{Balance: 100}

	err := acc.Withdraw(30)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if acc.Balance != 70 {
		t.Errorf("Expected balance 70, got %f", acc.Balance)
	}
}

func TestAccountWithdrawInsufficientFunds(t *testing.T) {
	acc := Account{Balance: 50}

	err := acc.Withdraw(100)
	if err == nil {
		t.Error("Expected insufficient funds error, got nil")
	}

	// Balance should be unchanged
	if acc.Balance != 50 {
		t.Errorf("Expected balance 50, got %f", acc.Balance)
	}
}

// ===== Tests for Constructor Functions =====

func TestNewValidatedUserValid(t *testing.T) {
	user, err := NewValidatedUser("Alice", "alice@example.com", 30)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if user.Name() != "Alice" || user.Email() != "alice@example.com" || user.Age() != 30 {
		t.Errorf("Constructor produced incorrect values: %+v", user)
	}
}

func TestNewValidatedUserInvalid(t *testing.T) {
	tests := []struct {
		name  string
		email string
		age   int
		valid bool
	}{
		{"", "test@example.com", 30, false},        // Empty name
		{"Alice", "invalid-email", 30, false},       // No @
		{"Bob", "bob@example.com", -1, false},       // Negative age
		{"Charlie", "charlie@example.com", 200, false}, // Age too high
		{"David", "david@example.com", 25, true},   // Valid
	}

	for _, tt := range tests {
		_, err := NewValidatedUser(tt.name, tt.email, tt.age)
		if tt.valid && err != nil {
			t.Errorf("Expected valid user for %q, got error: %v", tt.name, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("Expected error for invalid user %q", tt.name)
		}
	}
}

// ===== Tests for Embedding Code Reuse =====

func TestEmployeeEmbeddedMethods(t *testing.T) {
	emp := Employee{
		Named:  Named{Name: "Oliver"},
		Title:  "Engineer",
		Salary: 100000,
	}

	if emp.GetName() != "Oliver" {
		t.Errorf("Expected promoted GetName to work")
	}

	emp.SetName("Peter")
	if emp.Name != "Peter" {
		t.Errorf("Expected SetName to modify embedded field")
	}
}

// ===== Tests for Struct Comparison =====

func TestPointComparison(t *testing.T) {
	p1 := Point{X: 1.0, Y: 2.0}
	p2 := Point{X: 1.0, Y: 2.0}
	p3 := Point{X: 2.0, Y: 3.0}

	if p1 != p2 {
		t.Error("Expected equal points to be equal")
	}

	if p1 == p3 {
		t.Error("Expected different points to be different")
	}
}

// ===== Tests for Counter (Value vs Pointer Receivers) =====

func TestCounterValueReceiver(t *testing.T) {
	c := Counter{value: 5}
	result := c.Get()

	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}

func TestCounterPointerReceiver(t *testing.T) {
	c := Counter{value: 0}
	c.Increment()
	c.Increment()
	c.Increment()

	if c.Get() != 3 {
		t.Errorf("Expected 3, got %d", c.Get())
	}
}

// ===== Tests for Interface Implementation =====

func TestRectangleShape(t *testing.T) {
	rect := Rectangle{Width: 10, Height: 5}

	if rect.Area() != 50 {
		t.Errorf("Expected area 50, got %f", rect.Area())
	}

	if rect.Perimeter() != 30 {
		t.Errorf("Expected perimeter 30, got %f", rect.Perimeter())
	}
}

func TestCircleShape(t *testing.T) {
	circle := Circle{Radius: 3}

	expectedArea := 3.14159 * 3 * 3
	if circle.Area()-expectedArea > 0.01 {
		t.Errorf("Expected area ~%f, got %f", expectedArea, circle.Area())
	}
}

func TestShapeInterface(t *testing.T) {
	var shapes []Shape
	shapes = append(shapes, Rectangle{Width: 10, Height: 5})
	shapes = append(shapes, Circle{Radius: 3})

	if len(shapes) != 2 {
		t.Errorf("Expected 2 shapes, got %d", len(shapes))
	}

	// All shapes should have area > 0
	for _, shape := range shapes {
		if shape.Area() <= 0 {
			t.Errorf("Expected positive area, got %f", shape.Area())
		}
	}
}

// ===== Tests for Nil Receiver =====

func TestStringSetNilReceiver(t *testing.T) {
	var nilSet *StringSet
	nilSet.Add("test") // Should not panic
}

func TestStringSetValidReceiver(t *testing.T) {
	set := &StringSet{}
	set.Add("apple")

	if set.values == nil || !set.values["apple"] {
		t.Error("Expected apple to be in set")
	}
}

// ===== Benchmark Tests =====

func BenchmarkStructCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Person{
			Name:  "Alice",
			Age:   30,
			Email: "alice@example.com",
		}
	}
}

func BenchmarkMethodCall(b *testing.B) {
	p := Person{Name: "Alice", Age: 30}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Describe()
	}
}

func BenchmarkPointerMethodCall(b *testing.B) {
	p := &Person{Name: "Alice", Age: 30}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Birthday()
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	user := User{Name: "Alice", Age: 30, Email: "alice@example.com"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(user)
	}
}

// ===== Helper Functions =====

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
