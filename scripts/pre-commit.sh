#!/bin/sh

# Get the working directory (we don't know if the git push command happened in /dcim)
# So we cut the current path until sources.etixlabs.com

PROJECT_ADDRESS="github.com/7joe7/csvstreamtest"
SCRIPTS_FOLDER="scripts"

echo "----- checking formatting of modified and staged go source code -----"

pushd $GOPATH/src/$PROJECT_ADDRESS > /dev/null

if [[ $(git diff --name-only --cached) ]]; then DIFFERENCES=$(goimports -l $(git diff --name-only --cached | grep -v '/vendor/' | grep -v '.pb.go' | grep .go)); fi

if [[ $DIFFERENCES ]]; then
    echo "\nFormatting issues found in these files:\n";
    echo $DIFFERENCES;
    exit 1;
fi

if ! gosec ./...; then
    echo "\nSecurity issues found"
    exit 1;
fi

echo "\n----- no issues found -----"

popd > /dev/null