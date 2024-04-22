CREATE TABLE IF NOT EXISTS scooters (
  `id` BINARY(16) NOT NULL PRIMARY KEY,
  `latitude` FLOAT NOT NULL,
  `longitude` FLOAT NOT NULL,
  `is_available` BOOLEAN NOT NULL DEFAULT FALSE,
  `opt_lock_version` INT NOT NULL DEFAULT 0
);