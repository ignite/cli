#!/usr/bin/env bash

set -euo pipefail

for modfile in $(find . -name go.mod); do
 echo "Updating $modfile"
 DIR=$(dirname $modfile)
 if [[ $DIR == *"testdata"* ]]; then
   echo "Skipping testdata directory"
   continue
 fi
 (cd $DIR; go mod tidy)
done
