package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func main() {

	currentFolder, regex, targetfolder := validateArgs()

	allFiles := getFileNames(currentFolder)
	matchedFiles := getFilesForRegex(allFiles, regex)

	moveFiles(currentFolder, matchedFiles, targetfolder)

}

func moveFiles(currentFolder string, matchedFiles []string, targetFolder string) {

	for _, fileName := range matchedFiles {
		sourcePath := filepath.Join(currentFolder, fileName)
		destPath := filepath.Join(targetFolder, fileName)

		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("Warning: File '%s' already exists in the target folder, skipping \n ", fileName)
			continue
		}

		err := os.Rename(sourcePath, destPath)
		if err != nil {
			err = copyAndDelete(sourcePath, destPath)
			if err != nil {
				fmt.Printf("Error moving file '%s' : %v\n ", fileName, err)
				continue
			}
		}

		fmt.Printf("Moved: %s \n", fileName)

	}

}

func copyAndDelete(sourcePath string, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}

	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy the content: %v", err)
	}

	sourceFile.Close()
	destFile.Close()

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to remove source file : %v", err)
	}

	return nil
}

func validateArgs() (string, string, string) {

	flag.Parse()
	if flag.NArg() != 3 {
		fmt.Printf("please enter folder currently with files, files regex and target folder to move")
		os.Exit(1)
	}

	regex := flag.Arg(1)

	if re, err := regexp.Compile(regex); err != nil {
		fmt.Printf("the regex %s is not a valid regex", re)
		os.Exit(1)
	}

	currentFolder := flag.Arg(0)
	if !folderExists(currentFolder) {
		fmt.Println("The folder with the files does not exist")
		os.Exit(1)
	}

	targetFolder := flag.Arg(2)
	if !folderExists(targetFolder) {
		fmt.Println("target folder was not found creating one ")
		os.Mkdir(targetFolder, 0755)
		fmt.Printf("create target folder %s .\n", targetFolder)
	}

	return currentFolder, regex, targetFolder

}

func getFilesForRegex(files []string, re string) []string {
	var matchedFiles []string

	pattern, err := regexp.Compile(re)

	if err != nil {
		fmt.Printf("Invalid regex Pattern %v \n", err)
		return matchedFiles
	}

	for _, file := range files {
		if pattern.MatchString(file) {
			matchedFiles = append(matchedFiles, file)
		}
	}
	return matchedFiles
}

func getFileNames(folder string) []string {
	var files []string

	entries, err := os.ReadDir(folder)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return files
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files
}

func folderExists(path string) bool {
	info, error := os.Stat(path)
	if os.IsNotExist(error) {
		return false
	}

	return info.IsDir()
}
