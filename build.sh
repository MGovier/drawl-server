#!/bin/zsh
rm -rf ./builds/*
GOOS=linux GOARCH=amd64 go build -o ./builds/