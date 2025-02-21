
# STEP 1 build executable binary
FROM golang:1.24.0-alpine AS builder
RUN apk update && apk add git
RUN apk --update add ca-certificates
RUN apk add --no-cache tzdata
WORKDIR /weather-server
#get dependancies

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
COPY public .

RUN go mod download

COPY . .

#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/weather-server
# STEP 2 build a small image
# start from scratch
FROM ubuntu:jammy
ENV TZ=America/Denver
# Copy our static executable
COPY --from=builder /go/bin/weather-server /app/weather-server
COPY --from=builder /weather-server/public /app/public
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV GO_ENV=production
ENTRYPOINT ["/app/weather-server"]
