setup-env:
	echo "Setting up environment..."
	@go mod download

debug:
	echo "Debugging..."
	@go run main.go

build:
	echo "Building..."
	@go build -o bin/main main.go

docker-build:
	echo "Building docker image..."
	@docker build -t search-api:latest .

run-dependencies:
	echo "Running dependencies..."
	@docker-compose up -d

tests:
	echo "Running tests..."
	@go test -v ./...

