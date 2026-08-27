-- name: CreateAppointment :one
INSERT INTO appointments (client_id, appt_date, notes) VALUES ($1, $2, $3) RETURNING id, client_id, appt_date, appt_status, late_fee, payment_method, notes, receipt_url, loyalty_reward, tip, created_at;

-- name: GetAppointment :one
SELECT id, client_id, appt_date, appt_status, late_fee, payment_method, notes, receipt_url, loyalty_reward, tip, created_at FROM appointments WHERE id = $1;

-- name: ListAppointments :many
SELECT id, client_id, appt_date, appt_status, late_fee, payment_method, notes, receipt_url, loyalty_reward, tip, created_at FROM appointments;

-- name: ListAppointmentsByDateRange :many
SELECT id, client_id, appt_date, appt_status, late_fee, payment_method, notes, receipt_url, loyalty_reward, tip, created_at
FROM appointments
WHERE appt_date >= $1 AND appt_date < $2
ORDER BY appt_date ASC;

-- name: ListCompleteAppointmentsForPeriod :many
SELECT id, appt_date, tip
FROM appointments
WHERE appt_status = 'complete'
  AND appt_date >= $1 AND appt_date < $2
ORDER BY appt_date ASC;

-- name: ListUpcomingAppointments :many
SELECT id, client_id, appt_date, appt_status, late_fee, payment_method, notes, receipt_url, loyalty_reward, tip, created_at
FROM appointments
WHERE appt_date > now() AND appt_status != 'cancelled'
ORDER BY appt_date ASC
LIMIT 10;

-- name: GetAppointmentCountForPeriod :one
SELECT COUNT(*) FROM appointments
WHERE appt_status = 'complete'
  AND appt_date >= $1 AND appt_date < $2;

-- name: UpdateAppointment :one
UPDATE appointments
SET appt_date = $2, appt_status = $3, late_fee = $4, payment_method = $5, notes = $6, receipt_url = $7, loyalty_reward = $8, tip = $9
WHERE id = $1
RETURNING id, client_id, appt_date, appt_status, late_fee, payment_method, notes, receipt_url, loyalty_reward, tip, created_at;

-- name: DeleteAppointment :exec
DELETE FROM appointments WHERE id = $1;