module qonaqzhai-backend/services/realtime

go 1.25.0

require (
	github.com/gorilla/websocket v1.5.3
	github.com/jackc/pgx/v5 v5.9.2
	google.golang.org/grpc v1.68.0
	qonaqzhai-backend/gen/proto v0.0.0-00010101000000-000000000000
	qonaqzhai-backend/pkg v0.0.0-00010101000000-000000000000
)

replace (
	qonaqzhai-backend/gen/proto => ../../gen/proto
	qonaqzhai-backend/pkg => ../../pkg
)
