package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func parseGitignore(gitignorePath string) []string {
	ignorePatterns := []string{}

	file, err := os.Open(gitignorePath)
	if err != nil {
		return ignorePatterns
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			ignorePatterns = append(ignorePatterns, line)
		}
	}

	return ignorePatterns
}

func shouldIgnore(path string, ignoreDirs []string, ignoreExts []string, gitignorePatterns []string) bool {
	fileName := filepath.Base(path)
	if strings.HasPrefix(fileName, ".") {
		return true
	}

	pathParts := strings.Split(path, string(os.PathSeparator))
	for _, ignoredDir := range ignoreDirs {
		for _, part := range pathParts {
			if part == ignoredDir {
				return true
			}
		}
	}

	for _, ext := range ignoreExts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	for _, pattern := range gitignorePatterns {
		matched, _ := filepath.Match(pattern, path)
		if matched {
			return true
		}
	}

	return false
}

func generateDirectoryStructure(sourcePath string, outputPath string, ignoreDirs []string, ignoreExts []string) {
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		fmt.Printf("Error: Source path %s does not exist.\n", sourcePath)
		return
	}

	gitignorePath := filepath.Join(sourcePath, ".gitignore")
	gitignorePatterns := parseGitignore(gitignorePath)

	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		fmt.Printf("Error getting absolute path: %s\n", err)
		return
	}
	rootName := filepath.Base(absPath)
	if rootName == "." {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %s\n", err)
			return
		}
		rootName = filepath.Base(currentDir)
	}
	structure := []string{rootName}

	var buildTree func(path, prefix string)
	buildTree = func(path, prefix string) {
		items, _ := os.ReadDir(path)

		dirs := []string{}
		files := []string{}

		for _, item := range items {
			if item.IsDir() {
				dirs = append(dirs, item.Name())
			} else {
				files = append(files, item.Name())
			}
		}

		sort.Strings(dirs)
		sort.Strings(files)

		allItems := append(dirs, files...)

		for i, item := range allItems {
			fullPath := filepath.Join(path, item)

			if shouldIgnore(fullPath, ignoreDirs, ignoreExts, gitignorePatterns) {
				continue
			}

			isLast := (i == len(allItems)-1)
			connector := "└── "
			if !isLast {
				connector = "├── "
			}

			info, _ := os.Stat(fullPath)
			if info.IsDir() {
				structure = append(structure, fmt.Sprintf("%s%s%s", prefix, connector, item))
				if isLast {
					buildTree(fullPath, prefix+"    ")
				} else {
					buildTree(fullPath, prefix+"│   ")
				}
			} else {
				structure = append(structure, fmt.Sprintf("%s%s- %s", prefix, connector, item))
			}
		}
	}

	buildTree(sourcePath, "")

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error writing to %s: %s\n", outputPath, err)
		return
	}
	defer file.Close()

	file.WriteString(strings.Join(structure, "\n"))
	fmt.Printf("Project structure written to %s\n", outputPath)
}

func main() {
	source := flag.String("s", "/", "Absolute path to the project folder")
	output := flag.String("o", "project_structure", "Path to save the output markdown file")
	iDirectory := flag.String("id", "", "Directories to ignore, comma-separated")
	iExtension := flag.String("ie", "", "Extensions to ignore, comma-separated")

	flag.Parse()

	ignoreDirs := []string{}
	if *iDirectory != "" {
		for _, d := range strings.Split(*iDirectory, ",") {
			ignoreDirs = append(ignoreDirs, strings.TrimSpace(d))
		}
	}

	ignoreExts := []string{}
	if *iExtension != "" {
		for _, e := range strings.Split(*iExtension, ",") {
			e = strings.TrimSpace(e)
			if !strings.HasPrefix(e, ".") {
				e = "." + e
			}
			ignoreExts = append(ignoreExts, e)
		}
	}

	outputPath := *output
	if !strings.HasSuffix(outputPath, ".md") {
		outputPath += ".md"
	}

	generateDirectoryStructure(*source, outputPath, ignoreDirs, ignoreExts)
}
