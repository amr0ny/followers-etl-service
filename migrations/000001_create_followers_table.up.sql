CREATE TABLE IF NOT EXISTS followers (
     email VARCHAR(255) PRIMARY KEY,
     full_name VARCHAR(255) NOT NULL,
     modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
