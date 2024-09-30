FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY .. .
RUN go build -o app ./cmd/storage/main.go

FROM alpine

WORKDIR /app
COPY --from=builder /build/app /app/app
COPY . .

CMD ["./app"]