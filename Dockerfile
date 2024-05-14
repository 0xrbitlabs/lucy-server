FROM golang:1.22-alpine as builder
RUN update-ca-certificates
WORKDIR app/
COPY go.mod .
ENV GO111MODULE=on
RUN go mod download && go mod verify
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -extldflags "-static"' -o /app .
FROM alpine:latest
WORKDIR /usr/local/bin
COPY --from=builder /app /usr/local/bin/app
EXPOSE 8080
ENTRYPOINT ["./app"]
