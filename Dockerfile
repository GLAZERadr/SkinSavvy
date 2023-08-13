# Stage 1: Build the application
FROM golang:1.20.5-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Create a minimal container to run the application
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=0 /app/app .

EXPOSE 8080

CMD ["./app"]
