package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Controller struct {
	mux    *http.ServeMux
	sqld   *SQLDriver
	bearer string
}

func newController(driver *SQLDriver, token string) *Controller {
	c := &Controller{
		mux:    http.NewServeMux(),
		sqld:   driver,
		bearer: token,
	}
	c.mux.HandleFunc("/queue", c.queueLengthResponder)
	c.mux.HandleFunc("/held", c.heldMessageCountResponder)
	c.mux.HandleFunc("/ping", c.smokeTestResponder)
	c.mux.HandleFunc("/json", c.jsonResponder)
	return c
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
	auth := r.Header.Get("Authorization")
	expected := "Bearer " + c.bearer
	switch auth {
	case expected:
		c.mux.ServeHTTP(w, r)
	case "":
		w.Header().Set("WWW-Authenticate", "Bearer realm=\"Camdram email-webtools\"")
		http.Error(w, "401 unauthorized", http.StatusUnauthorized)
	default:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	}
}

func (c *Controller) queueLengthResponder(w http.ResponseWriter, r *http.Request) {
	queueLength, err := c.sqld.GetQueueLength()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "Postal queue length: %d\n", queueLength)
}

func (c *Controller) heldMessageCountResponder(w http.ResponseWriter, r *http.Request) {
	heldCount, err := c.sqld.GetHeldMessageCount()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "Held message count: %d\n", heldCount)
}

func (c *Controller) smokeTestResponder(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong\n")
}

func (c *Controller) jsonResponder(w http.ResponseWriter, r *http.Request) {
	queueLength, err := c.sqld.GetQueueLength()
	if err != nil {
		log.Println(err)
		return
	}
	heldCount, err := c.sqld.GetHeldMessageCount()
	if err != nil {
		log.Println(err)
		return
	}
	data := map[string]int{
		"PostalQueue":  queueLength,
		"HeldMessages": heldCount,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(&data); err != nil {
		log.Println(err)
	}
}
