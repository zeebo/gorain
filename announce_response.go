package main

import (
	"bytes"
	"encoding/binary"
	"net"
)

type CompactAnnounceResponse struct {
	Interval   int      `bencode:"interval"`
	Complete   int      `bencode:"complete"`
	Incomplete int      `bencode:"incomplete"`
	Peers      []string `bencode:"peers"`
}

type AnnounceResponse struct {
	Interval   int          `bencode:"interval"`
	Complete   int          `bencode:"complete"`
	Incomplete int          `bencode:"incomplete"`
	Peers      []StructPeer `bencode:"peers"`
}

type StructPeer struct {
	PeerId string `bencode:"peer id"`
	LastIP string `bencode:"ip"`
	Port   int16  `bencode:"port"`
}

func (s *StructPeer) Compact() string {
	var peer [6]byte

	//first parse the ip
	ip := net.ParseIP(s.LastIP).To4()
	if ip == nil {
		return ""
	}

	var buf bytes.Buffer //make a buffer
	binary.Write(&buf, binary.BigEndian, s.Port)

	copy(peer[:], ip)
	copy(peer[4:], buf.Bytes())

	return string(peer[:])
}
