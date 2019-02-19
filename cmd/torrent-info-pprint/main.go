package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/anacrolix/envpprof"
	"github.com/anacrolix/tagflag"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/bradfitz/iter"
)

var flags struct {
	JustName    bool
	PieceHashes bool
	Files       bool
	tagflag.StartPos
}

type FileExtended struct {
	Length  int64
	PathMap []map[string]string
}

func fromFiles(files []metainfo.FileInfo) (fe []FileExtended) {
	for _, file := range files {
		fe = append(fe, fromFileInfo(file))
	}

	return fe
}

func fromFileInfo(fi metainfo.FileInfo) (fe FileExtended) {
	fe.Length = fi.Length
	fe.PathMap = make([]map[string]string, len(fi.Path))
	for index, path := range fi.Path {
		fe.PathMap[index] = map[string]string{path: url.QueryEscape(path)}
	}

	return fe
}

func processReader(r io.Reader) error {
	var info metainfo.Info
	decoder := bencode.NewDecoder(r)
	err := decoder.Decode(&info)
	if err != nil {
		return err
	}

	if flags.JustName {
		fmt.Printf("%s\n", info.Name)
		return nil
	}
	d := map[string]interface{}{
		"Name":        info.Name,
		"NumPieces":   info.NumPieces(),
		"PieceLength": info.PieceLength,
		"NumFiles":    len(info.UpvertedFiles()),
		"TotalLength": info.TotalLength(),
	}
	if flags.Files {
		d["Files"] = fromFiles(info.Files)
		d["FilesOriginal"] = info.Files
	}
	if flags.PieceHashes {
		d["PieceHashes"] = func() (ret []string) {
			for i := range iter.N(info.NumPieces()) {
				ret = append(ret, hex.EncodeToString(info.Pieces[i*20:(i+1)*20]))
			}
			return
		}()
	}
	b, _ := json.MarshalIndent(d, "", "  ")
	_, err = os.Stdout.Write(b)
	return err
}

func main() {
	defer envpprof.Stop()
	tagflag.Parse(&flags)
	err := processReader(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	if !flags.JustName {
		os.Stdout.WriteString("\n")
	}
}
