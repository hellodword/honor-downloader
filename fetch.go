package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Anna struct {
	TorrentLink  string
	Filename     string
	RealFilename string
}

var (
	regexTorrentLink = regexp.MustCompile(`(?m)href="(/[^" \n]+\.torrent)"`)
	regexFilename    = regexp.MustCompile(`(?m)&nbsp;file&nbsp;.([^<]+).<`)
)

func fetchAnna(link string) (*Anna, error) {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}

	// TODO UA cookies fingerprinting for CF

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var anna Anna

	bs := regexTorrentLink.FindSubmatch(b)
	if len(bs) != 2 {
		return nil, errors.New("unable to find torrent link")
	}
	anna.TorrentLink = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, string(bs[1]))

	bs = regexFilename.FindSubmatch(b)
	if len(bs) != 2 {
		return nil, errors.New("unable to find filename")
	}
	anna.Filename = string(bs[1])
	if strings.Contains(anna.Filename, "&nbsp;(extract)") {
		return nil, fmt.Errorf("illegal filename: %s", anna.Filename)
	}

	return &anna, nil
}
