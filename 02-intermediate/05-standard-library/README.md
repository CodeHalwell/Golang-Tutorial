# Chapter 5: Standard Library Deep Dive

## Learning Objectives

1. Master essential string operations (`strings`, `strconv`)
2. Use the `fmt` package correctly for formatting
3. Work with time and dates (`time` package)
4. Understand `sync` package for synchronization
5. Use standard library data structures (`container/`)
6. Work with sort and search (`sort`)
7. Master regular expressions (`regexp`)
8. Use hashing and cryptography packages
9. Understand `encoding` packages
10. Learn about `log` and logging best practices

---

## Strings Package

### Common String Operations

```go
import "strings"

// Check if string contains substring
if strings.Contains("hello world", "world") {
    fmt.Println("Contains 'world'")
}

// Case conversion
lower := strings.ToLower("HELLO")      // "hello"
upper := strings.ToUpper("hello")      // "HELLO"

// Trim whitespace
trimmed := strings.TrimSpace("  hello  ")  // "hello"
trimmed = strings.Trim("##hello##", "#")  // "hello"

// Split and Join
parts := strings.Split("a,b,c", ",")  // []string{"a", "b", "c"}
joined := strings.Join([]string{"a", "b", "c"}, ",")  // "a,b,c"

// Replace
replaced := strings.Replace("hello hello", "hello", "hi", 1)  // "hi hello"
replaced = strings.ReplaceAll("hello hello", "hello", "hi")   // "hi hi"

// HasPrefix and HasSuffix
strings.HasPrefix("hello.txt", "hello")  // true
strings.HasSuffix("hello.txt", ".txt")   // true

// Index operations
idx := strings.Index("hello world", "world")      // 6
idx = strings.LastIndex("hello hello", "hello")   // 6

// Count occurrences
count := strings.Count("hello hello", "hello")    // 2

// Field splitting
fields := strings.Fields("one   two   three")  // []string{"one", "two", "three"}
```

### String Builder for Efficiency

```go
import "strings"

// String concatenation
var sb strings.Builder

for i := 0; i < 1000; i++ {
    sb.WriteString(fmt.Sprintf("Line %d\n", i))
}

result := sb.String()

// Why? Building strings in a loop with + is O(n²)
// StringBuilder is O(n)
```

---

## Strconv Package

### String to Number Conversions

```go
import "strconv"

// Parse integer
i, err := strconv.Atoi("123")       // int, simpler
i64, err := strconv.ParseInt("123", 10, 64)  // int64, more control
i32, err := strconv.ParseInt("123", 10, 32)

// Parse float
f, err := strconv.ParseFloat("3.14", 64)

// Parse boolean
b, err := strconv.ParseBool("true")  // true

// Number to string
s := strconv.Itoa(123)               // "123"
s = strconv.FormatInt(123, 10)       // "123" in base 10
s = strconv.FormatFloat(3.14, 'f', 2, 64)  // "3.14"
```

### Error Handling in Conversion

```go
import (
    "log"
    "strconv"
)

// With explicit error handling
i, err := strconv.Atoi("not_a_number")
if err != nil {
    log.Printf("Conversion error: %v", err)
}
```

---

## Fmt Package

### String Formatting Verbs

```go
import "fmt"

// General formatting
fmt.Printf("%v\n", value)      // Default format
fmt.Printf("%#v\n", value)     // Show type info
fmt.Printf("%T\n", value)      // Type only

// Strings
fmt.Printf("%s\n", "hello")    // String
fmt.Printf("%q\n", "hello")    // Quoted string
fmt.Printf("%x\n", "hello")    // Hex encoding

// Integers
fmt.Printf("%d\n", 42)         // Decimal
fmt.Printf("%o\n", 42)         // Octal
fmt.Printf("%x\n", 42)         // Hex
fmt.Printf("%b\n", 42)         // Binary
fmt.Printf("%c\n", 65)         // Character (ASCII 65 = 'A')

// Floats
fmt.Printf("%f\n", 3.14)       // Float
fmt.Printf("%e\n", 3.14)       // Scientific
fmt.Printf("%.2f\n", 3.14159)  // 2 decimal places

// Pointers
fmt.Printf("%p\n", &value)     // Pointer address

// Width and padding
fmt.Printf("%10s\n", "hello")  // Right-aligned, width 10
fmt.Printf("%-10s\n", "hello") // Left-aligned, width 10
fmt.Printf("%010d\n", 42)      // Zero-padded to width 10
```

---

## Time Package

### Basic Time Operations

```go
import (
    "fmt"
    "time"
)

// Current time
now := time.Now()
fmt.Println(now)  // 2024-01-15 14:30:45.123456789 +0000 UTC

// Time components
year := now.Year()
month := now.Month()
day := now.Day()
hour := now.Hour()
minute := now.Minute()
second := now.Second()

// Creating specific time
t := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)

// Time parsing
t, err := time.Parse("2006-01-02", "2024-01-15")
t, err = time.Parse("2006-01-02 15:04:05", "2024-01-15 14:30:45")

// Time formatting
formatted := now.Format("2006-01-02 15:04:05")
formatted = now.Format(time.RFC3339)
formatted = now.Format(time.RFC1123)

// Date only
dateOnly := now.Format("2006-01-02")  // "2024-01-15"
timeOnly := now.Format("15:04:05")    // "14:30:45"
```

### Duration and Intervals

```go
import "time"

// Create duration
d := 5 * time.Second
d = 2 * time.Hour
d = 3 * time.Minute + 30 * time.Second

// Work with duration
fmt.Println(d.Seconds())       // Total seconds as float64
fmt.Println(d.Minutes())       // Total minutes
fmt.Println(d.Hours())         // Total hours

// Time arithmetic
now := time.Now()
tomorrow := now.Add(24 * time.Hour)
yesterday := now.Add(-24 * time.Hour)

// Duration between times
elapsed := time.Since(startTime)     // Duration since startTime
until := time.Until(deadline)        // Duration until deadline

// Compare times
if now.Before(deadline) {
    fmt.Println("Not yet!")
}

if now.After(startTime) {
    fmt.Println("Started")
}
```

### Timers and Tickers

```go
import (
    "fmt"
    "time"
)

// Timer (one-shot)
timer := time.NewTimer(3 * time.Second)
<-timer.C  // Wait for timer to fire
fmt.Println("Timer fired")

// Ticker (repeating)
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

for i := 0; i < 5; i++ {
    <-ticker.C
    fmt.Println("Tick")
}

// After (convenience function)
<-time.After(2 * time.Second)
fmt.Println("2 seconds elapsed")
```

---

## Sync Package

### Mutex for Mutual Exclusion

```go
import "sync"

type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Read() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

### RWMutex for Read-Write Locks

```go
import "sync"

type Cache struct {
    mu    sync.RWMutex
    data  map[string]interface{}
}

func (c *Cache) Get(key string) interface{} {
    c.mu.RLock()  // Read lock (multiple readers OK)
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()   // Write lock (exclusive)
    defer c.mu.Unlock()
    c.data[key] = value
}
```

### WaitGroup for Goroutine Synchronization

```go
import (
    "sync"
    "time"
)

var wg sync.WaitGroup

// Add 3 goroutines to wait for
wg.Add(3)

for i := 0; i < 3; i++ {
    go func(id int) {
        defer wg.Done()
        // Do work
    }(i)
}

// Block until all Done() called
wg.Wait()
```

### Once for One-Time Initialization

```go
import "sync"

var once sync.Once
var instance *Singleton

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}

// GetInstance() always returns same instance
// Guaranteed to happen once, even with concurrent calls
```

### Atomic Operations

```go
import (
    "sync/atomic"
)

var counter int64

// Atomic increment (no lock needed)
atomic.AddInt64(&counter, 1)

// Atomic load
value := atomic.LoadInt64(&counter)

// Atomic store
atomic.StoreInt64(&counter, 100)

// Compare and swap
atomic.CompareAndSwapInt64(&counter, 100, 101)
```

---

## Sort Package

### Basic Sorting

```go
import "sort"

// Sort integers
ints := []int{3, 1, 4, 1, 5, 9}
sort.Ints(ints)  // [1, 1, 3, 4, 5, 9]

// Sort strings
strings := []string{"banana", "apple", "cherry"}
sort.Strings(strings)  // [apple, banana, cherry]

// Sort floats
floats := []float64{3.14, 2.71, 1.41}
sort.Float64s(floats)
```

### Custom Sorting

```go
import "sort"

type Person struct {
    Name string
    Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Sort
people := []Person{...}
sort.Sort(ByAge(people))

// Or with SortFunc (Go 1.21+)
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
```

### Searching

```go
import "sort"

nums := []int{1, 3, 4, 6, 8, 9}

// Binary search
idx := sort.SearchInts(nums, 6)  // returns 3
found := idx < len(nums) && nums[idx] == 6

// Custom search
idx = sort.Search(len(people), func(i int) bool {
    return people[i].Age >= 30
})
```

---

## Regexp Package

### Regular Expression Matching

```go
import "regexp"

re := regexp.MustCompile(`\d+`)  // One or more digits

// Check if matches
if re.MatchString("abc123def") {
    fmt.Println("Contains digits")
}

// Find first match
match := re.FindString("abc123def456")  // "123"

// Find all matches
matches := re.FindAllString("abc123def456", -1)  // ["123", "456"]

// Find with groups
re = regexp.MustCompile(`(\d+)-(\d+)-(\d+)`)
groups := re.FindStringSubmatch("2024-01-15")
// groups[0] = "2024-01-15"
// groups[1] = "2024"
// groups[2] = "01"
// groups[3] = "15"
```

### String Replacement

```go
import "regexp"

re := regexp.MustCompile(`\d+`)

// Replace first
result := re.ReplaceAllString("abc123def456", "X")  // "abcXdefX"

// Replace with function
result = re.ReplaceAllStringFunc("abc123def456", func(s string) string {
    return "[" + s + "]"
})  // "abc[123]def[456]"
```

---

## Hash and Crypto Packages

### MD5, SHA1, SHA256

```go
import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "fmt"
)

data := []byte("hello")

// MD5 (not for security)
md5sum := md5.Sum(data)
fmt.Printf("%x\n", md5sum)

// SHA1 (deprecated for security)
sha1sum := sha1.Sum(data)

// SHA256 (preferred)
sha256sum := sha256.Sum256(data)
fmt.Printf("%x\n", sha256sum)

// Or using hash writer interface
h := sha256.New()
h.Write(data)
hash := h.Sum(nil)
```

---

## Encoding Packages

### Base64

```go
import (
    "encoding/base64"
    "fmt"
)

// Encode
encoded := base64.StdEncoding.EncodeToString([]byte("hello"))
fmt.Println(encoded)  // "aGVsbG8="

// Decode
decoded, _ := base64.StdEncoding.DecodeString("aGVsbG8=")
fmt.Println(string(decoded))  // "hello"
```

### Hex

```go
import (
    "encoding/hex"
    "fmt"
)

// Encode
encoded := hex.EncodeToString([]byte("hello"))
fmt.Println(encoded)  // "68656c6c6f"

// Decode
decoded, _ := hex.DecodeString("68656c6c6f")
fmt.Println(string(decoded))  // "hello"
```

---

## Logging Best Practices

### Standard Library Logging

```go
import (
    "log"
    "os"
)

// Basic logging
log.Println("Info message")
log.Printf("Formatted: %s\n", "value")
log.Fatalf("Fatal error: %v", err)  // Prints and exits

// Custom logger
file, _ := os.Create("app.log")
logger := log.New(file, "[INFO] ", log.LstdFlags)
logger.Println("Message")

// Flags
log.SetFlags(log.LstdFlags | log.Lshortfile)  // Include filename/line
```

---

## Practical Examples

### Example 1: Parse Command-Line Arguments

```go
import (
    "flag"
    "fmt"
)

func main() {
    name := flag.String("name", "World", "Name to greet")
    count := flag.Int("count", 1, "Number of greetings")
    flag.Parse()

    for i := 0; i < *count; i++ {
        fmt.Printf("Hello, %s!\n", *name)
    }
}

// Usage: ./program -name Alice -count 3
```

### Example 2: CSV Processing

```go
import (
    "encoding/csv"
    "os"
    "strconv"
    "strings"
)

type Record struct {
    Name string
    Age  int
}

func readCSV(filename string) ([]Record, error) {
    file, _ := os.Open(filename)
    reader := csv.NewReader(file)

    records := []Record{}
    lines, _ := reader.ReadAll()

    for i := 1; i < len(lines); i++ {
        age, _ := strconv.Atoi(lines[i][1])
        records = append(records, Record{
            Name: lines[i][0],
            Age:  age,
        })
    }

    return records, nil
}
```

### Example 3: Rate Limiting with Ticker

```go
import (
    "fmt"
    "time"
)

func rateLimitedWork() {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for i := 0; i < 10; i++ {
        <-ticker.C
        fmt.Printf("Work %d\n", i)
    }
}
```

---

## Summary

- **strings**: Split, join, replace, trim, case conversion
- **strconv**: Parse and format numbers
- **fmt**: Formatting with verbs and width
- **time**: Dates, durations, timers, tickers
- **sync**: Mutex, RWMutex, WaitGroup, Once, Atomic
- **sort**: Sort slices, custom sorting, search
- **regexp**: Pattern matching, replacement
- **crypto**: Hashing (MD5, SHA1, SHA256)
- **encoding**: Base64, Hex, JSON, XML, CSV
- **log**: Standard logging

This chapter covered the most commonly used standard library packages. Next: Module 3 - Concurrency Patterns.
