#!/usr/bin/env bash

set -xe

go get "github.com/julienschmidt/httprouter"
go get "github.com/gorilla/mux"


env GOOS=linux go build -v -ldflags '-d -s -w' -a  -tags netgo -installsuffix netgo -o bin/mortgage-calculator mortgage-calculator.go