# wishmill-server

swag init --generalInfo .\internal\api\api.go

go-bindata db\migrations
move to db/migrations
change package to migrations

go run ./cmd/wishmill -config S:\wishmill.yaml