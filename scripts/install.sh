#!/bin/bash

install_celify() {
    # Determine OS and Architecture
    os=$(uname | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m | tr '[:upper:]' '[:lower:]' | sed -e s/x86_64/amd64/)

    # Define download URL
    tar="celify-${os}-${arch}.tar.gz"
    url="https://github.com/rdalbuquerque/celify/releases/latest/download"

    # Defined filename
    version=$(curl -s https://api.github.com/repos/rdalbuquerque/celify/releases/latest | jq -r '.tag_name')
    filename="celify_${version}_$os-$arch"

    # Download and install
    echo "Downloading version ${version} of celify-$os-$arch..."
    curl -sL "$url/$tar" -o "/tmp/$tar"
    if [ ! -f "/tmp/$tar" ]; then
      echo "Error downloading celify-$os-$arch."
      return 1
    fi
    echo

    tar -xzf "/tmp/$tar" -C /tmp
    # Optionally remove the tar file after extracting
    # rm "/tmp/$tar"

    echo "Moving /tmp/$filename to /usr/local/bin/celify (you might be asked for your password due to sudo)"
    if [ -x "$(command -v sudo)" ]; then
      sudo mv "/tmp/$filename" "/usr/local/bin/celify"
    else
      mv "/tmp/$filename" "/usr/local/bin/celify"
    fi
    /usr/local/bin/celify --version
    # test exit code
    if [ $? -eq 0 ]; then
      echo "celify installed successfully"
    else
      echo "Failed to install celify"
      return 1
    fi
    echo
}

# Execute the function
install_celify
