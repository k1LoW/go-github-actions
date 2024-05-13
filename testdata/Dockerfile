FROM golang:latest

ENTRYPOINT ["env", "GOPATH=/go", "GOROOT=", "go", "test", "./...", "-coverprofile=coverage.out", "-covermode=count"]
