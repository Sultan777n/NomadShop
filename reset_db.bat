@echo off
echo [1/5] Dropping database...
psql -U postgres -c "DROP DATABASE IF EXISTS nomadshop;"

echo [2/5] Creating database...
psql -U postgres -c "CREATE DATABASE nomadshop;"

echo [3/5] Forcing migration to version 0 (reset)...
migrate -path migrations -database "postgres://postgres:asd12345@localhost:5432/nomadshop?sslmode=disable" force 0

echo [4/5] Running migrations...
migrate -path migrations -database "postgres://postgres:asd12345@localhost:5432/nomadshop?sslmode=disable" up

echo [5/5] Done!
pause
