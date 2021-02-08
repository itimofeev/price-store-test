FROM golang:1.15 as builder

WORKDIR /build/price-store-test

COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on go mod download

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /bin/price-store github.com/itimofeev/price-store-test/cmd/price-store

FROM alpine

RUN apk add --update ca-certificates dmidecode # Certificates for SSL & dmidecode

WORKDIR /bin

COPY --from=builder /bin/price-store /bin/price-store

CMD ["price-store"]
