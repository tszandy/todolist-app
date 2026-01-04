FROM golang:1.25-alpine AS basic

FROM basic AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./main.go

# Dev stage for live reloading with Air
FROM basic AS dev
WORKDIR /src
RUN go install github.com/air-verse/air@latest
EXPOSE 8080
CMD ["air"]

# Deploy stage (production)
FROM alpine:3.18 AS deploy
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /src/backend /app/backend
EXPOSE 8080
CMD ["/app/backend"]
