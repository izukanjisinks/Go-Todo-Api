CREATE TABLE todo_tasks (
    id NVARCHAR(36) PRIMARY KEY,
    title NVARCHAR(200) NOT NULL,
    description NVARCHAR(500),
    assigned_to NVARCHAR(100) NOT NULL,
    status NVARCHAR(20) NOT NULL DEFAULT 'Draft',
    reviewed_by NVARCHAR(100),
    approved_by NVARCHAR(100),
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE(),
    CONSTRAINT chk_status CHECK (status IN ('Draft', 'Review', 'Approved'))
);

CREATE INDEX idx_todo_tasks_assigned_to ON todo_tasks(assigned_to);
CREATE INDEX idx_todo_tasks_status ON todo_tasks(status);
CREATE INDEX idx_todo_tasks_reviewed_by ON todo_tasks(reviewed_by);
