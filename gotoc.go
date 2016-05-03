package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Toc struct {
	// DocToc
	depth  []int
	spaces []string
	title  []string
	link   []string
}

func headerTitle() string {
	header := "\n**Table of Contents** (by [`gotoc`](https://github.com/Belekkk/gotoc))"
	return header
}

func headerCredits() string {
	return "\n<!-- Table of Contents generated by [gotoc](https://github.com/Belekkk/gotoc) -->"
}

func createHeader(title string) string {
	contentHeader := fmt.Sprintf("%s %s \n", headerCredits(), title)
	return contentHeader
}

func removeSpecialCharacters(str string) string {
	r, _ := regexp.Compile("([a-zA-Z]+)")
	cleanedStr := r.FindAllString(str, -1)
	return strings.Join(cleanedStr, " ")
}

func formatLink(str string) string {
	link := strings.Replace(str, " ", "-", -1)
	link = strings.ToLower(link)
	link = strings.Join([]string{"#", link}, "")
	return link
}

func computeSpaces(depth int) string {
	spaces := strings.Repeat(" ", depth)
	return spaces
}

func outputLine(n int) string {
	line := strings.Repeat("-", n)
	return line
}

func main() {
	var startsWith bool
	var line string
	var doctoc = new(Toc)
	var fileContent []string

	filename := flag.String("file", "", "a string")
	maxDepth := flag.Int("depth", 3, "an int")
	title := flag.String("title", headerTitle(), "a string")
	noTitle := flag.Bool("notitle", false, "a bool")
	flag.Parse()

	file, err := os.OpenFile(*filename, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		startsWith = strings.HasPrefix(line, "#")
		if startsWith == true {
			headerCount := strings.Count(line, "#")

			if headerCount <= *maxDepth {
				content := line[headerCount+1:]
				doctoc.depth = append(doctoc.depth, headerCount)
				indentation := (headerCount - 1) * 2
				nSpaces := computeSpaces(indentation)
				cleanedContent := removeSpecialCharacters(content)
				doctoc.spaces = append(doctoc.spaces, nSpaces)
				doctoc.title = append(doctoc.title, content)
				doctoc.link = append(doctoc.link, formatLink(cleanedContent))
			}
		}
		fileContent = append(fileContent, line)
	}

	err = os.Remove(*filename)
	if err != nil {
		log.Fatal(err)
	}

	if *noTitle == true {
		*title = ""
	}

	newFile, err := os.OpenFile(*filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	if _, err = newFile.WriteString(createHeader(*title)); err != nil {
		log.Fatal(err)
	}
	os.Stderr.WriteString("Generating Table of Contents\n")
	os.Stderr.WriteString(fmt.Sprintf("%s\n", outputLine(28)))
	for i := range doctoc.depth {
		doctocContent := fmt.Sprintf("%s- [%s](%s) \n", doctoc.spaces[i], doctoc.title[i], doctoc.link[i])
		if _, err = newFile.WriteString(doctocContent); err != nil {
			log.Fatal(err)
		}
	}

	os.Stderr.WriteString(fmt.Sprintf("'%s' will be updated\n", *filename))
	for i := range fileContent {
		originalContent := fmt.Sprintf("%s\n", fileContent[i])
		if _, err = newFile.WriteString(originalContent); err != nil {
			log.Fatal(err)
		}
	}
	os.Stderr.WriteString("Done\n")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
