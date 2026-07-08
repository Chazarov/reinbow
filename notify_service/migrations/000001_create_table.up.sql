CREATE TABLE IF NOT EXISTS client (
    id SERIAL PRIMARY KEY,
    client_name VARCHAR(300),
    client_token VARCHAR(4096),
    client_login VARCHAR(300),
    client_password VARCHAR(300),
    email VARCHAR(400)
);

CREATE TABLE IF NOT EXISTS notify_type_message (
    id SERIAL PRIMARY KEY,
    notify_type VARCHAR(200) UNIQUE,
    notify_description VARCHAR(200),
    want_telegram BOOLEAN,
    want_email BOOLEAN,
    want_webhook TEXT[]
)