FROM golang:1.25-alpine AS builder

WORKDIR /server

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /server/server .

EXPOSE 8080

CMD ["./server"]