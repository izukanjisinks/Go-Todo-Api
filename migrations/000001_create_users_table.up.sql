CREATE TABLE users (
                       id INT IDENTITY(1,1) PRIMARY KEY,
                       username NVARCHAR(255) NOT NULL UNIQUE,
                       hashed_password NVARCHAR(255) NOT NULL,
                       session_token NVARCHAR(255),
                       csrf_token NVARCHAR(255),
                       created_at DATETIME2 DEFAULT GETDATE(),
                       updated_at DATETIME2 DEFAULT GETDATE()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_session_token ON users(session_token);