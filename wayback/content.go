package wayback

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GetHTML handles url and after fetching HTML content tries to clean up
func GetHTML(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Unable to make an HTTP request -> %s", err.Error())
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	bodyBuf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(bodyBuf, res.Body)
	if err != nil {
		return "", fmt.Errorf("Unable fetch content from HTTP request -> %s", err.Error())
	}

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
	contentLines := append(lines[:htmlStartIndex+1], lines[headRewriteIndex+1:removeLineStart+1]...)
	contentLines = append(contentLines, lines[removeLineEnd+1:endOfHTML+1]...)
	return strings.Join(contentLines, "\n"), nil
}
