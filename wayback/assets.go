package wayback

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	urlParser "net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FetchAssets downloading and saving assets to specified output directory
func FetchAssets(domain, htmlContent, outputDir string) (string, error) {
	_, err := os.Stat(outputDir)
	if err != nil {
		os.Mkdir(outputDir, 0777)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(htmlContent))
	if err != nil {
		return "", err
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		s.SetAttr("href", "")
	})

	doc.Find("link[href], *[src]").Each(func(i int, s *goquery.Selection) {
		processAssetNode(domain, s, outputDir)
	})

	newContent, err := doc.Html()
	if err != nil {
		return "", err
	}

	return newContent, nil
}

func processAssetNode(domain string, s *goquery.Selection, outputDir string) {
	linkType, _ := s.Attr("type")
	hrefURL, isHref := s.Attr("href")
	srcURL, isSrc := s.Attr("src")

	tagName := strings.ToLower(goquery.NodeName(s))
	finalFileName := ""

	if isSrc && srcURL[0] == '/' {
		srcURL = fmt.Sprintf("http://web.archive.org%s", srcURL)
	} else if isHref && hrefURL[0] == '/' {
		hrefURL = fmt.Sprintf("http://web.archive.org%s", hrefURL)
	}

	if tagName == "script" || tagName == "img" {
		finalFileName = downloadAsset(srcURL, outputDir)
	} else if linkType == "text/css" {
		finalFileName = downloadAsset(hrefURL, outputDir)
	} else {
		return
	}

	if len(finalFileName) > 0 {
		attrName := "href"
		if isSrc {
			attrName = "src"
		}

		s.SetAttr(attrName, finalFileName)
	}

}

func downloadAsset(url, outputDir string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))

	up, err := urlParser.Parse(url)
	if err != nil {
		return ""
	}

	fileName := fmt.Sprintf("%s%s", hex.EncodeToString(hasher.Sum(nil)), filepath.Ext(up.Path))

	filePath := path.Join(outputDir, fileName)
	_, err = os.Stat(filePath)
	if err == nil {
		return filePath
	}

	res, err := http.Get(url)
	if err != nil {
		log.Println("ASSET FETCH: HTTP Request Error -> ", url, err.Error())
		return ""
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("ASSET FETCH: HTTP status error %d %s \n", res.StatusCode, res.Status)
		return ""
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("ASSET FETCH: ERROR creating a file -> ", url, err.Error())
		return ""
	}

	_, err = io.Copy(f, res.Body)

	if err != nil {
		log.Println("ASSET FETCH: ERROR downloading -> ", url, err.Error())
		return ""
	}

	return path.Join("assets", fileName)
}
