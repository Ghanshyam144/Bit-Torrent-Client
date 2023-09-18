package main

import (
	"encoding/binary"
	"encoding/hex"
)

func buildHandshake(infoHash string) []byte {
	buff := make([]byte, 68)
	copy(buff[0:], []byte{19})                      // pstrlen
	copy(buff[1:20], []byte("BitTorrent protocol")) // pstr
	copy(buff[20:28], []byte{0})                    // reserve

	infohash, _ := hex.DecodeString(infoHash) // info hash
	copy(buff[28:], infohash)

	copy(buff[48:], PEER_ID) //peer_id

	return buff
}

func sendKeepAlive(peerConn *PeerConnection) error {
	_, err := peerConn.conn.Write([]byte{0, 0, 0, 0})
	return err
}

func sendChoke(peerConn *PeerConnection) error {
	buff := make([]byte, 5)

	copy(buff[0:], retuByte(4, 1)) // length
	copy(buff[4:], []byte{0})      // id

	_, err := peerConn.conn.Write(buff)
	return err
}

func sendUnchoke(peerConn *PeerConnection) error {
	buff := make([]byte, 5)

	copy(buff[0:], retuByte(4, 1)) // length
	copy(buff[4:], []byte{1})      // id

	_, err := peerConn.conn.Write(buff)
	return err
}

func sendInterested(peerConn *PeerConnection) error {
	buff := make([]byte, 5)

	copy(buff[0:], retuByte(4, 1)) // length
	copy(buff[4:], []byte{2})      // id

	_, err := peerConn.conn.Write(buff)
	return err
}

func sendNotInterested(peerConn *PeerConnection) error {
	buff := make([]byte, 5)

	copy(buff[0:], retuByte(4, 1)) //length
	copy(buff[4:], []byte{3})      // id

	_, err := peerConn.conn.Write(buff)
	return err
}

func sendHave(peerConn *PeerConnection, index int) error {
	buff := make([]byte, 9)

	copy(buff[0:], retuByte(4, 5))             // length
	copy(buff[4:5], []byte{4})                 // id
	copy(buff[5:], retuByte(4, uint32(index))) // index

	_, err := peerConn.conn.Write(buff)
	return err
}

func sendRequest(peerConn *PeerConnection, index int, offset int, length int) error {
	buff := make([]byte, 17)

	copy(buff[0:], retuByte(4, 13))              // length
	copy(buff[4:], []byte{6})                    //id
	copy(buff[5:], retuByte(4, uint32(index)))   // index
	copy(buff[9:], retuByte(4, uint32(offset)))  // offset
	copy(buff[13:], retuByte(4, uint32(length))) // block length

	_, err := peerConn.conn.Write(buff)
	return err
}

// 									%TODO

func sendBitfield(peerConn *PeerConnection, field []byte) error {
	length := len(info.Info.Pieces) / 20
	buff := make([]byte, 5+length)

	copy(buff[0:], retuByte(4, uint32(1+length))) // length
	copy(buff[4:5], []byte{5})                    // id

	copy(buff[5:], field) // bitfield

	_, err := peerConn.conn.Write(buff)
	return err
}

func sendPiece(peerConn *PeerConnection, index int, offset int, block []byte) error {
	buff := make([]byte, 13+len(block))

	copy(buff[0:], retuByte(4, uint32(9+len(block)))) // length
	copy(buff[4:], []byte{7})                         // id
	copy(buff[5:], retuByte(4, uint32(index)))        // index
	copy(buff[9:], retuByte(4, uint32(offset)))       // offset
	copy(buff[13:], block)                            // block

	_, err := peerConn.conn.Write(buff)
	return err
}

func sendCancel(peerConn *PeerConnection, index int, offset int, length int) error {
	buff := make([]byte, 17)

	copy(buff[0:], retuByte(4, 13))              // length
	copy(buff[4:], []byte{8})                    //id
	copy(buff[5:], retuByte(4, uint32(index)))   // index
	copy(buff[9:], retuByte(4, uint32(offset)))  // offset
	copy(buff[13:], retuByte(4, uint32(length))) // block length

	_, err := peerConn.conn.Write(buff)
	return err
}

func sendPort(peerConn *PeerConnection, p uint16) error {
	buff := make([]byte, 7)

	copy(buff[0:], retuByte(4, 3)) // length
	copy(buff[4:], []byte{9})      // id

	// port
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, p)
	copy(buff[5:], port)

	_, err := peerConn.conn.Write(buff)
	return err
}
