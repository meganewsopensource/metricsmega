FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN apk update && apk add git

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
