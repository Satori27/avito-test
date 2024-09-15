FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go build -o avito /app/cmd/tender

FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /app/avito .

EXPOSE 8080

CMD ["./avito"]