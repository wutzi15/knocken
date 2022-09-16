FROM golang:latest as builder

WORKDIR /app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build knocken.go
RUN mkdir -p /app/html

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=builder /app/html /app/html
COPY --from=builder /app/knocken /app/knocken

CMD ["/app/knocken"]