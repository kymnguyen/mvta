FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.work go.sum* ./
COPY auth-svc ./auth-svc
COPY vehicle-svc ./vehicle-svc
COPY tracking-svc ./tracking-svc

WORKDIR /build/auth-svc

RUN go mod download
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /app/auth-svc ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/auth-svc /app/
EXPOSE 50000
CMD ["/app/auth-svc"]