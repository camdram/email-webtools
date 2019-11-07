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
			go checkJSON(port, token, serverName)
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func checkJSON(port string, token string, serverName string) {
	data := fetchFromServer(port, token, serverName)
	if data["PostalQueue"] > 10 {
		go func() {
			mailer := NewMailer()
			defer mailer.Teardown()
			messageBody := `The length of the Camdram Postal mail queue is alerting!

Summary:

The queue of mail waiting to be delivered by Postal is exceeding normal
operational levels. This could be due to server issues (eg. the message
broker or SMTP daemon has crashed) or simply because Camdram is sending
a very high rate of email (eg. spam).

Remedial action to take:

1. Login to the admin interface at https://mail.camdram.net and monitor
   the situation.
2. Ensure both MariaDB and RabbitMQ are functioning using 'sudo
   systemctl status mariadb rabbitmq-server'.
3. Run the following to check the status of the email system 'cd
   /home/postal/app && sudo -Hu postal procodile status'.
4. If necessary, restart the email system by typing 'sudo systemctl
   restart postal'.
5. If the queue is not going down then open a Postal console using
   the following command 'cd /home/postal/app && sudo -Hu postal
   bin/postal console' and type the following Ruby code inside:

=======================================================================
org = Organization.first
server = org.servers.first
db = server.message_db
begin
  sleep 3600
  messages = db.messages(where: {held: 1}, limit: 140)
  messages.each do |msg|
    msg.add_to_message_queue(manual: true)
  end
end until messages.length == 0
=======================================================================

   This will sleep for one hour before attempting to resend the first
   140 messages in the queue.`
			mailer.Send("camdram-admins@srcf.net", "charlie@charliejonas.co.uk", "Postal Queue Alert", messageBody)
		}()
	}
	if data["HeldMessages"] > 0 {
		go func() {
			mailer := NewMailer()
			defer mailer.Teardown()
			messageBody := `

Summary:

The number of messages being held back by Postal is non-zero. This
could be due to server issues (eg. the message broker or SMTP daemon
has crashed) or simply because Camdram is sending a high rate of email
(eg. spam).

Remedial action to take:

1. Login to the admin interface at https://mail.camdram.net and monitor
   the situation.
2. Ensure both MariaDB and RabbitMQ are functioning using 'sudo
   systemctl status mariadb rabbitmq-server'.
3. Run the following to check the status of the email system 'cd
   /home/postal/app && sudo -Hu postal procodile status'.
4. If necessary, restart the email system by typing 'sudo systemctl
   restart postal'.
5. If the queue is not going down then open a Postal console using
   the following command 'cd /home/postal/app && sudo -Hu postal
   bin/postal console' and type the following Ruby code inside:

=======================================================================
org = Organization.first
server = org.servers.first
db = server.message_db
begin
	sleep 3600
	messages = db.messages(where: {held: 1}, limit: 140)
	messages.each do |msg|
	msg.add_to_message_queue(manual: true)
	end
end until messages.length == 0
=======================================================================

   This will sleep for one hour before attempting to resend the first
   140 messages in the queue.`
			mailer.Send("camdram-admins@srcf.net", "charlie@charliejonas.co.uk", "Held Message Queue Alert", messageBody)
		}()
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
