# PostgreSQL Array Handling in Go

## Overview

PostgreSQL supports native array types (e.g., `TEXT[]`, `INTEGER[]`), but Go's `database/sql` package doesn't automatically handle them. You need the `github.com/lib/pq` driver to properly work with PostgreSQL arrays.

## The Problem

When you try to insert or scan a Go `[]string` directly with PostgreSQL's `TEXT[]` type, it won't work:

```go
// ❌ This FAILS - Go slice is not automatically converted
permissions := []string{"users:read", "users:create"}
db.Exec("INSERT INTO roles (permissions) VALUES ($1)", permissions)
// Error: cannot convert []string to PostgreSQL array
```

## The Solution: `pq.Array()`

The `lib/pq` package provides `pq.Array()` wrapper to handle conversions:

### Writing to Database (INSERT/UPDATE)

```go
import "github.com/lib/pq"

permissions := []string{"users:read", "users:create", "users:update"}

// ✅ Wrap with pq.Array() when inserting
db.Exec(
    "INSERT INTO roles (name, permissions) VALUES ($1, $2)",
    "Admin",
    pq.Array(permissions), // Converts Go []string to PostgreSQL TEXT[]
)
```

### Reading from Database (SELECT)

```go
var permissions []string

// ✅ Wrap the destination pointer with pq.Array() when scanning
db.QueryRow("SELECT permissions FROM roles WHERE name = $1", "Admin").
    Scan(pq.Array(&permissions)) // Converts PostgreSQL TEXT[] to Go []string
```

## How It Works

### When Writing (INSERT/UPDATE)
```go
pq.Array(role.Permissions)
```
- Takes your Go `[]string`
- Converts it to PostgreSQL array format: `{"users:read","users:create"}`
- PostgreSQL stores it as a native `TEXT[]` type

### When Reading (SELECT)
```go
pq.Array(&role.Permissions)
```
- PostgreSQL returns the array as: `{"users:read","users:create"}`
- `pq.Array()` parses it and populates your Go `[]string`

## Real Example from RoleService

### CreateRole - Writing Array
```go
func (s *RoleService) CreateRole(role *models.Role) error {
    query := `
        INSERT INTO roles (role_id, name, description, permissions)
        VALUES ($1, $2, $3, $4)
    `
    
    // Wrap permissions slice with pq.Array()
    _, err := s.db.Exec(
        query, 
        role.RoleId, 
        role.Name, 
        role.Description, 
        pq.Array(role.Permissions), // ✅ Convert []string to TEXT[]
    )
    
    return err
}
```

### GetRoleByID - Reading Array
```go
func (s *RoleService) GetRoleByID(roleID uuid.UUID) (*models.Role, error) {
    query := `
        SELECT role_id, name, description, permissions
        FROM roles
        WHERE role_id = $1
    `
    
    role := &models.Role{}
    err := s.db.QueryRow(query, roleID).Scan(
        &role.RoleId,
        &role.Name,
        &role.Description,
        pq.Array(&role.Permissions), // ✅ Convert TEXT[] to []string
    )
    
    return role, err
}
```

## Database Schema

```sql
CREATE TABLE roles (
    role_id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    permissions TEXT[] NOT NULL DEFAULT '{}', -- PostgreSQL array type
);

-- Insert with PostgreSQL ARRAY syntax
INSERT INTO roles (role_id, name, permissions) VALUES
    (gen_random_uuid(), 'Admin', ARRAY['users:read', 'users:create', 'users:update']);
```

## Common Patterns

### 1. Insert with Array
```go
_, err := db.Exec(
    "INSERT INTO roles (name, permissions) VALUES ($1, $2)",
    "Moderator",
    pq.Array([]string{"content:read", "content:create"}),
)
```

### 2. Update with Array
```go
_, err := db.Exec(
    "UPDATE roles SET permissions = $1 WHERE name = $2",
    pq.Array([]string{"users:read", "content:read"}),
    "User",
)
```

### 3. Query Single Row
```go
var permissions []string
err := db.QueryRow(
    "SELECT permissions FROM roles WHERE name = $1", 
    "Admin",
).Scan(pq.Array(&permissions))
```

### 4. Query Multiple Rows
```go
rows, _ := db.Query("SELECT name, permissions FROM roles")
defer rows.Close()

for rows.Next() {
    var name string
    var permissions []string
    
    rows.Scan(&name, pq.Array(&permissions))
    fmt.Printf("%s: %v\n", name, permissions)
}
```

## What Happens Without pq.Array()?

### Without pq.Array() on INSERT
```go
// ❌ WRONG
db.Exec("INSERT INTO roles (permissions) VALUES ($1)", []string{"read", "write"})

// Error: sql: converting argument $1 type: unsupported type []string, a slice of string
```

### Without pq.Array() on SELECT
```go
// ❌ WRONG
var permissions []string
db.QueryRow("SELECT permissions FROM roles").Scan(&permissions)

// Error: sql: Scan error on column index 0: unsupported Scan, storing driver.Value type []uint8 into type *[]string
```

## Key Takeaways

1. **Always use `pq.Array()`** when working with PostgreSQL array columns
2. **On INSERT/UPDATE**: Wrap the Go slice → `pq.Array(mySlice)`
3. **On SELECT**: Wrap the destination pointer → `pq.Array(&mySlice)`
4. **Import required**: `import "github.com/lib/pq"`
5. **Works with**: `[]string`, `[]int`, `[]int64`, `[]float64`, `[]bool`, etc.

## Alternative: JSON Storage

If you don't want to deal with PostgreSQL arrays, you can store as JSON:

```sql
-- Use JSONB instead of TEXT[]
CREATE TABLE roles (
    permissions JSONB NOT NULL DEFAULT '[]'
);
```

```go
import "encoding/json"

// Write
permJSON, _ := json.Marshal(permissions)
db.Exec("INSERT INTO roles (permissions) VALUES ($1)", permJSON)

// Read
var permJSON []byte
db.QueryRow("SELECT permissions FROM roles").Scan(&permJSON)
json.Unmarshal(permJSON, &permissions)
```

But PostgreSQL arrays are more efficient and provide better query capabilities (e.g., `ANY`, `@>` operators).

## References

- [lib/pq documentation](https://pkg.go.dev/github.com/lib/pq)
- [PostgreSQL Array Types](https://www.postgresql.org/docs/current/arrays.html)
