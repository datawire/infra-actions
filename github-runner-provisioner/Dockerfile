FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./github-runner-provisioner

FROM ubuntu:22.04

WORKDIR /root/

RUN apt-get update && apt-get install -y ca-certificates openssl
ENV SSL_CERT_DIR=/etc/ssl/certs

COPY --from=builder /app/main .

CMD ["./main"]
