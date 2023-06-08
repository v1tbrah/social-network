CREATE TABLE IF NOT EXISTS table_hashtag
(
    id   serial  PRIMARY KEY,
    name varchar UNIQUE NOT NULL
);