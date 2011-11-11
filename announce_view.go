package main

import (
	"fmt"
	"github.com/zeebo/bencode"
	"http"
	"launchpad.net/gobson/bson"
	"log"
)

func announceRecover(w http.ResponseWriter) {
	err := recover()
	if err == nil {
		return
	}
	w.WriteHeader(http.StatusInternalServerError)

	enc := bencode.NewEncoder(w)
	enc.Encode(M{"failure reason": err})

	//DEBUG: newline for debugging
	fmt.Fprintln(w, "")
}

func announce(c *Context) {
	defer announceRecover(c.w)

	log.Print(c.r.RawURL)
	a, err := ParseAnnounce(c.r)
	if err != nil {
		log.Panic("ERROR: ", err)
	}
	//grab the user by ip
	co := c.DB.C("users")

	var user *User
	co.Find(M{"ips": a.IP.String()}).One(&user)
	if user == nil {
		log.Panic("Unauthorized IP: ", a.IP)
	}

	//update the user in the database
	user.LastIP = a.IP.String()
	user.LastSeen = bson.Now()
	user.Update(a.InfoHash, a.Event, a.Left)
	co.Update(M{"_id": user.ID}, user)

	var peers []StructPeer
	query := co.Find(M{"info.infohash": a.InfoHash, "_id": M{"$ne": user.ID}})
	query.Limit(min(50, max(0, a.Numwant))) //bound between 0 <= n <= 50
	selector := M{"lastip": 1, "port": 1}
	if !a.NoPeerId {
		selector["peerid"] = 1
	}
	query.Select(selector)
	query.All(&peers)

	if query.Iter().Err() != nil {
		panic(query.Iter().Err())
	}

	enc := bencode.NewEncoder(c.w)
	//build the response
	if a.Compact {
		response := CompactAnnounceResponse{
			Interval:   30,
			Complete:   0,
			Incomplete: 0,
			Peers:      make([]string, len(peers)),
		}
		for i := range peers {
			response.Peers[i] = peers[i].Compact()
		}
		enc.Encode(response)
	} else {
		response := AnnounceResponse{
			Interval:   30,
			Complete:   0,
			Incomplete: 0,
			Peers:      peers,
		}
		enc.Encode(response)
	}
	fmt.Fprintln(c.w, "")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
