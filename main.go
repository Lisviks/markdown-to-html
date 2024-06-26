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
	regex := regexp.MustCompile(`(!?)\[(.*?)\]\((.*?)\)`)
	line = regex.ReplaceAllStringFunc(line, func(s string) string {
		parts := regex.FindStringSubmatch(s)

		if parts[1] == "!" {
			return s
		}
		return `<a href="` + parts[3] + `">` + parts[2] + `</a>`
	})

	return line
}

func inlineCode(line string) string {
	regex := regexp.MustCompile("`([^`]+)`")
	line = regex.ReplaceAllString(line, "<code>$1</code>")
	return line
}

func bold(line string) string {
	regex := regexp.MustCompile(`\*\*(.*?)\*\*|__(.*?)__`)
	line = regex.ReplaceAllString(line, "<strong>$1$2</strong>")
	return line
}

func italic(line string) string {
	regex := regexp.MustCompile(`\*(.*?)\*|_(.*?)_`)
	line = regex.ReplaceAllString(line, "<em>$1$2</em>")
	return line
}

func imageTag(line string) string {
	regex := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
	line = regex.ReplaceAllString(line, `<img src="$2" alt="$1">`)
	return line
}

func unorderedListItem(line string) string {
	li := strings.Replace(line, "- ", "<li>", 1)
	li += "</li>\n"
	return li
}

func orderedListItem(line string) string {
	regex := regexp.MustCompile(`^(?:\s*)(\d+)\.\s+(.+)`)
	li := regex.ReplaceAllString(line, `<li>$2`)
	li += "</li>\n"
	return li
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

	html := "<article>\n"
	inCodeBlock := false
	inUnorderedListBlock := false
	inOnrderedListBlock := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		orderedListRegex := regexp.MustCompile(`^(?:\s*)(\d+)\.\s+(.+)`)
		if inCodeBlock || inUnorderedListBlock || inOnrderedListBlock {
			if inCodeBlock {
				if strings.HasPrefix(line, "```") {
					inCodeBlock = false
					html += "</code></pre>\n"
					continue
				}
				html += line + "\n"
			} else if inUnorderedListBlock {
				if !strings.HasPrefix(line, "- ") {
					inUnorderedListBlock = false
					html += "</ul>\n"
					continue
				}
				html += unorderedListItem(line)
			} else if inOnrderedListBlock {
				if !orderedListRegex.MatchString(line) {
					inOnrderedListBlock = false
					html += "</ol>\n"
					continue
				}
				html += orderedListItem(line)
			}
		} else {
			if strings.HasPrefix(line, "```") {
				inCodeBlock = true
				html += "<pre><code>"
				continue
			}
			if strings.HasPrefix(line, "- ") {
				inUnorderedListBlock = true
				html += "<ul>\n" + unorderedListItem(line)
				continue
			}
			if orderedListRegex.MatchString(line) {
				inOnrderedListBlock = true
				html += "<ol>\n" + orderedListItem(line)
				continue
			}

			line = anchor(line)
			line = inlineCode(line)
			line = bold(line)
			line = italic(line)
			line = imageTag(line)
			if strings.HasPrefix(line, "#") {
				html += headings(line) + "\n"
			} else if len(line) >= 1 && !strings.HasPrefix(line, "<") {
				html += paragraph(line) + "\n"
			} else if strings.HasPrefix(line, "<") {
				html += line + "\n"
			}
		}
	}

	html += "</article>\n"
	data := []byte(html)

	var fileName string
	var path string
	if len(args) == 3 {
		path = filepath.Join(".", args[2])
	} else {
		path = filepath.Join(".", "out")
	}
	err = os.MkdirAll(path, 0755)

	if err != nil {
		log.Fatal(err)
	}

	if len(args) == 2 {
		fileName = args[1] + ".html"
	} else {
		fileName = inputFileName + ".html"
	}

	e := os.WriteFile(path+"/"+fileName, data, 0644)
	if err != nil {
		log.Fatal(e)
	}
	fmt.Println("HTML successfully written to", fileName)
}
