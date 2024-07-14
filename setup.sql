CREATE TABLE IF NOT EXISTS users
(
    id          SERIAL PRIMARY KEY,
    user_name   VARCHAR(50) UNIQUE NOT NULL,
    first_name  VARCHAR(255) NOT NULL,
    last_name   VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    user_status VARCHAR(1) NOT NULL,
    department  VARCHAR(255)
);