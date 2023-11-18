#!/bin/bash

# Determine OS and Architecture
os=$(uname | tr '[:upper:]' '[:lower:]')
arch=$(uname -m | tr '[:upper:]' '[:lower:]' | sed -e s/x86_64/amd64/)

# Define download URL
url="https://github.com/rdalbuquerque/celify/releases/latest/download/celify-${os}-${arch}"

# Download and install
curl -L "$url" -o celify
chmod +x celify

# Move to a location in PATH
sudo mv celify /usr/local/bin/

echo "Celify CLI installed successfully"
