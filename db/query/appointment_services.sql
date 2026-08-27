-- name: CreateAppointmentService :one
INSERT INTO appointment_services (appointment_id, service_name, service_price, design_price) VALUES ($1, $2, $3, $4) RETURNING id, appointment_id, service_name, service_price, design_price, created_at;

-- name: ListAppointmentServicesByAppointment :many
SELECT id, service_name, service_price, design_price
FROM appointment_services
WHERE appointment_id = $1;

-- name: DeleteAppointmentService :exec
DELETE FROM appointment_services WHERE id = $1;