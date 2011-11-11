package main

import (
	"gorilla.googlecode.com/hg/gorilla/mux"
	"http"
	"launchpad.net/mgo"
	"log"
	"os"
)

//helper type for encoding responses
type M map[string]interface{}

//Our mongo session
var session *mgo.Session

//Mongo info
const (
	mgoUrl      = "localhost"
	mgoDatabase = "rain"
)

//Helpers for setting up http
type Context struct {
	DB mgo.Database
	w  http.ResponseWriter
	r  *http.Request
}

//Our handler
type handler func(*Context)

//Wrap our handler with the context (copy of the mongo session)
func w(f handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(&Context{
			DB: session.Copy().DB(mgoDatabase),
			w:  w,
			r:  r,
		})
	}
}

//Setup dat http
func main() {
	var err os.Error
	session, err = mgo.Mongo(mgoUrl)
	if err != nil {
		panic(err)
	}

	router := new(mux.Router)
	router.HandleFunc("/a", w(announce)) //dont support scrapes
	router.HandleFunc("/ping", w(ping))
	if err := http.ListenAndServe(":9988", router); err != nil {
		log.Fatal(err)
	}
}
