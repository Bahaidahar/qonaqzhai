module qonaqzhai-backend/services/gateway

go 1.25.0

require (
	golang.org/x/time v0.11.0
	qonaqzhai-backend/pkg v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/grpc v1.68.0 // indirect
	google.golang.org/protobuf v1.36.0 // indirect
	qonaqzhai-backend/gen/proto v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	qonaqzhai-backend/gen/proto => ../../gen/proto
	qonaqzhai-backend/pkg => ../../pkg
)
