CREATE DATABASE SkinSavvy;
USE SkinSavvy;

SELECT * FROM users;

CREATE TABLE users (
	id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    fullname VARCHAR(255),
    age INT,
    password VARCHAR(255)
)
