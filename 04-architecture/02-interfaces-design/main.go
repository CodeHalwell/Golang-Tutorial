package main

import (
	"fmt"
	"io"
)

// === INTERFACE DESIGN PATTERNS ===

// Example 1: Small, focused interfaces (idiomatic Go)
//
// Go philosophy: Accept interfaces, return concrete types
// This is more testable and flexible than large interfaces

// Reader is a small, focused interface (from io package)
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer is another small interface
type Writer interface {
	Write(p []byte) (n int, err error)
}

// By combining small interfaces, we get composability
type ReadWriter interface {
	Reader
	Writer
}

// Closer manages resource cleanup
type Closer interface {
	Close() error
}

// ReadWriteCloser combines three interfaces
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// Example 2: Interface nil vs typed nil (critical Go concept!)
//
// This is a common gotcha that trips up many Go developers

type FileError struct {
	path string
	err  error
}

func (fe *FileError) Error() string {
	return fmt.Sprintf("file error: %s: %v", fe.path, fe.err)
}

// This function returns an interface{}
func getError(condition bool) error {
	var err *FileError // nil pointer
	if condition {
		err = &FileError{path: "test.txt", err: io.EOF}
	}
	return err // Returns interface{} with underlying nil
}

// Example 3: Design for testability using interfaces
//
// Instead of concrete dependencies, use interfaces

// Database is an interface for database operations
type Database interface {
	Save(key string, value []byte) error
	Load(key string) ([]byte, error)
	Delete(key string) error
}

// UserService depends on Database interface, not concrete implementation
type UserService struct {
	db Database
}

func (us *UserService) SaveUser(id string, data []byte) error {
	return us.db.Save("user:"+id, data)
}

func (us *UserService) GetUser(id string) ([]byte, error) {
	return us.db.Load("user:" + id)
}

// MockDatabase for testing
type MockDatabase struct {
	data map[string][]byte
}

func (md *MockDatabase) Save(key string, value []byte) error {
	md.data[key] = value
	return nil
}

func (md *MockDatabase) Load(key string) ([]byte, error) {
	if val, ok := md.data[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("not found")
}

func (md *MockDatabase) Delete(key string) error {
	delete(md.data, key)
	return nil
}

// Example 4: Composition over inheritance
//
// Go doesn't have inheritance; instead use composition

// Shape is a basic interface
type Shape interface {
	Area() float64
}

// Circle implements Shape
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

// ColoredShape composes Shape with color
type ColoredShape struct {
	shape Shape
	color string
}

func (cs ColoredShape) Area() float64 {
	return cs.shape.Area()
}

func (cs ColoredShape) Color() string {
	return cs.color
}

// Example 5: Strategy pattern using interfaces
//
// Swap algorithm implementations at runtime

type SortStrategy interface {
	Sort(items []int)
}

type QuickSort struct{}

func (qs QuickSort) Sort(items []int) {
	fmt.Println("Quick sort implementation")
}

type MergeSort struct{}

func (ms MergeSort) Sort(items []int) {
	fmt.Println("Merge sort implementation")
}

// Sorter uses the strategy interface
type Sorter struct {
	strategy SortStrategy
}

func (s *Sorter) SetStrategy(st SortStrategy) {
	s.strategy = st
}

func (s *Sorter) Sort(items []int) {
	s.strategy.Sort(items)
}

// Example 6: Dependency injection pattern
//
// Inject dependencies rather than creating them

type Logger interface {
	Log(msg string)
}

type ConsoleLogger struct{}

func (cl ConsoleLogger) Log(msg string) {
	fmt.Printf("[LOG] %s\n", msg)
}

type Application struct {
	logger Logger
}

// Constructor for dependency injection
func NewApplication(logger Logger) *Application {
	return &Application{logger: logger}
}

func (app *Application) Run() {
	app.logger.Log("Application started")
}

// Example 7: Observer pattern with interfaces
//
// Decouple event producers from consumers

type EventListener interface {
	OnEvent(event string, data interface{})
}

type EventEmitter struct {
	listeners []EventListener
}

func (ee *EventEmitter) Subscribe(listener EventListener) {
	ee.listeners = append(ee.listeners, listener)
}

func (ee *EventEmitter) Emit(event string, data interface{}) {
	for _, listener := range ee.listeners {
		listener.OnEvent(event, data)
	}
}

type ConsoleEventListener struct{}

func (cel ConsoleEventListener) OnEvent(event string, data interface{}) {
	fmt.Printf("Event: %s, Data: %v\n", event, data)
}

// Example 8: Adapter pattern
//
// Adapt incompatible interfaces

// We have this interface
type LegacyAPI interface {
	DoWork(input string) string
}

type OldAPI struct{}

func (oa OldAPI) DoWork(input string) string {
	return "Old: " + input
}

// But we want this newer interface
type ModernAPI interface {
	Process(input string) (string, error)
}

// Adapter bridges them
type APIAdapter struct {
	legacy LegacyAPI
}

func (aa APIAdapter) Process(input string) (string, error) {
	return aa.legacy.DoWork(input), nil
}

// Example 9: Interface embedding for composition
//
// Embed interfaces to create larger ones
// (Note: Writer is already defined earlier in this file)

type Flusher interface {
	Flush() error
}

type WriteFlusher interface {
	Writer
	Flusher
}

// Example 10: The empty interface{}
//
// Used sparingly for maximum flexibility

func PrintAnything(v interface{}) {
	switch val := v.(type) {
	case int:
		fmt.Printf("Integer: %d\n", val)
	case string:
		fmt.Printf("String: %s\n", val)
	case []int:
		fmt.Printf("Slice: %v\n", val)
	default:
		fmt.Printf("Unknown: %T\n", v)
	}
}

func main() {
	fmt.Println("=== Interface Design Patterns in Go ===\n")

	// Example 1: Small interfaces
	fmt.Println("1. Small, focused interfaces:")
	var reader Reader = nil // Can assign any Reader implementation
	_ = reader              // Would use for file operations

	// Example 2: Interface nil vs typed nil (CRITICAL)
	fmt.Println("\n2. Interface nil vs typed nil:")
	fmt.Println("This is a TRAP for Go developers!")

	err := getError(false)
	fmt.Printf("err is nil: %v\n", err == nil) // TRUE - actual nil
	if err != nil {
		fmt.Println(err)
	}

	err = getError(true)
	fmt.Printf("err is nil: %v\n", err == nil) // FALSE - interface wrapping nil
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example 3: Testability
	fmt.Println("\n3. Testable design with interfaces:")
	mockDB := &MockDatabase{data: make(map[string][]byte)}
	service := &UserService{db: mockDB}

	service.SaveUser("123", []byte("John Doe"))
	user, _ := service.GetUser("123")
	fmt.Printf("User: %s\n", string(user))

	// Example 4: Composition
	fmt.Println("\n4. Composition over inheritance:")
	circle := Circle{Radius: 5}
	colored := ColoredShape{shape: circle, color: "red"}
	fmt.Printf("Area: %.2f, Color: %s\n", colored.Area(), colored.Color())

	// Example 5: Strategy pattern
	fmt.Println("\n5. Strategy pattern:")
	sorter := &Sorter{}
	sorter.SetStrategy(QuickSort{})
	sorter.Sort([]int{3, 1, 2})
	sorter.SetStrategy(MergeSort{})
	sorter.Sort([]int{3, 1, 2})

	// Example 6: Dependency injection
	fmt.Println("\n6. Dependency injection:")
	app := NewApplication(ConsoleLogger{})
	app.Run()

	// Example 7: Observer pattern
	fmt.Println("\n7. Observer pattern:")
	emitter := &EventEmitter{}
	emitter.Subscribe(ConsoleEventListener{})
	emitter.Emit("user_registered", map[string]string{"user": "john"})

	// Example 8: Adapter pattern
	fmt.Println("\n8. Adapter pattern:")
	oldAPI := OldAPI{}
	adapter := APIAdapter{legacy: oldAPI}
	result, _ := adapter.Process("hello")
	fmt.Println("Result:", result)

	// Example 10: Empty interface
	fmt.Println("\n10. Type assertions:")
	PrintAnything(42)
	PrintAnything("hello")
	PrintAnything([]int{1, 2, 3})

	// === KEY TAKEAWAYS ===
	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("1. Use small interfaces (1-3 methods typically)")
	fmt.Println("2. Accept interfaces, return concrete types")
	fmt.Println("3. Watch out for interface nil vs typed nil!")
	fmt.Println("4. Design for testability from the start")
	fmt.Println("5. Use composition instead of inheritance")
	fmt.Println("6. Inject dependencies, don't create them")
	fmt.Println("7. Interfaces enable loose coupling")
}
