#!/bin/bash
set -e

LFS_THRESHOLD="50MB"

if [ -z "$1" ]; then
    echo "Usage: $0 <project_name>"
    exit 1
fi

PROJECT_NAME="$1"
TFS_BASE_URL="https://tfs.revenuepremier.com/tfs/RSI/${PROJECT_NAME}/_git"

echo "Using TFS project: ${PROJECT_NAME}"
echo "TFS_BASE_URL is set to: ${TFS_BASE_URL}"

echo "Adding GitHub to SSH known hosts..."
mkdir -p ~/.ssh
ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

REPO_LIST_FILE="repo_list.txt"
if [ ! -f "$REPO_LIST_FILE" ]; then
    echo "Error: Repository list file '$REPO_LIST_FILE' not found."
    exit 1
fi

echo "Reading repository list from file: $REPO_LIST_FILE"
REPO_NAMES=$(awk -v proj="$PROJECT_NAME" '
    BEGIN { flag=0 }
    $0 ~ "^"proj"=\\(" { flag=1; next }
    flag && /\)/ { flag=0; next }
    flag { print $0 }
' "$REPO_LIST_FILE")

if [[ -z "$REPO_NAMES" ]]; then
    echo "Error: No repositories found for project '$PROJECT_NAME' in file '$REPO_LIST_FILE'."
    exit 1
fi

echo "Found repositories for project $PROJECT_NAME:"
echo "$REPO_NAMES"

OUTPUT_FILE="$(pwd)/files_over_100MB.txt"
> "$OUTPUT_FILE"

for repo in $REPO_NAMES; do
    echo "Processing repository: $repo"

    TFS_REPO_URL="${TFS_BASE_URL}/$repo"
    REPO_PATH="$(pwd)/$repo.git"
    WORKING_PATH="$(pwd)/$repo-work"

    [ -d "$REPO_PATH" ] && rm -rf "$REPO_PATH"
    [ -d "$WORKING_PATH" ] && rm -rf "$WORKING_PATH"

    echo "Cloning repository from TFS..."
    mkdir -p "$REPO_PATH"
    expect << EOF
        set timeout -1
        spawn git clone --mirror "$TFS_REPO_URL" "$REPO_PATH"
        expect {
            "Username for" { 
                send "$TFS_USERNAME\r"
                exp_continue 
            }
            "Password for" { 
                send "$TFS_PASS\r"
                exp_continue 
            }
            eof { 
                wait
                exit 0 
            }
        }
EOF

    echo "Creating working repository..."
    git clone "$REPO_PATH" "$WORKING_PATH"
    cd "$WORKING_PATH"

    echo "Finding files over $LFS_THRESHOLD in repository $repo..."
    echo "Repository: $repo" >> "$OUTPUT_FILE"
    find . -type f -size +50M -exec du -m -h {} + | sort -rh >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"

    cd ..
    echo "Cleaning up working clone..."
    rm -rf "$WORKING_PATH"

    echo "Cleaning up mirrored repository..."
    rm -rf "$REPO_PATH"
done

echo "All repositories processed."
