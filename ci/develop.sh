#!/usr/bin/env sh

set -e -v

dagger develop

dirs="merge-dirs rust go bun"

for dir in $dirs; do
    cd "$dir"
    dagger develop
    cd -

    if [ -d "$dir/tests" ]; then
        cd "$dir/tests"
        dagger develop
        cd -
    fi
done
