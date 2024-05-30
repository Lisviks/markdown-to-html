package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func headings(line string) string {
	count := strings.Count(strings.SplitN(line, " ", 2)[0], "#")
	content := strings.TrimSpace(line[count:])
	return "<h" + strconv.Itoa(count) + ">" + content + "</h" + strconv.Itoa(count) + ">"
}

func paragraph(line string) string {
	return "<p>" + line + "</p>"
}

func anchor(line string) string {
	linkRegex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	convertedLine := linkRegex.ReplaceAllString(line, `<a href="$2">$1</a>`)
	return convertedLine
}

func inlineCode(line string) string {
	codeRegex := regexp.MustCompile("`([^`]+)`")
	convertedLine := codeRegex.ReplaceAllString(line, "<code>$1</code>")
	return convertedLine
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No markdown file provided.")
		return
	}

	filePath := args[0]
	fileNameWithExt := filepath.Base(filePath)
	ext := filepath.Ext(fileNameWithExt)
	inputFileName := strings.TrimSuffix(fileNameWithExt, ext)

	if ext != ".md" {
		fmt.Println("Provide a markdown file.")
		return
	}
	file, err := os.Open(filePath)

	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()

	html := ""
	inCodeBlock := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if inCodeBlock {
			if strings.HasPrefix(line, "```") {
				inCodeBlock = false
				html += "</code></pre>\n"
				continue
			}
			html += line + "\n"
		} else {
			if strings.HasPrefix(line, "```") {
				inCodeBlock = true
				html += "<pre><code>"
				continue
			}

			line = anchor(line)
			line = inlineCode(line)
			if strings.HasPrefix(line, "#") {
				html += headings(line) + "\n"
			} else if len(line) >= 1 {
				html += paragraph(line) + "\n"
			}
		}
	}

	data := []byte(html)

	fileName := inputFileName + ".html"
	e := os.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Fatal(e)
	}
	fmt.Println("HTML successfully written to", fileName)
}
