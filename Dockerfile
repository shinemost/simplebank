FROM golang:1.16-alpine3.13 AS builder
WORKDIR /app
COPY . .
RUN go build -mod=mod -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-arm64.tar.gz | tar xvz
        

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8081
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]