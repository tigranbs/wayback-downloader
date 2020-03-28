package wayback

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetAnnualURL calls Wayback API and tries to fetch latest website URL based on given Year
func GetAnnualURL(domain string, year int) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://archive.org/wayback/available?url=%s&timestamp=%d1231", domain, year))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("Invalid status code %d", res.StatusCode)
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	result := &WaybackMachineResult{}

	err = json.Unmarshal(responseData, result)
	if err != nil {
		return "", err
	}

	return result.ArchivedSnapshots.Closest.URL, nil
}
