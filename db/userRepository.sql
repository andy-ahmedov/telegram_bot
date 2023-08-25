CREATE TABLE user_repository (
	id SERIAL PRIMARY KEY,
	user_name  VARCHAR(40),
	client_name VARCHAR(50),
	phone_number VARCHAR(20),
	chatID INT NOT NULL,
	bonus INT,
	code VARCHAR(4),
	authorized VARCHAR(3)
);