.PHONY: proto proto-tools tidy build test

GOBIN := $(shell go env GOPATH)/bin

proto-tools:
	@command -v protoc >/dev/null || (echo "install protoc: brew install protobuf" && exit 1)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto: proto-tools
	PATH=$(GOBIN):$$PATH protoc \
	  -I proto \
	  --go_out=gen/proto --go_opt=module=qonaqzhai-backend/gen/proto \
	  --go-grpc_out=gen/proto --go-grpc_opt=module=qonaqzhai-backend/gen/proto \
	  proto/auth/v1/auth.proto \
	  proto/core/v1/core.proto \
	  proto/payment/v1/payment.proto \
	  proto/realtime/v1/realtime.proto

tidy:
	@for dir in . services/auth services/core services/payment services/realtime services/gateway pkg; do \
	  if [ -f $$dir/go.mod ]; then (cd $$dir && go mod tidy); fi; \
	done

build:
	@for svc in auth core payment realtime gateway; do \
	  echo "== building $$svc =="; \
	  (cd services/$$svc && go build ./...); \
	done

test:
	@for svc in auth core payment realtime gateway; do \
	  echo "== testing $$svc =="; \
	  (cd services/$$svc && go test -race ./...); \
	done
