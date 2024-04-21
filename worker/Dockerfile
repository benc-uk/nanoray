# ========================================================
# Stage 1: Build the Go binary
# ========================================================
FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./lib lib/
COPY ./worker worker/

RUN CGO_ENABLED=0 GOOS=linux go build -o worker nanoray/worker

# ========================================================
# Stage 2: Create a minimal container to run the Go binary
# ========================================================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/worker .

ENV PORT=4000

# Set these at runtime
#ENV CONTROLLER_ADDR=controller:5000
#MAX_JOBS=10
#HOSTNAME=worker1

EXPOSE 4000

CMD ["./worker"]