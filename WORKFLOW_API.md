# Todo Workflow API Documentation

## Overview
The workflow API manages todos through a state machine with three statuses: **Draft**, **Review**, and **Approved**.

## Endpoints

### 1. Create Todo Task
**POST** `/workflow/todos`

Creates a new todo in Draft status.

**Request Body:**
```json
{
  "title": "Complete project documentation",
  "description": "Write comprehensive API docs",
  "assigned_to": "user123"
}
```

**Response:** `201 Created`
```json
{
  "ID": "uuid-here",
  "Title": "Complete project documentation",
  "Description": "Write comprehensive API docs",
  "AssignedTo": "user123",
  "Status": "Draft",
  "CreatedAt": "2024-12-02T10:00:00Z",
  "UpdatedAt": "2024-12-02T10:00:00Z",
  "ReviewedBy": "",
  "ApprovedBy": ""
}
```

---

### 2. Submit for Review
**POST** `/workflow/todos/{id}/submit/{submitted_by}`

Submits a draft todo for review. Only the assigned user can submit.

**Example:** `POST /workflow/todos/550e8400-e29b-41d4-a716-446655440000/submit/user123`

**No Request Body Required**

**Response:** `200 OK`
```json
{
  "message": "Todo submitted for review successfully"
}
```

**Business Rules:**
- Todo must be in **Draft** status
- Only the assigned user can submit

---

### 3. Approve Todo
**POST** `/workflow/todos/{id}/approve/{approved_by}`

Approves a todo in review. Only the reviewer can approve.

**Example:** `POST /workflow/todos/550e8400-e29b-41d4-a716-446655440000/approve/user123`

**No Request Body Required**

**Response:** `200 OK`
```json
{
  "message": "Todo approved successfully"
}
```

**Business Rules:**
- Todo must be in **Review** status
- Only the reviewer (who submitted it) can approve

---

### 4. Reject Todo
**POST** `/workflow/todos/{id}/reject/{rejected_by}`

Rejects a todo and sends it back to Draft status.

**Example:** `POST /workflow/todos/550e8400-e29b-41d4-a716-446655440000/reject/user123`

**No Request Body Required**

**Response:** `200 OK`
```json
{
  "message": "Todo rejected successfully"
}
```

**Business Rules:**
- Todo must be in **Review** status
- Resets status to **Draft** and clears the reviewer

---

### 5. Get Todos by User
**GET** `/workflow/todos/user?user_id={userId}`

Retrieves all todos assigned to a specific user.

**Response:** `200 OK`
```json
[
  {
    "ID": "uuid-1",
    "Title": "Task 1",
    "Description": "Description 1",
    "AssignedTo": "user123",
    "Status": "Draft",
    "CreatedAt": "2024-12-02T10:00:00Z",
    "UpdatedAt": "2024-12-02T10:00:00Z",
    "ReviewedBy": "",
    "ApprovedBy": ""
  }
]
```

---

### 6. Get Todos by Status
**GET** `/workflow/todos/status?status={status}`

Retrieves all todos with a specific status.

**Valid Status Values:** `Draft`, `Review`, `Approved`

**Response:** `200 OK`
```json
[
  {
    "ID": "uuid-1",
    "Title": "Task in Review",
    "Description": "Description",
    "AssignedTo": "user123",
    "Status": "Review",
    "CreatedAt": "2024-12-02T10:00:00Z",
    "UpdatedAt": "2024-12-02T10:00:00Z",
    "ReviewedBy": "user123",
    "ApprovedBy": ""
  }
]
```

---

## Workflow State Machine

```
Draft → Review → Approved
  ↑       ↓
  └───────┘ (reject)
```

1. **Draft**: Initial state when todo is created
2. **Review**: Todo submitted for review by assigned user
3. **Approved**: Todo approved by reviewer
4. **Reject**: Sends todo back to Draft status

---

## Database Schema

**Table:** `todo_tasks`

| Column       | Type          | Description                    |
|--------------|---------------|--------------------------------|
| id           | NVARCHAR(36)  | Primary key (UUID)             |
| title        | NVARCHAR(200) | Todo title                     |
| description  | NVARCHAR(500) | Todo description               |
| assigned_to  | NVARCHAR(100) | User assigned to the todo      |
| status       | NVARCHAR(20)  | Draft, Review, or Approved     |
| reviewed_by  | NVARCHAR(100) | User who submitted for review  |
| approved_by  | NVARCHAR(100) | User who approved the todo     |
| created_at   | DATETIME2     | Creation timestamp             |
| updated_at   | DATETIME2     | Last update timestamp          |

---

## Migration

Run the migration to create the table:

```bash
# Navigate to your project directory
cd c:\Users\Admin\GolandProjects\todo-api

# Run migration (adjust connection string as needed)
migrate -path ./migrations -database "your_connection_string" up
```
