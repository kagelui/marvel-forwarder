CREATE TABLE "public"."server_synchronizations"
(
    last_sync_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_syncing   BOOL
);

CREATE TABLE "public"."characters"
(
    id          SERIAL PRIMARY KEY,
    external_id INTEGER NOT NULL,
    name        TEXT    NOT NULL,
    description TEXT    NOT NULL
);
