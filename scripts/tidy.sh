#!/bin/bash

set -e

function process_dir() {
    local dir=$1

    if [ -f "$dir/go.mod" ]; then 
        echo "--- Tidying $dir ---"
        cd "$dir"
        go mod tidy
        cd -
    fi

    for subdir in "$dir"/*; do
        if [ -d "$subdir" ]; then
            process_dir "$subdir"
        fi
    done
}

process_dir "."

cd "$(git rev-parse --show-toplevel)"
go work sync
