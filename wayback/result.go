package wayback

// WaybackMachineResult is the main type for and API response
type WaybackMachineResult struct {
	Timestamp         string          `json:"timestamp"`
	URL               string          `json:"url"`
	ArchivedSnapshots archiveSnapshot `json:"archived_snapshots"`
}

type archiveSnapshot struct {
	Closest struct {
		Timestamp string `json:"timestamp"`
		URL       string `json:"url"`
		Available bool   `json:"available"`
		Status    string `json:"status"`
	} `json:"closest"`
}
