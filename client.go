package client

import (
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var server, port, token string

func main() {
	log.Printf("Starting Camdram Email Web Tools client...")

	// Read in settings from config file.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: %s", err.Error())
	}
	server = os.Getenv("HTTP_SERVER")
	if server == "" {
		log.Fatalf("Server name not set in .env file, exiting...")
	}
	port = os.Getenv("HTTP_PORT")
	if port == "" {
		log.Fatalf("Server port not set in .env file, exiting...")
	}
	token = os.Getenv("HTTP_AUTH_TOKEN")
	if token == "" {
		log.Fatalf("Server HTTP auth token not set in .env file, exiting...")
	}

	// Check the endpoints on the remote server once every minute.
	ticker := time.NewTicker(1 * time.Minute)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	for {
		select {
		case <-ticker.C:
			go checkQueueLength()
			go checkHeldMessageCount()
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func checkQueueLength() {
	responseBody := makeRequest("queue")
	log.Printf("Queue: %s", responseBody)
}

func checkHeldMessageCount() {
	responseBody := makeRequest("held")
	log.Printf("Held: %s", responseBody)
}

func remoteUrl(endpoint string) string {
	return "http://" + server + ":" + port + "/" + endpoint
}

func makeRequest(endpoint string) string {
	url := remoteUrl(endpoint)
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
