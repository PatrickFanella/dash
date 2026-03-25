-- name: ListSections :many
SELECT * FROM sections ORDER BY sort_order ASC;

-- name: GetSection :one
SELECT * FROM sections WHERE id = $1;

-- name: CreateSection :one
INSERT INTO sections (name, icon, cols, collapsed, sort_order, section_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateSection :one
UPDATE sections
SET name = $2, icon = $3, cols = $4, collapsed = $5, sort_order = $6, section_type = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSection :exec
DELETE FROM sections WHERE id = $1;
