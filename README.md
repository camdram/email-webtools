# Camdram Email Web Tools

`email-webtools` is a small micro-service written in Go that we at Camdram use to keep track of the Email delivery queue length and more generally ensure that email receipt and delivery is functioning as expected.

## Development
You will need to install the Golang programming language (see [here](https://golang.org/doc/install#install) for details). Currently dependencies need to be install manually but this may change in future - run the following in a terminal:
```bash
go get github.com/joho/godotenv
go get -u github.com/go-sql-driver/mysql
```

Then get hacking! You can run your code with `go run server.go`.

## Compiling
We compile the project down to a single binary executable that gets uploaded to our server. This avoids having to install the entire Go toolchain which is simply unnecessary. From your development machine run the following:
```bash
GOARCH=amd64 GOOS=linux go tool dist install -v pkg/runtime
GOARCH=amd64 GOOS=linux go install -v -a std
GOARCH=amd64 GOOS=linux go build server.go
```

This will produce a single `server` binary file.

---

### Copyright

The code in this Git repository is released under the [MIT License](https://en.wikipedia.org/wiki/MIT_License).

Copyright (C) 2019 *The Association of Cambridge Theatre Societies* and contributing groups.
