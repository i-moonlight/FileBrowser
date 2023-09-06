#!/bin/bash

echo "This is a test script."

# Check if arguments are provided
if [ $# -gt 0 ]; then
    echo "Arguments:"
    for arg in "$@"; do
        echo "- $arg"
    done
else
    echo "No arguments provided."
fi

# Simulate a time-consuming task
echo "Simulating a time-consuming task..."
sleep 1

echo "Test script completed."