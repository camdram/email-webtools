package client

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartListner(port string, token string, serverName string) {
	log.Printf("Starting Camdram Email Web Tools client...")
	if serverName == "" {
		log.Fatalf("Server name not set in .env file, exiting...")
	}
	ticker := time.NewTicker(1 * time.Minute)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	for {
		select {
		case <-ticker.C:
			go checkQueueLength(port, token, serverName)
			go checkHeldMessageCount(port, token, serverName)
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func checkQueueLength(port string, token string, serverName string) {
	responseBody := makeRequest("queue", port, token, serverName)
	log.Printf("Queue: %s", responseBody)
}

func checkHeldMessageCount(port string, token string, serverName string) {
	responseBody := makeRequest("held", port, token, serverName)
	log.Printf("Held: %s", responseBody)
}

func makeRequest(endpoint string, port string, token string, serverName string) string {
	url := remoteURL(endpoint, port, serverName)
	client := &http.Client{}
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err.Error())
	}
	return string(body)
}

func remoteURL(endpoint string, port string, serverName string) string {
	return "http://" + serverName + ":" + port + "/" + endpoint
}
