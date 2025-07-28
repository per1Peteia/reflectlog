-- +goose Up
ALTER TABLE clips
    ADD COLUMN clip_title text NOT NULL;

-- +goose Down
ALTER TABLE clips
    DROP COLUMN clip_title;

