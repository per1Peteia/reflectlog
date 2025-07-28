-- name: createClip :one
INSERT INTO clips (id, created_at, updated_at, clip_text, clip_brief, clip_title)
    VALUES (gen_random_uuid (), now(), now(), $1, $2, $3)
RETURNING
    *;

-- name: getClips :many
SELECT
    *
FROM
    clips
ORDER BY
    created_at ASC;

