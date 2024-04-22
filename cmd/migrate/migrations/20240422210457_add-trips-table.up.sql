CREATE TABLE IF NOT EXISTS trips (
  `id` BINARY(16) NOT NULL PRIMARY KEY,
  `user_id` BINARY(16) NOT NULL,
  `scooter_id` BINARY(16) NOT NULL,
  `is_finished` BOOLEAN NOT NULL DEFAULT FALSE,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  FOREIGN KEY (`scooter_id`) REFERENCES scooters(`id`)
);