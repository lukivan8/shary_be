#!/bin/bash

# Database setup script for Shary BE
# This script helps set up the PostgreSQL database for the application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up database for Shary BE...${NC}"

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}PostgreSQL is not installed. Please install PostgreSQL first.${NC}"
    echo "Visit: https://www.postgresql.org/download/"
    exit 1
fi

# Check if golang-migrate is installed
MIGRATE_CMD="migrate"
if ! command -v migrate &> /dev/null; then
    echo -e "${YELLOW}golang-migrate is not installed. Installing...${NC}"
    if command -v go &> /dev/null; then
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        
        # Get Go bin directory and add to PATH
        GOPATH=$(go env GOPATH)
        GO_BIN="$GOPATH/bin"
        
        # Check if migrate was installed
        if [ -f "$GO_BIN/migrate" ]; then
            echo -e "${GREEN}migrate installed successfully at $GO_BIN/migrate${NC}"
            MIGRATE_CMD="$GO_BIN/migrate"
        else
            echo -e "${RED}Failed to install migrate. Please install manually.${NC}"
            echo "Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
            exit 1
        fi
    else
        echo -e "${RED}Go is not installed. Please install Go first, or install golang-migrate manually.${NC}"
        echo "Visit: https://golang.org/doc/install"
        exit 1
    fi
fi

# Database configuration
DB_NAME="shary_be"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"

echo -e "${YELLOW}Database configuration:${NC}"
echo "Database name: $DB_NAME"
echo "Database user: $DB_USER"
echo "Database host: $DB_HOST"
echo "Database port: $DB_PORT"
echo ""

# Check if database exists
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo -e "${YELLOW}Database '$DB_NAME' already exists.${NC}"
    read -p "Do you want to drop and recreate it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Dropping database '$DB_NAME'...${NC}"
        dropdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    else
        echo -e "${GREEN}Using existing database.${NC}"
    fi
fi

# Create database if it doesn't exist
if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo -e "${YELLOW}Creating database '$DB_NAME'...${NC}"
    createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
fi

# Run migrations using golang-migrate
echo -e "${YELLOW}Running database migrations...${NC}"
DATABASE_URL="postgres://$DB_USER:your_password@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Note: Replace 'your_password' with actual password or use environment variable
if [ -n "$DB_PASSWORD" ]; then
    DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
fi

echo -e "${YELLOW}Using migrate command: $MIGRATE_CMD${NC}"
$MIGRATE_CMD -path migrations -database "$DATABASE_URL" up

echo -e "${GREEN}Database setup completed successfully!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Set your database URL:"
echo "   export DATABASE_URL=\"postgres://$DB_USER:your_password@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable\""
echo ""
echo "2. Run the application:"
echo "   go run main.go"
echo ""
echo "3. Test the API:"
echo "   curl http://localhost:4000/health"
echo "   curl http://localhost:4000/api/items" 