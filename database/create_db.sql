CREATE DATABASE IF NOT EXISTS peru_db;

USE peru_db;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);

-- Insert users, ignoring any errors such as duplicate key errors (i.e. if the user already exists)
INSERT IGNORE INTO users (username, password, email) VALUES ('test1', 'password1', 'test1@example.com');
INSERT IGNORE INTO users (username, password, email) VALUES ('test2', 'password2', 'test2@example.com');
INSERT IGNORE INTO users (username, password, email) VALUES ('test3', 'password3', 'test3@example.com');


CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT,
    from_user VARCHAR(255) NOT NULL,
    to_user VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    PRIMARY KEY (id)
);