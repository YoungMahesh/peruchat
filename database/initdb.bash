#!/bin/bash

# Check if MYSQL_ROOT_PASSWORD is not set in the .env file
if ! grep -q "MYSQL_ROOT_PASSWORD" .env; then
    # create a random password for mysql database
    echo "MYSQL_ROOT_PASSWORD=$(openssl rand -hex 10)" >> .env
fi

# start mysql-server
docker compose -f docker-compose.yml up -d

# Load MYSQL_ROOT_PASSWORD from .env file
export $(grep -v '^#' ./.env | xargs)

# Connect to MySQL and execute the .sql file
docker exec -i peruchat_mysql mysql -u root -p$MYSQL_ROOT_PASSWORD < ./create_db.sql