package main

type TorrentInfo struct {
	InfoHash string
	State    TorrentState
}

type User struct {
	Password   string
	Port       int
	IPs        []string
	Info       []TorrentInfo
	Downloaded int64
	Uploaded   int64
}
