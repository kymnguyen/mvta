FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum* ./
COPY tracking-svc/go.mod tracking-svc/go.sum ./tracking-svc/

WORKDIR /app/tracking-svc
RUN go mod download

COPY tracking-svc/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/tracking-svc/server .

EXPOSE 8002
CMD ["./server"]
