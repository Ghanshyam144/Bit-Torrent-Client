package main

import (
	"crypto/sha1"
	"io"
	"net"
	"strconv"
	"time"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
)

func handleBufferMessage(peerConn *PeerConnection, t int) bool {
	for {
		msgLength, msgId, err := messageType(peerConn, t)
		if msgId == -2 {
			return true
		}
		if err != nil {
			return false
		}
		if handleMessage(peerConn, msgId, msgLength) != nil {
			return false
		}
	}
}

func rebuildHandShake(torrent *gotorrentparser.Torrent, peer Peer, workQueue chan *Piece, peerConn *PeerConnection) bool {
	conn, err := net.DialTimeout("tcp", peer.ip+":"+strconv.Itoa(int(peer.port)), 5*time.Second)
	if err != nil {
		return false
	}

	conn.Write(buildHandshake(torrent.InfoHash))
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	recieve := make([]byte, 68)
	_, err = io.ReadFull(conn, recieve)
	if err != nil {
		return false
	}
	peerConn.conn = conn
	return true
}

func handShake(torrent *gotorrentparser.Torrent, peer Peer, workQueue chan *Piece) {
	conn, err := net.DialTimeout("tcp", peer.ip+":"+strconv.Itoa(int(peer.port)), 5*time.Second)
	if err != nil {
		removePeer(peer)
		return
	}

	conn.Write(buildHandshake(torrent.InfoHash))
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	recieve := make([]byte, 68)
	_, err = io.ReadFull(conn, recieve)
	if err != nil {
		removePeer(peer)
		return
	}

	field := make([]bool, len(info.Info.Pieces)/20)
	go download(PeerConnection{conn, peer, PEER_ID, true, false, &field}, torrent, workQueue)
}

func validation(index int) bool {
	if piece[index].data == nil {
		return false
	}

	if !(sha1.Sum(*piece[index].data) == piece[index].hash) {
		println("Invalid piece with index :", index)
		return false
	}
	return true
}

func requestPiece(index int, peerConn *PeerConnection) bool {
	len := int64(16384)
	for i := int64(0); i < piece[index].length; i += 16384 {
		if len > int64(piece[index].length)-i {
			len = int64(piece[index].length) - i
		}
		sendRequest(peerConn, int(index), int(i), int(len))
	}
	return getPiece(index, peerConn)
}

func getPiece(index int, peerConn *PeerConnection) bool {
	if piece[index].data == nil {
		data := make([]byte, piece[index].length)
		piece[index].data = &data
	}

	numBlock := (piece[index].length + 16383) / 16384
	for numBlock > 0 {
		length, id, err := messageType(peerConn, 1200)
		if err != nil {
			return false
		}
		if id == 7 {
			if handleMessage(peerConn, id, length) != nil {
				return false
			}
			numBlock--
		} else {
			if handleMessage(peerConn, id, length) != nil {
				return false
			}
		}
	}
	return validation(index)
}

func download(peerConn PeerConnection, torrent *gotorrentparser.Torrent, workQueue chan *Piece) {
	defer removePeer(peerConn.peer)

	sendUnchoke(&peerConn)
	sendInterested(&peerConn)
	handleBufferMessage(&peerConn, 5)

	for reqPiece := range workQueue {
		if peerConn.choked || !(*peerConn.bitfield)[reqPiece.index] {
			workQueue <- reqPiece
			if peerConn.choked {
				active := handleBufferMessage(&peerConn, 2)
				if !active {
					peerConn.conn.Close()
					println("Connection Ended: ", peerConn.peer.ip)
					if !rebuildHandShake(torrent, peerConn.peer, workQueue, &peerConn) {
						return
					}
					println("Connection Rebuilt: ", peerConn.peer.ip)
				}
			}
			continue
		}
		mutex.Lock()

		//  end game ... when we require only one last piece : to ensure fast download!!
		if len(workQueue) == 0 {
			for ii := range piece {
				if !piecedone[ii] {
					workQueue <- piece[ii]
				}
			}
		}
		// if last piece done..
		if piecedone[reqPiece.index] {
			mutex.Unlock()
			continue
		}
		mutex.Unlock()

		println("Requesting piece: " + strconv.Itoa(reqPiece.index))
		valid := requestPiece(reqPiece.index, &peerConn)
		if valid {
			write(reqPiece.index)
			println("recieved piece: ", reqPiece.index, " ", len(piecedone))
			sendHave(&peerConn, int(reqPiece.index))
			reqPiece.data = nil
		} else {
			workQueue <- reqPiece
			peerConn.conn.Close()
			println("Connection Closed: ", peerConn.peer.ip)
			if !rebuildHandShake(torrent, peerConn.peer, workQueue, &peerConn) {
				return
			}
			println("Connection Rebuilt: ", peerConn.peer.ip)
		}

	}
}

func startDownload(workQueue chan *Piece, torrent *gotorrentparser.Torrent) {
	for {
		peerList := getPeers(torrent)
		for _, peer := range peerList {
			go handShake(torrent, peer, workQueue)
		}
		time.Sleep(60 * time.Second)
	}
}
