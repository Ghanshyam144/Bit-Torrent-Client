package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
)

// message format to send while connecting ...
func buildConnect() []byte {
	buff := make([]byte, 16)

	bb := make([]byte, 8)
	binary.BigEndian.PutUint64(bb, 0x41727101980)
	copy(buff[0:], bb) // connection id

	copy(buff[8:], retuByte(4, 0))
	copy(buff[12:], retuByte(4, rand.Uint32()))
	return buff
}

// decrypt message recieved while connecting ...
func decryptConnect(buffer []byte) connResp {
	var res connResp
	res.action = binary.BigEndian.Uint32(buffer[0:4])
	res.transaction_id = binary.BigEndian.Uint32(buffer[4:8])
	res.connection_id = binary.BigEndian.Uint64(buffer[8:])
	return res
}

// message format to send while announcing ...
func buildAnnounce(cnnctid uint64, torrent *gotorrentparser.Torrent, port uint16) []byte {
	buff := make([]byte, 98)
	bb := make([]byte, 8)
	binary.BigEndian.PutUint64(bb, cnnctid)
	copy(buff[0:], bb) // connection_id

	copy(buff[8:], retuByte(4, 1))              // action
	copy(buff[12:], retuByte(4, rand.Uint32())) // transaction_id

	infohash, _ := hex.DecodeString(torrent.InfoHash)
	copy(buff[16:], infohash)
	copy(buff[36:], PEER_ID)

	//  download
	// left
	// upload
	// event
	// ip

	copy(buff[88:], retuByte(4, rand.Uint32())) // key
	copy(buff[92:], retuByte(4, 4294967295))    // num_want

	// port
	por := make([]byte, 2)
	binary.BigEndian.PutUint16(por, port)
	copy(buff[96:], por) // port

	return buff
}

// message decrypting of announce response ...
func decryptAnnounce(buffer []byte, n int) annResp {
	var res annResp
	res.action = binary.BigEndian.Uint32(buffer[0:4])
	res.transactionId = binary.BigEndian.Uint32(buffer[4:8])
	res.leechers = binary.BigEndian.Uint32(buffer[12:16])
	res.seeders = binary.BigEndian.Uint32(buffer[16:20])

	temp := buffer[20:]
	for i := 0; i < (n - 20); i += 6 {
		var k Peer
		for j := i; j < i+4; j++ {
			k.ip += strconv.Itoa(int(temp[j]))
			if j < i+3 {
				k.ip += "."
			}
		}
		k.port = binary.BigEndian.Uint16(temp[i+4 : i+6])
		res.peers = append(res.peers, k)
	}
	return res
}

// socket messaging ....
func buildConnection(index int, buff []byte, torrent *gotorrentparser.Torrent, Peers *[]Peer) {
	trackerUrl, err := url.Parse(torrent.Announce[index])
	if err != nil {
		println("Error in url: ", err.Error())
		return
	}

	conn, err := net.Dial("udp", trackerUrl.Host)
	if err != nil {
		println("Error in establishing connection: ", err.Error())
		return
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		println("Connection SetReadDeadline, Error = ", err.Error())
		return
	}

	conn.Write(buff)

	buffer := make([]byte, 16)
	_, err = conn.Read(buffer)
	if err != nil {
		println("Error in reading msg: ", err.Error())
		return
	}
	println("Connection Successful")
	resp := decryptConnect(buffer)

	if resp.action == 0 {
		bufAnnc := buildAnnounce(resp.connection_id, torrent, 6881)
		conn.Write(bufAnnc)
		response := make([]byte, 1048576)
		n, err := conn.Read(response)
		if err != nil {
			println("Error in Announce: ", err.Error())
			return
		}
		println("Announce response size: ", n)

		npeer := decryptAnnounce(response, n)
		*Peers = append(*Peers, npeer.peers...)
		if len(*Peers) < len(npeer.peers) {
			*Peers = npeer.peers
		}
	}
}

// getting peers through various in between methods ...
func getPeers(torrent *gotorrentparser.Torrent) []Peer {
	buff := buildConnect()

	var Peers []Peer
	for i := range torrent.Announce {
		if torrent.Announce[i][0:3] == "udp" {
			buildConnection(i, buff, torrent, &Peers)
		}
	}

	newPeers := make([]Peer, 0)
	for _, i := range Peers {
		np := i.ip + fmt.Sprintf("%v", i.port)
		if !listOfPeers[np] {
			newPeers = append(newPeers, i)
			listOfPeers[np] = true
		}
	}
	return newPeers
}
