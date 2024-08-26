package main

import (
	"fmt"
	"net/http"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
)

func main() {
	var output_dir = "data"

	// https://annas-archive.li/md5/08b0f97b98c977da93cd5e5623686af5
	var only_files = "08b0f97b98c977da93cd5e5623686af5"
	var url = "https://annas-archive.li/dyn/small_file/torrents/external/libgen_rs_non_fic/r_4319000.torrent"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	mi, err := metainfo.Load(resp.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", mi)

	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.NoDHT = true
	clientConfig.NoDefaultPortForwarding = true
	// clientConfig.DisableWebseeds = true
	// clientConfig.DisableWebtorrent = true
	clientConfig.DataDir = output_dir

	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	t, err := client.AddTorrent(mi)
	if err != nil {
		panic(err)
	}

	files := t.Files()

	var pending = make(map[int]struct{})
	for _, f := range files {
		if f.DisplayPath() == only_files {
			begin, end := f.BeginPieceIndex(), f.EndPieceIndex()
			for i := begin; i < end; i++ {
				pending[i] = struct{}{}
			}
			t.DownloadPieces(begin, end)
			break
		}
	}
	fmt.Println(pending)

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
			// fmt.Printf("%s: %d\t%s: %+v\n", "index", ev.Index, "state", ev.PieceState)
			if ev.PieceState.Completion == expected {
				delete(pending, ev.Index)
				if len(pending) == 0 {
					break
				}
			}
		}
	}

}
