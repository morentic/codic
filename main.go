package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const ignoreFileName = ".filecountignore"

var extension = map[string][]string{
	"YAML": {".yaml", ".yml"},
}

var types = map[string]string{
	".go":   "Go",
	".py":   "Python",
	".js":   "JavaScript",
	".php":  "PHP",
	".css":  "CSS",
	".html": "HTML",
	".scss": "SCSS",
	".vue":  "Vue",
	".json": "JSON",
	".yaml": "YAML",
	".xml":  "XML",
}

var fileCount = map[string]int{
	".go":   0,
	".py":   0,
	".js":   0,
	".php":  0,
	".css":  0,
	".html": 0,
	".scss": 0,
	".vue":  0,
	".json": 0,
	".yaml": 0,
	".xml":  0,
}

var totalFiles int

func main() {
	root := "./testdata"

	if len(os.Args) >= 2 {
		root = os.Args[1]
	}

	ignoreFileContent, err := os.ReadFile(ignoreFileName)
	var ignoreList []string
	if err == nil {
		ignoreLines := string(ignoreFileContent)
		for line := range strings.SplitSeq(ignoreLines, "\n") {
			if line == "" {
				continue
			}
			ignoreList = append(ignoreList, line)
		}
	}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		ext := filepath.Ext(path)

		if slices.Contains(ignoreList, fileName) {
			return nil
		}

		for fileExt, fileType := range types {
			if extension[fileType] == nil {
				if fileExt == ext {
					fileCount[fileExt]++
				}
				continue
			}

			for _, polymorphicExtension := range extension[fileType] {
				if polymorphicExtension == ext {
					fileCount[fileExt]++
				}
			}
		}

		totalFiles++

		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	table := table.New()
	for ext, count := range fileCount {
		table.Row(types[ext], fmt.Sprintf("%d", count))
	}
	style := lipgloss.NewStyle().Bold(true)
	table.Row(style.Render("Total"), style.Render(fmt.Sprintf("%d", totalFiles)))
	table.BorderRow(true)
	fmt.Println(table.Render())
}
