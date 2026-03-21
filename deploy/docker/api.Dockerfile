FROM golang:1.25 AS builder
WORKDIR /workspace
COPY apps/backend/go.mod apps/backend/go.sum ./apps/backend/
WORKDIR /workspace/apps/backend
RUN go mod download
WORKDIR /workspace
COPY . .
WORKDIR /workspace/apps/backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /out/api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]
