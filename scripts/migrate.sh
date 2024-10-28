#!/bin/bash

# Get the directory of the script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Default values
COMMAND="up"
STEP=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--command)
            COMMAND="$2"
            shift
            shift
            ;;
        -s|--step)
            STEP="$2"
            shift
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Validate command
if [[ "$COMMAND" != "up" && "$COMMAND" != "down" ]]; then
    echo "Invalid command. Use 'up' or 'down'"
    exit 1
fi

# Run the migration
if [ -n "$STEP" ]; then
    echo "Running migration ${COMMAND} ${STEP} steps..."
    go run cmd/migrate/main.go -command=${COMMAND} -steps=${STEP}
else
    echo "Running migration ${COMMAND}..."
    go run cmd/migrate/main.go -command=${COMMAND}
fi