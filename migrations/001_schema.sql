CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    hash       TEXT NOT NULL,
    created    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE videos (
    id           BIGSERIAL PRIMARY KEY,
    public_id    UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    title        TEXT NOT NULL,
    author_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    views BIGINT NOT NULL DEFAULT 0,

    s3_key       TEXT NOT NULL UNIQUE, -- "videos/{key}/original.mp4"
    size_bytes   BIGINT,
    duration_ms  INTEGER
);
