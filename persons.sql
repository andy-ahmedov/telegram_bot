CREATE TABLE persons (
	id SERIAL PRIMARY KEY,
	first_name VARCHAR(20) NOT NULL,
	last_name VARCHAR(25) NOT NULL,
	phone_number VARCHAR(20) NOT NULL,
	bonus INT NOT NULL
);