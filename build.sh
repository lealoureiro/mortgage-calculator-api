#!/usr/bin/env bash

set -xe

echo "Getting project dependencies..."
dep ensure

echo "Compiling application..."
env GOOS=linux go build -v -ldflags '-d -s -w' -a  -tags netgo -installsuffix netgo -o bin/mortgage-calculator-api mortgage-calculator-api.go