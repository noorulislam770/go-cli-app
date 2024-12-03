package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// Test file structure setup helper
func setupTestFiles(t *testing.T) (string, string, func()) {
	// Create temporary test directories
	sourceDir, err := os.MkdirTemp("", "source")
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	targetDir, err := os.MkdirTemp("", "target")
	if err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	// Create test files
	testFiles := []string{
		"image001.jpg",
		"doc123.pdf",
		"image002.png",
		"test.txt",
		"IMG_20230101.jpg",
	}

	for _, file := range testFiles {
		path := filepath.Join(sourceDir, file)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(sourceDir)
		os.RemoveAll(targetDir)
	}

	return sourceDir, targetDir, cleanup
}

// Test getFileNames function
func TestGetFileNames(t *testing.T) {
	sourceDir, _, cleanup := setupTestFiles(t)
	defer cleanup()

	files := getFileNames(sourceDir)
	if len(files) != 5 {
		t.Errorf("Expected 5 files, got %d", len(files))
	}
}

// Test getFilesForRegex function
func TestGetFilesForRegex(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		pattern  string
		expected int
	}{
		{
			name:     "Match JPG files",
			files:    []string{"test.jpg", "test.txt", "image.jpg"},
			pattern:  `\.jpg$`,
			expected: 2,
		},
		{
			name:     "Match numbered files",
			files:    []string{"file1.txt", "file2.txt", "test.txt"},
			pattern:  `file\d+\.txt$`,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := getFilesForRegex(tt.files, tt.pattern)
			if len(matched) != tt.expected {
				t.Errorf("Expected %d matches, got %d", tt.expected, len(matched))
			}
		})
	}
}

// Test moveFiles function
func TestMoveFiles(t *testing.T) {
	sourceDir, targetDir, cleanup := setupTestFiles(t)
	defer cleanup()

	// Get all jpg files
	allFiles := getFileNames(sourceDir)
	matchedFiles := getFilesForRegex(allFiles, `\.jpg$`)

	// Move the files
	moveFiles(sourceDir, matchedFiles, targetDir)

	// Check target directory
	movedFiles := getFileNames(targetDir)
	if len(movedFiles) != 2 { // Should have 2 .jpg files
		t.Errorf("Expected 2 files in target directory, got %d", len(movedFiles))
	}
}

// Test invalid regex pattern
func TestInvalidRegex(t *testing.T) {
	invalidPattern := "["
	if isValidRegex(invalidPattern) {
		t.Error("Expected invalid regex pattern to return false")
	}
}

func isValidRegex(regex string) bool {

	if re, err := regexp.Compile(regex); err != nil {
		fmt.Printf("the regex %s is not a valid regex", re)
		return false
	}
	return true
}
