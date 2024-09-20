FROM golang:1.23.1-alpine AS builder

RUN apk --no-cache add ca-certificates git

WORKDIR /app/
COPY . .

ENV CGO_ENABLED=0

RUN go mod download
RUN go build -o app cmd/server/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app /app

EXPOSE 9090/tcp

ENTRYPOINT ["/app"]
