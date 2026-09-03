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

-- name: ListUpcomingAppointmentsForDashboard :many
SELECT
    appointments.id,
    appointments.appt_date,
    clients.client_name
FROM appointments
JOIN clients
    ON appointments.client_id = clients.id
WHERE appointments.appt_date > now()
  AND appointments.appt_status != 'cancelled'
  AND clients.deleted_at IS NULL
ORDER BY appointments.appt_date ASC
LIMIT 5;

-- name: ListAppointmentsWithClient :many
SELECT
  appointments.id,
  appointments.appt_date,
  appointments.appt_status,
  clients.client_name
FROM appointments
JOIN clients
  ON appointments.client_id = clients.id
ORDER BY appointments.appt_date DESC;

-- name: GetAppointmentWithClient :one
SELECT
    appointments.id,
    appointments.client_id,
    appointments.appt_date,
    appointments.appt_status,
    appointments.late_fee,
    appointments.payment_method,
    appointments.notes,
    appointments.receipt_url,
    appointments.loyalty_reward,
    appointments.tip,
    appointments.created_at,

    clients.client_name,
    clients.contact_method,

    (
        SELECT COUNT(*)
        FROM appointments AS client_appointments
        WHERE client_appointments.client_id = appointments.client_id
        AND client_appointments.appt_status = 'complete'
        AND (
            client_appointments.appt_date < appointments.appt_date
            OR (
                client_appointments.appt_date = appointments.appt_date
                AND client_appointments.created_at < appointments.created_at
            )
        )
    ) + 1 AS appointment_rank

FROM appointments

JOIN clients
ON clients.id = appointments.client_id

WHERE appointments.id = $1;


-- name: ListAppointmentsForCalendar :many
SELECT
    appointments.id,
    appointments.appt_date,
    appointments.appt_status,
    appointments.notes,
    clients.client_name
FROM appointments
JOIN clients
    ON appointments.client_id = clients.id
WHERE appointments.appt_date >= $1
    AND appointments.appt_date < $2
    AND appointments.appt_status != 'cancelled'
    AND clients.deleted_at IS NULL
ORDER BY appointments.appt_date ASC;