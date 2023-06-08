CREATE TABLE IF NOT EXISTS table_user
(
    id      serial  PRIMARY KEY,
    name    varchar NOT NULL,
    surname varchar NOT NULL,
    city_id bigint,
        CONSTRAINT fk_city
        FOREIGN KEY (city_id)
        REFERENCES table_city(id)
);
