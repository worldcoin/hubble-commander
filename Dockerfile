FROM golang:1.16-alpine as build-env
WORKDIR /go/src/app

# Install tools
RUN apk update && apk add --no-cache git gcc libc-dev

# Fetch dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

# Build hubble executable
COPY . .
RUN go build -ldflags="-w -s" -o build/hubble ./main

# Fetch latest certificates
RUN update-ca-certificates --verbose

################################################################################
# Create minimal docker image for our app
FROM scratch

# Drop priviliges
USER 10001:10001

# Configure SSL CA certificates
COPY --from=build-env --chown=0:10001 --chmod=040 \
    /etc/ssl/certs/ca-certificates.crt /
ENV SSL_CERT_FILE="/ca-certificates.crt"

# Executable
COPY --from=build-env --chown=0:10001 --chmod=010 /go/src/app/build/hubble /bin
STOPSIGNAL SIGTERM
HEALTHCHECK NONE
ENTRYPOINT ["/bin"]
