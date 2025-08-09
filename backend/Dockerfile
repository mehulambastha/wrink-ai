# Stage 1
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Stage 2
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/main .

FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 5000

CMD ["./main"]




