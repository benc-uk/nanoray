# ========================================================#
# Stage 1: Build the Go binary
# ========================================================
FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./lib lib/
COPY ./frontend frontend/

RUN CGO_ENABLED=0 GOOS=linux go build -o frontend nanoray/frontend

# ========================================================
# Stage 2: Create a minimal container to run the Go binary
# ========================================================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/frontend .

ENV PORT=8000
# Set these at runtime
#ENV CONTROLLER_ADDR=controller:5000

EXPOSE 8000

CMD ["./frontend"]