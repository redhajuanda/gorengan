
-- +migrate Up
CREATE TABLE users (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    first_name VARCHAR(32),
    last_name VARCHAR(32),
    email VARCHAR(100),
    password VARCHAR(255),
    address TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO `users` (`id`, `first_name`, `last_name`, `email`, `password`, `address`, `created_at`, `updated_at`) VALUES
('c7a2df29-047c-4674-a553-0416d4325e6c',	'Super',	'Admin',	'super@admin.com',	'$2a$04$VdPk/HVxCz0ncH.QbPCRyOZCyp90ZAQjEfst3tCQS5pb5Riszl8c.',	'Jakarta',	'2020-08-09 11:30:25',	'2020-08-09 11:30:25');

-- +migrate Down
DROP TABLE users;