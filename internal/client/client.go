package client

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/camdram/email-webtools/internal/assets"
	"github.com/cbroglie/mustache"
)

var errorAlertLastSent, queueAlertLastSent, heldAlertLastSent time.Time
var errorAlertExponent, queueAlertExponent, heldAlertExponent int

func StartListner(port string, token string, serverName string, to string) {
	if serverName == "" {
		log.Fatalln("Server name not set")
	}
	if to == "" {
		log.Fatalln("Mail recipient address(es) not set")
	}

	log.Println("Starting Email Web Tools in client mode")
	ticker := time.NewTicker(1 * time.Minute)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	for {
		select {
		case <-ticker.C:
			go checkJSON(port, token, serverName, to)
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func checkJSON(port string, token string, serverName string, to string) {
	data, err := fetchFromServer(port, token, serverName)
	if err != nil {
		if time.Now().UTC().Sub(errorAlertLastSent).Minutes() > calcTimeDiff(errorAlertExponent) {
			go sendErrorAlert(to, err)
		}
		log.Println("Failed to make request to remote server:", err)
		return
	}
	errorAlertExponent = 0
	if data["PostalQueue"] > 10 {
		if time.Now().UTC().Sub(queueAlertLastSent).Minutes() > calcTimeDiff(queueAlertExponent) {
			go sendQueueAlert(to, data)
		}
	} else {
		queueAlertExponent = 0
	}
	if data["HeldMessages"] > 0 {
		if time.Now().UTC().Sub(heldAlertLastSent).Minutes() > calcTimeDiff(heldAlertExponent) {
			go sendHeldAlert(to, data)
		}
	} else {
		heldAlertExponent = 0
	}
}

func calcTimeDiff(exponent int) float64 {
	num := math.Pow(1.6, float64(exponent)) + float64(exponent)*3.2
	return math.Min(num, 45)
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

func sendErrorAlert(to string, err error) {
	mailer, err := NewMailer()
	if err != nil {
		log.Println("Failed to initialise alert system:", err)
		return
	}
	defer mailer.Teardown()
	messageBody := fmt.Sprintln("Failed to make request to remote email-webtools server:", err)
	err = mailer.Send("camdram-admins@srcf.net", to, "Postal Queue Alert", messageBody)
	if err != nil {
		log.Fatalln("Failed to send alert:", err)
	} else {
		errorAlertLastSent = time.Now().UTC()
		errorAlertExponent++
	}
}

func sendQueueAlert(to string, data map[string]int) {
	log.Println("Sending Postal queue alert")
	mailer, err := NewMailer()
	if err != nil {
		log.Println("Failed to initialise alert system:", err)
		return
	}
	defer mailer.Teardown()
	buf, err := assets.Asset("assets/postal-queue.txt")
	if err != nil {
		log.Fatalln("Failed to load alert message:", err)
		return
	}
	messageBody, err := mustache.Render(string(buf), data)
	if err != nil {
		log.Fatalln("Failed to render alert message:", err)
		return
	}
	err = mailer.Send("camdram-admins@srcf.net", to, "Postal Queue Alert", messageBody)
	if err != nil {
		log.Fatalln("Failed to send alert:", err)
	} else {
		queueAlertLastSent = time.Now().UTC()
		queueAlertExponent++
	}
}

func sendHeldAlert(to string, data map[string]int) {
	log.Println("Sending held message queue alert")
	mailer, err := NewMailer()
	if err != nil {
		log.Println("Failed to initialise alert system:", err)
		return
	}
	defer mailer.Teardown()
	buf, err := assets.Asset("assets/held-messages.txt")
	if err != nil {
		log.Fatalln("Failed to load alert message:", err)
		return
	}
	messageBody, err := mustache.Render(string(buf), data)
	if err != nil {
		log.Fatalln("Failed to render alert message:", err)
		return
	}
	err = mailer.Send("camdram-admins@srcf.net", to, "Held Message Queue Alert", messageBody)
	if err != nil {
		log.Fatalln("Failed to send alert:", err)
	} else {
		heldAlertLastSent = time.Now().UTC()
		heldAlertExponent++
	}
}
