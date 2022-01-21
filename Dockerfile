FROM golang:1.16-alpine as build-env
WORKDIR /src

# Install tools
RUN apk update && apk add --no-cache git gcc libc-dev

# Fetch dependencies
COPY go.mod go.sum .
RUN go mod download && go mod verify

# Build static hubble executable
COPY . .
RUN go build -ldflags="-w -s -linkmode external -extldflags -static" -o hubble ./main

# Fetch latest certificates
RUN update-ca-certificates --verbose

################################################################################
# Create minimal docker image for our app
FROM scratch
WORKDIR /

# Drop priviliges
USER 10001:10001

# Configure SSL CA certificates
COPY --from=build-env --chown=0:10001 --chmod=040 \
    /etc/ssl/certs/ca-certificates.crt /
ENV SSL_CERT_FILE="/ca-certificates.crt"

# Executable
COPY --from=build-env --chown=0:10001 --chmod=010 /src/hubble .
STOPSIGNAL SIGTERM
HEALTHCHECK NONE
ENTRYPOINT ["/hubble", "start"]
