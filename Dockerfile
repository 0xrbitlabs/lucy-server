FROM golang:latest as builder
RUN update-ca-certificates
WORKDIR app/
COPY go.mod .
ENV GO111MODULE=on
RUN go mod download && go mod verify
COPY . .
RUN go build -o /app .
FROM debian:latest
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /usr/local/bin/app
WORKDIR /usr/local/bin
EXPOSE 1000
ENTRYPOINT ["app"]
