#!/bin/bash

# This script automates the installation of envcontainer from its GitHub repository.
# It downloads the latest binary release, unpacks it,
# and moves the 'envcontainer' executable to /usr/local/bin.

# --- Configuration ---
GITHUB_REPO="erickmaria/envcontainer"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="envcontainer" # The primary script we want to install

# --- Function to check for required commands ---
check_command() {
    if ! command -v "$1" &> /dev/null; then
        echo -e "\n\033[0;31mError:\033[0m The command '$1' is not installed."
        echo "Please install '$1' and try again. For example, on Debian/Ubuntu: sudo apt install $1" >&2
        echo "For 'jq' on Debian/Ubuntu: sudo apt install jq" >&2
        exit 1
    fi
}

# --- Main Installation Script ---
echo "🚀 Starting envcontainer installation script..."
echo "----------------------------------------------"

# 1. Check for necessary tools
echo "1. Checking for required tools (curl, unzip, jq)..."
check_command "curl"
check_command "unzip"
check_command "jq"
echo "All required tools are present. ✅"

# Create a temporary directory for download and unpacking
TMP_DIR=$(mktemp -d -t envcontainer_install_XXXXXX)
if [ -z "$TMP_DIR" ]; then
    echo -e "\n\033[0;31mError:\033[0m Failed to create a temporary directory." >&2
    exit 1
fi
echo "Using temporary directory: $TMP_DIR"

# Ensure the temporary directory is cleaned up on script exit (even if errors occur)
trap "echo 'Cleaning up temporary files...' && rm -rf \"$TMP_DIR\"" EXIT

# 2. Download the latest release
echo -e "\n2. Fetching latest release information from GitHub (https://api.github.com/repos/$GITHUB_REPO/releases/latest)..."
LATEST_RELEASE_INFO=$(curl -sL "https://api.github.com/repos/$GITHUB_REPO/releases/latest")

if [ -z "$LATEST_RELEASE_INFO" ]; then
    echo -e "\n\033[0;31mError:\033[0m Could not retrieve latest release information for $GITHUB_REPO." >&2
    echo "Please check your internet connection or the repository URL." >&2
    exit 1
fi

# Extract the download URL and tag_name using jq
# IMPORTANT: For 'erickmaria/envcontainer', the 'assets'.
# We must use '.browser_download_url' to download.
DOWNLOAD_URL=$(echo "$LATEST_RELEASE_INFO" | jq -r .assets[0].browser_download_url)
TAG_NAME=$(echo "$LATEST_RELEASE_INFO" | jq -r .tag_name)

if [ -z "$DOWNLOAD_URL" ] || [ "$DOWNLOAD_URL" == "null" ]; then
    echo -e "\n\033[0;31mError:\033[0m Could not find a suitable download URL for the latest release." >&2
    echo "This repository does not appear to publish specific downloadable assets with 'browser_download_url'." >&2
    exit 1
fi

echo -e "Found latest release: \033[0;32m$TAG_NAME\033[0m"
echo "Download URL: $DOWNLOAD_URL"

DOWNLOADED_FILE="$TMP_DIR/envcontainer_release.zip"
echo "Downloading $TAG_NAME source code to $DOWNLOADED_FILE..."
if ! curl -sL "$DOWNLOAD_URL" -o "$DOWNLOADED_FILE"; then
    echo -e "\n\033[0;31mError:\033[0m Failed to download the release zip file." >&2
    exit 1
fi
echo "Download complete. ✅"

# 3. Unpack the file
echo -e "\n3. Unpacking the release archive..."
if ! unzip -q "$DOWNLOADED_FILE" -d "$TMP_DIR"; then
    echo -e "\n\033[0;31mError:\033[0m Failed to unpack the zip file '$DOWNLOADED_FILE'." >&2
    exit 1
fi

# Find the unpacked source directory (e.g., erickmaria-envcontainer-HASH)
UNPACKED_SOURCE_DIR=$(find "$TMP_DIR" -maxdepth 1 -mindepth 1 -type f -name "$BINARY_NAME" | head -n 1)

if [ -z "$UNPACKED_SOURCE_DIR" ]; then
    echo -e "\n\033[0;31mError:\033[0m Could not find the unpacked source directory in '$TMP_DIR'." >&2
    exit 1
fi
echo "Unpacked to: $UNPACKED_SOURCE_DIR"

# Define source and target paths for the binary
SOURCE_SCRIPT="$UNPACKED_SOURCE_DIR"
TARGET_SCRIPT="$INSTALL_DIR/$BINARY_NAME"


# Check if the binary exists in the unpacked directory
if [ ! -f "$SOURCE_SCRIPT" ]; then
    echo -e "\n\033[0;31mError:\033[0m binary '$BINARY_NAME' not found in '$UNPACKED_SOURCE_DIR'." >&2
    echo "Please verify the repository structure or update BINARY_NAME in this script." >&2
    exit 1
fi

# 4. Move the unpacked file to /usr/local/bin
echo -e "\n4. Moving '$BINARY_NAME' to '$INSTALL_DIR' and setting permissions..."

# Create INSTALL_DIR if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
    echo "Creating installation directory: $INSTALL_DIR (requires sudo)..."
    if ! sudo mkdir -p "$INSTALL_DIR"; then
        echo -e "\n\033[0;31mError:\033[0m Failed to create directory '$INSTALL_DIR'. Do you have sufficient permissions?" >&2
        exit 1
    fi
fi

echo "Moving $BINARY_NAME (requires sudo)..."
if ! sudo mv "$SOURCE_SCRIPT" "$TARGET_SCRIPT"; then
    echo -e "\n\033[0;31mError:\033[0m Failed to move '$BINARY_NAME' to '$INSTALL_DIR'. Do you have sufficient permissions?" >&2
    exit 1
fi

echo "Making '$TARGET_SCRIPT' executable (requires sudo)..."
if ! sudo chmod +x "$TARGET_SCRIPT"; then
    echo -e "\n\033[0;31mError:\033[0m Failed to make '$TARGET_SCRIPT' executable. Do you have sufficient permissions?" >&2
    exit 1
fi

echo -e "Successfully installed \0  33[0;32menvcontainer\033[0m to $INSTALL_DIR. ✅"
echo "----------------------------------------------"

echo -e "\n🎉 Installation of envcontainer is complete!"
