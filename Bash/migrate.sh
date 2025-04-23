#!/bin/bash
set -e

repo=${1:-$REPO_NAME}
LFS_THRESHOLD="1MB"
LFS_TRACK_LIMIT="50MB"
step_size=${2:-1000}
project=${3:-"Products"}

if [[ -z "$repo" || -z "$TFS_USERNAME" || -z "$TFS_PASS" || -z "$MIGRATION_PAT" ]]; then
    echo "Error: Missing required environment variables (REPO_NAME, TFS_USERNAME, TFS_PASS, MIGRATION_PAT)."
    exit 1
fi

echo "Adding GitHub to SSH known hosts..."
mkdir -p ~/.ssh
ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

TFS_REPO_URL="https://tfs.revenuepremier.com/tfs/RSI/$project/_git/$repo"
REPO_PATH="$(pwd)/$repo.git"
WORKING_PATH="$(pwd)/$repo-work"
GITHUB_URL="https://$MIGRATION_PAT@github.com/revenue-solutions-inc/$repo.git"

echo "Cleaning up existing directories..."
[ -d "$REPO_PATH" ] && rm -rf "$REPO_PATH"
[ -d "$WORKING_PATH" ] && rm -rf "$WORKING_PATH"

echo "Cloning repository from TFS (mirror)..."
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

echo "Creating working repository with all branches..."
git clone --no-single-branch "$REPO_PATH" "$WORKING_PATH"
cd "$WORKING_PATH"

echo "Ensuring all branches are tracked locally..."
git branch -r | grep -v HEAD | sed 's/origin\///' | while read branch; do
    if ! git show-ref --quiet --verify "refs/heads/$branch"; then
        git checkout -f -b "$branch" "origin/$branch"
    else
        echo "Branch $branch already exists. Skipping creation."
        git checkout -f "$branch"
    fi
done

echo "Initializing Git LFS..."
GIT_LFS_TRACK_NO_INSTALL=1 git lfs install --skip-repo

git config --global lfs.concurrenttransfers 10
git config --global lfs.activitytimeout 3600
git config --global lfs.httptransfertimeout 3600
git config --global lfs.url "https://$MIGRATION_PAT@github.com/revenue-solutions-inc/$repo.git/info/lfs"
git config lfs.https://github.com/revenue-solutions-inc/$repo.git/info/lfs.locksverify false

echo "Git LFS initialized."

echo "Identifying large files (>$LFS_THRESHOLD)..."
mapfile -d '' large_files < <(
    git rev-list --all --objects | \
    git cat-file --batch-check='%(objecttype) %(objectname) %(objectsize) %(rest)' | \
    awk -v threshold="${LFS_THRESHOLD%MB}" '
    BEGIN { threshold_bytes = threshold * 1024 * 1024 }
    $1 == "blob" && $3 > threshold_bytes { print $4 "\0" }
    ' | sort -zu
)

if [ ${#large_files[@]} -gt 0 ]; then
    echo "Found large files:"
    printf '%s\n' "${large_files[@]}"
    
    echo "Configuring LFS tracking for files >$LFS_TRACK_LIMIT..."
    for file in "${large_files[@]}"; do
        if [ -f "$file" ]; then
            size=$(git cat-file -s "$(git rev-parse HEAD:"$file")" 2>/dev/null || continue)
            if [ "$size" -gt $(( ${LFS_TRACK_LIMIT%MB} * 1024 * 1024 )) ]; then
                echo "Tracking $file with LFS (size: $((size/1024/1024))MB"
                git lfs track "$file"
            fi
        else
            echo "Skipping invalid path: $file"
        fi
    done
    
    if [[ -n $(git diff --name-only) ]]; then
        git add .gitattributes
        git commit -m "Track large files with Git LFS"
    fi
    
    echo "Rewriting history to apply LFS changes..."
    git lfs migrate import --yes --everything --above="$LFS_TRACK_LIMIT" --verbose
else
    echo "No files above $LFS_THRESHOLD found."
fi

echo "Configuring GitHub remote..."
if git remote | grep -q origin; then
    echo "Removing existing origin remote..."
    git remote remove origin
fi

git remote add origin "$GITHUB_URL"
git remote set-url origin "$GITHUB_URL"
echo "GitHub remote configured successfully."

echo "Configuring Git buffers for large pushes..."
git config --global http.postBuffer 5242880000
git config --global pack.windowMemory 256m
git config --global pack.packSizeLimit 2g
git config --global core.sshCommand "ssh -o UserKnownHostsFile=~/.ssh/known_hosts"

echo "Creating GitHub repository..."
response=$(curl -s -o /dev/null -w "%{http_code}" \
    -H "Authorization: token $MIGRATION_PAT" \
    -H "Accept: application/vnd.github.v3+json" \
    "https://api.github.com/repos/revenue-solutions-inc/$repo")

if [[ "$response" != "200" ]]; then
    echo "Creating GitHub repository..."
    create_response=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
        "https://api.github.com/orgs/revenue-solutions-inc/repos" \
        -H "Authorization: token $MIGRATION_PAT" \
        -H "Accept: application/vnd.github.v3+json" \
        -d '{
            "name": "'"$repo"'",
            "private": true,
            "has_issues": false,
            "has_projects": false,
            "has_wiki": false
        }')
    
    if [[ "$create_response" != "201" ]]; then
        echo "Error: Failed to create repository (HTTP $create_response)"
        exit 1
    fi
    echo "GitHub repository created successfully."
else
    echo "GitHub repository already exists."
fi

echo "Configuring Git LFS credentials..."
echo "https://$MIGRATION_PAT@github.com" > .git-credentials
git config --global credential.helper "store --file=$(pwd)/.git-credentials"

echo "Git Auth with PAT..."
git ls-remote "$GITHUB_URL"

# echo "Starting incremental push to GitHub..."

# mapfile -t branches < <(git for-each-ref --format='%(refname:short)' refs/heads/ refs/remotes/origin/)
# branches=("${branches[@]#origin/}")
# mapfile -t branches < <(printf '%s\n' "${branches[@]}" | grep -v '^$' | sort -u)

# for branch in "${branches[@]}"; do
#     if [[ -z "$branch" || "$branch" == "HEAD" ]]; then
#         continue
#     fi

#     echo "Pushing branch: $branch"
#     git checkout -f "$branch"
#     git push origin "$branch" --force --progress
# done

echo "Starting incremental push to GitHub..."
mapfile -t branches < <(git for-each-ref --format='%(refname:short)' refs/heads/ refs/remotes/origin/)
branches=("${branches[@]#origin/}")
mapfile -t branches < <(printf '%s\n' "${branches[@]}" | grep -v '^$' | sort -u)

for branch in "${branches[@]}"; do
    if [[ -z "$branch" || "$branch" == "HEAD" ]]; then
        continue
    fi

    echo "Pushing branch: $branch"
    git checkout "$branch"
    step_commits=$(git log --oneline --reverse refs/heads/"$branch" | awk "NR % $step_size == 0")

    if [[ -z "$step_commits" ]]; then
        echo "No commits found for incremental push. Pushing all at once..."
        git push origin "$branch" --force --progress
    else
        echo "Pushing commits incrementally for branch: $branch"
        echo "$step_commits" | while read commit message; do
            echo "Pushing commit: $commit"
            git push origin +"$commit":refs/heads/"$branch"
        done
    fi
done

echo "Pushing all branches to GitHub..."
git push origin --all --force --progress

echo "Pushing tags to GitHub..."
git push origin --tags --force --progress

echo "Cleaning up..."
cd ..
rm -rf "$WORKING_PATH"
rm -f .git-credentials

echo "Migration of $repo completed successfully!"