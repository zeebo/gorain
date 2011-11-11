package main

import (
	"fmt"
	"launchpad.net/gobson/bson"
	"os"
)

type TStatus string

const (
	TStatusSeed  = "seed"
	TStatusLeech = "leech"
)

func getStatus(left int64) TStatus {
	if left == 0 {
		return TStatusSeed
	}
	return TStatusLeech
}

type TInfo struct {
	InfoHash string
	Status   TStatus
}

type User struct {
	ID         bson.ObjectId `bson:"_id"`
	PeerId     string
	LastSeen   bson.Timestamp
	Password   string
	Port       int
	LastIP     string
	IPs        []string
	Info       []TInfo
	Downloaded int64
	Uploaded   int64
}

func (u *User) Update(hash string, ev TEvent, left int64) os.Error {
	switch ev {
	//if the event is stopped, we need to remove the infohash from the list
	case TEventStopped:
		return u.remove(hash)

	//completed and started add the info hash if it isn't there, with status
	//determined by the number of bytes left
	case TEventCompleted:
		fallthrough
	case TEventStarted:
		status := getStatus(left)
		u.update(hash, status)
	}

	return nil
}

func (u *User) remove(hash string) os.Error {
	for idx := range u.Info {
		if u.Info[idx].InfoHash == hash {
			u.Info = append(u.Info[:idx], u.Info[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Unable to find info_hash: %x", hash)
}

func (u *User) update(hash string, status TStatus) {
	//update
	for idx := range u.Info {
		if u.Info[idx].InfoHash == hash {
			u.Info[idx].Status = status
			return
		}
	}
	//otherwise append
	u.Info = append(u.Info, TInfo{
		InfoHash: hash,
		Status:   status,
	})
}
