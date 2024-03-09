
# STEP 1 build executable binary
FROM golang:1.22.1-alpine3.19 as builder
RUN apk update && apk add git
RUN apk --update add ca-certificates
RUN apk add --no-cache tzdata
WORKDIR /weather-server
#get dependancies

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

#build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/weather-server
# STEP 2 build a small image
# start from scratch
FROM scratch
ENV TZ America/Denver
# Copy our static executable
COPY --from=builder /go/bin/weather-server /go/bin/weather-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV GO_ENV production
ENTRYPOINT ["/go/bin/weather-server"]
