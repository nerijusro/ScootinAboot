CREATE TABLE IF NOT EXISTS users (
  `id` BINARY(16) NOT NULL PRIMARY KEY,
  `full_name` VARCHAR(255) NOT NULL,
  `is_eligible_to_travel` BOOLEAN NOT NULL DEFAULT TRUE,
  `opt_lock_version` INT NOT NULL DEFAULT 0
);
