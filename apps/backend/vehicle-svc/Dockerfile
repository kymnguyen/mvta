FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.work go.sum* ./
COPY auth-svc ./auth-svc
COPY vehicle-svc ./vehicle-svc
COPY tracking-svc ./tracking-svc

WORKDIR /build/vehicle-svc

RUN go mod download
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o /app/vehicle-svc ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/vehicle-svc /app/

EXPOSE 50001
CMD ["/app/vehicle-svc"]