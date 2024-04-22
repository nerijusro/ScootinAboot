ALTER TABLE scooters
ADD CONSTRAINT fk_occupied_by
FOREIGN KEY (occupied_by)
REFERENCES users(id);