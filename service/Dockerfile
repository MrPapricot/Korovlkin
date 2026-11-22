FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /build/main \
    ./src/main.go

FROM alpine:3.21

RUN adduser -D -H appuser

WORKDIR /app

COPY --from=builder /build/main ./main

USER appuser

EXPOSE 8080

ENTRYPOINT ["./main"]
