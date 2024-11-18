# Build Stage
FROM golang:1.23.1 AS buildstage

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build for Linux with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /build/guild-bot ./src/main.go

# Final Stage
FROM alpine:latest

WORKDIR /app

# Copy the statically built binary
COPY --from=buildstage /build/guild-bot /app/guild-bot

# Use the binary as the entrypoint
ENTRYPOINT ["/app/guild-bot"]
