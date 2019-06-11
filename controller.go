package main

import (
	"fmt"
	"log"
	"net/http"
)

type Controller struct {
	mux    *http.ServeMux
	sqld   *SqlDriver
	bearer string
}

func newController(driver *SqlDriver, token string) *Controller {
	c := &Controller{
		mux:    http.NewServeMux(),
		sqld:   driver,
		bearer: token,
	}
	c.mux.HandleFunc("/", c.queueLengthResponder)
	c.mux.HandleFunc("/ping", c.smokeTestResponder)
	return c
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "Bearer "+c.bearer {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		c.mux.ServeHTTP(w, r)
	}
}

func (c *Controller) queueLengthResponder(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Postal queue length: %d", c.sqld.GetQueueLength())
}

func (c *Controller) smokeTestResponder(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong!")
}
