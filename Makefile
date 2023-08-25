all: run

run: 
	go run cmd/bot/main.go

build_up:
	sudo docker-compose up --build bot


docker_clean:
	docker stop postgres_boston boston
	docker rm boston postgres_boston
	docker rmi telegram_bot-bot telegram_bot-db

fclean: docker_clean
		sudo rm -rf .database

create_dbs:
	psql -c "CREATE DATABASE telegram_bot" 	
	psql -d telegram_bot -c "\i /home/andy/github.com/andy-ahmedov/telegram_bot/clientRepository.sql"
	psql -d telegram_bot -c "\i /home/andy/github.com/andy-ahmedov/telegram_bot/userRepository.sql"
	psql -d telegram_bot -c "\d"
	psql -d telegram_bot -c "\d client_repository"

del_db:
	psql -d telegram_bot -c "DROP TABLE client_repository"
	psql -d telegram_bot -c "DROP TABLE userRepository"
	psql -c "DROP DATABASE telegram_bot"

del_client_table:
	psql -d telegram_bot -c "DROP TABLE client_repository"

del_boltdb:
	rm -rf bot.db

del_user_table:
	psql -d telegram_bot -c "DROP TABLE user_repository"

client_table:
	psql -d telegram_bot -c "\i /home/andy/github.com/andy-ahmedov/telegram_bot/clientRepository.sql"

user_table:
	psql -d telegram_bot -c "\i /home/andy/github.com/andy-ahmedov/telegram_bot/userRepository.sql"

show_client_table:
	psql -d telegram_bot -c "SELECT * FROM client_repository"

show_user_table:
	psql -d telegram_bot -c "SELECT * FROM user_repository"

reset_client_table:
	psql -d telegram_bot -c "DROP TABLE client_repository"
	psql -d telegram_bot -c "\i /home/andy/github.com/andy-ahmedov/telegram_bot/clientRepository.sql"

reset_user_table:
	psql -d telegram_bot -c "DROP TABLE user_repository"
	psql -d telegram_bot -c "\i /home/andy/github.com/andy-ahmedov/telegram_bot/userRepository.sql"

reset:
	psql -d telegram_bot -c "DROP TABLE user_repository"
	psql -d telegram_bot -c "\i /home/andy/github.com/andy-ahmedov/telegram_bot/userRepository.sql"
	rm -rf bot.db

rr:
	psql -d telegram_bot -c "DROP TABLE user_repository"
	psql -d telegram_bot -c "\i /home/andy/github.com/andy-ahmedov/telegram_bot/userRepository.sql"
	rm -rf bot.db
	go run cmd/bot/main.go

edit chatid:
	psql -d telegram_bot -c "UPDATE user_repository SET chatid = 346653896 Where chatid = 437499454"

export_client_table:
	psql -d telegram_bot -c "\copy client_repository to 'clients' csv;"
