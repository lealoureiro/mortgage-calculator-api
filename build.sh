#!/usr/bin/env bash

set -xe

go get "github.com/gorilla/mux"

env GOOS=linux go build -v -ldflags '-d -s -w' -a  -tags netgo -installsuffix netgo -o bin/mortgage-calculator-api mortgage-calculator-api.go