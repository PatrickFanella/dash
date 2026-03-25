-- name: ListServices :many
SELECT * FROM services ORDER BY sort_order ASC;

-- name: GetService :one
SELECT * FROM services WHERE id = $1;

-- name: CreateService :one
INSERT INTO services (title, url, description, icon, status_check, status_check_url, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateService :one
UPDATE services
SET title = $2, url = $3, description = $4, icon = $5, status_check = $6,
    status_check_url = $7, sort_order = $8, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteService :exec
DELETE FROM services WHERE id = $1;

-- name: ListServicesBySection :many
SELECT s.*
FROM services s
JOIN service_section_mappings m ON m.service_id = s.id
WHERE m.section_id = $1
ORDER BY m.sort_order ASC;
