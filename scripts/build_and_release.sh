#!/bin/bash

# Exit on error
set -e


GIT_SHORT_REV=$(git rev-parse --short HEAD)

# Variables (Update as needed)
REPO_OWNER="heimdahl-xyz"           # Replace with your GitHub username or organization
REPO_NAME="heimdahl-cli"          # Replace with your repository name
TAG=$GIT_SHORT_REV                      # Replace with your release tag (e.g., v1.0.0)
GITHUB_TOKEN=$1     # Replace with your GitHub token or set as an environment variable
BUILD_DIR="build"                    # Directory to store builds
BINARY_NAME="heimdahl-cli"           # Replace with the desired binary name
PLATFORMS=("linux/amd64"
"linux/arm64"
"darwin/amd64"
"darwin/arm64"
"windows/amd64"
"windows/arm64") # Platforms to build for

# Ensure dependencies
command -v jq >/dev/null 2>&1 || { echo "jq is required. Install it and try again."; exit 1; }
command -v zip >/dev/null 2>&1 || { echo "zip is required. Install it and try again."; exit 1; }

# Step 1: Cross-compile binaries
echo "Cross-compiling binaries..."
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

for PLATFORM in "${PLATFORMS[@]}"; do
  OS=$(echo $PLATFORM | cut -d'/' -f1)
  ARCH=$(echo $PLATFORM | cut -d'/' -f2)
  OUTPUT="$BUILD_DIR/${BINARY_NAME}-${OS}-${ARCH}"

  if [ "$OS" == "windows" ]; then
    OUTPUT="${OUTPUT}.exe"
  fi

  echo "Building for $OS/$ARCH..."
  GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT
done

# Step 2: Create tarballs and zip archives
echo "Creating tarballs and zip archives..."
for FILE in $BUILD_DIR/*; do
  BASENAME=$(basename $FILE)

  if [[ $FILE == *.exe ]]; then
    zip -j $BUILD_DIR/${BASENAME}.zip $FILE
  else
    tar -czvf $BUILD_DIR/${BASENAME}.tar.gz -C $BUILD_DIR $BASENAME
  fi
done

## Step 3: Get Release ID
#echo "Fetching release ID for tag $TAG..."
RELEASE_RESPONSE=$(curl -s \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github+json" \
  "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/tags/$TAG")

RELEASE_ID=$(echo "$RELEASE_RESPONSE" | jq -r .id)

if [ "$RELEASE_ID" == "null" ]; then
  echo "Release with tag $TAG not found. Creating a new release..."
  RELEASE="heimdahl-cli-$TAG"

  RELEASE_RESPONSE=$(curl -s \
    -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github+json" \
    https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases \
    -d "{\"tag_name\": \"$RELEASE\", \"name\": \"$RELEASE\", \"body\": \"Release created via script.\"}")

  RELEASE_ID=$(echo "$RELEASE_RESPONSE" | jq -r .id)
fi

if [ "$RELEASE_ID" == "null" ]; then
  echo "Failed to create or retrieve release. Exiting."
  exit 1
fi

echo "Release ID: $RELEASE_ID"

# Step 4: Upload assets
echo "Uploading assets to GitHub release"

UPLOAD_URL="https://uploads.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/$RELEASE_ID/assets"
echo "Uploading assets to GitHub release..."

for FILE in $BUILD_DIR/*.{tar.gz,zip}; do
  BASENAME=$(basename $FILE)
  echo "Uploading $BASENAME..."

  curl -s \
    -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Content-Type: $(file --mime-type -b $FILE)" \
    --data-binary @$FILE \
    "$UPLOAD_URL?name=$BASENAME"

  echo "$BASENAME uploaded."
done
