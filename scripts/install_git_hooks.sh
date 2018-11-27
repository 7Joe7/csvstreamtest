#!/bin/sh

PROJECT_ADDRESS="github.com/7joe7/csvstreamtest"
SCRIPTS_FOLDER="scripts"

# install pre-commit hook
cat <<EOT > $GOPATH/src/$PROJECT_ADDRESS/.git/hooks/pre-commit
#!/bin/sh
$GOPATH/src/$PROJECT_ADDRESS/$SCRIPTS_FOLDER/pre-commit.sh
EOT

chmod 755 $GOPATH/src/$PROJECT_ADDRESS/.git/hooks/pre-commit

echo " ----- pre-commit git hook installed -----"
