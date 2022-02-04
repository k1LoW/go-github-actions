FROM golang:1.17

ENTRYPOINT ["go", "test", "./...", "-coverprofile=coverage.out", "-covermode=count"]
