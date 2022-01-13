CREATE USER 'api'@'localhost';

CREATE DATABASE CandleDB;

GRANT ALL PRIVILEGES ON CandleDB . * TO 'api'@'localhost';

USE CandleDB;

CREATE TABLE Users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(20) NOT NULL,
    first_name VARCHAR(20) NOT NULL,
    last_name VARCHAR(20) NOT NULL,
    institution VARCHAR(50) NOT NULL,
    pfp_url VARCHAR(2048)
);

CREATE TABLE Orgs (
    id  INT AUTO_INCREMENT PRIMARY KEY,
    org_name VARCHAR(20) NOT NULL,
    institution VARCHAR(50) NOT NULL,
    org_pic_url VARCHAR(2048)
);

CREATE TABLE Members (
    user_id INT NOT NULL,
    org_id INT NOT NULL,
    PRIMARY KEY (user_id, org_id),
    member_role VARCHAR(20) NOT NULL,
    title VARCHAR(20) NOT NULL
);


