FROM golang:alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rhttpbin ./cmd/rhttpbin

# runtime
FROM scratch
COPY --from=builder /app/rhttpbin /app/rhttpbin
COPY --from=builder /app/web ./web

EXPOSE 8000
CMD ["/app/rhttpbin"]