
# IMS - Inventory Management System Microservice

This microservice manages **hubs**, **SKUs**, and **inventory operations** for the Omniful Inventory Management Platform. It provides robust CRUD APIs, validation endpoints, Redis-based caching, GORM-based database migrations, and internationalized responses.

---

## Project Structure

```
IMS/
├── cache/              # Redis cache setup and usage
├── configs/            # Configuration and environment management
├── db/                 # Database connection and migrations
├── handler/            # HTTP handlers for hubs, SKUs, and inventory
├── localstack/         # LocalStack testing support
├── middleware/         # Authentication and logging middleware
├── models/             # GORM models for Hub, SKU, Inventory, etc.
├── router/             # API route definitions
├── server/             # Gin server setup
├── docker-compose.yml  # Redis and PostgreSQL containers
├── go.mod / go.sum     # Go module dependencies
└── main.go             # Application entry point
```

---

## Features

- CRUD operations for **Hubs** and **SKUs**
- Inventory **upsert** with conflict resolution
- Validation endpoints for hub and SKU existence
- Redis-based **caching** for validation lookups
- Token-based **authentication middleware**
- Structured logging via middleware
- **i18n** support for multi-language response messages
- Go-native DB migration tooling
- Local development with **Docker Compose**

---

## Authentication

All APIs require an authorization header:

```
Authorization: Bearer my-secret-token
```

---

## Project Workflow

### 1. Service Initialization

- Loads configs, connects to DB, initializes Redis, sets up i18n, and starts the Gin server.

### 2. Request Handling

- Middleware checks for a valid bearer token.
- Logging and i18n are handled uniformly across routes.

### 3. Validation Endpoints

- `GET /validate/hub/:id` and `GET /validate/sku/:id` check Redis for cached results first.
- If not found in cache, fallback to DB lookup.
- Returns `{ "is_valid": true/false }`.

### 4. Inventory Upsert

- `POST /inventory/upsert` allows inserting or updating inventory quantities.
- Conflict resolution is handled at the DB level.

### 5. CRUD Endpoints

- Full POST, GET, PUT, DELETE for **Hub** and **SKU** models.

---

## API Endpoints Example

### Validate Hub

**GET** `/validate/hub/{hub_id}`

**Headers:**
```
Authorization: Bearer my-secret-token
```

**Response:**
```json
{
  "data": {
    "is_valid": true
  }
}
```
