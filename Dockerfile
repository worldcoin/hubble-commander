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

# Create (emtpy) config file
RUN touch empty-file

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

# Add empty config to avoid warning
COPY --from=build-env --chown=0:10001 --chmod=040 /src/empty-file /commander-config.yaml

# Create badger data dir
COPY --from=build-env --chown=10001:0 --chmod=700 /var/empty /badger
VOLUME ["/badger"]
ENV HUBBLE_BADGER_PATH=/badger

# Configure logging
ENV HUBBLE_LOG_FORMAT=json
ENV HUBBLE_LOG_LEVEL=info

# Executable
COPY --from=build-env --chown=0:10001 --chmod=010 /src/hubble .
STOPSIGNAL SIGTERM
HEALTHCHECK NONE
ENTRYPOINT ["/hubble"]
CMD ["start"]
