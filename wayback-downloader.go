package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"wayback-downloader/wayback"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("USAGE: %s <domain name> <year> \n", os.Args[0])
		os.Exit(0)
	}

	domain := os.Args[1]
	yearStr := os.Args[2]

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		fmt.Printf("Unable to parse Year argument, it's should be a number -> %s \n", err.Error())
		fmt.Printf("USAGE: %s <domain name> <year> \n", os.Args[0])
		os.Exit(0)
	}

	waybackURL, err := wayback.GetAnnualURL(domain, year)

	// Request the HTML page.
	res, err := http.Get(waybackURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	bodyBuf := bytes.NewBuffer([]byte{})
	io.Copy(bodyBuf, res.Body)

	lines := strings.Split(bodyBuf.String(), "\n")
	removeLineStart := -1
	removeLineEnd := -1
	htmlStartIndex := -1
	headRewriteIndex := -1
	endOfHTML := -1

	for i, line := range lines {
		if strings.Contains(line, "<html") {
			htmlStartIndex = i
		}

		if strings.Contains(line, "</html>") {
			endOfHTML = i
		}

		if strings.Contains(line, "<!-- End Wayback Rewrite JS Include -->") {
			headRewriteIndex = i
		}

		if strings.Contains(line, "<!-- BEGIN WAYBACK TOOLBAR INSERT -->") {
			removeLineStart = i
		}

		if strings.Contains(line, "<!-- END WAYBACK TOOLBAR INSERT -->") {
			removeLineEnd = i
		}
	}

	lines[htmlStartIndex] = lines[htmlStartIndex] + "<head>"
	tmpLine := append(lines[:htmlStartIndex+1], lines[headRewriteIndex+1:removeLineStart+1]...)

	tmpLine = append(tmpLine, lines[removeLineEnd+1:endOfHTML+1]...)
	ioutil.WriteFile("fb.html", []byte(strings.Join(tmpLine, "\n")), 0644)
}
