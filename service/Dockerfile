FROM golang:1.27-alpine AS docs-generator

WORKDIR /src

COPY ./src ./src

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g ./src/main.go -o ./src/docs --parseInternal


FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod ./

RUN go mod download

COPY --from=docs-generator /src .

RUN go mod tidy

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
