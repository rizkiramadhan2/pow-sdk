FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rbac-server ./cmd/server

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/rbac-server /app/rbac-server

EXPOSE 10010

CMD ["/app/rbac-server"]