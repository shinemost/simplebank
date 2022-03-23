FROM golang:1.16-alpine3.13 AS builder
WORKDIR /app
COPY . .
RUN go build -mod=mod -o main main.go

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8081
CMD ["/app/main"]