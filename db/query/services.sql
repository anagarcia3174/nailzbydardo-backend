-- name: CreateService :one
INSERT INTO services (service_name, price) VALUES ($1, $2) RETURNING id, service_name, price, created_at;

-- name: ListServices :many
SELECT id, service_name, price, created_at FROM services;

-- name: UpdateService :one
UPDATE services
SET service_name = $2, price = $3
WHERE id = $1
RETURNING id, service_name, price, created_at;

-- name: DeleteService :exec
DELETE FROM services WHERE id = $1;