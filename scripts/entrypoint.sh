#!/bin/sh
cd /data/migrations

# Run Goose migrations
goose up
echo "PATH is $PATH"

cd /root/

# Start your main application
exec ./main
