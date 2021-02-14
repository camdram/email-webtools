package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Controller struct {
	router *mux.Router
	server *http.Server
	sqld   *SQLDriver
	bearer string
}

func newController(driver *SQLDriver, token string, addr string) *Controller {
	c := &Controller{
		router: mux.NewRouter(),
		sqld:   driver,
		bearer: token,
	}
	c.router.HandleFunc("/ping", c.smokeTestResponder)
	c.router.HandleFunc("/json", c.jsonResponder)
	c.router.HandleFunc("/queue", c.queueLengthResponder)
	c.router.HandleFunc("/held", c.heldMessageCountResponder)
	c.server = &http.Server{
		Handler: c.router,
		Addr:    addr,
	}
	return c
}

func (c *Controller) ListenAndServe() error {
	c.router.Use(logMiddleware)
	c.router.Use(serverHeaderMiddleware)
	c.router.Use(c.authMiddleware)
	return c.server.ListenAndServe()
}

func (c *Controller) Shutdown(ctx context.Context) error {
	return c.server.Shutdown(ctx)
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
