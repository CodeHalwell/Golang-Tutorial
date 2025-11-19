package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Custom error type for validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}

// Custom error type for file operations
type FileError struct {
	Operation string
	Filename  string
	Err       error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s failed for file '%s': %v", e.Operation, e.Filename, e.Err)
}

// Unwrap allows errors.Is and errors.As to work with FileError
func (e *FileError) Unwrap() error {
	return e.Err
}

// Example 1: Simple error checking
func readFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		// Wrap the error with additional context
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

// Example 2: Custom error type
func validateEmail(email string) error {
	if !strings.Contains(email, "@") {
		return &ValidationError{
			Field:   "email",
			Message: "must contain @",
		}
	}
	return nil
}

// Example 3: Sentinel error
var ErrUserNotFound = errors.New("user not found")

func findUser(id string) (string, error) {
	if id == "" {
		return "", ErrUserNotFound
	}
	return "User: " + id, nil
}

// Example 4: Error wrapping with context
func processData(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("processData: input is empty")
	}

	if len(input) < 3 {
		return "", fmt.Errorf("processData: input too short (got %d, need 3)", len(input))
	}

	return strings.ToUpper(input), nil
}

// Example 5: Multiple return values for error handling
func divideWithValidation(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Example 6: Handling different error types
func handleUserOperation(userID string) error {
	// Try to find user
	user, err := findUser(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// Specific handling for not found
			return fmt.Errorf("user lookup failed: %w", err)
		}
		// Generic error
		return fmt.Errorf("unexpected error: %w", err)
	}

	fmt.Println("Found:", user)
	return nil
}

// Example 7: Panic recovery (rare, use with caution)
func safeDivide(a, b float64) (result float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("division panicked: %v", r)
		}
	}()

	if b == 0 {
		panic("division by zero")
	}

	return a / b, nil
}

// Example 8: Deferred cleanup with error handling
func processFile(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("open failed: %w", err)
	}

	// In real code with writable files:
	// defer func() {
	//     if err := file.Close(); err != nil {
	//         log.Printf("close error: %v", err)
	//     }
	// }()

	fmt.Printf("File contents (%d bytes)\n", len(file))
	return nil
}

// Example 9: Error checking in loops
func processMultipleFiles(filenames []string) error {
	for _, filename := range filenames {
		if err := processFile(filename); err != nil {
			return fmt.Errorf("processing %s: %w", filename, err)
		}
	}
	return nil
}

// Example 10: Type assertion for custom errors
func handleFileError(err error) {
	var fileErr *FileError
	if errors.As(err, &fileErr) {
		log.Printf("File operation '%s' failed on file '%s'",
			fileErr.Operation, fileErr.Filename)
		return
	}

	// Not a FileError
	log.Printf("Unknown error: %v", err)
}

func main() {
	fmt.Println("=== Error Handling Examples ===\n")

	// Example 1: Simple error
	fmt.Println("Example 1: Simple error handling")
	if _, err := readFile("/nonexistent/file.txt"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 2: Custom error type
	fmt.Println("\nExample 2: Custom validation error")
	if err := validateEmail("invalid-email"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if err := validateEmail("valid@example.com"); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Email is valid")
	}

	// Example 3: Sentinel error
	fmt.Println("\nExample 3: Sentinel error")
	if _, err := findUser(""); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			fmt.Println("User not found (sentinel error match)")
		}
	}

	// Example 4: Error wrapping
	fmt.Println("\nExample 4: Error wrapping")
	if _, err := processData("ok"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 5: Division with error
	fmt.Println("\nExample 5: Operation with error")
	result, err := divideWithValidation(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 2 = %f\n", result)
	}

	result, err = divideWithValidation(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 6: Handling different errors
	fmt.Println("\nExample 6: Handling different error types")
	if err := handleUserOperation("test"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 7: Safe divide with panic recovery
	fmt.Println("\nExample 7: Panic recovery")
	result, err = safeDivide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 8: File processing
	fmt.Println("\nExample 8: File processing")
	// Create a test file
	testFile := "/tmp/test.txt"
	_ = ioutil.WriteFile(testFile, []byte("test content"), 0644)
	if err := processFile(testFile); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Example 9: Errors.As with custom types
	fmt.Println("\nExample 9: Using errors.As")
	fileErr := &FileError{
		Operation: "read",
		Filename:  "data.txt",
		Err:       errors.New("permission denied"),
	}
	handleFileError(fileErr)

	fmt.Println("\n=== Best Practices ===")
	fmt.Println("1. Check errors immediately")
	fmt.Println("2. Wrap errors with context using %w")
	fmt.Println("3. Use sentinel errors (var ErrXxx = errors.New())")
	fmt.Println("4. Create custom error types for domain-specific errors")
	fmt.Println("5. Use errors.Is and errors.As for comparison")
	fmt.Println("6. Don't ignore error return values")
	fmt.Println("7. Log or return errors, don't silently ignore them")
}
