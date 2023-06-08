CREATE TABLE IF NOT EXISTS table_interest
(
    id   serial  PRIMARY KEY,
    name varchar UNIQUE NOT NULL
);