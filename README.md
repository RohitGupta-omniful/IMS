
Inventory Management System (IMS) Microservice
=============================================

This microservice manages hubs, SKUs, and inventory operations. It provides robust CRUD APIs, validation endpoints, caching, database migrations, and internationalized responses.

Project Structure
-----------------
cache/         # Redis cache setup and usage
configs/       # Application configuration and environment management
db/            # Database connection and migration logic using GORM
handler/       # HTTP handlers for hubs, SKUs, and inventory operations
localstack/    # Configuration for localstack-based testing
middleware/    # Authentication and logging middleware
models/        # ORM models for Hub, SKU, Inventory, etc.
router/        # API route definitions
server/        # Server and Gin router initialization
main.go        # Application entry point
docker-compose.yml  # Setup for Redis, PostgreSQL, and other services
go.mod / go.sum     # Go module dependencies

Features
--------
- Full CRUD operations for Hubs and SKUs
- Inventory upsert with intelligent conflict resolution
- Validation endpoints for checking the existence of hubs and SKUs
- Redis caching for optimized validation checks
- Middleware for token-based authentication and structured logging
- Internationalization (i18n) support for response messages and logs
- Database migrations using Go-native tooling
- Local development ready with Docker Compose support

Application Workflow
--------------------
1. On startup, the service initializes configurations, database, Redis, i18n settings, and the Gin HTTP server.
2. Authentication middleware checks for a valid bearer token in the Authorization header.
3. Clients can verify if a hub or SKU exists using /validate/hub/:id or /validate/sku/:id.
4. The service checks Redis first for validation data, falling back to the database if not cached.
5. Inventory records can be inserted or updated via the /inventory/upsert endpoint.
6. CRUD APIs allow management of hubs and SKUs with support for POST, GET, PUT, and DELETE.
7. Docker Compose facilitates seamless local development with Redis and PostgreSQL containers.

API Endpoints Example
-------------
Validate Hub
GET /validate/hub/{hub_id}
#### Required Headers

| Key            | Value               |
|----------------|---------------------|
| Authorization  | Bearer secret-token |

Response:
{
  "data": {
    "is_valid": true
  }
}

