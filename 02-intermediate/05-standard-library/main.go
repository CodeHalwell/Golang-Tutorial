package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ===== EXAMPLE 1: Strings Package =====

func example1_StringsPackage() {
	fmt.Println("\n=== EXAMPLE 1: Strings Package ===")

	// Contains
	if strings.Contains("hello world", "world") {
		fmt.Println("Contains 'world'")
	}

	// Case conversion
	fmt.Printf("Lower: %s\n", strings.ToLower("HELLO"))
	fmt.Printf("Upper: %s\n", strings.ToUpper("hello"))

	// Trim
	trimmed := strings.TrimSpace("  hello  ")
	fmt.Printf("Trimmed: '%s'\n", trimmed)

	// Split and Join
	parts := strings.Split("a,b,c", ",")
	fmt.Printf("Split: %v\n", parts)
	fmt.Printf("Join: %s\n", strings.Join(parts, "-"))

	// Replace
	replaced := strings.ReplaceAll("hello hello", "hello", "hi")
	fmt.Printf("Replaced: %s\n", replaced)

	// HasPrefix/HasSuffix
	fmt.Printf("HasPrefix: %v\n", strings.HasPrefix("hello.txt", "hello"))
	fmt.Printf("HasSuffix: %v\n", strings.HasSuffix("hello.txt", ".txt"))

	// Index
	fmt.Printf("Index: %d\n", strings.Index("hello world", "world"))

	// Count
	fmt.Printf("Count: %d\n", strings.Count("hello hello", "hello"))

	// Fields
	fields := strings.Fields("one   two   three")
	fmt.Printf("Fields: %v\n", fields)
}

// ===== EXAMPLE 2: String Builder =====

func example2_StringBuilder() {
	fmt.Println("\n=== EXAMPLE 2: String Builder ===")

	var sb strings.Builder

	// Build string efficiently
	for i := 1; i <= 3; i++ {
		sb.WriteString(fmt.Sprintf("Line %d\n", i))
	}

	result := sb.String()
	fmt.Printf("Result:\n%s", result)
	fmt.Printf("Length: %d\n", sb.Len())
}

// ===== EXAMPLE 3: Strconv Package =====

func example3_StrconvPackage() {
	fmt.Println("\n=== EXAMPLE 3: Strconv Package ===")

	// String to Int
	i, err := strconv.Atoi("123")
	if err == nil {
		fmt.Printf("Atoi: %d (type: %T)\n", i, i)
	}

	// String to Int64
	i64, _ := strconv.ParseInt("999", 10, 64)
	fmt.Printf("ParseInt: %d\n", i64)

	// String to Float
	f, _ := strconv.ParseFloat("3.14", 64)
	fmt.Printf("ParseFloat: %f\n", f)

	// String to Bool
	b, _ := strconv.ParseBool("true")
	fmt.Printf("ParseBool: %v\n", b)

	// Int to String
	s := strconv.Itoa(456)
	fmt.Printf("Itoa: %s (type: %T)\n", s, s)

	// Format Int
	s = strconv.FormatInt(789, 10)
	fmt.Printf("FormatInt: %s\n", s)

	// Format Float
	s = strconv.FormatFloat(2.71828, 'f', 2, 64)
	fmt.Printf("FormatFloat: %s\n", s)
}

// ===== EXAMPLE 4: Fmt Package =====

func example4_FmtPackage() {
	fmt.Println("\n=== EXAMPLE 4: Fmt Package ===")

	value := 42
	text := "hello"

	// General formatting
	fmt.Printf("%%v: %v\n", value)
	fmt.Printf("%%T: %T\n", value)

	// Integer formatting
	fmt.Printf("%%d: %d\n", value)
	fmt.Printf("%%x: %x\n", value)
	fmt.Printf("%%o: %o\n", value)
	fmt.Printf("%%b: %b\n", value)

	// String formatting
	fmt.Printf("%%s: %s\n", text)
	fmt.Printf("%%q: %q\n", text)
	fmt.Printf("%%v: %v\n", []byte(text))

	// Float formatting
	f := 3.14159
	fmt.Printf("%%f: %f\n", f)
	fmt.Printf("%%.2f: %.2f\n", f)
	fmt.Printf("%%e: %e\n", f)

	// Width and padding
	fmt.Printf("%%10s: '%10s'\n", text)
	fmt.Printf("%%-10s: '%-10s'\n", text)
	fmt.Printf("%%010d: '%010d'\n", 42)
}

// ===== EXAMPLE 5: Time Package - Basics =====

func example5_TimeBasics() {
	fmt.Println("\n=== EXAMPLE 5: Time Package - Basics ===")

	// Current time
	now := time.Now()
	fmt.Printf("Now: %v\n", now)

	// Time components
	fmt.Printf("Year: %d, Month: %d, Day: %d\n", now.Year(), now.Month(), now.Day())
	fmt.Printf("Hour: %d, Minute: %d, Second: %d\n", now.Hour(), now.Minute(), now.Second())

	// Create specific time
	t := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)
	fmt.Printf("Created: %v\n", t)

	// Parse time
	parsed, _ := time.Parse("2006-01-02", "2024-01-15")
	fmt.Printf("Parsed: %v\n", parsed)

	// Format time
	fmt.Printf("Formatted: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("RFC3339: %s\n", now.Format(time.RFC3339))

	// Weekday and Month names
	fmt.Printf("Weekday: %s\n", now.Weekday())
}

// ===== EXAMPLE 6: Time Package - Duration =====

func example6_TimeDuration() {
	fmt.Println("\n=== EXAMPLE 6: Time Package - Duration ===")

	// Create durations
	d := 5 * time.Second
	fmt.Printf("Duration: %v\n", d)
	fmt.Printf("Seconds: %.0f\n", d.Seconds())

	// Arithmetic
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	fmt.Printf("Tomorrow: %s\n", tomorrow.Format("2006-01-02"))

	// Time.Since and Until
	start := now.Add(-5 * time.Second)
	elapsed := time.Since(start)
	fmt.Printf("Elapsed: %v\n", elapsed)

	// Compare times
	if now.Before(tomorrow) {
		fmt.Println("Now is before tomorrow")
	}
}

// ===== EXAMPLE 7: Time Package - Timers and Tickers =====

func example7_TimersAndTickers() {
	fmt.Println("\n=== EXAMPLE 7: Time Timers and Tickers ===")

	// Timer (one-shot)
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C
	fmt.Println("Timer fired!")

	// Ticker (repeating)
	ticker := time.NewTicker(100 * time.Millisecond)
	count := 0
	for range ticker.C {
		count++
		if count >= 3 {
			break
		}
	}
	ticker.Stop()
	fmt.Printf("Ticker fired %d times\n", count)

	// After convenience
	<-time.After(50 * time.Millisecond)
	fmt.Println("After fired!")
}

// ===== EXAMPLE 8: Sync - Mutex =====

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

func example8_SyncMutex() {
	fmt.Println("\n=== EXAMPLE 8: Sync - Mutex ===")

	counter := &Counter{}

	// Sequential increments
	for i := 0; i < 5; i++ {
		counter.Increment()
	}

	fmt.Printf("Counter value: %d\n", counter.Read())

	// Concurrent increments
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("After concurrent increments: %d\n", counter.Read())
}

// ===== EXAMPLE 9: Sync - RWMutex =====

type Cache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func (c *Cache) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func example9_RWMutex() {
	fmt.Println("\n=== EXAMPLE 9: Sync - RWMutex ===")

	cache := &Cache{data: make(map[string]interface{})}

	// Write
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Read
	fmt.Printf("key1: %v\n", cache.Get("key1"))
	fmt.Printf("key2: %v\n", cache.Get("key2"))

	// Concurrent reads and writes
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Set(fmt.Sprintf("key%d", id), fmt.Sprintf("value%d", id))
		}(i)
	}

	wg.Wait()
	fmt.Println("Concurrent operations complete")
}

// ===== EXAMPLE 10: Sync - WaitGroup =====

func example10_WaitGroup() {
	fmt.Println("\n=== EXAMPLE 10: Sync - WaitGroup ===")

	var wg sync.WaitGroup

	// Add 3 goroutines
	wg.Add(3)

	for i := 1; i <= 3; i++ {
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d working\n", id)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("Goroutine %d done\n", id)
		}(i)
	}

	fmt.Println("Waiting for all goroutines...")
	wg.Wait()
	fmt.Println("All goroutines complete")
}

// ===== EXAMPLE 11: Sync - Once =====

type Singleton struct {
	Value string
}

var once sync.Once
var instance *Singleton

func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{Value: "initialized"}
	})
	return instance
}

func example11_Once() {
	fmt.Println("\n=== EXAMPLE 11: Sync - Once ===")

	// Multiple calls return same instance
	s1 := GetInstance()
	s2 := GetInstance()

	fmt.Printf("s1 == s2: %v\n", s1 == s2)
	fmt.Printf("Value: %s\n", s1.Value)

	// Concurrent calls
	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			GetInstance()
		}()
	}

	wg.Wait()
	fmt.Println("All concurrent calls got same instance")
}

// ===== EXAMPLE 12: Atomic Operations =====

func example12_Atomic() {
	fmt.Println("\n=== EXAMPLE 12: Atomic Operations ===")

	var counter int64

	// Atomic operations
	atomic.AddInt64(&counter, 5)
	fmt.Printf("After AddInt64(5): %d\n", atomic.LoadInt64(&counter))

	atomic.StoreInt64(&counter, 100)
	fmt.Printf("After StoreInt64(100): %d\n", atomic.LoadInt64(&counter))

	// CompareAndSwap
	swapped := atomic.CompareAndSwapInt64(&counter, 100, 200)
	fmt.Printf("CompareAndSwap: %v, Value: %d\n", swapped, atomic.LoadInt64(&counter))

	// Concurrent increments with atomic
	var wg sync.WaitGroup
	var atomicCounter int64

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&atomicCounter, 1)
		}()
	}

	wg.Wait()
	fmt.Printf("Final atomic counter: %d\n", atomic.LoadInt64(&atomicCounter))
}

// ===== EXAMPLE 13: Sort Package =====

func example13_Sort() {
	fmt.Println("\n=== EXAMPLE 13: Sort Package ===")

	// Sort integers
	ints := []int{3, 1, 4, 1, 5, 9, 2, 6}
	sort.Ints(ints)
	fmt.Printf("Sorted ints: %v\n", ints)

	// Sort strings
	strings := []string{"banana", "apple", "cherry"}
	sort.Strings(strings)
	fmt.Printf("Sorted strings: %v\n", strings)

	// Sort floats
	floats := []float64{3.14, 2.71, 1.41, 1.73}
	sort.Float64s(floats)
	fmt.Printf("Sorted floats: %v\n", floats)

	// Check if sorted
	if sort.IntsAreSorted([]int{1, 2, 3, 4}) {
		fmt.Println("Already sorted")
	}
}

// ===== EXAMPLE 14: Custom Sort =====

type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func example14_CustomSort() {
	fmt.Println("\n=== EXAMPLE 14: Custom Sort ===")

	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}

	// Sort by age
	sort.Sort(ByAge(people))
	fmt.Println("Sorted by age:")
	for _, p := range people {
		fmt.Printf("  %s: %d\n", p.Name, p.Age)
	}
}

// ===== EXAMPLE 15: Regexp Package =====

func example15_Regexp() {
	fmt.Println("\n=== EXAMPLE 15: Regexp Package ===")

	// Match digits
	re := regexp.MustCompile(`\d+`)

	if re.MatchString("abc123def") {
		fmt.Println("Contains digits")
	}

	// Find patterns
	match := re.FindString("abc123def456")
	fmt.Printf("First match: %s\n", match)

	matches := re.FindAllString("abc123def456", -1)
	fmt.Printf("All matches: %v\n", matches)

	// Find with groups
	re = regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
	groups := re.FindStringSubmatch("Date: 2024-01-15")

	fmt.Println("Date groups:")
	if groups != nil {
		fmt.Printf("  Full: %s\n", groups[0])
		fmt.Printf("  Year: %s\n", groups[1])
		fmt.Printf("  Month: %s\n", groups[2])
		fmt.Printf("  Day: %s\n", groups[3])
	}
}

// ===== EXAMPLE 16: Regex Replace =====

func example16_RegexReplace() {
	fmt.Println("\n=== EXAMPLE 16: Regex Replace ===")

	re := regexp.MustCompile(`\d+`)

	// Replace all
	result := re.ReplaceAllString("abc123def456", "X")
	fmt.Printf("Replaced: %s\n", result)

	// Replace with function
	result = re.ReplaceAllStringFunc("abc123def456", func(s string) string {
		return "[" + s + "]"
	})
	fmt.Printf("With function: %s\n", result)
}

// ===== EXAMPLE 17: Hashing =====

func example17_Hashing() {
	fmt.Println("\n=== EXAMPLE 17: Hashing ===")

	data := []byte("hello world")

	// MD5
	md5sum := md5.Sum(data)
	fmt.Printf("MD5: %x\n", md5sum)

	// SHA256
	sha256sum := sha256.Sum256(data)
	fmt.Printf("SHA256: %x\n", sha256sum)

	// Using hash writer interface
	h := sha256.New()
	h.Write([]byte("hello"))
	h.Write([]byte(" "))
	h.Write([]byte("world"))
	fmt.Printf("SHA256 (writer): %x\n", h.Sum(nil))
}

// ===== EXAMPLE 18: Base64 and Hex Encoding =====

func example18_Encoding() {
	fmt.Println("\n=== EXAMPLE 18: Encoding ===")

	text := "hello world"

	// Base64
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	fmt.Printf("Base64 encoded: %s\n", encoded)

	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	fmt.Printf("Base64 decoded: %s\n", string(decoded))

	// Hex
	hexEncoded := hex.EncodeToString([]byte(text))
	fmt.Printf("Hex encoded: %s\n", hexEncoded)

	hexDecoded, _ := hex.DecodeString(hexEncoded)
	fmt.Printf("Hex decoded: %s\n", string(hexDecoded))
}

// ===== EXAMPLE 19: JSON Marshaling =====

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func example19_JSONMarshal() {
	fmt.Println("\n=== EXAMPLE 19: JSON Marshaling ===")

	user := User{Name: "Alice", Age: 30}

	// Marshal to JSON
	data, _ := json.Marshal(user)
	fmt.Printf("Marshaled: %s\n", string(data))

	// Unmarshal from JSON
	var u User
	json.Unmarshal([]byte(`{"name":"Bob","age":25}`), &u)
	fmt.Printf("Unmarshaled: %+v\n", u)

	// Pretty print
	data, _ = json.MarshalIndent(user, "", "  ")
	fmt.Printf("Pretty:\n%s\n", string(data))
}

// ===== EXAMPLE 20: CSV Operations =====

func example20_CSV() {
	fmt.Println("\n=== EXAMPLE 20: CSV Operations ===")

	// Write CSV
	records := [][]string{
		{"Name", "Age", "City"},
		{"Alice", "30", "NYC"},
		{"Bob", "25", "LA"},
	}

	fmt.Println("CSV records:")
	for i, record := range records {
		fmt.Printf("  Row %d: %v\n", i+1, record)
	}
}

// ===== MAIN FUNCTION =====

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Chapter 5: Standard Library - Comprehensive Examples     ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	example1_StringsPackage()
	example2_StringBuilder()
	example3_StrconvPackage()
	example4_FmtPackage()
	example5_TimeBasics()
	example6_TimeDuration()
	example7_TimersAndTickers()
	example8_SyncMutex()
	example9_RWMutex()
	example10_WaitGroup()
	example11_Once()
	example12_Atomic()
	example13_Sort()
	example14_CustomSort()
	example15_Regexp()
	example16_RegexReplace()
	example17_Hashing()
	example18_Encoding()
	example19_JSONMarshal()
	example20_CSV()

	// Command-line flag example
	fmt.Println("\n=== BONUS: Command-Line Flags ===")
	name := flag.String("name", "World", "Name to greet")
	count := flag.Int("count", 1, "Number of greetings")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.Parse()

	if *verbose {
		fmt.Printf("Processing with name=%s, count=%d\n", *name, *count)
	}

	for i := 0; i < *count; i++ {
		fmt.Printf("Hello, %s!\n", *name)
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Key Takeaways:                                           ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Println("║ 1. strings package for string operations                   ║")
	fmt.Println("║ 2. StringBuilder for efficient string building              ║")
	fmt.Println("║ 3. strconv for string↔number conversions                   ║")
	fmt.Println("║ 4. fmt verbs for precise formatting                        ║")
	fmt.Println("║ 5. time package for dates and durations                    ║")
	fmt.Println("║ 6. sync.Mutex for thread safety                            ║")
	fmt.Println("║ 7. sync.RWMutex for read-heavy workloads                   ║")
	fmt.Println("║ 8. WaitGroup for goroutine coordination                    ║")
	fmt.Println("║ 9. Atomic operations for lock-free programming             ║")
	fmt.Println("║ 10. regexp for pattern matching                            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
