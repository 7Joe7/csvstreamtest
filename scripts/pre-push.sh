#!/bin/sh

PROJECT_ADDRESS="github.com/7joe7/csvstreamtest"
SCRIPTS_FOLDER="scripts"

echo "----- running unit tests -----"

pushd $GOPATH/src/$PROJECT_ADDRESS > /dev/null

for service_dir in $(git diff --name-only HEAD~1 | grep -Eo '(cloud|site)/([^/]+)' | sort | uniq) ; do
  cd $service_dir
    service_name=$(echo $service_dir | rev | cut -d / -f 1 | rev)
    echo "\n$service_name\n--------------------------------"

    go list # check that we are in a Golang service
    if [ $? -eq 0 ]; then
        echo "\nChecking imports...\n"
        if [ -d "./vendor/sources.etixlabs.com/dcim" ]; then
          if [ "$(find ./vendor/sources.etixlabs.com/dcim -maxdepth 1 -type d | wc -w)" -eq 2 ]; then
            if ! [ -d "./vendor/sources.etixlabs.com/dcim/libellula" ]; then
              echo "Too many directories in ./$service_name/vendor/sources.etixlabs.com/dcim"
              exit 1
            else
              echo "Imports ok"
            fi
          else
            echo "Too many directories in ./$service_name/vendor/sources.etixlabs.com/dcim"
            exit 1
          fi
        fi
        echo "\nRunning tests...\n"
        go test $(go list ./... | grep -v contracts)
        if [ $? -ne 0 ]; then
          echo "Tests failed for service:" $service_name
          exit 1
        fi
    fi

  cd $WORKING_DIR
done

popd > /dev/null
