FROM golang:1.25-alpine AS builder

WORKDIR /server

COPY go.mod go.sum ./

RUN go mod download && go mod tidy

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN go build -o server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /server/server .

COPY --from=builder /go/bin/goose /bin/goose

COPY /server/db/migrations /migrations

EXPOSE 8080

CMD ["/bin/sh", "-c", "goose -dir /migrations postgres \"$DATABASE_URL\" up && ./server"]