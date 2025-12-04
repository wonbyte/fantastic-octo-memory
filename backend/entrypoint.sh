#!/bin/sh

# Wait for database to be ready
echo "Waiting for database to be ready..."
until pg_isready -h ${DB_HOST:-postgres} -p ${DB_PORT:-5432} -U ${DB_USER:-platform_user}; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is ready!"

# Construct DATABASE_URL if not provided
if [ -z "$DATABASE_URL" ]; then
  DB_HOST=${DB_HOST:-postgres}
  DB_PORT=${DB_PORT:-5432}
  DB_NAME=${DB_NAME:-construction_platform}
  DB_USER=${DB_USER:-platform_user}
  DB_PASSWORD=${DB_PASSWORD:-change_me_in_production}
  export DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
fi

echo "Running migrations..."
echo "DATABASE_URL: ${DATABASE_URL}"

# Run migrations using golang-migrate
migrate -path /app/migrations -database "$DATABASE_URL" up

if [ $? -eq 0 ]; then
  echo "Migrations completed successfully!"
else
  echo "Migration failed!"
  exit 1
fi

# Start the application
echo "Starting application..."
exec /app/main
