#!/bin/bash
cd /data/migrations

# Run Goose migrations
goose up

cd /root/

# Start your main application
exec ./main
