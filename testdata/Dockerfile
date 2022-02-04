FROM golang:1.17

RUN pwd && ls

COPY ../ /workspace/

WORKDIR /workspace/

ENTRYPOINT ["go", "test", "./...", "-coverprofile=coverage.out", "-covermode=count"]
