package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type CLIOptions struct {
	SourceFolder    string
	RegexPattern    string
	TargetFolder    string
	InteractiveMode bool
}

func main() {

	options, err := RunCli()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := ProcessFiles(options); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing files %v\n", err)
		os.Exit(1)
	}

	// currentFolder, regex, targetfolder := validateArgs()

	// allFiles := getFileNames(currentFolder)
	// matchedFiles := getFilesForRegex(allFiles, regex)

	// moveFiles(currentFolder, matchedFiles, targetfolder)

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

func RunCli() (*CLIOptions, error) {
	interactive := flag.Bool("i", false, "Run in interactive mode.")
	help := flag.Bool("h", false, "Show help Messages")

	flag.Usage = func() {
		printUsageInfo()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	options := &CLIOptions{
		InteractiveMode: *interactive,
	}

	if len(os.Args) == 1 {
		handleNakedMode(options)
	}

	if *interactive {
		return handleInteractiveMode(options)
	}
	return handleDirectMode(options)
}

func printUsageInfo() {
	fmt.Fprintf(os.Stderr, "Usage of File Mover:\n\n")
	fmt.Fprintf(os.Stderr, "  Used for moving files based on a regex syntax:\n")
	fmt.Fprintf(os.Stderr, "  How to use :\n")
	fmt.Fprintf(os.Stderr, "  Direct mode: %s [source folder] [regex pattern] [target folder]\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  Interactive mode: %s -i\n\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExample:\n")
	fmt.Fprintf(os.Stderr, "  %s ./source \".*\\.txt$\" ./target\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "  %s -i\n", filepath.Base(os.Args[0]))

}

func handleNakedMode(options *CLIOptions) {
	displayModeSelection(options)

	// printUsageInfo()
	// return handleInteractiveMode(options)
}

func displayModeSelection(options *CLIOptions) {
	fmt.Println("Welcome to FileMover!")
	fmt.Println("Select a mode:")
	fmt.Println("  i - Interactive mode")
	fmt.Println("  h - Help")
	fmt.Println("  q - Quit")
	fmt.Print("Enter your choice: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "i":
			fmt.Println("Entering interactive mode...")
			handleInteractiveMode(options)
		case "h":
			printUsageInfo()
			displayModeSelection(options)
		case "q":
			fmt.Println("Goodbye!")
			os.Exit(1)
		default:
			fmt.Println("Invalid option. Please try again.")
			displayModeSelection(options)
		}
	}
}

func handleInteractiveMode(options *CLIOptions) (*CLIOptions, error) {
	reader := bufio.NewReader(os.Stdin)

	// Get source folder
	fmt.Print("Enter source folder path: ")
	sourceFolder, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading source folder: %v", err)
	}
	options.SourceFolder = strings.TrimSpace(sourceFolder)

	// Validate source folder
	if !folderExists(options.SourceFolder) {
		return nil, fmt.Errorf("source folder does not exist: %s", options.SourceFolder)
	}

	// Get regex pattern
	fmt.Print("Enter regex pattern (e.g., .*\\.txt$ for txt files): ")
	regexPattern, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading regex pattern: %v", err)
	}
	options.RegexPattern = strings.TrimSpace(regexPattern)

	// Get target folder
	fmt.Print("Enter target folder path: ")
	targetFolder, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading target folder: %v", err)
	}
	options.TargetFolder = strings.TrimSpace(targetFolder)

	// Create target folder if it doesn't exist
	if !folderExists(options.TargetFolder) {
		fmt.Printf("Target folder doesn't exist. Creating: %s\n", options.TargetFolder)
		if err := os.MkdirAll(options.TargetFolder, 0755); err != nil {
			return nil, fmt.Errorf("failed to create target folder: %v", err)
		}
	}
	fmt.Println(options)
	return options, nil
}

func handleDirectMode(options *CLIOptions) (*CLIOptions, error) {
	args := flag.Args()
	if len(args) != 3 {
		flag.Usage()
		fmt.Println(len(args))
		return nil, fmt.Errorf("incorrect number of arguments")
	}

	options.SourceFolder = args[0]
	options.RegexPattern = args[1]
	options.TargetFolder = args[2]

	// Validate inputs
	if !folderExists(options.SourceFolder) {
		return nil, fmt.Errorf("source folder does not exist: %s", options.SourceFolder)
	}

	if !folderExists(options.TargetFolder) {
		fmt.Printf("Target folder doesn't exist. Creating: %s\n", options.TargetFolder)
		if err := os.MkdirAll(options.TargetFolder, 0755); err != nil {
			return nil, fmt.Errorf("failed to create target folder: %v", err)
		}
	}

	return options, nil
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
	fmt.Printf("validation successful\nsource folder :%v \ndest folder : %v \nregex : %v \n", currentFolder, targetFolder, regex)

	return currentFolder, regex, targetFolder

}

func getFilesForRegex(files []string, re string) []string {
	var matchedFiles []string

	pattern, err := regexp.Compile(re)

	if err != nil {
		fmt.Printf("Invalid regex Pattern %v \n", err)
		return matchedFiles
	}

	fmt.Printf("checking files for regex %v \n", re)
	for _, file := range files {
		fmt.Printf("checking for file %v \n", file)
		if pattern.MatchString(file) {
			matchedFiles = append(matchedFiles, file)
			fmt.Printf("file matched, appending to the list\n")
		} else {
			fmt.Printf("file didnot match the regex\n")
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
	fmt.Println("read files for the folder")
	fmt.Println(files)
	return files
}

func folderExists(path string) bool {
	info, error := os.Stat(path)
	fmt.Printf("\nChecking for folders Existance \n")
	if os.IsNotExist(error) {
		println("folder does not exist")
		return false
	}
	isDir := info.IsDir()
	if isDir {
		fmt.Println("Folder Exists")
	}
	return isDir
}

func ProcessFiles(options *CLIOptions) error {
	allFiles := getFileNames(options.SourceFolder)
	if len(allFiles) == 0 {
		fmt.Errorf("no files matches the regular expressions ")
	}

	matchedFiles := getFilesForRegex(allFiles, options.RegexPattern)
	if len(matchedFiles) == 0 {
		return fmt.Errorf("no files matched the regex pattern")
	}

	moveFiles(options.SourceFolder, matchedFiles, options.TargetFolder)
	return nil
}
