
Inventory Management System (IMS) Microservice
=============================================

This microservice manages hubs, SKUs, and inventory operations. It provides robust CRUD APIs, validation endpoints, caching, database migrations, and internationalized responses.

Project Structure
-----------------
cache/             # Redis cache setup and usage
configs/           # Configuration and environment management
db/                # Database connection and GORM-based migrations
handler/           # HTTP handlers for hubs, SKUs, and inventory
localstack/        # Localstack testing setup
middleware/        # Token authentication and request logging
models/            # GORM models for Hub, SKU, Inventory, etc.
router/            # API route definitions
server/            # Gin router and server setup
main.go            # Entry point of the application
docker-compose.yml # Redis and PostgreSQL setup
go.mod / go.sum    # Go module dependencies

Features
--------
- Full CRUD operations for Hubs and SKUs
- Inventory upsert with conflict resolution
- Hub and SKU validation endpoints
- Redis-based caching for validations
- Token-based authentication middleware
- Structured logging of requests
- i18n support for localizable messages and logs
- Go-based database migration tooling
- Docker Compose support for local development

Application Workflow
--------------------
1. Service initializes configs, DB, Redis, i18n, and Gin HTTP server
2. Token-based middleware validates the Authorization header
3. Validation endpoints check existence of hubs or SKUs
4. Redis is used as a cache; DB is checked on cache miss
5. /inventory/upsert allows inserting or updating inventory records
6. Standard CRUD operations for hub and SKU resources
7. Development environment runs with docker-compose

API Endpoints (Examples)
------------------------

Validate Hub
------------
GET /validate/hub/{hub_id}

Headers:
Authorization: Bearer secret-token

Response:
{
  "data": {
    "is_valid": true
  }
}

