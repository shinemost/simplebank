FROM golang:1.16-alpine3.13 AS builder
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -mod=mod -o main main.go

        

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8081
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]