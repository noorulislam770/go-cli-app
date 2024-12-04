package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetFileNames(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{"test1.txt", "test2.txt", "apples.jpg", "accounts.pdf", "cv.doc", "sheets.xls"}
	for _, fname := range testFiles {
		if err := os.WriteFile(filepath.Join(tmpDir, fname), []byte("test"), 0666); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	files := getFileNames(tmpDir)
	if len(files) != len(testFiles) {
		t.Errorf("Expected %d files, got %d", len(testFiles), len(files))
	}
}

func TestGetFilesForRegex(t *testing.T) {
	testCases := []struct {
		name     string
		files    []string
		regex    string
		expected int
	}{
		{
			name:     "Match txt files",
			files:    []string{"test1.txt", "test2.txt", "test.pdf"},
			regex:    `\.txt$`,
			expected: 2,
		},
		{
			name:     "Match numbered files",
			files:    []string{"test1.txt", "test2.txt", "test.pdf"},
			regex:    `\d`,
			expected: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matched := getFilesForRegex(tc.files, tc.regex)
			if len(matched) != tc.expected {
				t.Errorf("Expected %d matches, got %d", tc.expected, len(matched))
			}
		})
	}
}

func TestMoveFiles(t *testing.T) {
	// Create source and destination directories
	srcDir, err := os.MkdirTemp("", "src")
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	destDir, err := os.MkdirTemp("", "dest")
	if err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}
	defer os.RemoveAll(destDir)

	// Create test file
	testFile := "test.txt"
	if err := os.WriteFile(filepath.Join(srcDir, testFile), []byte("test"), 0666); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test moving file
	moveFiles(srcDir, []string{testFile}, destDir)

	// Verify file was moved
	if _, err := os.Stat(filepath.Join(destDir, testFile)); os.IsNotExist(err) {
		t.Errorf("Expected file to exist in destination")
	}
}

func TestFolderExists(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if !folderExists(tmpDir) {
		t.Error("Expected folder to exist")
	}

	if folderExists("/nonexistent/path") {
		t.Error("Expected folder to not exist")
	}
}

func TestCopyAndDelete(t *testing.T) {
	// Create source and destination directories
	srcDir, err := os.MkdirTemp("", "src")
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	destDir, err := os.MkdirTemp("", "dest")
	if err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}
	defer os.RemoveAll(destDir)

	// Create test file
	content := []byte("test content")
	srcPath := filepath.Join(srcDir, "test.txt")
	destPath := filepath.Join(destDir, "test.txt")

	if err := os.WriteFile(srcPath, content, 0666); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test copy and delete
	if err := copyAndDelete(srcPath, destPath); err != nil {
		t.Fatalf("copyAndDelete failed: %v", err)
	}

	// Verify destination file exists and has correct content
	destContent, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(destContent) != string(content) {
		t.Error("Destination file content doesn't match source")
	}

	// Verify source file was deleted
	if _, err := os.Stat(srcPath); !os.IsNotExist(err) {
		t.Error("Source file should have been deleted")
	}
}
