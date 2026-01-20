-- Simple migration to change user ID from INT to UUID
-- Prerequisites: users, todos, and shared_tasks tables must be empty

-- Step 1: Drop and recreate users table with UUID id
DROP TABLE IF EXISTS shared_tasks;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;

-- Recreate users table with UUID id
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    session_token VARCHAR(255),
    csrf_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_session_token ON users(session_token);

-- Recreate todos table with UUID user_id
CREATE TABLE todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_name VARCHAR(100) NOT NULL,
    task_description VARCHAR(100) NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_todos_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_completed ON todos(completed);

-- Recreate shared_tasks table with UUID user_id columns
CREATE TABLE shared_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    shared_with_id UUID NOT NULL,
    todo_id UUID NOT NULL,
    comment VARCHAR(255) NOT NULL,
    CONSTRAINT fk_sharedtasks_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_sharedtasks_shared_with FOREIGN KEY (shared_with_id) REFERENCES users(id) ON DELETE NO ACTION,
    CONSTRAINT fk_sharedtasks_todo FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE NO ACTION
);

CREATE INDEX idx_sharedtasks_owner ON shared_tasks(owner_id);
CREATE INDEX idx_sharedtasks_shared_with ON shared_tasks(shared_with_id);
CREATE INDEX idx_sharedtasks_todo ON shared_tasks(todo_id);
