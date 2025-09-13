#!/bin/bash

# local dev script used to easily test Goose DB changes to a locally hosted db. No prod data here.

echo "Running goose up migrations to local db!"

# SCHEMA=$1
# if [ -z "$SCHEMA" ]; then
#     echo "missing required schema arg"
#     exit 1
# fi

GOOSE_CMD=$1
if [ -z "$GOOSE_CMD" ]; then
    echo "missing required goose command"
    exit 1
fi

# goose postgres "user=postgres password=postgres host=127.0.0.1 port=5432 dbname=postgres sslmode=disable searchpath=$SCHEMA" $GOOSE_CMD
goose postgres "user=postgres password=postgres host=127.0.0.1 port=5432 dbname=postgres sslmode=disable" $GOOSE_CMD