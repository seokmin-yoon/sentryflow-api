FROM golang:alpine AS builder

RUN apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o sentryflow-api .

EXPOSE 9090

CMD ["./sentryflow-api"]