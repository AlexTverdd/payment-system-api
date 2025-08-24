FROM golang:1.24-alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o payment_api main.go

FROM alpine:latest
WORKDIR /app

COPY --from=build /app/payment_api .
COPY wait-for-it.sh .

RUN apk add --no-cache bash
RUN chmod +x wait-for-it.sh

EXPOSE 8080

CMD ["./wait-for-it.sh", "db:5432", "--timeout=30", "--", "./payment_api"]