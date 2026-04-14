set shell := ["bash", "-uc"]


# Run database migrations up
migrate:
    go run cmd/migrate/main.go up

# Run database migrations down
migrate-down:
    go run cmd/migrate/main.go down

# Generate Protobuf files
generate-proto:
    protoc --proto_path=pkg/proto --go_out=internal/pb --go_opt=paths=source_relative --go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative pkg/proto/*.proto

# Generate SQLC query handlers
sqlc-generate:
    sqlc generate

# Init Swagger documentation
generate-swagger:
    swag init -g cmd/client/main.go

# --- Formatting & Quality ---

# Format Go code
fmt:
    go fmt ./...

# Run go vet
vet:
    go vet ./...

# Run linter
lint:
    golangci-lint run

# --- Build & Run ---

# Run the client
run-client:
    go run cmd/client/main.go

# Run the server
run-server:
    go run cmd/server/main.go

# Build the client binary
build-client:
    go build -ldflags="-s -w" -o client cmd/client/main.go

# Build the server binary
build-server:
    go build -ldflags="-s -w" -o server cmd/server/main.go

# Build the migration binary
build-migrate:
    go build -ldflags="-s -w" -o migrate cmd/migrate/main.go

# --- Testing ---
test-go:
    go test -race -covermode=atomic -coverprofile=coverage.txt -v ./...

# Run Hurl integration tests
test-hurl:
    hurl --test --variable baseUrl=http://localhost:5000 tests/hurl/*.hurl

# Run all tests (Go and Hurl)
test: test-all

# Run all tests (Go and Hurl)
test-all: test-go test-hurl

# Show test coverage in browser
coverage: test-go
    go tool cover -html=coverage.txt

# --- Documentation ---

# Build database documentation
db-docs:
    dbdocs build doc/db.dbml

# Export DBML to SQL schema
db-schema:
    dbml2sql --postgres -o doc/schema.sql doc/db.dbml

# Import SQL schema to DBML
db-sql-to-dbml:
    sql2dbml --postgres doc/schema.sql -o doc/db.dbml

# --- Docker ---

# Spin up containers
docker-up:
    docker compose up -d --build

# Tear down containers
docker-down:
    docker compose down
