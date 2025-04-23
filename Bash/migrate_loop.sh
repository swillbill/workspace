#!/bin/bash
set -e

project=${1:? "Missing project name"}
step_size=${2:-1000}

repos_list_file="repo_list.txt"

if [[ ! -f "$repos_list_file" ]]; then
    echo "Error: $repos_list_file not found."
    exit 1
fi

echo "Processing project: '$project' from $repos_list_file"

repos=$(awk -v proj="$project" '
    $0 ~ "^"proj"=\\(" {flag=1; next}
    flag && /^\)/ {flag=0}
    flag {print $1}
' "$repos_list_file")

if [[ -z "$repos" ]]; then
    echo "No repositories found for project: $project"
    exit 1
fi

for repo in $repos; do
    echo "--------------------------------------"
    echo "Migrating repository: $repo"
    echo "--------------------------------------"
    chmod +x ./migrate.sh
    ./migrate.sh "$repo" "$step_size" "$project"
done

echo "All repositories for project '$project' have been processed."
