package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	torrentParser "github.com/j-muller/go-torrent-parser"
	bencode "github.com/jackpal/bencode-go"
)

func main() {
	arg := os.Args[1:]

	rand.Read(PEER_ID)
	file, err := os.Open(arg[0])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	path = arg[1]
	info = bencodeTorrent{}
	err = bencode.Unmarshal(file, &info)
	if err != nil {
		panic(err)
	}
	pieceLength = info.Info.PieceLength

	for i := range info.Info.Files {
		info.Info.Length += info.Info.Files[i].Length
	}
	lastLen := info.Info.Length % pieceLength
	if lastLen == 0 {
		lastLen = pieceLength
	}

	piece = make([]*Piece, len(info.Info.Pieces)/20)
	for i := 0; i < len(info.Info.Pieces); i += 20 {
		var v Piece
		v.index = i / 20
		v.length = info.Info.PieceLength
		v.data = nil
		for j := 0; j < 20; j++ {
			v.hash[j] = info.Info.Pieces[i+j]
		}
		if i+20 == len(info.Info.Pieces) {
			v.length = lastLen
		}
		piece[i/20] = &v
	}

	torrent, err := torrentParser.ParseFromFile(arg[0])
	if err != nil {
		panic(err)
	}

	workQueue := make(chan *Piece, len(piece))
	for i := range piece {
		workQueue <- piece[i]
	}

	go startDownload(workQueue, torrent)
	for len(piecedone) != len(piece) {
		fmt.Println("download = ", float64(len(piecedone))/float64(len(piece))*100, "%")
		fmt.Println("active peers = ", len(listOfPeers))
		time.Sleep(10 * time.Second)
	}
}
