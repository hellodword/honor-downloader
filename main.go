package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	pwd, _ := os.Getwd()
	if pwd == "" {
		pwd = "."
	}

	annaUrl := flag.String("anna", "", "anna url")
	realFilename := flag.String("name", "", "real filename")
	dataDir := flag.String("dir", pwd, "data dir")
	verbose := flag.Bool("verbose", true, "show verbose log")
	version := flag.Bool("version", false, "show version")
	flag.Parse()

	if *version {
		printVersion()
		os.Exit(0)
	}

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if *annaUrl == "" {
		flag.Usage()
		os.Exit(1)
	}

	slog.Debug("fetching anna", "url", *annaUrl)
	anna, err := fetchAnna(*annaUrl)
	if err != nil {
		slog.Error(
			"error occurred while fetching anna",
			"error", err,
		)
		os.Exit(1)
	}
	anna.RealFilename = *realFilename
	if anna.RealFilename == "" {
		anna.RealFilename = filepath.Base(anna.Filename)
	}

	slog.Debug("fetching torrent file", "TorrentLink", anna.TorrentLink)
	resp, err := http.Get(anna.TorrentLink)
	if err != nil {
		slog.Error(
			"error occurred while fetching torrent file",
			"error", err,
		)
		os.Exit(1)
	}
	defer resp.Body.Close()

	targetPath := filepath.Join(*dataDir, anna.RealFilename)
	slog.Debug("downloading anna", "targetPath", targetPath, "Filename", anna.Filename)
	err = downloadTorrent(resp.Body, "", targetPath, func(path string) bool {
		return filepath.Base(path) == anna.Filename
	})
	if err != nil {
		slog.Error(
			"error occurred while downloading the torrent",
			"error", err,
		)
		os.Exit(1)
	}
}
