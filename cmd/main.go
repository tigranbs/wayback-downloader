package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"wayback-downloader/wayback"

	"github.com/spf13/cobra"
)

// BaseCMD is the main entry point for CLI execution
var BaseCMD = &cobra.Command{
	Use:   fmt.Sprintf("%s <domain> <year>", os.Args[0]),
	Short: "Archive.org Wayback Machine downloader",
	Run:   mainCMD,
}

var (
	domainVar string
	outputVar string
	yearVar   int
)

func init() {
	BaseCMD.Flags().StringVarP(&domainVar, "domain", "d", "", "Domain name of target (required)")
	BaseCMD.Flags().IntVarP(&yearVar, "year", "y", 2006, "Year for downloading target (required)")
	BaseCMD.Flags().StringVarP(&outputVar, "output", "o", "", "Downloaded output directory (required)")

	BaseCMD.MarkFlagRequired("domain")
	BaseCMD.MarkFlagRequired("year")
}

func mainCMD(cmd *cobra.Command, args []string) {
	waybackURL, err := wayback.GetAnnualURL(domainVar, yearVar)

	// Request the HTML page.
	content, err := wayback.GetHTML(waybackURL)
	if err != nil {
		log.Fatal(err)
	}

	if len(outputVar) == 0 {
		outputVar = path.Join(domainVar, strconv.Itoa(yearVar))
	}

	_, err = os.Stat(outputVar)
	if err != nil {
		os.MkdirAll(outputVar, 0777)
	}

	content, err = wayback.FetchAssets(domainVar, content, path.Join(outputVar, "assets"))
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(path.Join(outputVar, "index.html"), []byte(content), 0644)
}
