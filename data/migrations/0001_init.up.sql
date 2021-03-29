CREATE TABLE "public"."characters"
(
    id          SERIAL PRIMARY KEY,
    external_id INTEGER UNIQUE NOT NULL,
    name        TEXT    NOT NULL,
    description TEXT    NOT NULL
);
