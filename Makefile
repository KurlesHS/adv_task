create_db:
	export PGPASSWORD=${DB_PASS} && psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -c 'create database ${DB_NAME};'

drop_db:
	export PGPASSWORD=${DB_PASS} && psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -c 'drop database ${DB_NAME};'

migrate_up:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down -all

test_repo: create_db drop_db
