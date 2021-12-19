#!/usr/bin/env bash
set -xe

# get all of the dependencies needed
go get

# create the application binary that EB uses
go build -o bin/application main.go