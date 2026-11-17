# STEP 1 build executable binary
FROM golang:1.27.4-alpine3.24 AS builder
WORKDIR /sdr

# Populate the module cache from go.{mod,sum} first so it is reused.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/sdr .

# STEP 2 runtime image. Unlike the other services this one cannot start from
# scratch: it runs rtl_sdr (librtlsdr + libusb) to read the dongle.
FROM alpine:3.24
RUN apk add --no-cache rtl-sdr tzdata ca-certificates
ENV TZ=America/Denver
COPY --from=builder /go/bin/sdr /go/bin/sdr
ENTRYPOINT ["/go/bin/sdr"]
