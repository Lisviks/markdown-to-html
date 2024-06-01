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
	linkRegex := regexp.MustCompile(`(!?)\[(.*?)\]\((.*?)\)`)
	convertedLine := linkRegex.ReplaceAllStringFunc(line, func(s string) string {
		parts := linkRegex.FindStringSubmatch(s)

		if parts[1] == "!" {
			return s
		}
		return `<a href="` + parts[3] + `">` + parts[2] + `</a>`
	})

	return convertedLine
}

func inlineCode(line string) string {
	codeRegex := regexp.MustCompile("`([^`]+)`")
	convertedLine := codeRegex.ReplaceAllString(line, "<code>$1</code>")
	return convertedLine
}

func bold(line string) string {
	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*|__(.*?)__`)
	convertedLine := boldRegex.ReplaceAllString(line, "<strong>$1$2</strong>")
	return convertedLine
}

func italic(line string) string {
	italicRegex := regexp.MustCompile(`\*(.*?)\*|_(.*?)_`)
	convertedLine := italicRegex.ReplaceAllString(line, "<em>$1$2</em>")
	return convertedLine
}

func imageTag(line string) string {
	imageRegex := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	convertedLine := imageRegex.ReplaceAllString(line, `<img src="$2" alt="$1">`)
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
			line = bold(line)
			line = italic(line)
			line = imageTag(line)
			if strings.HasPrefix(line, "#") {
				html += headings(line) + "\n"
			} else if strings.HasPrefix(line, "![") {
				html += imageTag(line) + "\n"
			} else if len(line) >= 1 {
				html += paragraph(line) + "\n"
			}
		}
	}

	data := []byte(html)

	var fileName string
	path := filepath.Join(".", "out")
	err = os.MkdirAll(path, 0755)

	if err != nil {
		log.Fatal(err)
	}

	if len(args) == 2 {
		fileName = args[1] + ".html"
	} else {
		fileName = inputFileName + ".html"
	}

	e := os.WriteFile("out/"+fileName, data, 0644)
	if err != nil {
		log.Fatal(e)
	}
	fmt.Println("HTML successfully written to", fileName)
}
