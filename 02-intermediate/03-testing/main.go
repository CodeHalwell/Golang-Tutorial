package main

import (
	"fmt"
	"strings"
)

// === TABLE-DRIVEN TESTS AND TESTING PATTERNS ===

// Calculator demonstrates functions we'll test
type Calculator struct{}

// Add returns the sum of two numbers
func (c *Calculator) Add(a, b int) int {
	return a + b
}

// Subtract returns the difference of two numbers
func (c *Calculator) Subtract(a, b int) int {
	return a - b
}

// Multiply returns the product of two numbers
func (c *Calculator) Multiply(a, b int) int {
	return a * b
}

// Divide returns the quotient and any error
func (c *Calculator) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// StringProcessor provides string manipulation functions
type StringProcessor struct{}

// ReverseString reverses a string
func (sp *StringProcessor) ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CountWords counts the number of words in a string
func (sp *StringProcessor) CountWords(s string) int {
	if strings.TrimSpace(s) == "" {
		return 0
	}
	return len(strings.Fields(s))
}

// ToUpperCase converts string to uppercase
func (sp *StringProcessor) ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// ToLowerCase converts string to lowercase
func (sp *StringProcessor) ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// Validate checks if a value is valid
type Validator interface {
	Validate() error
}

// Email represents an email address
type Email struct {
	address string
}

// NewEmail creates a new email
func NewEmail(address string) *Email {
	return &Email{address: address}
}

// Validate checks if email is valid
func (e *Email) Validate() error {
	if e.address == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if !strings.Contains(e.address, "@") {
		return fmt.Errorf("email must contain @")
	}
	parts := strings.Split(e.address, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("invalid email format")
	}
	if !strings.Contains(parts[1], ".") {
		return fmt.Errorf("email domain must contain .")
	}
	return nil
}

// User represents a user
type User struct {
	Name  string
	Email string
}

// Validate validates user data
func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if u.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	email := NewEmail(u.Email)
	return email.Validate()
}

// DataStore defines an interface for storing and retrieving data
type DataStore interface {
	Save(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

// MockDataStore is a test double for DataStore
type MockDataStore struct {
	data map[string]interface{}
}

// NewMockDataStore creates a new mock data store
func NewMockDataStore() *MockDataStore {
	return &MockDataStore{
		data: make(map[string]interface{}),
	}
}

// Save stores a value
func (m *MockDataStore) Save(key string, value interface{}) error {
	m.data[key] = value
	return nil
}

// Get retrieves a value
func (m *MockDataStore) Get(key string) (interface{}, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("key not found")
}

// Delete removes a value
func (m *MockDataStore) Delete(key string) error {
	delete(m.data, key)
	return nil
}

// Service uses a DataStore
type Service struct {
	store DataStore
}

// NewService creates a new service
func NewService(store DataStore) *Service {
	return &Service{store: store}
}

// SaveUser saves a user
func (s *Service) SaveUser(key string, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return s.store.Save(key, user)
}

// GetUser retrieves a user
func (s *Service) GetUser(key string) (*User, error) {
	data, err := s.store.Get(key)
	if err != nil {
		return nil, err
	}
	user, ok := data.(*User)
	if !ok {
		return nil, fmt.Errorf("value is not a user")
	}
	return user, nil
}

// Example test demonstrations
func demonstrateTableDrivenTests() {
	fmt.Println("=== Table-Driven Tests ===\n")

	// Example: Testing Add function
	testCases := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed signs", 5, -3, 2},
		{"zero", 0, 0, 0},
		{"large numbers", 1000000, 2000000, 3000000},
	}

	fmt.Println("Calculator Add tests:")
	calc := &Calculator{}
	for _, tc := range testCases {
		result := calc.Add(tc.a, tc.b)
		passed := result == tc.expected
		symbol := "✓"
		if !passed {
			symbol = "✗"
		}
		fmt.Printf("%s %s: %d + %d = %d\n", symbol, tc.name, tc.a, tc.b, result)
	}
}

func demonstrateStringProcessing() {
	fmt.Println("\n=== String Processing Tests ===\n")

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "hello", "olleh"},
		{"with spaces", "hello world", "dlrow olleh"},
		{"single char", "a", "a"},
		{"empty string", "", ""},
		{"palindrome", "racecar", "racecar"},
	}

	fmt.Println("StringProcessor ReverseString tests:")
	sp := &StringProcessor{}
	for _, tc := range testCases {
		result := sp.ReverseString(tc.input)
		passed := result == tc.expected
		symbol := "✓"
		if !passed {
			symbol = "✗"
		}
		fmt.Printf("%s %s: \"%s\" -> \"%s\"\n", symbol, tc.name, tc.input, result)
	}

	// Word counting tests
	fmt.Println("\nStringProcessor CountWords tests:")
	wordTests := []struct {
		name     string
		input    string
		expected int
	}{
		{"single word", "hello", 1},
		{"multiple words", "hello world foo", 3},
		{"with extra spaces", "hello  world", 2},
		{"empty string", "", 0},
		{"just spaces", "   ", 0},
	}

	for _, tc := range wordTests {
		result := sp.CountWords(tc.input)
		passed := result == tc.expected
		symbol := "✓"
		if !passed {
			symbol = "✗"
		}
		fmt.Printf("%s %s: \"%s\" -> %d words\n", symbol, tc.name, tc.input, result)
	}
}

func demonstrateErrorHandling() {
	fmt.Println("\n=== Error Handling Tests ===\n")

	fmt.Println("Validator tests:")
	testCases := []struct {
		name      string
		validator Validator
		shouldErr bool
	}{
		{"valid email", NewEmail("test@example.com"), false},
		{"missing @", NewEmail("testexample.com"), true},
		{"empty email", NewEmail(""), true},
		{"no domain", NewEmail("test@"), true},
		{"no local part", NewEmail("@example.com"), true},
		{"no dot in domain", NewEmail("test@example"), true},
		{"valid user", &User{Name: "John", Email: "john@example.com"}, false},
		{"invalid user email", &User{Name: "John", Email: "invalid"}, true},
		{"empty name", &User{Name: "", Email: "test@example.com"}, true},
	}

	for _, tc := range testCases {
		err := tc.validator.Validate()
		passed := (err != nil) == tc.shouldErr
		symbol := "✓"
		if !passed {
			symbol = "✗"
		}
		errStr := "nil"
		if err != nil {
			errStr = err.Error()
		}
		fmt.Printf("%s %s: %s\n", symbol, tc.name, errStr)
	}
}

func demonstrateMocking() {
	fmt.Println("\n=== Mocking Tests ===\n")

	// Create mock store
	store := NewMockDataStore()

	// Create service with mock
	service := NewService(store)

	// Test SaveUser
	user := &User{Name: "Alice", Email: "alice@example.com"}
	err := service.SaveUser("user1", user)
	fmt.Printf("SaveUser result: %v\n", err)

	// Test GetUser
	retrieved, err := service.GetUser("user1")
	fmt.Printf("GetUser result: %v (user: %v)\n", err, retrieved)

	// Test with invalid user
	invalidUser := &User{Name: "Bob", Email: "invalid"}
	err = service.SaveUser("user2", invalidUser)
	fmt.Printf("SaveUser with invalid user: %v\n", err)
}

func main() {
	fmt.Println("=== TESTING PATTERNS IN GO ===\n")

	demonstrateTableDrivenTests()
	demonstrateStringProcessing()
	demonstrateErrorHandling()
	demonstrateMocking()

	fmt.Println("\n=== TESTING BEST PRACTICES ===")
	fmt.Println("1. Use table-driven tests for multiple cases")
	fmt.Println("2. Test error cases, not just happy path")
	fmt.Println("3. Use test doubles (mocks, stubs, fakes)")
	fmt.Println("4. Organize tests in test files (*_test.go)")
	fmt.Println("5. Use subtests for related test cases")
	fmt.Println("6. Check both output and errors")
	fmt.Println("7. Test edge cases (empty, nil, zero, max values)")
	fmt.Println("8. Use meaningful test names")
}
