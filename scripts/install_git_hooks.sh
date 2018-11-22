#!/bin/sh

PROJECT_ADDRESS="github.com/7joe7/csvstreamtest"
SCRIPTS_FOLDER="scripts"

# install pre-commit hook
cat <<EOT > $GOPATH/src/$PROJECT_ADDRESS/.git/hooks/pre-commit
#!/bin/sh
$GOPATH/src/$PROJECT_ADDRESS/$SCRIPTS_FOLDER/pre-commit.sh
EOT

# install pre-push hook
cat <<EOT > $GOPATH/src/$PROJECT_ADDRESS/.git/hooks/pre-push
#!/bin/sh
cd $GOPATH/src/$PROJECT_ADDRESS
$GOPATH/src/$PROJECT_ADDRESS/$SCRIPTS_FOLDER/pre-push.sh
EOT

chmod 755 $GOPATH/src/$PROJECT_ADDRESS/.git/hooks/pre-commit
chmod 755 $GOPATH/src/$PROJECT_ADDRESS/.git/hooks/pre-push

echo " ----- pre-commit and pre-push git hooks installed -----"
