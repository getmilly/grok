#! /bin/bash

# $1 >> Git commit hash

COMMIT_MESSAGE=$(git log --format=%B -n 1 $1)
CURRENT_TAG=$(git describe --tags)

# Check if the 'git describe' result matches the 1.0.0-123-g1234567 pattern

if [[ $CURRENT_TAG =~ [0-9]+.[0-9]+.[0-9]+ ]]; then
    IFS='-' read -ra TAG_DATA <<< $CURRENT_TAG
    CURR_VERSION=${TAG_DATA[0]}

    IFS='.' read -ra VERSION_DATA <<< $CURR_VERSION

    MAJOR=${VERSION_DATA[0]}
    MINOR=${VERSION_DATA[1]}
    PATCH=${VERSION_DATA[2]}

    if [[ $COMMIT_MESSAGE =~ "[major]" ]]; then
        MAJOR=$(expr $MAJOR + 1)
        MINOR=0
        PATCH=0
    elif [[ $COMMIT_MESSAGE =~ "[minor]" ]]; then
        MINOR=$(expr $MINOR + 1)
        PATCH=0
    elif [[ $COMMIT_MESSAGE =~ "[patch]" ]]; then
        PATCH=$(expr $PATCH + 1)
    else
        echo "Commit type [major | minor | patch] is required"
        exit 1
    fi

    echo "$MAJOR.$MINOR.$PATCH"
else
    echo "0.0.0"
fi