#!/usr/bin/env sh

set -e -v

dagger develop

for dir in merge-dirs rust; do
    cd "$dir"
    dagger develop
    cd -

    if [ -d "$dir/tests" ]; then
        cd "$dir/tests"
        dagger develop
        cd -
    fi
done
