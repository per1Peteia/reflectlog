-- +goose Up
CREATE TABLE clips (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    clip_text text NOT NULL,
    clip_brief text NOT NULL
);

-- +goose Down
DROP TABLE clips;

