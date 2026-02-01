FROM golang:alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rhttpbin ./cmd/rhttpbin

# runtime
FROM scratch
COPY --from=builder /app/rhttpbin /rhttpbin
COPY --from=builder /app/web /web

ENV GO_PORT=80
EXPOSE 80

CMD ["/rhttpbin"]