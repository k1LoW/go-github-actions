FROM golang:1.17

ENTRYPOINT ["env", "GOPATH=/go", "GOROOT=", "go", "test", "./...", "-coverprofile=coverage.out", "-covermode=count"]
