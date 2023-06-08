CREATE TABLE IF NOT EXISTS table_hashtag_per_post
(
    post_id bigint NOT NULL,
    hashtag_id bigint NOT NULL,
    PRIMARY KEY (post_id, hashtag_id),
    CONSTRAINT fk_post
    FOREIGN KEY (post_id)
    REFERENCES table_post (id),
    CONSTRAINT fk_hashtag
    FOREIGN KEY (hashtag_id)
    REFERENCES table_hashtag (id)
);
