# Camdram Email Web Tools

`email-webtools` is a small micro-service written in Go that we at Camdram use to monitor our Email systems and ensure that email receipt and delivery is functioning as expected.

## How does it work?
At Camdram we use [Postal](https://postal.atech.media/) for the sending and receiving of emails. This service connects to the Postal database in MySQL (technically MariaDB) and, when a correctly authenticated HTTP request is made, executes queries to determine the length of the mail queue and the number of held messaged.

Queued messages are messages that are actively awaiting delivery. Held messages are those that have been put to one side and will not be delivered, needing manual intervention.

## Compiling
We compile the project down to a single executable that gets uploaded to our server via SFTP or similar. This avoids having to install the entire Go toolchain which is simply unnecessary. Both of these methods produce a single `email-webtools` binary file in your working directory.

### Docker
First [install Docker](https://docs.docker.com/install/) and then run the following in a terminal window:
```bash
docker build -t camdram/email-webtools .
docker run -v `pwd`:/go/src/github.com/camdram/email-webtools camdram/email-webtools
```

### Old School
You will need to install the Golang programming language (see [here](https://golang.org/doc/install#install) for details). Then run the build using the included Makefile:
```bash
GOARCH=amd64 GOOS=linux go tool dist install -v pkg/runtime
GOARCH=amd64 GOOS=linux go install -v -a std
make all
```

## Installing
You'll need to create a `.env` config file. This should contain something along the following lines:
```
HTTP_PORT: 8080
MYSQL_USER: username
MYSQL_PASSWORD: password
MYSQL_DBL postal
```

---

### Copyright

The code in this Git repository is released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License).

Copyright (C) 2019 *The Association of Cambridge Theatre Societies* and contributing groups.
