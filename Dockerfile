FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev
COPY go.mod go.sum ./
RUN go mod download

# Copy and build the source code
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache sqlite
RUN mkdir -p /app/data
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"] 