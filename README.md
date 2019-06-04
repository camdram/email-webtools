# Camdram Email Web Tools

`email-webtools` is a small micro-service written in Go that we at Camdram use to keep track of the Email delivery queue length and more generally ensure that email receipt and delivery is functioning as expected.

## Compiling
We compile the project down to a single executable that gets uploaded to our server via SFTP or similar. This avoids having to install the entire Go toolchain which is simply unnecessary. Both of these methods produce a single `email-webtools` binary file in your working directory.

### Docker
First [install Docker](https://docs.docker.com/install/) and then run the following in a terminal window:
```bash
docker build -t camdram/email-webtools .
docker run -v `pwd`:/go/src/github.com/camdram/email-webtools camdram/email-webtools
```

### Old School
You will need to install the Golang programming language (see [here](https://golang.org/doc/install#install) for details). Currently the project dependencies need to be install manually but this may change in future.
```bash
go get github.com/joho/godotenv
go get -u github.com/go-sql-driver/mysql
GOARCH=amd64 GOOS=linux go tool dist install -v pkg/runtime
GOARCH=amd64 GOOS=linux go install -v -a std
GOARCH=amd64 GOOS=linux go build build -o email-webtools main.go controller.go driver.go
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
