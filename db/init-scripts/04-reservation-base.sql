CREATE TABLE IF NOT EXISTS reservations (    
    id BIGSERIAL PRIMARY KEY,
    property_id BIGINT NOT NULL,    
    user_id BIGINT NOT NULL,   
    owner_id BIGINT NOT NULL,    
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,   
    status VARCHAR(50) DEFAULT 'RESERVATION_STATUS_PENDING', -- Ex: PENDING, CONFIRMED, CANCELLED
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_reservation_property 
        FOREIGN KEY (property_id) 
        REFERENCES properties (id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_reservation_user 
        FOREIGN KEY (user_id) 
        REFERENCES users (id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_reservation_owner
        FOREIGN KEY (owner_id) 
        REFERENCES users (id) 
        ON DELETE CASCADE,

    CONSTRAINT chk_valid_dates 
        CHECK (check_out_date > check_in_date)
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_reservations_updated_at
BEFORE UPDATE ON reservations
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();