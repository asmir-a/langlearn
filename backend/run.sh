#!/bin/zsh
nodemon -e go --watch "./**/**.go" --signal SIGTERM --exec go run .
