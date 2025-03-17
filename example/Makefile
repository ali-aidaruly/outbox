# Set environment variables
POSTGRES_CONTAINER=outbox_postgres
NATS_CONTAINER=outbox_nats
POSTGRES_PORT=5432
POSTGRES_USER=testuser
POSTGRES_PASSWORD=testpass
POSTGRES_DB=testdb
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5433/$(POSTGRES_DB)?sslmode=disable

# Start all services
up:
	docker-compose up -d

# Stop and remove all services
down:
	docker-compose down -v

# Adjusted migration path (one level up)
MIGRATIONS_PATH=../migrations/postgres

# Run database migrations inside PostgreSQL container
migrate:
	docker cp $(MIGRATIONS_PATH) $(POSTGRES_CONTAINER):/migrations
	docker exec -i $(POSTGRES_CONTAINER) sh -c \
		"PGPASSWORD=$(POSTGRES_PASSWORD) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -h outbox_postgres -p $(POSTGRES_PORT) -f /migrations/0001_create_uuid_extension.sql && \
		 PGPASSWORD=$(POSTGRES_PASSWORD) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -h outbox_postgres -p $(POSTGRES_PORT) -f /migrations/0002_create_outbox_table.sql"

# Build the example service
build:
	go build -o outbox-service ./cmd/main.go

# Set default values (can be overridden)
HTTP_PORT ?= 9090
POLL_INTERVAL ?= 10s
POLL_BATCH_SIZE ?= 5
LEASE_DURATION ?= 5s

# make run HTTP_PORT=9090 POLL_INTERVAL=10s POLL_BATCH_SIZE=3 LEASE_DURATION=10s
# Run the service with configurable arguments
run: build
	./outbox-service \
		-http-port=$(HTTP_PORT) \
		-poll-interval=$(POLL_INTERVAL) \
		-poll-batch-size=$(POLL_BATCH_SIZE) \
		-with-lease-duration=$(LEASE_DURATION)

# Run NATS consumer
run-consumer:
	go run ./consumer/main.go

# Clean build artifacts
clean:
	rm -f outbox-example

# Full setup (start services, migrate, build, and run)
setup: up migrate
