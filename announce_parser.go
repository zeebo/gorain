package main

import (
	"fmt"
	"http"
	"net"
	"os"
	"strconv"
	"strings"
)

type parser struct {
	r   *http.Request
	err os.Error
}

func (p *parser) Int(v *int, s string, r bool) {
	if p.err != nil || (!r && p.r.FormValue("") == "") {
		return
	}
	*v, p.err = strconv.Atoi(p.r.FormValue(s))
	if p.err != nil {
		p.err = fmt.Errorf("Parsing %s [int]: %q", s, p.err)
	}
}

func (p *parser) Int64(v *int64, s string, r bool) {
	if p.err != nil || (!r && p.r.FormValue("") == "") {
		return
	}
	*v, p.err = strconv.Atoi64(p.r.FormValue(s))
	if p.err != nil {
		p.err = fmt.Errorf("Parsing %s [int64]: %q", s, p.err)
	}
}

func (p *parser) Bool(v *bool, s string, r bool) {
	if p.err != nil || (!r && p.r.FormValue("") == "") {
		return
	}
	*v, p.err = strconv.Atob(p.r.FormValue(s))
	if p.err != nil {
		p.err = fmt.Errorf("Parsing %s [bool]: %q", s, p.err)
	}
}

func (p *parser) IP(v *net.IP) {
	if p.err != nil {
		return
	}

	//grab the part before the port
	chunks := strings.Split(p.r.RemoteAddr, ":")
	noPort := strings.Join(chunks[:len(chunks)-1], ":")

	//check for ipv6 in []'s and remove them
	if noPort[0] == '[' {
		noPort = noPort[1 : len(noPort)-1]
	}

	*v = net.ParseIP(noPort)
	if *v == nil {
		p.err = fmt.Errorf("Parsing IP [ip]: %q", p.r.RemoteAddr)
	}
}
