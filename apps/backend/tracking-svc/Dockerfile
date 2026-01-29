FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.work go.sum* ./
COPY auth-svc ./auth-svc
COPY vehicle-svc ./vehicle-svc
COPY tracking-svc ./tracking-svc

WORKDIR /build/tracking-svc

RUN go mod download
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /app/tracking-svc ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/tracking-svc /app/
EXPOSE 50002
CMD ["/app/tracking-svc"]