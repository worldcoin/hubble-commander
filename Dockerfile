FROM golang:1.16

LABEL org.opencontainers.image.source="https://github.com/Worldcoin/hubble-commander"

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make build

RUN mkdir -p /db/data/hubble

HEALTHCHECK --interval=3s --timeout=3s  CMD curl --fail -L -X POST 'localhost:8080' -H 'Content-Type: application/json' --data-raw '{"jsonrpc": "2.0","method": "hubble_getVersion","params": [],"id": 1}'

ENTRYPOINT ["build/hubble"]
CMD ["start"]
