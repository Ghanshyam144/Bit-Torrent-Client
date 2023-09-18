package main

import (
	"net"
)

type PeerConnection struct {
	conn       net.Conn
	peer       Peer
	peerId     []byte
	choked     bool
	interested bool
	bitfield   *[]bool
}

type connResp struct {
	action         uint32
	transaction_id uint32
	connection_id  uint64
}

type Peer struct {
	ip   string
	port uint16
}

type annResp struct {
	action        uint32
	transactionId uint32
	seeders       uint32
	leechers      uint32
	peers         []Peer
}

type file struct {
	Length int64    `bencode:"length"`
	Path   []string `bencode:"path"`
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int64  `bencode:"piece length"`
	Length      int64  `bencode:"length"`
	Files       []file `bencode:"files"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type Piece struct {
	index  int
	length int64
	hash   [20]byte
	data   *[]byte
}
