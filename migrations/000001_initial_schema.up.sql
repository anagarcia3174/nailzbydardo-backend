CREATE TYPE payment_method AS ENUM ('cash', 'zelle', 'cash_app', 'other');
CREATE TYPE appointment_status AS ENUM ('booked', 'complete', 'no_show', 'cancelled');
CREATE TYPE discount_type AS ENUM ('amount', 'percent');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT sessions_user_id_users_id_fk
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    client_name VARCHAR(255) NOT NULL,
    contact_method VARCHAR(255),
    notes TEXT,
    birthday DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ

);
CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    service_name VARCHAR(255) NOT NULL,
    price NUMERIC(10,  2) NOT NULL CHECK (price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()

);
CREATE TABLE IF NOT EXISTS expenses (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    expense_name VARCHAR(255) NOT NULL,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    date_purchased DATE NOT NULL,
    receipt_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    client_id UUID NOT NULL,
    appt_date TIMESTAMPTZ NOT NULL,
    appt_status appointment_status DEFAULT 'booked' NOT NULL,
    late_fee NUMERIC(10, 2) CHECK (late_fee >= 0),
    payment_method payment_method,
    notes TEXT,
    receipt_url TEXT,
    loyalty_reward BOOLEAN DEFAULT false,
    tip NUMERIC(10, 2) CHECK (tip >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT appointments_client_id_clients_id_fk
        FOREIGN KEY (client_id)
        REFERENCES clients(id)
        ON DELETE RESTRICT
);
CREATE INDEX idx_appointments_client_id ON appointments(client_id);
CREATE INDEX idx_appointments_appt_date ON appointments(appt_date);
CREATE TABLE IF NOT EXISTS appointment_services (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    appointment_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    service_price NUMERIC(10,2) NOT NULL CHECK (service_price >= 0),
    design_price NUMERIC(10,2) DEFAULT 0 CHECK (design_price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT appointment_services_appointment_id_appointments_id_fk
        FOREIGN KEY (appointment_id)
        REFERENCES appointments(id)
        ON DELETE CASCADE
);
CREATE INDEX idx_appointment_services_appointment_id ON appointment_services(appointment_id);
CREATE TABLE IF NOT EXISTS appointment_discounts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    appointment_id UUID NOT NULL,
    discount_name TEXT NOT NULL,
    discount_type discount_type NOT NULL,
    discount_value NUMERIC(10, 2) NOT NULL CHECK (discount_value >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT appointment_discounts_appointment_id_appointments_id_fk
        FOREIGN KEY (appointment_id)
        REFERENCES appointments(id)
        ON DELETE CASCADE
);
CREATE INDEX idx_appointment_discounts_appointment_id ON appointment_discounts(appointment_id);
