FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /agent-os ./cmd/server

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /agent-os .

EXPOSE 8080

ENV PORT=8080

CMD ["./agent-os"]
