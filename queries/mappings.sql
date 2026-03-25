-- name: AddServiceToSection :one
INSERT INTO service_section_mappings (service_id, section_id, sort_order)
VALUES ($1, $2, $3)
ON CONFLICT (service_id, section_id) DO NOTHING
RETURNING *;

-- name: RemoveServiceFromSection :exec
DELETE FROM service_section_mappings
WHERE service_id = $1 AND section_id = $2;

-- name: ListMappingsBySection :many
SELECT * FROM service_section_mappings
WHERE section_id = $1
ORDER BY sort_order ASC;

-- name: UpdateMappingSortOrder :exec
UPDATE service_section_mappings
SET sort_order = $3
WHERE service_id = $1 AND section_id = $2;

-- name: DeleteMappingsByService :exec
DELETE FROM service_section_mappings WHERE service_id = $1;

-- name: ListSectionIDsByService :many
SELECT section_id FROM service_section_mappings WHERE service_id = $1;
