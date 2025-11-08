CREATE TABLE todos (
    id NVARCHAR(36) PRIMARY KEY,
    task_name NVARCHAR(100) NOT NULL,
    task_description NVARCHAR(100) NOT NULL,
    completed BIT DEFAULT 0,
    user_id INT NOT NULL,
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE(),
    CONSTRAINT fk_todos_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_completed ON todos(completed);
