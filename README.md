# wayback-downloader
Archive.org Wayback Machine downloader CLI written in GO

## Installation

```
go install github.com/tigranbs/wayback-downloader@latest
```

## Usage

```bash
~> wayback-downloader --help
Archive.org Wayback Machine downloader

Usage:
  wayback-downloader <domain> <year> [flags]

Flags:
  -d, --domain string   Domain name of target (required)
  -e, --endYear int     End Year for downloading target (default 2006)
  -h, --help            help for /tmp/go-build3312683796/b001/exe/wayback-downloader
  -o, --output string   Downloaded output directory
  -s, --startYear int   Start Year for downloading target (default 2006)
  -y, --year int        Year for downloading target
```
