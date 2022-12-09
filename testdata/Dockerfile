FROM golang:1.19

ENTRYPOINT ["env", "GOPATH=/go", "GOROOT=", "go", "test", "./...", "-coverprofile=coverage.out", "-covermode=count"]
