# Camdram Email Web Tools

`email-webtools` is a small micro-service written in Go that we at Camdram use to monitor our Email systems and ensure that email receipt and delivery is functioning as expected.

## How does it work?

At Camdram we use [Postal](https://postal.atech.media/) for the sending and receiving of emails. This service connects to the Postal database in MySQL (technically MariaDB) and, when a correctly authenticated HTTP request is made, executes queries to determine the length of the mail queue and the volume of held mail.

Queued messages are messages that are actively awaiting delivery. Held messages are those that have been put to one side and will not be delivered, needing manual intervention.

## Compiling

We compile the project down to a single statically-linked executable which avoids having to install the entire Go toolchain on our server. Both of the methods detailed below produce a single `email-webtools` binary file in your working directory.

### Docker

First [install Docker](https://docs.docker.com/install/) and then run the following in a terminal window:

```bash
docker build -t camdram/email-webtools:latest .
docker run --rm -v ${PWD}:/app camdram/email-webtools:latest
```

### Old School

You will need to install version 1.21.0 of the Go programming language (see [here](https://golang.org/doc/install#install) for details). Then run the build using the included Makefile:

```bash
make clean
make all
```

## Deploying

You'll need to create a `.env` config file to house the authentication settings. This should contain something along the following lines:

```
HTTP_SERVER: hostname
HTTP_PORT: 8080
HTTP_AUTH_TOKEN: yourauthtoken
MYSQL_USER: username
MYSQL_PASSWORD: password
MAIN_DB: postal
SERVER_DB: postal-server-1
SMTP_TO: address1@example.com,address2@example.com
```

---

### Copyright

The code in this Git repository is released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License).

Copyright (c) 2019 various members of the Camdram Web Team and other contributors.
