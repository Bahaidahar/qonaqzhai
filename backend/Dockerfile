# Multi-stage build: compile all five service binaries from the Go workspace
# into one small image. Each compose service runs a different binary via its
# command, so we build once and reuse the image five times.
FROM golang:1.25 AS build
WORKDIR /src

# Copy the whole workspace (go.work ties the modules together with replace
# directives, so a single module copy isn't enough).
COPY . .

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -trimpath -o /out/auth     ./services/auth/cmd/auth     \
 && go build -trimpath -o /out/core     ./services/core/cmd/core     \
 && go build -trimpath -o /out/payment  ./services/payment/cmd/payment \
 && go build -trimpath -o /out/realtime ./services/realtime/cmd/realtime \
 && go build -trimpath -o /out/gateway  ./services/gateway/cmd

# Minimal runtime — static binaries, no libc needed.
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/ /usr/local/bin/
# Default command overridden per service in docker-compose.
CMD ["gateway"]
