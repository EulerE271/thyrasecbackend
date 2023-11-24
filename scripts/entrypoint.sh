#!/bin/sh

# Navigate to migrations directory
cd /data/migrations

# Run Goose migrations
goose up
echo "Migrations completed."

# Navigate to the root directory
cd /root/

# Start your main application
exec ./main
