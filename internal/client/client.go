package client

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/camdram/email-webtools/internal/assets"
)

func StartListner(port string, token string, serverName string) {
	log.Println("Starting Email Web Tools in client mode")
	if serverName == "" {
		log.Fatalln("Server name not set in .env file")
	}
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
	data, err := fetchFromServer(port, token, serverName)
	if err != nil {
		log.Println("Failed to make request to remote server:", err)
		return
	}
	if data["PostalQueue"] > 10 {
		go func() {
			log.Println("Sending Postal queue alert")
			mailer, err := NewMailer()
			if err != nil {
				log.Println("Failed to initialise alert system:", err)
				return
			}
			defer mailer.Teardown()
			data, err := assets.Asset("assets/postal-queue.txt")
			if err != nil {
				log.Fatalln("Failed to load alert message:", err)
				return
			}
			messageBody := string(data)
			err = mailer.Send("camdram-admins@srcf.net", "charlie@charliejonas.co.uk", "Postal Queue Alert", messageBody)
			if err != nil {
				log.Fatalln("Failed to send alert:", err)
			}
		}()
	}
	if data["HeldMessages"] > 0 {
		go func() {
			log.Println("Sending held message queue alert")
			mailer, err := NewMailer()
			if err != nil {
				log.Println("Failed to initialise alert system:", err)
				return
			}
			defer mailer.Teardown()
			data, err := assets.Asset("assets/held-messages.txt")
			if err != nil {
				log.Fatalln("Failed to load alert message:", err)
				return
			}
			messageBody := string(data)
			err = mailer.Send("camdram-admins@srcf.net", "charlie@charliejonas.co.uk", "Held Message Queue Alert", messageBody)
			if err != nil {
				log.Fatalln("Failed to send alert:", err)
			}
		}()
	}
}

func fetchFromServer(port string, token string, serverName string) (map[string]int, error) {
	url := remoteURL("json", port, serverName)
	client := &http.Client{}
	var data map[string]int
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func remoteURL(endpoint string, port string, serverName string) string {
	return "http://" + serverName + ":" + port + "/" + endpoint
}
