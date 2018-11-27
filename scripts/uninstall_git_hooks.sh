#!/bin/sh

PROJECT_ADDRESS="github.com/7joe7/csvstreamtest"
SCRIPTS_FOLDER="scripts"

rm $GOPATH/src/$PROJECT_ADDRESS/.git/hooks/pre-commit
rm $GOPATH/src/$PROJECT_ADDRESS/.git/hooks/pre-push
