package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Request the HTML page.
	res, err := http.Get("http://web.archive.org/web/20060101014129/http://facebook.com:80/")
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
	tmpLine := append(lines[:htmlStartIndex+1], lines[headRewriteIndex+1:removeLineStart]...)

	tmpLine2 := append(lines[headRewriteIndex+2:removeLineStart], lines[removeLineEnd+1:endOfHTML+1]...)
	lines = append(tmpLine, tmpLine2...)
	ioutil.WriteFile("fb.html", []byte(strings.Join(lines, "\n")), 0644)
}
