CREATE TABLE IF NOT EXISTS table_city
(
    id   serial  PRIMARY KEY,
    name varchar UNIQUE NOT NULL
);