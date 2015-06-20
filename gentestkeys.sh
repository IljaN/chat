#!/usr/bin/env bash
openssl genrsa -out keys/app.rsa 4096
openssl rsa -in keys/app.rsa -outform PEM -pubout -out keys/app.rsa.pub


