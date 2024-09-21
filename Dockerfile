
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# RUN go test ./...

RUN go build -o scanner .

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/scanner .

COPY ./files .

RUN chmod +x run.sh

RUN apt-get update && \
    apt-get install -y rtl-sdr rtl-433 mosquitto-clients && \
    rm -rf /var/lib/apt/lists/*

#CMD ["tail", "-f", "/dev/null"]
CMD ["/bin/sh", "-c", "./run.sh && ./scanner"]