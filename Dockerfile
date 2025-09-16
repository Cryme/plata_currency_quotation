# ---- Build stage ----
FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN make initial-setup

RUN make build

# ---- Run stage ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/server /app/server

ENTRYPOINT ["/app/server"]