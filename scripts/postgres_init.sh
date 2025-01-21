#!/bin/bash

echo 'Spinning up Postgres DB...'

PASSWORD=postgres
USER=postgres
DB=postgres
HOST=127.0.0.1
PORT=5432

docker run -e POSTGRES_PASSWORD=$PASSWORD -e POSTGRES_USER=$USER -e POSTGRES_DB=$DB -p $HOST:$PORT:5432 -d postgres
if [ $? -ne 0 ]; then
    echo 'failed to create local db'
else
    echo "Successfully created local db! User=$USER Password=$PASSWORD DB=$DB Host=$HOST Port=$PORT"
fi
