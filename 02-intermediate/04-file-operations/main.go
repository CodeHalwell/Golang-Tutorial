package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ===== EXAMPLE 1: Read Entire File =====

func example1_ReadEntireFile() {
	fmt.Println("\n=== EXAMPLE 1: Read Entire File ===")

	// Create a test file first
	testFile := "test1.txt"
	content := "Hello, World!\nThis is line 2.\nAnd line 3."
	os.WriteFile(testFile, []byte(content), 0644)

	// Read entire file
	data, err := os.ReadFile(testFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("File size: %d bytes\n", len(data))
	fmt.Printf("Content:\n%s\n", string(data))

	// Cleanup
	os.Remove(testFile)
}

// ===== EXAMPLE 2: Read Line by Line =====

func example2_ReadLineByLine() {
	fmt.Println("\n=== EXAMPLE 2: Read Line by Line ===")

	// Create test file
	testFile := "test2.txt"
	lines := []string{"Line 1", "Line 2", "Line 3", "Line 4"}
	content := strings.Join(lines, "\n")
	os.WriteFile(testFile, []byte(content), 0644)

	// Open file
	file, err := os.Open(testFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file.Close()

	// Read line by line
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		fmt.Printf("Line %d: %s\n", lineNum, line)
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v\n", err)
	}

	// Cleanup
	os.Remove(testFile)
}

// ===== EXAMPLE 3: Read in Chunks =====

func example3_ReadInChunks() {
	fmt.Println("\n=== EXAMPLE 3: Read in Chunks ===")

	// Create test file with known content
	testFile := "test3.txt"
	content := "This is some binary-like content for chunk reading."
	os.WriteFile(testFile, []byte(content), 0644)

	file, err := os.Open(testFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file.Close()

	// Read in 16-byte chunks
	buffer := make([]byte, 16)
	chunkNum := 0
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		chunkNum++
		fmt.Printf("Chunk %d (%d bytes): %s\n", chunkNum, n, string(buffer[:n]))
	}

	// Cleanup
	os.Remove(testFile)
}

// ===== EXAMPLE 4: Write Entire File =====

func example4_WriteEntireFile() {
	fmt.Println("\n=== EXAMPLE 4: Write Entire File ===")

	testFile := "test4.txt"
	content := []byte("Hello, this is written content!")

	// Write file with 0644 permissions
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("File written: %s\n", testFile)

	// Verify by reading back
	readBack, _ := os.ReadFile(testFile)
	fmt.Printf("Verified content: %s\n", string(readBack))

	// Cleanup
	os.Remove(testFile)
}

// ===== EXAMPLE 5: Write with Buffering =====

func example5_WriteWithBuffering() {
	fmt.Println("\n=== EXAMPLE 5: Write with Buffering ===")

	testFile := "test5.txt"

	// Create file
	file, err := os.Create(testFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file.Close()

	// Write with buffering
	writer := bufio.NewWriter(file)

	lines := []string{"Buffered line 1", "Buffered line 2", "Buffered line 3"}
	for i, line := range lines {
		_, err := writer.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			fmt.Printf("Error writing: %v\n", err)
		}
		fmt.Printf("Wrote line %d (buffered)\n", i+1)
	}

	// IMPORTANT: Flush to disk
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Flush error: %v\n", err)
	}

	fmt.Println("All data flushed to disk")

	// Verify
	data, _ := os.ReadFile(testFile)
	fmt.Printf("File content:\n%s", string(data))

	// Cleanup
	os.Remove(testFile)
}

// ===== EXAMPLE 6: Append to File =====

func example6_AppendToFile() {
	fmt.Println("\n=== EXAMPLE 6: Append to File ===")

	testFile := "test6.txt"

	// Create initial file
	os.WriteFile(testFile, []byte("Initial content\n"), 0644)

	// Open for appending
	file, err := os.OpenFile(testFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file.Close()

	// Append lines
	_, err = file.WriteString("Appended line 1\n")
	_, err = file.WriteString("Appended line 2\n")

	fmt.Println("Lines appended")

	// Verify
	data, _ := os.ReadFile(testFile)
	fmt.Printf("File content:\n%s", string(data))

	// Cleanup
	os.Remove(testFile)
}

// ===== EXAMPLE 7: File Information =====

func example7_FileInfo() {
	fmt.Println("\n=== EXAMPLE 7: File Information ===")

	testFile := "test7.txt"
	os.WriteFile(testFile, []byte("Test content for stat"), 0644)

	// Get file info
	info, err := os.Stat(testFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Name:       %s\n", info.Name())
	fmt.Printf("Size:       %d bytes\n", info.Size())
	fmt.Printf("Modified:   %v\n", info.ModTime())
	fmt.Printf("Is Dir:     %v\n", info.IsDir())
	fmt.Printf("Permissions:%v\n", info.Mode())

	// Cleanup
	os.Remove(testFile)
}

// ===== EXAMPLE 8: File Permissions =====

func example8_FilePermissions() {
	fmt.Println("\n=== EXAMPLE 8: File Permissions ===")

	testFile := "test8.txt"
	os.WriteFile(testFile, []byte("Test"), 0644)

	// Get info
	info, _ := os.Stat(testFile)
	mode := info.Mode()

	// Check permissions
	fmt.Printf("Permissions: %o\n", mode.Perm())

	isReadable := (mode & 0444) != 0
	isWritable := (mode & 0222) != 0
	isExecutable := (mode & 0111) != 0

	fmt.Printf("Readable:   %v\n", isReadable)
	fmt.Printf("Writable:   %v\n", isWritable)
	fmt.Printf("Executable: %v\n", isExecutable)

	// Change permissions
	os.Chmod(testFile, 0600) // rw-------
	info, _ = os.Stat(testFile)
	fmt.Printf("After chmod: %o\n", info.Mode().Perm())

	// Cleanup
	os.Remove(testFile)
}

// ===== EXAMPLE 9: Directory Operations =====

func example9_DirectoryOperations() {
	fmt.Println("\n=== EXAMPLE 9: Directory Operations ===")

	testDir := "testdir"

	// Create directory
	err := os.Mkdir(testDir, 0755)
	if err != nil {
		fmt.Printf("Error creating dir: %v\n", err)
	} else {
		fmt.Printf("Directory created: %s\n", testDir)
	}

	// Create nested directories
	nestedDir := filepath.Join(testDir, "subdir", "nested")
	os.MkdirAll(nestedDir, 0755)
	fmt.Printf("Nested directory created: %s\n", nestedDir)

	// Create test files
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test1"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("test2"), 0644)

	// List directory
	entries, _ := os.ReadDir(testDir)
	fmt.Printf("\nContents of %s:\n", testDir)
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("  [DIR]  %s\n", entry.Name())
		} else {
			info, _ := entry.Info()
			fmt.Printf("  [FILE] %s (%d bytes)\n", entry.Name(), info.Size())
		}
	}

	// Cleanup
	os.RemoveAll(testDir)
	fmt.Println("\nDirectory removed")
}

// ===== EXAMPLE 10: Path Manipulation =====

func example10_PathManipulation() {
	fmt.Println("\n=== EXAMPLE 10: Path Manipulation ===")

	path := "dir/subdir/file.txt"

	fmt.Printf("Original:    %s\n", path)
	fmt.Printf("Base:        %s\n", filepath.Base(path))
	fmt.Printf("Dir:         %s\n", filepath.Dir(path))
	fmt.Printf("Ext:         %s\n", filepath.Ext(path))

	// Join paths (platform-independent)
	joined := filepath.Join("dir", "subdir", "file.txt")
	fmt.Printf("Joined:      %s\n", joined)

	// Clean path
	messy := "path/./to/../file.txt"
	clean := filepath.Clean(messy)
	fmt.Printf("Messy:       %s\n", messy)
	fmt.Printf("Clean:       %s\n", clean)

	// Absolute path
	abs, _ := filepath.Abs("file.txt")
	fmt.Printf("Absolute:    %s\n", abs)

	// Current directory
	cwd, _ := os.Getwd()
	fmt.Printf("Current dir: %s\n", cwd)
}

// ===== EXAMPLE 11: Copy File =====

func example11_CopyFile() {
	fmt.Println("\n=== EXAMPLE 11: Copy File ===")

	srcFile := "source.txt"
	dstFile := "destination.txt"

	// Create source
	os.WriteFile(srcFile, []byte("This is the source file content."), 0644)

	// Copy using io.Copy
	src, _ := os.Open(srcFile)
	defer src.Close()

	dst, _ := os.Create(dstFile)
	defer dst.Close()

	bytes, err := io.Copy(dst, src)
	if err != nil {
		fmt.Printf("Error copying: %v\n", err)
	} else {
		fmt.Printf("Copied %d bytes\n", bytes)
	}

	// Verify
	dstContent, _ := os.ReadFile(dstFile)
	fmt.Printf("Destination contains: %s\n", string(dstContent))

	// Cleanup
	os.Remove(srcFile)
	os.Remove(dstFile)
}

// ===== EXAMPLE 12: Rename File =====

func example12_RenameFile() {
	fmt.Println("\n=== EXAMPLE 12: Rename File ===")

	oldName := "oldname.txt"
	newName := "newname.txt"

	// Create file
	os.WriteFile(oldName, []byte("Test"), 0644)
	fmt.Printf("Created: %s\n", oldName)

	// Rename
	err := os.Rename(oldName, newName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Renamed to: %s\n", newName)
	}

	// Verify
	info, _ := os.Stat(newName)
	fmt.Printf("New file exists: %v\n", info != nil)

	// Cleanup
	os.Remove(newName)
}

// ===== EXAMPLE 13: CSV File Operations =====

type Person struct {
	Name string
	Age  int
	City string
}

func example13_CSVOperations() {
	fmt.Println("\n=== EXAMPLE 13: CSV File Operations ===")

	csvFile := "data.csv"

	// Write CSV
	file, _ := os.Create(csvFile)
	writer := csv.NewWriter(file)

	records := [][]string{
		{"Name", "Age", "City"},
		{"Alice", "30", "New York"},
		{"Bob", "25", "Los Angeles"},
		{"Charlie", "35", "Chicago"},
	}

	for _, record := range records {
		writer.Write(record)
	}
	writer.Flush()
	file.Close()

	fmt.Println("CSV file written")

	// Read CSV
	file, _ = os.Open(csvFile)
	reader := csv.NewReader(file)
	readRecords, _ := reader.ReadAll()
	file.Close()

	fmt.Println("\nCSV content:")
	for i, record := range readRecords {
		fmt.Printf("Row %d: %v\n", i+1, record)
	}

	// Cleanup
	os.Remove(csvFile)
}

// ===== EXAMPLE 14: JSON File Operations =====

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

func example14_JSONOperations() {
	fmt.Println("\n=== EXAMPLE 14: JSON File Operations ===")

	jsonFile := "config.json"

	// Write JSON
	config := Config{
		Host:     "localhost",
		Port:     8080,
		Database: "mydb",
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(jsonFile, data, 0644)

	fmt.Println("JSON written:")
	fmt.Println(string(data))

	// Read JSON
	readData, _ := os.ReadFile(jsonFile)
	var readConfig Config
	json.Unmarshal(readData, &readConfig)

	fmt.Printf("\nParsed config: %+v\n", readConfig)

	// Cleanup
	os.Remove(jsonFile)
}

// ===== EXAMPLE 15: Process Large File =====

func example15_ProcessLargeFile() {
	fmt.Println("\n=== EXAMPLE 15: Process Large File ===")

	testFile := "large.txt"

	// Create a "large" file with many lines
	file, _ := os.Create(testFile)
	writer := bufio.NewWriter(file)
	for i := 1; i <= 1000; i++ {
		writer.WriteString(fmt.Sprintf("Line %d\n", i))
	}
	writer.Flush()
	file.Close()

	// Process line by line with progress
	file, _ = os.Open(testFile)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Increase buffer size for longer lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lineCount := 0
	sumCount := 0

	for scanner.Scan() {
		lineCount++

		// Every 250 lines, print progress
		if lineCount%250 == 0 {
			fmt.Printf("Processed %d lines\n", lineCount)
		}

		sumCount += lineCount
	}

	fmt.Printf("Total lines: %d\n", lineCount)

	// Cleanup
	os.Remove(testFile)
}

// ===== EXAMPLE 16: Safe File Writing =====

func example16_SafeFileWriting() {
	fmt.Println("\n=== EXAMPLE 16: Safe File Writing ===")

	filename := "important.txt"
	data := []byte("Important data that must not be lost!")

	// Write to temp file first
	tmpfile, _ := os.CreateTemp(filepath.Dir(filename), ".tmp")
	fmt.Printf("Writing to temp: %s\n", tmpfile.Name())

	_, _ = tmpfile.Write(data)
	tmpfile.Close()

	// Atomic rename
	os.Rename(tmpfile.Name(), filename)
	fmt.Printf("Renamed to: %s\n", filename)

	// Verify
	readBack, _ := os.ReadFile(filename)
	fmt.Printf("Verified: %s\n", string(readBack))

	// Cleanup
	os.Remove(filename)
}

// ===== EXAMPLE 17: Walk Directory Tree =====

func example17_WalkDirectoryTree() {
	fmt.Println("\n=== EXAMPLE 17: Walk Directory Tree ===")

	testDir := "treedir"

	// Create test directory structure
	os.MkdirAll(filepath.Join(testDir, "dir1"), 0755)
	os.MkdirAll(filepath.Join(testDir, "dir2", "subdir"), 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "dir1", "file2.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "dir2", "subdir", "file3.txt"), []byte("test"), 0644)

	fmt.Println("Walking directory tree:")

	filepath.WalkDir(testDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			fmt.Printf("[DIR]  %s\n", path)
		} else {
			info, _ := d.Info()
			fmt.Printf("[FILE] %s (%d bytes)\n", path, info.Size())
		}

		return nil
	})

	// Cleanup
	os.RemoveAll(testDir)
}

// ===== MAIN FUNCTION =====

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Chapter 4: File Operations - Comprehensive Examples      ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")

	example1_ReadEntireFile()
	example2_ReadLineByLine()
	example3_ReadInChunks()
	example4_WriteEntireFile()
	example5_WriteWithBuffering()
	example6_AppendToFile()
	example7_FileInfo()
	example8_FilePermissions()
	example9_DirectoryOperations()
	example10_PathManipulation()
	example11_CopyFile()
	example12_RenameFile()
	example13_CSVOperations()
	example14_JSONOperations()
	example15_ProcessLargeFile()
	example16_SafeFileWriting()
	example17_WalkDirectoryTree()

	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Key Takeaways:                                           ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Println("║ 1. Use os.ReadFile for small, complete files              ║")
	fmt.Println("║ 2. Use bufio.Scanner for streaming large files             ║")
	fmt.Println("║ 3. Always defer Close() to prevent leaks                  ║")
	fmt.Println("║ 4. Flush buffered writers before closing                   ║")
	fmt.Println("║ 5. Check scanner.Err() after reading                       ║")
	fmt.Println("║ 6. Use filepath for platform-independent paths             ║")
	fmt.Println("║ 7. Use atomic renames for critical files                   ║")
	fmt.Println("║ 8. Handle partial I/O (not all data written)               ║")
	fmt.Println("║ 9. Use appropriate buffer sizes for performance            ║")
	fmt.Println("║ 10. Check errors for permissions and file existence        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
