package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

var (
	errTorrentNoInfo = errors.New("torrent no info")
	errUnknown       = errors.New("unknown error")
)

func downloadTorrent(r io.Reader, tmpDir, targetPath string, fileFilter func(path string) bool) error {
	var err error
	var srcPath string

	if fileFilter == nil {
		return errors.New("fileFilter is nil")
	}

	mi, err := metainfo.Load(r)
	if err != nil {
		return err
	}

	if len(mi.InfoBytes) == 0 {
		return errTorrentNoInfo
	}

	hash := mi.HashInfoBytes().HexString()

	if tmpDir == "" {
		tmpDir, err = os.MkdirTemp(os.TempDir(), hash)
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)
	}

	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.NoDHT = true
	clientConfig.NoDefaultPortForwarding = true
	// clientConfig.DisableWebseeds = true
	// clientConfig.DisableWebtorrent = true
	clientConfig.DataDir = tmpDir
	clientConfig.DefaultStorage = NewStorage(clientConfig.DataDir, func(path string) bool {
		b := fileFilter(path)
		if b {
			srcPath = path
		}
		return b
	})

	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	t, err := client.AddTorrent(mi)
	if err != nil {
		return err
	}

	var pending = make(map[int]struct{})

	files := t.Files()
	for _, f := range files {
		if fileFilter(f.Path()) {
			begin, end := f.BeginPieceIndex(), f.EndPieceIndex()
			for i := begin; i < end; i++ {
				pending[i] = struct{}{}
			}
			t.DownloadPieces(begin, end)
			break
		}
	}

	sub := t.SubscribePieceStateChanges()
	defer sub.Close()

	expected := storage.Completion{
		Complete: true,
		Ok:       true,
	}

	for i := range pending {
		if t.PieceState(i).Completion == expected {
			delete(pending, i)
		}
	}

	if len(pending) > 0 {
		for ev := range sub.Values {
			if _, ok := pending[ev.Index]; !ok {
				continue
			}
			fmt.Printf("%s: %d\t%s: %+v\n", "piece", ev.Index, "state", ev.PieceState)
			if ev.PieceState.Completion == expected {
				delete(pending, ev.Index)
				if len(pending) == 0 {
					break
				}
			}
		}
	}

	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return errUnknown
	}

	return mv(srcPath, targetPath)
}

func mv(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't remove source file: %v", err)
	}
	return nil
}
