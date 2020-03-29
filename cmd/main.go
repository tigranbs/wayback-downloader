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
	domainVar    string
	outputVar    string
	yearVar      int
	startYearVar int
	endYearVar   int
)

func init() {
	BaseCMD.Flags().StringVarP(&domainVar, "domain", "d", "", "Domain name of target (required)")
	BaseCMD.Flags().IntVarP(&yearVar, "year", "y", 0, "Year for downloading target")
	BaseCMD.Flags().IntVarP(&startYearVar, "startYear", "s", 2006, "Start Year for downloading target")
	BaseCMD.Flags().IntVarP(&endYearVar, "endYear", "e", 2006, "End Year for downloading target")
	BaseCMD.Flags().StringVarP(&outputVar, "output", "o", "", "Downloaded output directory")

	BaseCMD.MarkFlagRequired("domain")
}

func mainCMD(cmd *cobra.Command, args []string) {
	if len(outputVar) == 0 {
		outputVar = domainVar
	}

	if yearVar > 0 {
		startYearVar = yearVar
		endYearVar = yearVar
	}

	for currentYear := startYearVar; currentYear <= endYearVar; currentYear++ {
		waybackURL, err := wayback.GetAnnualURL(domainVar, currentYear)

		// Request the HTML page.
		content, err := wayback.GetHTML(waybackURL)
		if err != nil {
			continue
		}

		log.Println("Downloading Domain:", domainVar, "Year:", currentYear)
		outputDirectory := path.Join(outputVar, strconv.Itoa(currentYear))

		_, err = os.Stat(outputDirectory)
		if err != nil {
			os.MkdirAll(outputDirectory, 0777)
		}

		content, err = wayback.FetchAssets(domainVar, content, path.Join(outputDirectory, "assets"))
		if err != nil {
			log.Fatal(err)
		}

		ioutil.WriteFile(path.Join(outputDirectory, "index.html"), []byte(content), 0644)
		log.Println("Completed Domain:", domainVar, "Year:", currentYear)
	}
}
