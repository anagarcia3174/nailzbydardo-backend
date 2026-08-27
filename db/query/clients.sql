-- name: CreateClient :one
INSERT INTO clients (client_name, contact_method, notes, birthday) VALUES ($1, $2, $3, $4) RETURNING id, client_name, contact_method, notes, birthday, created_at, deleted_at;

-- name: GetClient :one
SELECT id, client_name, contact_method, notes, birthday, created_at, deleted_at FROM clients WHERE id = $1;

-- name: ListClients :many
SELECT id, client_name, contact_method FROM clients WHERE deleted_at IS NULL;

-- name: UpdateClient :one
UPDATE clients
SET client_name = $2, contact_method = $3, notes = $4, birthday = $5
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, client_name, contact_method, notes, birthday, created_at, deleted_at;

-- name: SoftDeleteClient :exec
UPDATE clients SET deleted_at = now() WHERE id = $1;

-- name: ListClientAppointments :many
SELECT id, appt_date, appt_status, late_fee, payment_method, notes, loyalty_reward, tip, created_at
FROM appointments
WHERE client_id = $1
ORDER BY appt_date DESC;

-- name: ListCompleteClientAppointments :many
SELECT id, appt_date, tip
FROM appointments
WHERE client_id = $1 AND appt_status = 'complete'
ORDER BY appt_date DESC;
