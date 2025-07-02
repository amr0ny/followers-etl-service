CREATE SCHEMA IF NOT EXISTS etl_info_schema;

CREATE TABLE IF NOT EXISTS etl_info_schema.table_status (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(64) UNIQUE,
    last_activity TIMESTAMP DEFAULT '-infinity'
)
