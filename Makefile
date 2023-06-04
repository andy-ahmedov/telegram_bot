all: run

run: server_run

server_run:
	tilix -e 'bash -c "nats-streaming-server -cid prod | make publish_run | make subscribe_run; exec bash"'

publish_run:
	tilix -e 'bash -c "go run publish/publish.go; exec bash"'

subscribe_run:
	tilix -e 'bash -c "go run main_sub.go; exec bash"'

kill:
	pgrep nats-streaming- | xargs kill -KILL
	pgrep publish | xargs kill -KILL
	pgrep main_sub | xargs kill -KILL

postgre:
	psql -c "CREATE DATABASE telegram_bot" 	
	psql -d telegram_bot -c "\i /home/andy/tester/persons.sql"
	psql -d telegram_bot -c "\d"
	psql -d telegram_bot -c "\d persons"

del_db:
	psql -d telegram_bot -c "DROP TABLE persons"
	psql -c "DROP DATABASE telegram_bot"

del_table:
	psql -d telegram_bot -c "DROP TABLE persons"

table:
	psql -d telegram_bot -c "\i /home/andy/tester/persons.sql"

show_table:
	psql -d telegram_bot -c "SELECT * FROM persons"

reset_table:
	psql -d telegram_bot -c "DROP TABLE persons"
	psql -d telegram_bot -c "\i /home/andy/tester/persons.sql"
