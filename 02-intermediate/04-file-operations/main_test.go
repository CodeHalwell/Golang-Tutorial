package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// ===== Tests for File Reading =====

func TestReadEntireFile(t *testing.T) {
	// Create test file
	testFile := "test_read_entire.txt"
	content := "Hello, World!\nLine 2"
	os.WriteFile(testFile, []byte(content), 0644)
	defer os.Remove(testFile)

	// Read file
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(data) != content {
		t.Errorf("Expected %q, got %q", content, string(data))
	}
}

func TestReadNonExistentFile(t *testing.T) {
	_, err := os.ReadFile("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestReadFileSize(t *testing.T) {
	testFile := "test_size.txt"
	content := "1234567890"
	os.WriteFile(testFile, []byte(content), 0644)
	defer os.Remove(testFile)

	data, _ := os.ReadFile(testFile)
	if len(data) != 10 {
		t.Errorf("Expected size 10, got %d", len(data))
	}
}

// ===== Tests for File Writing =====

func TestWriteFile(t *testing.T) {
	testFile := "test_write.txt"
	content := []byte("Test content")
	defer os.Remove(testFile)

	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify
	readBack, _ := os.ReadFile(testFile)
	if string(readBack) != string(content) {
		t.Errorf("Content mismatch")
	}
}

func TestWriteFileOverwrite(t *testing.T) {
	testFile := "test_overwrite.txt"
	defer os.Remove(testFile)

	os.WriteFile(testFile, []byte("original"), 0644)
	os.WriteFile(testFile, []byte("new"), 0644)

	data, _ := os.ReadFile(testFile)
	if string(data) != "new" {
		t.Errorf("Expected overwrite, got %q", string(data))
	}
}

func TestAppendFile(t *testing.T) {
	testFile := "test_append.txt"
	defer os.Remove(testFile)

	// Initial write
	os.WriteFile(testFile, []byte("Line 1\n"), 0644)

	// Append
	file, _ := os.OpenFile(testFile, os.O_APPEND|os.O_WRONLY, 0644)
	file.WriteString("Line 2\n")
	file.Close()

	data, _ := os.ReadFile(testFile)
	expected := "Line 1\nLine 2\n"
	if string(data) != expected {
		t.Errorf("Expected %q, got %q", expected, string(data))
	}
}

// ===== Tests for File Information =====

func TestFileInfo(t *testing.T) {
	testFile := "test_info.txt"
	os.WriteFile(testFile, []byte("content"), 0644)
	defer os.Remove(testFile)

	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if info.Name() != "test_info.txt" {
		t.Errorf("Expected name test_info.txt, got %s", info.Name())
	}

	if info.Size() != 7 {
		t.Errorf("Expected size 7, got %d", info.Size())
	}

	if info.IsDir() {
		t.Error("Expected file, got directory")
	}
}

func TestFilePermissions(t *testing.T) {
	testFile := "test_perms.txt"
	os.WriteFile(testFile, []byte("test"), 0644)
	defer os.Remove(testFile)

	os.Chmod(testFile, 0600)
	info, _ := os.Stat(testFile)

	mode := info.Mode()
	expected := os.FileMode(0600)

	if mode.Perm() != expected {
		t.Errorf("Expected 0600, got %o", mode.Perm())
	}
}

// ===== Tests for Directory Operations =====

func TestCreateDirectory(t *testing.T) {
	testDir := "test_create_dir"
	defer os.RemoveAll(testDir)

	err := os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	info, _ := os.Stat(testDir)
	if !info.IsDir() {
		t.Error("Expected directory")
	}
}

func TestCreateNestedDirectory(t *testing.T) {
	testDir := "test_nested/deep/path"
	defer os.RemoveAll("test_nested")

	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	info, _ := os.Stat(testDir)
	if !info.IsDir() {
		t.Error("Expected directory to exist")
	}
}

func TestListDirectory(t *testing.T) {
	testDir := "test_list_dir"
	os.Mkdir(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Create files
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("test"), 0644)
	os.Mkdir(filepath.Join(testDir, "subdir"), 0755)

	entries, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}
}

func TestRemoveDirectory(t *testing.T) {
	testDir := "test_remove_dir"
	os.Mkdir(testDir, 0755)

	os.Remove(testDir)

	_, err := os.Stat(testDir)
	if err == nil {
		t.Error("Expected error, directory should be removed")
	}
}

// ===== Tests for Path Manipulation =====

func TestPathBase(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"dir/file.txt", "file.txt"},
		{"file.txt", "file.txt"},
		{"dir/subdir/file.txt", "file.txt"},
	}

	for _, tt := range tests {
		result := filepath.Base(tt.path)
		if result != tt.expected {
			t.Errorf("Base(%q) = %q, expected %q", tt.path, result, tt.expected)
		}
	}
}

func TestPathDir(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"dir/file.txt", "dir"},
		{"dir/subdir/file.txt", "dir/subdir"},
		{"file.txt", "."},
	}

	for _, tt := range tests {
		result := filepath.Dir(tt.path)
		if result != tt.expected {
			t.Errorf("Dir(%q) = %q, expected %q", tt.path, result, tt.expected)
		}
	}
}

func TestPathExt(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"file.txt", ".txt"},
		{"file.tar.gz", ".gz"},
		{"file", ""},
		{"dir/file.txt", ".txt"},
	}

	for _, tt := range tests {
		result := filepath.Ext(tt.path)
		if result != tt.expected {
			t.Errorf("Ext(%q) = %q, expected %q", tt.path, result, tt.expected)
		}
	}
}

func TestPathJoin(t *testing.T) {
	result := filepath.Join("dir", "subdir", "file.txt")
	if !hasComponents(result, []string{"dir", "subdir", "file.txt"}) {
		t.Errorf("Join failed: %s", result)
	}
}

func TestPathClean(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"path/./file.txt", "path/file.txt"},
		{"path/../file.txt", "file.txt"},
		{"path//file.txt", "path/file.txt"},
	}

	for _, tt := range tests {
		result := filepath.Clean(tt.path)
		if result != tt.expected {
			t.Errorf("Clean(%q) = %q, expected %q", tt.path, result, tt.expected)
		}
	}
}

// ===== Tests for File Copy =====

func TestCopyFile(t *testing.T) {
	srcFile := "test_copy_src.txt"
	dstFile := "test_copy_dst.txt"
	content := "Copy test content"

	os.WriteFile(srcFile, []byte(content), 0644)
	defer os.Remove(srcFile)
	defer os.Remove(dstFile)

	src, _ := os.Open(srcFile)
	dst, _ := os.Create(dstFile)

	_, err := io.Copy(dst, src)
	src.Close()
	dst.Close()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify
	readBack, _ := os.ReadFile(dstFile)
	if string(readBack) != content {
		t.Errorf("Content mismatch after copy")
	}
}

// ===== Tests for File Rename =====

func TestRenameFile(t *testing.T) {
	oldName := "test_old.txt"
	newName := "test_new.txt"

	os.WriteFile(oldName, []byte("test"), 0644)
	defer os.Remove(newName)

	err := os.Rename(oldName, newName)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Old file should not exist
	_, err = os.Stat(oldName)
	if err == nil {
		t.Error("Old file should not exist")
	}

	// New file should exist
	_, err = os.Stat(newName)
	if err != nil {
		t.Error("New file should exist")
	}
}

// ===== Tests for CSV Operations =====

func TestReadCSV(t *testing.T) {
	csvFile := "test_data.csv"
	defer os.Remove(csvFile)

	file, _ := os.Create(csvFile)
	writer := csv.NewWriter(file)

	records := [][]string{
		{"Name", "Age"},
		{"Alice", "30"},
		{"Bob", "25"},
	}

	for _, record := range records {
		writer.Write(record)
	}
	writer.Flush()
	file.Close()

	// Read back
	file, _ = os.Open(csvFile)
	reader := csv.NewReader(file)
	readRecords, _ := reader.ReadAll()
	file.Close()

	if len(readRecords) != 3 {
		t.Errorf("Expected 3 records, got %d", len(readRecords))
	}

	if readRecords[1][0] != "Alice" {
		t.Errorf("Expected Alice, got %s", readRecords[1][0])
	}
}

// ===== Tests for JSON Operations =====

func TestWriteJSON(t *testing.T) {
	jsonFile := "test_config.json"
	defer os.Remove(jsonFile)

	type TestConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	config := TestConfig{Host: "localhost", Port: 8080}
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(jsonFile, data, 0644)

	// Read back
	readData, _ := os.ReadFile(jsonFile)
	var readConfig TestConfig
	json.Unmarshal(readData, &readConfig)

	if readConfig.Host != "localhost" || readConfig.Port != 8080 {
		t.Errorf("Config mismatch: %+v", readConfig)
	}
}

// ===== Tests for Directory Walking =====

func TestWalkDirectory(t *testing.T) {
	testDir := "test_walk"
	os.MkdirAll(filepath.Join(testDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file2.txt"), []byte("test"), 0644)

	defer os.RemoveAll(testDir)

	fileCount := 0
	dirCount := 0

	filepath.WalkDir(testDir, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			dirCount++
		} else {
			fileCount++
		}
		return nil
	})

	if fileCount != 2 {
		t.Errorf("Expected 2 files, got %d", fileCount)
	}

	if dirCount != 2 { // testDir and subdir
		t.Errorf("Expected 2 dirs, got %d", dirCount)
	}
}

// ===== Benchmark Tests =====

func BenchmarkReadFile(b *testing.B) {
	testFile := "bench_read.txt"
	os.WriteFile(testFile, []byte("test content"), 0644)
	defer os.Remove(testFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.ReadFile(testFile)
	}
}

func BenchmarkWriteFile(b *testing.B) {
	testFile := "bench_write.txt"
	defer os.Remove(testFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.WriteFile(testFile, []byte("content"), 0644)
	}
}

func BenchmarkStatFile(b *testing.B) {
	testFile := "bench_stat.txt"
	os.WriteFile(testFile, []byte("test"), 0644)
	defer os.Remove(testFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.Stat(testFile)
	}
}

// ===== Helper Functions =====

func hasComponents(path string, components []string) bool {
	for _, comp := range components {
		if !contains(path, comp) {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
