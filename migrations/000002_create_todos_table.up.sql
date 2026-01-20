CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_name VARCHAR(100) NOT NULL,
    task_description VARCHAR(100) NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_todos_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_completed ON todos(completed);
