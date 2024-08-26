#!/bin/bash

# This script just runs curl requests continuously to the server
# It takes the number of iterations as an argument

# Check if the number of arguments is correct
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <number of iterations>"
    exit 1
fi

# Get the number of iterations
iterations=$1

# Run the curl requests
for i in $(seq 1 $iterations); do
    curl localhost:8090/worker/
done