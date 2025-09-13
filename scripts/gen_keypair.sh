#!/bin/bash

OVERWRITE=$1

if [[ -f "./private-key.pem"  && -f "./public-key.pem" && "${OVERWRITE}" !=  "overwrite" ]]; then
  echo "key pair already exists"
  exit 0
fi

echo "generating new local public and private keys..."
openssl genrsa -out private-key.pem 4096
openssl rsa -in private-key.pem -pubout -out public-key.pem
