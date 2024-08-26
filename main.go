package main

import (
	"net/http"
	"path/filepath"
)

func main() {

	anna, err := fetchAnna("https://annas-archive.li/md5/08b0f97b98c977da93cd5e5623686af5")
	if err != nil {
		panic(err)
	}
	anna.RealFilename = "The Embodied Soul: Aristotelian Psychology and Physiology in Medieval Europe between 1200 and 1420.epub"

	var dataDir = "data"

	resp, err := http.Get(anna.TorrentLink)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = downloadTorrent(resp.Body, "", filepath.Join(dataDir, anna.RealFilename), func(path string) bool {
		return filepath.Base(path) == anna.Filename
	})
	if err != nil {
		panic(err)
	}

}
