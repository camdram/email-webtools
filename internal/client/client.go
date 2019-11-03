package client

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var mailer *Mailer

func StartListner(port string, token string, serverName string) {
	log.Printf("Starting Camdram Email Web Tools client...")
	if serverName == "" {
		log.Fatalf("Server name not set in .env file, exiting...")
	}
	mailer = NewMailer()
	defer mailer.Teardown()
	ticker := time.NewTicker(1 * time.Minute)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	for {
		select {
		case <-ticker.C:
			go checkJSON(port, token, serverName)
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func checkJSON(port string, token string, serverName string) {
	data := fetchFromServer(port, token, serverName)
	if data["PostalQueue"] > 0 {
		go mailer.SendMail("camdram-admins@srcf.net", "charlie@charliejonas.co.uk", "Camdram Postal Queue", "Check Postal queue length!")
	}
	if data["HeldMessages"] > 0 {
		go mailer.SendMail("camdram-admins@srcf.net", "charlie@charliejonas.co.uk", "Camdram Held Message Queue", "Check Held message queue length!")
	}
}

func fetchFromServer(port string, token string, serverName string) map[string]int {
	url := remoteURL("json", port, serverName)
	client := &http.Client{}
	var data map[string]int
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error constructing new request: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %s", err.Error())
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		log.Fatalf("Error decoding JSON: %s", err.Error())
	}
	return data
}

func remoteURL(endpoint string, port string, serverName string) string {
	return "http://" + serverName + ":" + port + "/" + endpoint
}