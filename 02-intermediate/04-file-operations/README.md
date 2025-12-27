# Chapter 4: File Operations

## Learning Objectives

1. Read and write files with various methods
2. Work with file metadata and permissions
3. Handle directories and file paths
4. Use buffers for efficient I/O
5. Work with CSV and JSON files
6. Process large files efficiently
7. Handle file errors gracefully
8. Copy, rename, and delete files
9. Monitor file changes
10. Work with file permissions and ownership

---

## Basic File Operations

### Reading Files

#### Method 1: Read Entire File at Once

```go
import "os"

data, err := os.ReadFile("data.txt")
if err != nil {
    // Handle: file not found, permission denied, etc.
    log.Fatal(err)
}

// data is []byte
content := string(data)  // Convert to string
fmt.Println(content)
```

**Use Case**: Small files that fit in memory

**Pros**:
- Simple, one-liner
- Entire file available at once

**Cons**:
- Loads entire file into memory
- Bad for large files

#### Method 2: Read Line by Line

```go
import (
    "bufio"
    "os"
)

file, err := os.Open("data.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
}

if err := scanner.Err(); err != nil {
    log.Fatal(err)
}
```

**Use Case**: Large files that must be processed line by line

**Advantages**:
- Memory efficient (only one line at a time)
- Can process while reading

#### Method 3: Read in Chunks

```go
import (
    "io"
    "os"
)

file, err := os.Open("data.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

buffer := make([]byte, 1024)
for {
    n, err := file.Read(buffer)
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }

    // Process first n bytes
    process(buffer[:n])
}
```

**Use Case**: Binary files, network I/O, streaming data

---

### Writing Files

#### Method 1: Write Entire File

```go
import "os"

data := []byte("Hello, World!")
err := os.WriteFile("output.txt", data, 0644)
if err != nil {
    log.Fatal(err)
}
```

**File Permission**: `0644` = rw-r--r--
- Owner: read + write (6)
- Group: read only (4)
- Others: read only (4)

#### Method 2: Write with Buffering

```go
import (
    "bufio"
    "os"
)

file, err := os.Create("output.txt") // Truncates if exists
if err != nil {
    log.Fatal(err)
}
defer file.Close()

writer := bufio.NewWriter(file)
writer.WriteString("Line 1\n")
writer.WriteString("Line 2\n")

// IMPORTANT: Flush buffered data to disk
err = writer.Flush()
if err != nil {
    log.Fatal(err)
}
```

#### Method 3: Append to File

```go
import "os"

file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
    log.Fatal(err)
}
defer file.Close()

_, err = file.WriteString("New line\n")
if err != nil {
    log.Fatal(err)
}
```

**Flags**:
- `os.O_RDONLY` - Read only
- `os.O_WRONLY` - Write only
- `os.O_RDWR` - Read and write
- `os.O_CREATE` - Create if doesn't exist
- `os.O_TRUNC` - Truncate to zero length
- `os.O_APPEND` - Append to end
- `os.O_EXCL` - Used with CREATE, file must not exist

---

## File Metadata and Information

### FileInfo Interface

```go
import "os"

info, err := os.Stat("data.txt")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", info.Name())           // Filename
fmt.Printf("Size: %d bytes\n", info.Size())     // File size
fmt.Printf("Modified: %v\n", info.ModTime())    // Last modification time
fmt.Printf("Is Dir: %v\n", info.IsDir())        // Is directory?
fmt.Printf("Mode: %v\n", info.Mode())           // File permissions
```

### File Permissions

```go
import "os"

info, _ := os.Stat("data.txt")
mode := info.Mode()

// Check permissions
if mode&0444 != 0 {
    fmt.Println("Readable")
}

if mode&0222 != 0 {
    fmt.Println("Writable")
}

if mode&0111 != 0 {
    fmt.Println("Executable")
}

// Change permissions
os.Chmod("data.txt", 0600) // rw-------
```

---

## Directory Operations

### List Directory Contents

```go
import (
    "fmt"
    "os"
)

entries, err := os.ReadDir(".")
if err != nil {
    log.Fatal(err)
}

for _, entry := range entries {
    if entry.IsDir() {
        fmt.Printf("DIR:  %s\n", entry.Name())
    } else {
        info, _ := entry.Info()
        fmt.Printf("FILE: %s (%d bytes)\n", entry.Name(), info.Size())
    }
}
```

### Create Directories

```go
import "os"

// Create single directory
err := os.Mkdir("mydir", 0755)
if err != nil {
    log.Fatal(err)
}

// Create nested directories
err = os.MkdirAll("path/to/deep/dir", 0755)
if err != nil {
    log.Fatal(err)
}
```

### Remove Directories

```go
import "os"

// Remove empty directory
err := os.Remove("mydir")
if err != nil {
    log.Fatal(err)
}

// Remove directory and contents
err = os.RemoveAll("path/to/dir")
if err != nil {
    log.Fatal(err)
}
```

---

## Path Manipulation

### Working with Paths

```go
import (
    "filepath"
    "path"
)

// filepath: platform-specific (handles \ on Windows, / on Linux)
filepath.Join("dir", "subdir", "file.txt")          // dir/subdir/file.txt
filepath.Base("dir/subdir/file.txt")                // file.txt
filepath.Dir("dir/subdir/file.txt")                 // dir/subdir
filepath.Ext("file.txt")                            // .txt
filepath.IsAbs("/path/to/file")                     // true
filepath.Abs("file.txt")                            // /current/working/dir/file.txt

// path: URL-style paths (always use /)
path.Join("dir", "subdir", "file.txt")              // dir/subdir/file.txt
```

### Working with File Paths

```go
import (
    "os"
    "filepath"
)

// Get working directory
cwd, _ := os.Getwd()
fmt.Println(cwd)

// Change directory
os.Chdir("/tmp")

// Get home directory
home, _ := os.UserHomeDir()

// Clean path (remove . and ..)
clean := filepath.Clean("path/./to/../file.txt")  // path/file.txt
```

---

## Copying and Moving Files

### Copy File

```go
import (
    "io"
    "os"
)

// Method 1: Using io.Copy
src, _ := os.Open("source.txt")
defer src.Close()

dst, _ := os.Create("dest.txt")
defer dst.Close()

io.Copy(dst, src)  // Copy all data

// Method 2: With buffer for control
func copyFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    dest, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dest.Close()

    if _, err := io.Copy(dest, source); err != nil {
        return err
    }

    return dest.Sync()  // Ensure written to disk
}
```

### Move/Rename File

```go
import "os"

// Rename (works across filesystems on most systems)
err := os.Rename("old_name.txt", "new_name.txt")
if err != nil {
    log.Fatal(err)
}
```

---

## Working with Structured Data

### CSV Files

```go
import (
    "encoding/csv"
    "fmt"
    "os"
)

// Read CSV
file, _ := os.Open("data.csv")
defer file.Close()

reader := csv.NewReader(file)
records, err := reader.ReadAll()  // Read all records
if err != nil {
    log.Fatal(err)
}

for _, record := range records {
    fmt.Printf("Name: %s, Age: %s\n", record[0], record[1])
}

// Write CSV
file, _ = os.Create("output.csv")
defer file.Close()

writer := csv.NewWriter(file)
writer.Write([]string{"Name", "Age"})
writer.Write([]string{"Alice", "30"})
writer.Write([]string{"Bob", "25"})
writer.Flush()
```

### JSON Files

```go
import (
    "encoding/json"
    "os"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// Read JSON
data, _ := os.ReadFile("users.json")
var users []User
json.Unmarshal(data, &users)

// Write JSON
users := []User{
    {Name: "Alice", Age: 30},
    {Name: "Bob", Age: 25},
}

data, _ := json.MarshalIndent(users, "", "  ")
os.WriteFile("output.json", data, 0644)
```

---

## Efficient File Processing

### Process Large Files Line by Line

```go
import (
    "bufio"
    "os"
)

func processLargeFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Optional: increase buffer size for long lines
    buf := make([]byte, 0, 64*1024)  // 64KB buffer
    scanner.Buffer(buf, 1024*1024)   // Max 1MB per line

    lineNum := 0
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()

        // Process line
        processLine(line)

        // Show progress
        if lineNum%10000 == 0 {
            fmt.Printf("Processed %d lines\n", lineNum)
        }
    }

    return scanner.Err()
}
```

### Stream Processing with Channels

```go
import (
    "bufio"
    "os"
)

func readLinesChannel(filename string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        file, _ := os.Open(filename)
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            out <- scanner.Text()
        }
    }()
    return out
}

// Usage
for line := range readLinesChannel("data.txt") {
    processLine(line)
}
```

---

## Error Handling for File Operations

### Common File Errors

```go
import (
    "errors"
    "os"
)

file, err := os.Open("data.txt")
if err != nil {
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("File not found")
    } else if errors.Is(err, os.ErrPermission) {
        fmt.Println("Permission denied")
    } else {
        fmt.Println("Other error:", err)
    }
}
```

### Safe File Writing with Temporary Files

```go
import (
    "os"
    "path/filepath"
)

func safeWrite(filename string, data []byte) error {
    // Write to temporary file first
    tmpfile, err := os.CreateTemp(filepath.Dir(filename), ".tmp")
    if err != nil {
        return err
    }

    if _, err := tmpfile.Write(data); err != nil {
        tmpfile.Close()
        os.Remove(tmpfile.Name())
        return err
    }

    if err := tmpfile.Close(); err != nil {
        os.Remove(tmpfile.Name())
        return err
    }

    // Atomic rename
    return os.Rename(tmpfile.Name(), filename)
}
```

---

## Best Practices

### 1. Always Defer Close()

```go
// ✓ Good
file, err := os.Open("data.txt")
if err != nil {
    return err
}
defer file.Close()  // Guaranteed to run

// ✗ Bad
file, _ := os.Open("data.txt")
// If error happens here, file never closes
```

### 2. Flush Buffered Writers

```go
import "bufio"

writer := bufio.NewWriter(file)
writer.WriteString("data")
writer.Flush()  // Don't forget this!
```

### 3. Check Scanner Errors

```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    processLine(scanner.Text())
}

// ✓ Check for errors after loop
if err := scanner.Err(); err != nil {
    log.Fatal(err)
}
```

### 4. Use Atomic File Operations

```go
import "os"

// Bad: File partially written if crash
ioutil.WriteFile("important.txt", data, 0644)

// Good: Write to temp, then rename
tmpfile, _ := os.CreateTemp(".", ".tmp")
tmpfile.Write(data)
tmpfile.Close()
os.Rename(tmpfile.Name(), "important.txt")
```

### 5. Handle Partial I/O

```go
import "io"

n, err := file.Write(data)
if err != nil {
    // Not all data was written if err != nil
    log.Printf("Wrote %d bytes, error: %v", n, err)
}
```

### 6. Use Efficient Buffering

```go
// ✓ Good: Custom buffer size
scanner := bufio.NewScanner(file)
buf := make([]byte, 0, 64*1024)  // 64KB
scanner.Buffer(buf, 1024*1024)   // Max 1MB

// ✗ Bad: Default 64KB buffer for large lines
scanner := bufio.NewScanner(file)
```

### 7. Walk Directory Tree

```go
import "filepath"

filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err
    }

    if !d.IsDir() {
        fmt.Println(path)
    }
    return nil
})
```

---

## Practical Examples

### Example 1: Count Lines in File

```go
import (
    "bufio"
    "os"
)

func countLines(filename string) (int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    count := 0
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        count++
    }

    return count, scanner.Err()
}
```

### Example 2: Parse Configuration File

```go
import (
    "bufio"
    "os"
    "strings"
)

func parseConfig(filename string) (map[string]string, error) {
    config := make(map[string]string)
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "#") {
            continue  // Skip comments
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])
            config[key] = value
        }
    }

    return config, scanner.Err()
}
```

### Example 3: File Statistics

```go
import (
    "fmt"
    "os"
    "path/filepath"
)

type DirStats struct {
    FileCount      int
    DirCount       int
    TotalSize      int64
    LargestFile    string
    LargestSize    int64
}

func getStats(root string) (*DirStats, error) {
    stats := &DirStats{}

    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            stats.DirCount++
        } else {
            stats.FileCount++
            info, _ := d.Info()
            size := info.Size()
            stats.TotalSize += size

            if size > stats.LargestSize {
                stats.LargestSize = size
                stats.LargestFile = path
            }
        }

        return nil
    })

    return stats, err
}
```

---

## Summary

- **Read entire**: `os.ReadFile()` for small files
- **Read streams**: `bufio.Scanner` for large files
- **Write files**: `os.WriteFile()` or `bufio.Writer`
- **Append**: `os.OpenFile()` with `O_APPEND`
- **Metadata**: `os.Stat()` for info
- **Paths**: `filepath` package for platform-independent paths
- **Structured data**: `encoding/csv`, `encoding/json`
- **Always defer close**: Prevent resource leaks
- **Flush buffered writers**: Ensure data reaches disk
- **Handle errors properly**: Check both I/O and scanner errors

This chapter provides the foundation for real-world file operations. Next: Module 2 Chapter 5 - Standard Library Deep Dive.
