package main

import (
	"encoding/binary"
	"fmt"
	"sync"
)

var mutex = &sync.Mutex{}
var info bencodeTorrent
var piecedone = make(map[int]bool)
var pieceLength = int64(0)
var piece []*Piece
var PEER_ID = make([]byte, 20)
var listOfPeers = make(map[string]bool)
var path string

func removePeer(peer Peer) {
	mutex.Lock()
	delete(listOfPeers, peer.ip+fmt.Sprintf("%v", peer.port))
	mutex.Unlock()
}

func deleteData(index int) {
	mutex.Lock()
	piece[index].data = nil
	mutex.Unlock()
}

func markPiece(index int) {
	mutex.Lock()
	piecedone[index] = true
	mutex.Unlock()
}

func retuByte(sz int, num uint32) []byte {
	bb := make([]byte, sz)
	binary.BigEndian.PutUint32(bb, num)
	return bb
}
