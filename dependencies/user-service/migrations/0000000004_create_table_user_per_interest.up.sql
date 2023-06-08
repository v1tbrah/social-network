CREATE TABLE IF NOT EXISTS table_user_per_interest
(
    user_id     bigint  NOT NULL,
    interest_id bigint NOT NULL,
    PRIMARY KEY (user_id, interest_id)
);
