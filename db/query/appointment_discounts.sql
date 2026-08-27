-- name: CreateAppointmentDiscount :one
INSERT INTO appointment_discounts (appointment_id, discount_name, discount_type, discount_value) VALUES ($1, $2, $3, $4) RETURNING id, appointment_id, discount_name, discount_type, discount_value, created_at;

-- name: ListAppointmentDiscountsByAppointment :many
SELECT id, discount_name, discount_type, discount_value
FROM appointment_discounts
WHERE appointment_id = $1;

-- name: DeleteAppointmentDiscount :exec
DELETE FROM appointment_discounts WHERE id = $1;