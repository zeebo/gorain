package main

import (
	"fmt"
	"http"
	"net"
	"os"
)

type TEvent string

const (
	TEventNone      = ""
	TEventStarted   = "started"
	TEventStopped   = "stopped"
	TEventCompleted = "completed"
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
	Event      TEvent
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
	a.Event = TEvent(r.FormValue("event"))
	switch a.Event {
	case TEventStarted:
	case TEventStopped:
	case TEventCompleted:
	case TEventNone:
	default:
		return nil, fmt.Errorf("Unknown event type: %q", a.Event)
	}

	return a, nil
}
