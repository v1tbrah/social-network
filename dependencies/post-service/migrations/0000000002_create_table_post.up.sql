CREATE TABLE IF NOT EXISTS table_post
(
    id serial PRIMARY KEY,
    user_id bigint NOT NULL,
    description varchar,
    created_at timestamp NOT NULL
);
