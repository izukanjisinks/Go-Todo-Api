# Role-Based Access Control (RBAC) System

## Overview

This Todo API implements a comprehensive Role-Based Access Control (RBAC) system that replaces the simple `is_admin` boolean field with a flexible permission-based system.

## Architecture

### Components

1. **Role Model** (`internal/models/role.go`)
   - Defines roles with associated permissions
   - Permission constants for easy reference
   - Predefined roles (Super Admin, Admin, Moderator, User)

2. **User Model** (`internal/models/user.go`)
   - Updated with `RoleID` and `Role` fields
   - `HasPermission()` method for permission checking
   - Backward compatible with `IsAdmin` field

3. **Role Service** (`internal/services/role_service.go`)
   - CRUD operations for roles
   - Role assignment to users
   - Permission checking
   - Predefined role initialization

4. **RBAC Middleware** (`internal/middleware/rbac.go`)
   - `RequirePermission()` - Single permission check
   - `RequireAnyPermission()` - At least one permission
   - `RequireAllPermissions()` - All permissions required
   - `RequireRole()` - Specific role required

## Permission Categories

### User Management
- `users:read` - View user information
- `users:create` - Create new users
- `users:update` - Update user information
- `users:delete` - Delete users

### Content Management
- `content:read` - View content
- `content:create` - Create new content
- `content:update` - Update existing content
- `content:delete` - Delete content

### System Settings
- `settings:read` - View system settings
- `settings:update` - Modify system settings

### Reports
- `reports:view` - View reports
- `reports:export` - Export reports

### Roles
- `roles:manage` - Manage roles and permissions

## Predefined Roles

### Super Admin
**Description:** Full system access with all permissions

**Permissions:**
- All user management permissions
- All content management permissions
- All system settings permissions
- All reports permissions
- Role management

**Use Case:** System administrators with complete control

### Admin
**Description:** Administrative access with most permissions

**Permissions:**
- All user management permissions
- All content management permissions
- Settings read (no update)
- All reports permissions

**Use Case:** Department managers or team leads

### Moderator
**Description:** Content management and user viewing access

**Permissions:**
- User read only
- All content management permissions

**Use Case:** Content moderators or editors

### User
**Description:** Basic user access with read permissions

**Permissions:**
- Content read
- Reports view

**Use Case:** Regular users of the system

## Database Schema

### Roles Table
```sql
CREATE TABLE roles (
    role_id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    permissions TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Users Table Update
```sql
ALTER TABLE users ADD COLUMN role_id UUID;
ALTER TABLE users ADD CONSTRAINT fk_users_role
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE SET NULL;
```

## Usage Examples

### 1. Using Middleware

```go
import (
    "todo-api/internal/middleware"
    "todo-api/internal/models"
)

// Require a specific permission
http.Handle("/api/users", 
    middleware.RequirePermission(models.PermUsersRead)(
        http.HandlerFunc(getUsersHandler),
    ),
)

// Require any of multiple permissions
http.Handle("/api/content", 
    middleware.RequireAnyPermission(
        models.PermContentRead,
        models.PermContentCreate,
    )(
        http.HandlerFunc(contentHandler),
    ),
)

// Require all permissions
http.Handle("/api/admin/settings", 
    middleware.RequireAllPermissions(
        models.PermSettingsRead,
        models.PermSettingsUpdate,
    )(
        http.HandlerFunc(settingsHandler),
    ),
)

// Require specific role
http.Handle("/api/admin/roles", 
    middleware.RequireRole(models.RoleSuperAdmin)(
        http.HandlerFunc(roleManagementHandler),
    ),
)
```

### 2. Programmatic Permission Checking

```go
func handler(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*models.User)
    
    if !user.HasPermission(models.PermUsersRead) {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    
    // Handle request
}
```

### 3. Using Role Service

```go
roleService := services.NewRoleService(db)

// Initialize predefined roles (run once at startup)
err := roleService.InitializePredefinedRoles()

// Assign role to user
err = roleService.AssignRoleToUser(userID, roleID)

// Check permission
hasPermission, err := roleService.CheckPermission(userID, models.PermUsersRead)

// Create custom role
customRole := &models.Role{
    Name:        "Content Editor",
    Description: "Can edit content but not delete",
    Permissions: []string{
        models.PermContentRead,
        models.PermContentCreate,
        models.PermContentUpdate,
    },
}
err = roleService.CreateRole(customRole)
```

## Migration Guide

### From `is_admin` to RBAC

1. **Run the migration:**
   ```bash
   migrate -path ./migrations -database "postgres://..." up
   ```

2. **Initialize predefined roles:**
   ```go
   roleService := services.NewRoleService(db)
   err := roleService.InitializePredefinedRoles()
   ```

3. **Migrate existing users:**
   ```sql
   -- Assign Super Admin role to existing admins
   UPDATE users u
   SET role_id = (SELECT role_id FROM roles WHERE name = 'Super Admin')
   WHERE u.is_admin = true;
   
   -- Assign User role to non-admins
   UPDATE users u
   SET role_id = (SELECT role_id FROM roles WHERE name = 'User')
   WHERE u.is_admin = false;
   ```

4. **Update your routes:**
   - Replace `RequireAdmin()` middleware with `RequirePermission()` or `RequireRole()`
   - Use appropriate permission constants

5. **Backward Compatibility:**
   - The `is_admin` field is kept for backward compatibility
   - The `User.HasPermission()` method falls back to `is_admin` if no role is assigned
   - You can gradually migrate and eventually remove the `is_admin` field

## Creating Custom Roles

```go
// Example: Create a "Report Analyst" role
reportAnalyst := &models.Role{
    Name:        "Report Analyst",
    Description: "Can view and export reports",
    Permissions: []string{
        models.PermReportsView,
        models.PermReportsExport,
        models.PermContentRead,
    },
}

err := roleService.CreateRole(reportAnalyst)
```

## Best Practices

1. **Use Permission Constants:** Always use the predefined permission constants from `models` package
2. **Principle of Least Privilege:** Assign minimum permissions needed for a role
3. **Regular Audits:** Periodically review role permissions
4. **Custom Roles:** Create custom roles for specific use cases rather than modifying predefined ones
5. **Testing:** Test permission checks thoroughly, especially for sensitive operations
6. **Logging:** Log permission denials for security auditing

## Security Considerations

1. **Permission Checks:** Always check permissions on the server side, never trust client-side checks
2. **Role Assignment:** Only Super Admins should be able to assign roles
3. **Sensitive Operations:** Use `RequireAllPermissions()` for operations requiring multiple permissions
4. **Database Constraints:** Foreign key constraints ensure data integrity
5. **Null Handling:** Users without roles should have minimal or no permissions

## API Endpoints (Implemented)

```
GET    /roles                  - List all roles (requires authentication)
GET    /roles/{id}             - Get role details (requires authentication)
POST   /roles                  - Create role (requires authentication + roles:manage recommended)
PUT    /roles/{id}             - Update role (requires authentication + roles:manage recommended)
DELETE /roles/{id}             - Delete role (requires authentication + roles:manage recommended)
POST   /users/{id}/role        - Assign role to user (requires authentication + roles:manage recommended)
GET    /users/{id}/permissions - Get user permissions (requires authentication)
```

**Note:** All routes are protected with JWT authentication. For production, you should add RBAC middleware to restrict role management endpoints to users with `roles:manage` permission.

## Troubleshooting

### User has no permissions
- Check if user has a role assigned (`role_id` is not null)
- Verify the role has the required permissions
- Check if the role relationship is loaded in the query

### Permission check always fails
- Ensure the user's role is loaded (use JOIN in query)
- Verify permission constant spelling
- Check if middleware is applied in correct order (auth before RBAC)

### Migration fails
- Ensure all existing users have valid data
- Check database permissions
- Verify PostgreSQL version supports array types

## Future Enhancements

- [ ] Permission inheritance/hierarchies
- [ ] Time-based permissions (temporary access)
- [ ] Resource-level permissions (e.g., "can edit own content")
- [ ] Permission groups/categories
- [ ] Audit log for permission changes
- [ ] API for dynamic permission management
