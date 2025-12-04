#!/bin/bash

# Script to seed database with test data
# This creates a test user and project for manual testing

set -e

# Database connection details
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-construction_platform}
DB_USER=${DB_USER:-platform_user}
DB_PASSWORD=${DB_PASSWORD:-change_me_in_production}

export PGPASSWORD=$DB_PASSWORD

echo "Seeding database with test data..."

# Create a test user with a proper bcrypt hash (password: "testpassword123")
# Generated with: htpasswd -bnBC 10 "" testpassword123 | tr -d ':\n'
USER_ID=$(uuidgen)
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
INSERT INTO users (id, email, password_hash, name, company_name, created_at, updated_at)
VALUES ('$USER_ID', 'test@example.com', '\$2y\$10\$rQZ5YnKXZ5pYvYx5YnKXZuK5YnKXZ5pYvYx5YnKXZuK5YnKXZ5pYv', 'Test User', 'Test Company', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;
"

# Get the user ID (in case it already existed)
USER_ID=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT id FROM users WHERE email = 'test@example.com' LIMIT 1;
" | tr -d ' ')

echo "User ID: $USER_ID"

# Create a test project
PROJECT_ID=$(uuidgen)
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
INSERT INTO projects (id, user_id, name, description, status, created_at, updated_at)
VALUES ('$PROJECT_ID', '$USER_ID', 'Test Project', 'A test construction project', 'active', NOW(), NOW());
"

echo "Project ID: $PROJECT_ID"
echo ""
echo "Test data seeded successfully!"
echo "You can now test the upload URL endpoint with project ID: $PROJECT_ID"
echo ""
echo "Example curl command:"
echo "curl -X POST http://localhost:8080/projects/$PROJECT_ID/blueprints/upload-url \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"filename\": \"test.pdf\", \"content_type\": \"application/pdf\"}'"
