FROM golang:1.22-alpine AS builder
LABEL authors="Victor Assunção<hugotex00@gmail.com>"

ADD ../ /app
WORKDIR /app

RUN go mod download
RUN go build -o main ./cmd/main.go

FROM alpine:latest
LABEL authors="Victor Assunção<hugotex00@gmail.com>"

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/base.txt .

CMD ["/app/main"]