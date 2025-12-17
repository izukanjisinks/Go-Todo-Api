-- Revert user ID back to INT

-- Drop tables
DROP TABLE IF EXISTS shared_tasks;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;

-- Recreate users table with INT IDENTITY id
CREATE TABLE users (
    id INT IDENTITY(1,1) PRIMARY KEY,
    username NVARCHAR(255) NOT NULL UNIQUE,
    email NVARCHAR(255) NOT NULL UNIQUE,
    password NVARCHAR(255) NOT NULL,
    is_admin BIT NOT NULL DEFAULT 0,
    session_token NVARCHAR(255),
    csrf_token NVARCHAR(255),
    created_at DATETIME2 DEFAULT GETDATE(),
    updated_at DATETIME2 DEFAULT GETDATE()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Recreate todos table with INT user_id
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

-- Recreate shared_tasks table with INT user_id columns
CREATE TABLE shared_tasks (
    id NVARCHAR(36) PRIMARY KEY,
    owner_id INT NOT NULL,
    shared_with_id INT NOT NULL,
    todo_id NVARCHAR(36) NOT NULL,
    comment NVARCHAR(255) NOT NULL,
    CONSTRAINT fk_sharedtasks_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_sharedtasks_shared_with FOREIGN KEY (shared_with_id) REFERENCES users(id) ON DELETE NO ACTION,
    CONSTRAINT fk_sharedtasks_todo FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE NO ACTION 
);

CREATE INDEX idx_sharedtasks_owner ON shared_tasks(owner_id);
CREATE INDEX idx_sharedtasks_shared_with ON shared_tasks(shared_with_id);
CREATE INDEX idx_sharedtasks_todo ON shared_tasks(todo_id);
