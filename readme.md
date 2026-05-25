```md
# RBAC Service

A simple RBAC authorization service written in Go.

This project provides:

- RBAC HTTP service using Gin
- PostgreSQL storage
- Workspace-scoped authorization
- User, Role, Policy, Workspace entities
- Allow/Deny policy effect
- Wildcard policy matcher
- Go client SDK
- Simple shared secret authentication
- Docker Compose support
- Air hot reload for development
- Database auto migration

---

## Concept

This RBAC service uses `workspace_id` as the isolation boundary.

`workspace_id` can represent:

- tenant
- application
- organization
- workspace
- project
- environment

Authorization flow:

```text
User
  ↓
User roles in workspace
  ↓
Policies attached to roles
  ↓
Check resource + action
  ↓
Allow or deny
```

Policy priority:

```text
explicit deny > explicit allow > default deny
```

---

## Project Structure

```text
rbac/
 ├── cmd/
 │    ├── server/
 │    │    └── main.go
 │    └── cli/
 │         └── main.go
 │
 ├── app/
 │    └── http/
 │         ├── server.go
 │         ├── routes.go
 │         └── middleware.go
 │
 ├── handler/
 │    └── http/
 │         ├── authz_handler.go
 │         ├── workspace_handler.go
 │         ├── user_handler.go
 │         ├── role_handler.go
 │         └── policy_handler.go
 │
 ├── internal/
 │    ├── entity/
 │    ├── usecase/
 │    ├── repository/
 │    │    └── postgres/
 │    └── utils/
 │
 ├── client/
 │    ├── client.go
 │    ├── types.go
 │    ├── authz.go
 │    ├── workspace.go
 │    ├── user.go
 │    ├── role.go
 │    └── policy.go
 │
 ├── Dockerfile
 ├── docker-compose.yml
 ├── .air.toml
 └── go.mod
```

---

## Requirements

For local development:

- Go 1.25+
- PostgreSQL 16+
- Docker and Docker Compose

Optional:

- Air for hot reload

```bash
go install github.com/air-verse/air@latest
```

---

## Environment Variables

### Server

| Variable | Description | Default |
|---|---|---|
| `DATABASE_URL` | PostgreSQL connection URL | `postgres://postgres:postgres@localhost:5454/rbac?sslmode=disable` |
| `PORT` | HTTP server port | `10010` |
| `RBAC_SECRET` | Shared secret for API access | empty |
| `RBAC_AUTO_MIGRATE` | Auto run DB migration on startup | `true` |

Example:

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5454/rbac?sslmode=disable"
PORT=10010
RBAC_SECRET="super-secret-key"
RBAC_AUTO_MIGRATE=true
```

### Client

| Variable | Description |
|---|---|
| `RBAC_HOST` | RBAC service host |
| `RBAC_SECRET` | Shared secret |

Example:

```bash
RBAC_HOST="http://localhost:10010"
RBAC_SECRET="super-secret-key"
```

---

## Run with Docker Compose

### Normal mode

```bash
docker compose up --build
```

Service will be available at:

```text
http://localhost:10010
```

PostgreSQL will be available from host at:

```text
localhost:5454
```

Default secret:

```text
super-secret-key
```

---

## Run with Docker Compose Development Mode

Development mode uses Air for hot reload.

```bash
docker compose -f docker-compose.dev.yml up --build
```

When you edit Go files, the service will rebuild and restart automatically.

---

## Run Locally Without Docker

Start PostgreSQL first.

Example using Docker only for PostgreSQL:

```bash
docker run --name rbac-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=rbac \
  -p 5454:5454 \
  -d postgres:16
```

Run server:

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5454/rbac?sslmode=disable" \
RBAC_SECRET="super-secret-key" \
RBAC_AUTO_MIGRATE=true \
PORT=10010 \
go run ./cmd/server
```

---

## Run Locally with Air

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5454/rbac?sslmode=disable" \
RBAC_SECRET="super-secret-key" \
RBAC_AUTO_MIGRATE=true \
PORT=10010 \
air
```

---

## Database Migration

Migration can run automatically when server starts:

```bash
RBAC_AUTO_MIGRATE=true
```

Or manually using CLI:

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5454/rbac?sslmode=disable" \
go run ./cmd/cli migrate
```

---

## Authentication

Every API request must include:

```http
Authorization: Bearer <RBAC_SECRET>
```

Example:

```bash
-H "Authorization: Bearer super-secret-key"
```

If `RBAC_SECRET` is empty on the server, authentication middleware is disabled.

---

## API Usage with Curl

Base URL:

```text
http://localhost:10010
```

Secret:

```text
super-secret-key
```

---

## 1. Create Workspace

```bash
curl -X POST http://localhost:10010/v1/workspaces \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "id": "app_1",
    "name": "Application One"
  }'
```

---

## 2. Create User

```bash
curl -X POST http://localhost:10010/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "id": "user_1",
    "workspace_id": "app_1",
    "email": "john@example.com",
    "name": "John"
  }'
```

---

## 3. Create Role

```bash
curl -X POST http://localhost:10010/v1/roles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "id": "admin",
    "workspace_id": "app_1",
    "name": "Admin",
    "description": "Administrator role"
  }'
```

---

## 4. Create Allow Policy

This policy allows all project actions.

```bash
curl -X POST http://localhost:10010/v1/policies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "id": "project_all",
    "workspace_id": "app_1",
    "resource": "project",
    "action": "*",
    "effect": "allow"
  }'
```

---

## 5. Attach Policy to Role

```bash
curl -X POST http://localhost:10010/v1/roles/policies/attach \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "workspace_id": "app_1",
    "role_id": "admin",
    "policy_id": "project_all"
  }'
```

---

## 6. Assign Role to User

```bash
curl -X POST http://localhost:10010/v1/roles/assign \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "workspace_id": "app_1",
    "user_id": "user_1",
    "role_id": "admin"
  }'
```

---

## 7. Check Permission

```bash
curl -X POST http://localhost:10010/v1/authz/check \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "workspace_id": "app_1",
    "user_id": "user_1",
    "resource": "project",
    "action": "delete"
  }'
```

Expected response:

```json
{
  "allowed": true,
  "reason": "allow policy matched"
}
```

---

# Deny Policy Example

Create a deny policy:

```bash
curl -X POST http://localhost:10010/v1/policies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "id": "project_delete_denied",
    "workspace_id": "app_1",
    "resource": "project",
    "action": "delete",
    "effect": "deny"
  }'
```

Attach it to the role:

```bash
curl -X POST http://localhost:10010/v1/roles/policies/attach \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "workspace_id": "app_1",
    "role_id": "admin",
    "policy_id": "project_delete_denied"
  }'
```

Now check delete again:

```bash
curl -X POST http://localhost:10010/v1/authz/check \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer super-secret-key" \
  -d '{
    "workspace_id": "app_1",
    "user_id": "user_1",
    "resource": "project",
    "action": "delete"
  }'
```

Expected response:

```json
{
  "allowed": false,
  "reason": "explicit deny policy matched"
}
```

Deny always wins over allow.

---

## API Routes

```text
POST   /v1/authz/check

POST   /v1/workspaces
GET    /v1/workspaces/:workspace_id

POST   /v1/users
GET    /v1/workspaces/:workspace_id/users/:user_id

POST   /v1/roles
GET    /v1/workspaces/:workspace_id/roles
GET    /v1/workspaces/:workspace_id/roles/:role_id
POST   /v1/roles/assign
DELETE /v1/workspaces/:workspace_id/users/:user_id/roles/:role_id

POST   /v1/policies
GET    /v1/workspaces/:workspace_id/policies
GET    /v1/workspaces/:workspace_id/policies/:policy_id
POST   /v1/roles/policies/attach
DELETE /v1/workspaces/:workspace_id/roles/:role_id/policies/:policy_id
```

---

## Go Client SDK Usage

Install/import the client from your module:

```go
import rbacclient "github.com/your-org/rbac/client"
```

Example:

```go
package main

import (
	"context"
	"fmt"
	"log"

	rbacclient "github.com/your-org/rbac/client"
)

func main() {
	ctx := context.Background()

	rbac := rbacclient.New(rbacclient.Config{
		Host:   "http://localhost:10010",
		Secret: "super-secret-key",
	})

	result, err := rbac.Can(ctx, rbacclient.CheckRequest{
		WorkspaceID: "app_1",
		UserID:      "user_1",
		Resource:    "project",
		Action:      "delete",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("allowed:", result.Allowed)
	fmt.Println("reason:", result.Reason)
}
```

---

## Go Client Full Setup Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	rbacclient "github.com/your-org/rbac/client"
)

func main() {
	ctx := context.Background()

	rbac := rbacclient.New(rbacclient.Config{
		Host:   "http://localhost:10010",
		Secret: "super-secret-key",
	})

	_ = rbac.CreateWorkspace(ctx, rbacclient.Workspace{
		ID:   "app_1",
		Name: "Application One",
	})

	_ = rbac.CreateUser(ctx, rbacclient.User{
		ID:          "user_1",
		WorkspaceID: "app_1",
		Email:       "john@example.com",
		Name:        "John",
	})

	_ = rbac.CreateRole(ctx, rbacclient.Role{
		ID:          "admin",
		WorkspaceID: "app_1",
		Name:        "Admin",
		Description: "Administrator role",
	})

	_ = rbac.CreatePolicy(ctx, rbacclient.Policy{
		ID:          "project_all",
		WorkspaceID: "app_1",
		Resource:    "project",
		Action:      "*",
		Effect:      "allow",
	})

	_ = rbac.AttachPolicyToRole(ctx, rbacclient.RolePolicy{
		WorkspaceID: "app_1",
		RoleID:      "admin",
		PolicyID:    "project_all",
	})

	_ = rbac.AssignRoleToUser(ctx, rbacclient.UserRoleAssignment{
		WorkspaceID: "app_1",
		UserID:      "user_1",
		RoleID:      "admin",
	})

	result, err := rbac.Can(ctx, rbacclient.CheckRequest{
		WorkspaceID: "app_1",
		UserID:      "user_1",
		Resource:    "project",
		Action:      "delete",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("allowed:", result.Allowed)
	fmt.Println("reason:", result.Reason)
}
```

---

## CLI Usage

Run migration:

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5454/rbac?sslmode=disable" \
go run ./cmd/cli migrate
```

Seed test data:

```bash
RBAC_HOST="http://localhost:10010" \
RBAC_SECRET="super-secret-key" \
WORKSPACE_ID="app_1" \
USER_ID="user_1" \
go run ./cmd/cli seed
```

Check permission:

```bash
RBAC_HOST="http://localhost:10010" \
RBAC_SECRET="super-secret-key" \
WORKSPACE_ID="app_1" \
USER_ID="user_1" \
RESOURCE="project" \
ACTION="delete" \
go run ./cmd/cli check
```

Expected output:

```text
allowed: true
reason: allow policy matched
```

---

## Policy Matching

A policy contains:

```json
{
  "resource": "project",
  "action": "read",
  "effect": "allow"
}
```

The matcher supports wildcard:

```text
resource = "*"
action   = "*"
```

Examples:

| Policy Resource | Policy Action | Request | Result |
|---|---|---|---|
| `project` | `read` | `project:read` | match |
| `project` | `*` | `project:delete` | match |
| `*` | `read` | `invoice:read` | match |
| `*` | `*` | `billing:update` | match |
| `project` | `read` | `project:delete` | no match |

---

## Authorization Rules

```text
No role                  → denied
No matching allow policy → denied
Matching allow policy    → allowed
Matching deny policy     → denied
Deny always wins         → denied
```

Priority:

```text
deny > allow > default deny
```

---

## Example Roles and Policies

### Viewer

```json
{
  "resource": "project",
  "action": "read",
  "effect": "allow"
}
```

### Editor

```json
{
  "resource": "project",
  "action": "*",
  "effect": "allow"
}
```

### Restricted Editor

Allow all project actions:

```json
{
  "resource": "project",
  "action": "*",
  "effect": "allow"
}
```

Deny project deletion:

```json
{
  "resource": "project",
  "action": "delete",
  "effect": "deny"
}
```

### Admin

```json
{
  "resource": "*",
  "action": "*",
  "effect": "allow"
}
```

---

## Docker Compose Files

### `docker-compose.yml`

Use for normal local service:

```bash
docker compose up --build
```

### `docker-compose.dev.yml`

Use for development with hot reload:

```bash
docker compose -f docker-compose.dev.yml up --build
```

---

## Notes on IDs

This service uses `TEXT` IDs.

That is intentional because consuming applications may already have their own IDs.

Examples:

```text
workspace_id = "app_1"
user_id      = "auth0|abc123"
role_id      = "admin"
policy_id    = "project_read"
```

Recommended constraints:

```sql
CHECK (length(id) > 0 AND length(id) <= 128)
```

Recommended casing:

```text
Use lowercase IDs for roles, resources, and actions.
```

Example:

```text
admin
project
read
delete
```

---

## Security Notes

Current authentication uses a shared secret:

```http
Authorization: Bearer super-secret-key
```

This is fine for MVP/internal service usage.

For production, consider:

- per-application API keys
- hashed API keys in DB
- key rotation
- mTLS
- JWT validation
- audit logs
- rate limiting

---

## Production Recommendations

Before production, consider adding:

- versioned DB migrations using Goose, Atlas, or golang-migrate
- health endpoint `/healthz`
- readiness endpoint `/readyz`
- structured logging
- request ID middleware
- audit logs
- caching for permission checks
- integration tests
- CI pipeline
- better HTTP error mapping

---