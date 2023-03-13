#!/bin/sh

# exit immediately if there is an error
set -e
echo "running db migration...."
/app/migrate -path /app/migration -database "$DSN" -verbose up

echo "migration ran successfully"

#remove migration files from the container
rm -rf /app/migrate
rm -rf /app/migration

echo "starting the app 🚀🚀"
# exec "$@"    