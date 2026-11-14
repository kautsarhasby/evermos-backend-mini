
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main .

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]