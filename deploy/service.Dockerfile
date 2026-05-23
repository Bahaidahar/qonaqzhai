# Shared multi-stage Dockerfile for every Go microservice in backend/services.
# Pick the target with --build-arg SERVICE=<name>, e.g.
#
#   docker build -f deploy/service.Dockerfile --build-arg SERVICE=auth -t qonaqzhai-auth .
#
# The build context must be the repo root so we can copy backend/ in full.

ARG SERVICE
FROM golang:1.25-alpine AS builder
ARG SERVICE
WORKDIR /src
COPY backend ./backend
WORKDIR /src/backend
RUN go work sync
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" \
      -o /out/svc ./services/${SERVICE}/cmd/${SERVICE}

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/svc /svc
USER nonroot
ENTRYPOINT ["/svc"]
