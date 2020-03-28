package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	content, err := wayback.GetHTML(waybackURL)
	if err != nil {
		log.Fatal(err)
	}

	content, err = wayback.FetchAssets(content, "./assets")
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile("fb.html", []byte(content), 0644)
}
