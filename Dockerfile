# syntax=docker/dockerfile:1

# Build
FROM --platform=linux/amd64 golang:1.22
WORKDIR /usr/src/app

COPY . .
RUN go build -o=./bin/main ./cmd

EXPOSE 8080

CMD ["./bin/main"]