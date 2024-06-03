if [ -z "$(git diff --name-only remotes/origin/HEAD -- README.md)" ]; then
    echo "The CHANGELOG.md file needs to be modified in this PR before doing Merge."
    exit 1
fi  