# Build Stage

FROM golang:1.23.1 AS BuildStage

WORKDIR /

COPY . .

RUN go mod download

RUN go build -o guild-bot ./src/main.go

FROM alpine:latest

WORKDIR /

COPY --from=BuildStage /guild-bot /guild-bot

USER nonroo:nonroot

ENTRYPOINT [ "/guild-bot" ]
