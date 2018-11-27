#!/bin/sh

# Get the working directory (we don't know if the git push command happened in /dcim)
# So we cut the current path until sources.etixlabs.com

PROJECT_ADDRESS="github.com/7joe7/csvstreamtest"
SCRIPTS_FOLDER="scripts"

echo "----- checking formatting of modified and staged go source code -----\n"

pushd $GOPATH/src/$PROJECT_ADDRESS > /dev/null

if [[ $(git diff --name-only --cached) ]]; then DIFFERENCES=$(goimports -l $(git diff --name-only --cached | grep -v '/vendor/' | grep -v '.pb.go' | grep .go)); fi

if [[ $DIFFERENCES ]]; then
    echo "Formatting issues found in these files:\n";
    echo $DIFFERENCES;
    exit 1;
fi

echo "----- running gosec -----\n"
if ! gosec ./...; then
    echo "Security issues found\n"
    exit 1;
fi

echo "----- running go vet -----\n"
if ! go vet ./...; then
    echo "Code quality is compromised\n"
    exit 1;
fi

echo "----- running unit tests -----\n"
if ! go test ./...; then
    echo "Unit tests are failing\n"
    exit 1;
fi

echo "----- no issues found -----\n"

popd > /dev/null