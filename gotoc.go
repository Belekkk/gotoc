package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"path/filepath"
	"github.com/fatih/color"
)

type Toc struct {
	// DocToc
	depth  []int
	spaces []string
	title  []string
	link   []string
}

func getCurrentDir() string {
	pwd, err := os.Getwd()
    if err != nil {
		log.Fatal(err)
    }
    return pwd
}

func findAllMdFiles(searchDir string) []string{
	fileList := []string{}
	    filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			if filepath.Ext(f.Name()) == ".md" {
	        	fileList = append(fileList, path)
			}
	        return nil
	    })
		return fileList
}

func headerTitle() string {
	header := "\n**Table of Contents** *(by [`gotoc`](https://github.com/axelbellec/gotoc))*"
	return header
}

func headerCredits() string {
	return "\n<!-- Table of Contents generated by [gotoc](https://github.com/axelbellec/gotoc) -->"
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

func formatFilename(filename string, wd string) string {
	return strings.Replace(filename, wd, "", -1)
}

func computeToc(filename string, maxDepth int, title string, noTitle bool, wd string) {
	var startsWith bool
	var line string
	var doctoc = new(Toc)
	var fileContent []string
	var tocLevel int

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tocScanner := bufio.NewScanner(file)
	if err := tocScanner.Err(); err != nil {
		log.Fatal(err)
	}

	for tocScanner.Scan() {
		line = tocScanner.Text()
		startsWith = strings.HasPrefix(line, "#")
		if startsWith == true {
			tocLevel = tocLevel + 1
		}
	}

	file, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if tocLevel == 0 {
		os.Stderr.WriteString(fmt.Sprintf("File : '%s'\n", formatFilename(filename, wd)))
		color.Yellow("No headings founded\n")
		return
	} else {
		os.Stderr.WriteString(fmt.Sprintf("File : '%s'\n", formatFilename(filename, wd)))
		fileScanner := bufio.NewScanner(file)
		if err := fileScanner.Err(); err != nil {
			log.Fatal(err)
		}
		for fileScanner.Scan() {
			line = fileScanner.Text()
			startsWith = strings.HasPrefix(line, "#")
			if startsWith == true {
				headerCount := strings.Count(line, "#")

		 		if headerCount <= maxDepth {
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
	}

	err = os.Remove(filename)
	if err != nil {
		log.Fatal(err)
	}

	if noTitle == true {
		title = ""
	}

	newFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	if _, err = newFile.WriteString(createHeader(title)); err != nil {
		log.Fatal(err)
	}

	for i := range doctoc.depth {
		doctocContent := fmt.Sprintf("%s- [%s](%s) \n", doctoc.spaces[i], doctoc.title[i], doctoc.link[i])
		if _, err = newFile.WriteString(doctocContent); err != nil {
			log.Fatal(err)
		}
	}

	color.Green(fmt.Sprintf("'%s' will be updated\n", formatFilename(filename, wd)))
	for i := range fileContent {
		originalContent := fmt.Sprintf("%s\n", fileContent[i])
		if _, err = newFile.WriteString(originalContent); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func main() {

	allDir := flag.Bool("dir", false, "a bool")
	filename := flag.String("file", "", "a string")
	maxDepth := flag.Int("depth", 3, "an int")
	title := flag.String("title", headerTitle(), "a string")
	noTitle := flag.Bool("notitle", false, "a bool")
	flag.Parse()

	os.Stderr.WriteString("Generating Table of Contents\n")
	os.Stderr.WriteString(fmt.Sprintf("%s\n", outputLine(28)))

	currDir := getCurrentDir()
	color.Magenta(fmt.Sprintf("Working Directory : %s\n", currDir))

	if *allDir == true {
		mdFiles := findAllMdFiles(currDir)
		for _, md := range mdFiles {
			computeToc(md, *maxDepth, *title, *noTitle, currDir)
		}
	} else {
		computeToc(*filename, *maxDepth, *title, *noTitle, currDir)
	}
	os.Stderr.WriteString("Done\n")
}
