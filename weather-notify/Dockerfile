
# STEP 1 build executable binary
FROM golang:1.20.1-alpine3.17 as builder
RUN apk update && apk add git
RUN apk --update add ca-certificates
RUN apk add --no-cache tzdata
WORKDIR /weather-notify
#get dependancies

# We want to populate the module cache based on the go.{mod,sum} files.
COPY . .
RUN go mod download
RUN ls -l
#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/weather-notify
# STEP 2 build a small image
# start from scratch
FROM ubuntu:latest
ENV TZ America/Denver
# Copy our static executable
COPY --from=builder /weather-notify/fonts/ /go/bin/fonts/
COPY --from=builder /weather-notify/icons/ /go/bin/icons/
COPY --from=builder /go/bin/weather-notify /go/bin/weather-notify
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV GO_ENV production
ENTRYPOINT ["/go/bin/weather-notify"]
