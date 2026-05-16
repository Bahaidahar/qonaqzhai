# syntax=docker/dockerfile:1.7
FROM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache build-base
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/qonaqzhai ./cmd/qonaqzhai

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/qonaqzhai /app/qonaqzhai
USER nonroot:nonroot
ENV ADDR=:8080
EXPOSE 8080
ENTRYPOINT ["/app/qonaqzhai"]
