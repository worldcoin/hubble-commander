FROM golang:1.15

LABEL org.opencontainers.image.source="https://github.com/Worldcoin/hubble-commander"

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN make build

CMD ["build/hubble"]
