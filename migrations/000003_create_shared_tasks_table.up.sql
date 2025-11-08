CREATE TABLE shared_tasks (
            id NVARCHAR(36) PRIMARY KEY,
            owner_id INT NOT NULL,
            shared_with_id INT NOT NULL,
            todo_id NVARCHAR(36) NOT NULL,
            comment NVARCHAR(255) NOT NULL,
            CONSTRAINT fk_sharedtasks_owner FOREIGN KEY (owner_id)
                REFERENCES users(id) ON DELETE CASCADE,
            CONSTRAINT fk_sharedtasks_shared_with FOREIGN KEY (shared_with_id)
                REFERENCES users(id) ON DELETE NO ACTION,
            CONSTRAINT fk_sharedtasks_todo FOREIGN KEY (todo_id)
                REFERENCES todos(id) ON DELETE NO ACTION 
);

CREATE INDEX idx_sharedtasks_owner ON shared_tasks(owner_id);
CREATE INDEX idx_sharedtasks_shared_with ON shared_tasks(shared_with_id);
CREATE INDEX idx_sharedtasks_todo ON shared_tasks(todo_id);