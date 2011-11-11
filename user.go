package main

import "launchpad.net/gobson/bson"

type TorrentInfo struct {
	InfoHash string
	State    TorrentState
}

type User struct {
	ID         bson.ObjectId `bson:"_id"`
	LastSeen   bson.Timestamp
	Password   string
	Port       int
	IPs        []string
	Info       []TorrentInfo
	Downloaded int64
	Uploaded   int64
}
