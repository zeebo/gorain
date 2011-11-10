package main

import (
	"fmt"
	"github.com/zeebo/bencode"
	"http"
	"log"
	"net"
	"os"
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

	fmt.Fprintln(c.w, user)
}

type TorrentState string

const (
	TorrentStarted   = TorrentState("started")
	TorrentStopped   = TorrentState("stopped")
	TorrentCompleted = TorrentState("completed")
)

type Announce struct {
	InfoHash   string
	PeerId     string
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Compact    bool
	NoPeerId   bool
	Event      TorrentState
	Numwant    int
	IP         net.IP
	Key        string
	TrackerId  string
}

func ParseAnnounce(r *http.Request) (*Announce, os.Error) {
	//Terms we dont have to sanity check/convert
	a := &Announce{
		InfoHash:  r.FormValue("info_hash"),
		PeerId:    r.FormValue("peer_id"),
		Key:       r.FormValue("key"),
		TrackerId: r.FormValue("trackerid"),
	}

	//create a parser that stops after the first error
	p := parser{r: r}

	//load the values in
	p.Int(&a.Port, "port", false)
	p.Int(&a.Numwant, "numwant", false)
	p.Int64(&a.Uploaded, "uploaded", false)
	p.Int64(&a.Downloaded, "downloaded", false)
	p.Int64(&a.Left, "left", false)
	p.Bool(&a.Compact, "compact", false)
	p.Bool(&a.NoPeerId, "no_peer_id", false)
	p.IP(&a.IP)

	//return any conversion errors
	if p.err != nil {
		return nil, p.err
	}

	//need to parse in event as enum type
	a.Event = TorrentState(r.FormValue("event"))
	switch a.Event {
	case TorrentStarted:
	case TorrentStopped:
	case TorrentCompleted:
	default:
		return nil, fmt.Errorf("Unknown event type: %q", a.Event)
	}

	return a, nil
}
