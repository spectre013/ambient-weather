# STEP 1 build executable binary
FROM golang:1.21.6-alpine as builder
RUN apk update && apk add git
RUN apk --update add ca-certificates
RUN apk add --no-cache tzdata
WORKDIR /weather-transfer
#get dependancies

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/weather-transfer
# STEP 2 build a small image
# start from scratch
FROM ubuntu:24.04
ENV TZ America/Denver
# Copy our static executable
COPY --from=builder /go/bin/weather-transfer /go/bin/weather-transfer
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV GO_ENV production
ENTRYPOINT ["/go/bin/weather-transfer"]
