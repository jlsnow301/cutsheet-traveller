#!/bin/bash

# Get the full path to the file passed as an argument
file_path="$1"

# Change to the directory where the script is located
cd "$(dirname "$0")"

# Ensure the binary is executable
chmod +x ./src/cutsheet-traveller.mac

# Execute the binary within the src directory with the full path to the dragged file
./src/cutsheet-traveller.mac "$file_path"

read -p "Press enter to continue"