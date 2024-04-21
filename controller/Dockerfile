# ========================================================
# Stage 1: Build the Go binary
# ========================================================
FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./lib lib/
COPY ./controller controller/

RUN CGO_ENABLED=0 GOOS=linux go build -o controller nanoray/controller

# ========================================================
# Stage 2: Create a minimal container to run the Go binary
# ========================================================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/controller .

ENV PORT=5000

EXPOSE 5000

CMD ["./controller"]