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

	"github.com/camdram/email-webtools/assets"
	"github.com/cbroglie/mustache"
)

var errorAlertLastSent, queueAlertLastSent, heldAlertLastSent time.Time
var errorAlertExponent, queueAlertExponent, heldAlertExponent int

func StartListner(port string, token string, serverName string, userAgent string, to string) {
	ensureConfig(serverName, to)
	log.Println("Starting Email Web Tools in client mode")
	ticker := time.NewTicker(1 * time.Minute)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	for {
		select {
		case <-ticker.C:
			go checkJSON(port, token, serverName, userAgent, to)
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func ensureConfig(serverName string, to string) {
	if serverName == "" {
		log.Fatalln("Server name not set")
	}
	if to == "" {
		log.Fatalln("Mail recipient address(es) not set")
	}
}

func checkJSON(port string, token string, serverName string, userAgent string, to string) {
	data, err := fetchFromServer(port, token, serverName, userAgent)
	if err != nil {
		msg := fmt.Sprintf("Failed to make request to remote email-webtools server: %v", err)
		log.Println(msg)
		if time.Now().UTC().Sub(errorAlertLastSent).Minutes() > calcTimeDiff(errorAlertExponent) {
			go sendErrorAlert(to, msg)
		}
		return
	}
	errorAlertExponent = 0
	if data["PostalQueue"] > 10 {
		log.Println("Alert firing: Postal queue length is greater than ten")
		if time.Now().UTC().Sub(queueAlertLastSent).Minutes() > calcTimeDiff(queueAlertExponent) {
			go sendQueueAlert(to, data)
		}
	} else {
		queueAlertExponent = 0
	}
	if data["HeldMessages"] > 0 {
		log.Println("Alert firing: Number of held messages is greater than zero")
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

func fetchFromServer(port string, token string, serverName string, userAgent string) (map[string]int, error) {
	url := remoteURL("json", port, serverName)
	client := &http.Client{}
	var data map[string]int
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cache-Control", "no-store, max-age=0")
	req.Header.Set("Authorization", "Bearer "+token)
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
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

func sendErrorAlert(to string, msg string) {
	mailer, err := NewMailer()
	if err != nil {
		log.Fatalln("Failed to initialise email alert system:", err)
		return
	}
	defer mailer.Teardown()
	err = mailer.Send("camdram-admins@srcf.net", to, "email-webtools Error Alert", msg)
	if err != nil {
		log.Fatalln("Failed to send email alert:", err)
	} else {
		errorAlertLastSent = time.Now().UTC()
		errorAlertExponent++
	}
}

func sendQueueAlert(to string, data map[string]int) {
	log.Println("Sending Postal queue alert")
	mailer, err := NewMailer()
	if err != nil {
		log.Fatalln("Failed to initialise email alert system:", err)
		return
	}
	defer mailer.Teardown()
	template, err := assets.ReadFile("postal-queue.txt.mustache")
	if err != nil {
		log.Fatalln("Failed to load email alert template:", err)
		return
	}
	messageBody, err := mustache.Render(string(template), data)
	if err != nil {
		log.Fatalln("Failed to render email alert message body:", err)
		return
	}
	err = mailer.Send("camdram-admins@srcf.net", to, "Postal Queue Alert", messageBody)
	if err != nil {
		log.Fatalln("Failed to send email alert:", err)
	} else {
		queueAlertLastSent = time.Now().UTC()
		queueAlertExponent++
	}
}

func sendHeldAlert(to string, data map[string]int) {
	log.Println("Sending held message queue alert")
	mailer, err := NewMailer()
	if err != nil {
		log.Fatalln("Failed to initialise email alert system:", err)
		return
	}
	defer mailer.Teardown()
	template, err := assets.ReadFile("held-messages.txt.mustache")
	if err != nil {
		log.Fatalln("Failed to load email alert template:", err)
		return
	}
	messageBody, err := mustache.Render(string(template), data)
	if err != nil {
		log.Fatalln("Failed to render email alert message body:", err)
		return
	}
	err = mailer.Send("camdram-admins@srcf.net", to, "Held Message Queue Alert", messageBody)
	if err != nil {
		log.Fatalln("Failed to send email alert:", err)
	} else {
		heldAlertLastSent = time.Now().UTC()
		heldAlertExponent++
	}
}
