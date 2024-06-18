FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY ./ ./
RUN cd cmd/server && CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:latest AS final
WORKDIR /app
COPY --from=builder /app/cmd/server/server .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/resolv.conf /etc/resolv.conf
CMD [ "./server" ]
