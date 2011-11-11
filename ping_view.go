package main

import (
	"fmt"
	"http"
	"launchpad.net/gobson/bson"
	"net"
)

func pingRecover(w http.ResponseWriter) {
	err := recover()
	if err == nil {
		return
	}
	w.WriteHeader(http.StatusInternalServerError)

	fmt.Fprintln(w, "Error processing ping: %s", err)
}

func ping(c *Context) {
	defer pingRecover(c.w)

	p := parser{r: c.r}
	var ip net.IP
	p.IP(&ip)
	if p.err != nil {
		panic(p.err)
	}

	co := c.DB.C("users")

	var user *User
	co.Find(M{"ips": ip.String()}).One(&user)
	if user == nil {
		//insert a new user!
		co.Insert(M{
			"lastseen": bson.Now(),
			"lastip":   ip.String(),
			"ips":      []interface{}{ip.String()},
		})
		fmt.Fprintln(c.w, "You have been added to the user database")
		return
	}

	fmt.Fprintln(c.w, "You were already in here")
}
