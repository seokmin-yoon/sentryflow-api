FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o sentryflow-api .

EXPOSE 9090

CMD ["./sentryflow-api"]