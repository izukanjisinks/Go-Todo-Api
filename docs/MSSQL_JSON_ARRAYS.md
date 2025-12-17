# SQL Server Array Handling with JSON

## The Problem

**SQL Server does NOT have native array types** like PostgreSQL. You cannot use `TEXT[]` or similar array columns in SQL Server.

## The Solution: Store as JSON

Since SQL Server doesn't support arrays, we store the permissions list as a **JSON string** in an `NVARCHAR(MAX)` column.

### Database Schema

```sql
CREATE TABLE roles (
    role_id UNIQUEIDENTIFIER PRIMARY KEY,
    name NVARCHAR(100) NOT NULL,
    permissions NVARCHAR(MAX) NOT NULL DEFAULT '[]', -- JSON array as string
);
```

### Example Data in Database

```
role_id: 123e4567-e89b-12d3-a456-426614174000
name: Admin
permissions: ["users:read","users:create","content:read","content:update"]
```

The `permissions` column stores a JSON array as a string.

## Go Implementation

### Writing to Database (INSERT/UPDATE)

```go
import "encoding/json"

// Convert Go []string to JSON string
permissions := []string{"users:read", "users:create", "users:update"}
permissionsJSON, err := json.Marshal(permissions)
// Result: ["users:read","users:create","users:update"]

// Insert into database
db.Exec(
    "INSERT INTO roles (name, permissions) VALUES (@p1, @p2)",
    "Admin",
    string(permissionsJSON), // Store as JSON string
)
```

### Reading from Database (SELECT)

```go
import "encoding/json"

var permissionsJSON string
var permissions []string

// Read JSON string from database
db.QueryRow("SELECT permissions FROM roles WHERE name = @p1", "Admin").
    Scan(&permissionsJSON)

// Convert JSON string to Go []string
err := json.Unmarshal([]byte(permissionsJSON), &permissions)
// Result: []string{"users:read", "users:create", "users:update"}
```

## Complete Example from RoleService

### CreateRole - Writing JSON Array

```go
func (s *RoleService) CreateRole(role *models.Role) error {
    query := `
        INSERT INTO roles (role_id, name, description, permissions)
        VALUES (@p1, @p2, @p3, @p4)
    `
    
    // Step 1: Convert []string to JSON
    permissionsJSON, err := json.Marshal(role.Permissions)
    if err != nil {
        return fmt.Errorf("failed to marshal permissions: %w", err)
    }
    
    // Step 2: Insert JSON string into database
    _, err = s.db.Exec(
        query, 
        role.RoleId, 
        role.Name, 
        role.Description, 
        string(permissionsJSON), // Convert []byte to string
    )
    
    return err
}
```

### GetRoleByID - Reading JSON Array

```go
func (s *RoleService) GetRoleByID(roleID uuid.UUID) (*models.Role, error) {
    query := `
        SELECT role_id, name, description, permissions
        FROM roles
        WHERE role_id = @p1
    `
    
    role := &models.Role{}
    var permissionsJSON string // Temporary variable for JSON string
    
    // Step 1: Read JSON string from database
    err := s.db.QueryRow(query, roleID).Scan(
        &role.RoleId,
        &role.Name,
        &role.Description,
        &permissionsJSON, // Scan into string
    )
    if err != nil {
        return nil, err
    }
    
    // Step 2: Convert JSON string to []string
    err = json.Unmarshal([]byte(permissionsJSON), &role.Permissions)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal permissions: %w", err)
    }
    
    return role, nil
}
```

## How It Works

### Flow Diagram

```
Go Application                    SQL Server Database
-------------                     -------------------

[]string                          NVARCHAR(MAX)
["users:read",      --Marshal-->  '["users:read",
 "users:create"]                   "users:create"]'

                    <-Unmarshal--
```

### Step-by-Step

1. **In Go**: Permissions are `[]string{"users:read", "users:create"}`
2. **Marshal**: `json.Marshal()` converts to `[]byte` → `["users:read","users:create"]`
3. **Convert**: Cast to `string` for SQL Server
4. **Store**: SQL Server stores as `NVARCHAR(MAX)` text
5. **Retrieve**: SQL Server returns the JSON string
6. **Unmarshal**: `json.Unmarshal()` converts back to `[]string`

## Migration File (SQL Server)

```sql
-- Create roles table
CREATE TABLE roles (
    role_id UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWID(),
    name NVARCHAR(100) UNIQUE NOT NULL,
    permissions NVARCHAR(MAX) NOT NULL DEFAULT '[]',
    created_at DATETIME2 DEFAULT GETDATE()
);

-- Insert with JSON array
INSERT INTO roles (role_id, name, permissions) VALUES
    (NEWID(), 'Admin', '["users:read","users:create","users:update"]');
```

## Querying JSON in SQL Server

SQL Server 2016+ has JSON functions you can use:

### Check if permission exists

```sql
-- Check if array contains a value
SELECT * FROM roles
WHERE JSON_VALUE(permissions, '$[0]') = 'users:read'
   OR JSON_VALUE(permissions, '$[1]') = 'users:read'
   OR JSON_VALUE(permissions, '$[2]') = 'users:read';

-- Better: Use OPENJSON (SQL Server 2016+)
SELECT r.name
FROM roles r
CROSS APPLY OPENJSON(r.permissions) AS perms
WHERE perms.value = 'users:read';
```

### Get all permissions for a role

```sql
SELECT r.name, perms.value AS permission
FROM roles r
CROSS APPLY OPENJSON(r.permissions) AS perms;
```

## Comparison: PostgreSQL vs SQL Server

| Feature | PostgreSQL | SQL Server |
|---------|-----------|------------|
| **Native Arrays** | ✅ `TEXT[]` | ❌ No native arrays |
| **Storage** | Native array type | JSON string in `NVARCHAR(MAX)` |
| **Go Handling** | `pq.Array()` wrapper | `json.Marshal/Unmarshal` |
| **Query Support** | Array operators (`@>`, `&&`) | JSON functions (`OPENJSON`) |
| **Performance** | Faster (native type) | Slightly slower (parsing) |
| **Size** | More efficient | Larger (JSON overhead) |

## Key Differences from PostgreSQL

### PostgreSQL (What We Can't Use)
```go
// ❌ This only works with PostgreSQL
import "github.com/lib/pq"
db.Exec(query, pq.Array(permissions))
```

### SQL Server (What We Use Instead)
```go
// ✅ This works with SQL Server
import "encoding/json"
permJSON, _ := json.Marshal(permissions)
db.Exec(query, string(permJSON))
```

## Advantages of JSON Approach

1. **Database Agnostic**: Works with any SQL database (MySQL, SQLite, etc.)
2. **Flexible**: Can store complex nested structures if needed
3. **Standard**: JSON is universally supported
4. **Queryable**: SQL Server has JSON functions for querying

## Disadvantages

1. **Slower**: Requires JSON parsing on every read/write
2. **Larger**: JSON has overhead (`["`, `","`, `"]`)
3. **Type Safety**: Database doesn't enforce array element types
4. **Indexing**: Can't directly index array elements (need computed columns)

## Best Practices

1. **Always validate JSON**: Check `json.Marshal/Unmarshal` errors
2. **Use NVARCHAR(MAX)**: Don't limit JSON string length
3. **Default to empty array**: Use `'[]'` not `NULL`
4. **Consider caching**: Parse JSON once, cache in memory
5. **Use JSON functions**: Leverage SQL Server's `OPENJSON` for queries

## Alternative: Separate Table

For better performance and queryability, consider a junction table:

```sql
CREATE TABLE role_permissions (
    role_id UNIQUEIDENTIFIER,
    permission NVARCHAR(100),
    PRIMARY KEY (role_id, permission),
    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);
```

**Pros**: Better indexing, easier queries, normalized  
**Cons**: More complex queries, more tables to manage

For the RBAC system, **JSON is simpler and sufficient** for most use cases.

## Summary

- **SQL Server doesn't have arrays** → Use JSON strings instead
- **Write**: `json.Marshal([]string)` → `string` → SQL Server
- **Read**: SQL Server → `string` → `json.Unmarshal()` → `[]string`
- **Storage**: `NVARCHAR(MAX)` column with JSON array string
- **Works perfectly** for the RBAC permissions system
