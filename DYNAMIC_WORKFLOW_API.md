# Dynamic Workflow System API Documentation

## Overview
The dynamic workflow system allows users to create custom workflows with configurable steps and transitions. Tasks can then be run through these workflows with automatic state management and audit trails.

---

## Architecture

### Core Components:
1. **Workflows** - Templates that define the flow
2. **Steps** - States/stages in a workflow
3. **Transitions** - Rules for moving between steps
4. **Instances** - Tasks running through a workflow
5. **History** - Audit trail of all transitions

---

## Admin API (Workflow Configuration)

### 1. Create Workflow
**POST** `/api/workflows`

Creates a new workflow template.

**Request Body:**
```json
{
  "name": "Standard Approval Flow",
  "description": "Three-step approval process",
  "created_by": "admin_user"
}
```

**Response:** `201 Created`
```json
{
  "id": "workflow-uuid",
  "name": "Standard Approval Flow",
  "description": "Three-step approval process",
  "is_active": true,
  "created_by": "admin_user",
  "created_at": "2024-12-02T14:00:00Z",
  "updated_at": "2024-12-02T14:00:00Z"
}
```

---

### 2. Get All Workflows
**GET** `/api/workflows`

Retrieves all active workflows.

**Response:** `200 OK`
```json
[
  {
    "id": "workflow-uuid",
    "name": "Standard Approval Flow",
    "description": "Three-step approval process",
    "is_active": true,
    "created_by": "admin_user",
    "created_at": "2024-12-02T14:00:00Z",
    "updated_at": "2024-12-02T14:00:00Z"
  }
]
```

---

### 3. Get Workflow
**GET** `/api/workflows/{id}`

Retrieves a specific workflow.

---

### 4. Create Step
**POST** `/api/workflows/{workflow_id}/steps`

Creates a new step in a workflow.

**Request Body:**
```json
{
  "step_name": "Draft",
  "step_order": 1,
  "is_start_step": true,
  "is_end_step": false,
  "allowed_roles": ["employee", "manager"]
}
```

**Response:** `201 Created`
```json
{
  "id": "step-uuid",
  "workflow_id": "workflow-uuid",
  "step_name": "Draft",
  "step_order": 1,
  "is_start_step": true,
  "is_end_step": false,
  "allowed_roles": ["employee", "manager"],
  "created_at": "2024-12-02T14:00:00Z"
}
```

---

### 5. Get Workflow Steps
**GET** `/api/workflows/{workflow_id}/steps`

Retrieves all steps for a workflow, ordered by `step_order`.

---

### 6. Create Transition
**POST** `/api/workflows/{workflow_id}/transitions`

Creates a transition between two steps.

**Request Body:**
```json
{
  "from_step_id": "draft-step-uuid",
  "to_step_id": "review-step-uuid",
  "action_name": "submit",
  "condition_type": "assigned_user_only",
  "condition_value": ""
}
```

**Condition Types:**
- `assigned_user_only` - Only the assigned user can perform this action
- `creator_only` - Only the task creator can perform this action
- `not_assigned_user` - Anyone except the assigned user (for approvals)
- `any_user` - Any authenticated user
- `""` (empty) - No restrictions

**Response:** `201 Created`
```json
{
  "id": "transition-uuid",
  "workflow_id": "workflow-uuid",
  "from_step_id": "draft-step-uuid",
  "to_step_id": "review-step-uuid",
  "action_name": "submit",
  "condition_type": "assigned_user_only",
  "condition_value": "",
  "created_at": "2024-12-02T14:00:00Z"
}
```

---

### 7. Get Workflow Transitions
**GET** `/api/workflows/{workflow_id}/transitions`

Retrieves all transitions for a workflow.

---

## Task API (Workflow Execution)

### 1. Start Task
**POST** `/api/tasks`

Creates a new task and starts it in a workflow.

**Request Body:**
```json
{
  "workflow_id": "workflow-uuid",
  "title": "Complete project documentation",
  "description": "Write comprehensive API docs",
  "task_data": "{\"priority\": \"high\", \"deadline\": \"2024-12-31\"}",
  "assigned_to": "user123",
  "created_by": "manager456"
}
```

**Response:** `201 Created`
```json
{
  "id": "instance-uuid",
  "workflow_id": "workflow-uuid",
  "current_step_id": "draft-step-uuid",
  "title": "Complete project documentation",
  "description": "Write comprehensive API docs",
  "task_data": "{\"priority\": \"high\"}",
  "assigned_to": "user123",
  "created_by": "manager456",
  "created_at": "2024-12-02T14:00:00Z",
  "updated_at": "2024-12-02T14:00:00Z"
}
```

---

### 2. Get Task
**GET** `/api/tasks/{instance_id}?user_id={user_id}`

Retrieves a task with current step info and available actions for the user.

**Response:** `200 OK`
```json
{
  "id": "instance-uuid",
  "workflow_id": "workflow-uuid",
  "current_step_id": "draft-step-uuid",
  "title": "Complete project documentation",
  "description": "Write comprehensive API docs",
  "task_data": "{\"priority\": \"high\"}",
  "assigned_to": "user123",
  "created_by": "manager456",
  "created_at": "2024-12-02T14:00:00Z",
  "updated_at": "2024-12-02T14:00:00Z",
  "current_step_name": "Draft",
  "workflow_name": "Standard Approval Flow",
  "available_actions": [
    {
      "action_name": "submit",
      "to_step_name": "Review",
      "transition_id": "transition-uuid"
    }
  ]
}
```

---

### 3. Execute Action
**POST** `/api/tasks/{instance_id}/execute`

Executes a workflow action (moves task to next step).

**Request Body:**
```json
{
  "action_name": "submit",
  "user_id": "user123",
  "comments": "Ready for review"
}
```

**Response:** `200 OK`
```json
{
  "message": "Action executed successfully"
}
```

---

### 4. Get Available Actions
**GET** `/api/tasks/{instance_id}/actions?user_id={user_id}`

Retrieves actions the user can currently perform on the task.

**Response:** `200 OK`
```json
[
  {
    "action_name": "submit",
    "to_step_name": "Review",
    "transition_id": "transition-uuid"
  }
]
```

---

### 5. Get Task History
**GET** `/api/tasks/{instance_id}/history`

Retrieves the complete audit trail for a task.

**Response:** `200 OK`
```json
[
  {
    "id": "history-uuid-2",
    "instance_id": "instance-uuid",
    "from_step_id": "draft-step-uuid",
    "to_step_id": "review-step-uuid",
    "action_taken": "submit",
    "performed_by": "user123",
    "comments": "Ready for review",
    "timestamp": "2024-12-02T15:00:00Z"
  },
  {
    "id": "history-uuid-1",
    "instance_id": "instance-uuid",
    "from_step_id": null,
    "to_step_id": "draft-step-uuid",
    "action_taken": "created",
    "performed_by": "manager456",
    "comments": "Workflow instance created",
    "timestamp": "2024-12-02T14:00:00Z"
  }
]
```

---

### 6. Get Tasks by User
**GET** `/api/tasks/user?user_id={user_id}`

Retrieves all tasks assigned to a user.

---

### 7. Get Tasks by Workflow
**GET** `/api/workflows/{workflow_id}/tasks`

Retrieves all tasks running in a specific workflow.

---

## Example: Complete Workflow Setup

### Step 1: Create Workflow
```bash
POST /api/workflows
{
  "name": "Standard Approval",
  "description": "Draft → Review → Approved",
  "created_by": "admin"
}
# Returns: workflow_id
```

### Step 2: Create Steps
```bash
POST /api/workflows/{workflow_id}/steps
{
  "step_name": "Draft",
  "step_order": 1,
  "is_start_step": true,
  "is_end_step": false
}
# Returns: draft_step_id

POST /api/workflows/{workflow_id}/steps
{
  "step_name": "Review",
  "step_order": 2,
  "is_start_step": false,
  "is_end_step": false
}
# Returns: review_step_id

POST /api/workflows/{workflow_id}/steps
{
  "step_name": "Approved",
  "step_order": 3,
  "is_start_step": false,
  "is_end_step": true
}
# Returns: approved_step_id
```

### Step 3: Create Transitions
```bash
# Draft → Review (submit)
POST /api/workflows/{workflow_id}/transitions
{
  "from_step_id": "{draft_step_id}",
  "to_step_id": "{review_step_id}",
  "action_name": "submit",
  "condition_type": "assigned_user_only"
}

# Review → Approved (approve)
POST /api/workflows/{workflow_id}/transitions
{
  "from_step_id": "{review_step_id}",
  "to_step_id": "{approved_step_id}",
  "action_name": "approve",
  "condition_type": "not_assigned_user"
}

# Review → Draft (reject)
POST /api/workflows/{workflow_id}/transitions
{
  "from_step_id": "{review_step_id}",
  "to_step_id": "{draft_step_id}",
  "action_name": "reject",
  "condition_type": "not_assigned_user"
}
```

### Step 4: Use the Workflow
```bash
# Create task
POST /api/tasks
{
  "workflow_id": "{workflow_id}",
  "title": "Complete docs",
  "assigned_to": "user123",
  "created_by": "manager"
}
# Returns: instance_id, starts at Draft step

# User submits for review
POST /api/tasks/{instance_id}/execute
{
  "action_name": "submit",
  "user_id": "user123"
}
# Task moves to Review step

# Manager approves
POST /api/tasks/{instance_id}/execute
{
  "action_name": "approve",
  "user_id": "manager"
}
# Task moves to Approved step (end)
```

---

## Migration

Run migrations to create all tables:
```bash
migrate -path ./migrations -database "your_connection_string" up
```

This will create:
- `workflows`
- `workflow_steps`
- `workflow_transitions`
- `workflow_instances`
- `workflow_history`

---

## Benefits

✅ **Fully Configurable** - Create any workflow without code changes  
✅ **Role-Based** - Control who can perform actions  
✅ **Audit Trail** - Complete history of all transitions  
✅ **Flexible Conditions** - Support for complex business rules  
✅ **Parallel Workflows** - Run multiple workflows simultaneously  
✅ **Backward Compatible** - Old hardcoded workflow still works  

---

## Next Steps

1. Run migrations to create tables
2. Create your first workflow via API
3. Define steps and transitions
4. Start running tasks through the workflow
5. Build frontend UI for visual workflow designer
