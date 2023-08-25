CREATE TABLE client_repository (
	id SERIAL PRIMARY KEY,
	client_name VARCHAR(60),
	phone_number VARCHAR(20) NOT NULL,
	bonus INT NOT NULL
);