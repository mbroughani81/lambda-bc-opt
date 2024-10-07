#!/bin/bash

# Number of activations to retrieve, passed as the first argument
N=$1

# Check if N is provided
if [ -z "$N" ]; then
    echo "Usage: $0 <number_of_activations>"
    exit 1
fi

# Initialize variables
limit=200
skip=0
remaining=$N

# Fetch activations in batches until we get the last N activations
while [ "$remaining" -gt 0 ]; do
    # Calculate the limit for the next batch
    if [ "$remaining" -lt "$limit" ]; then
        limit=$remaining
    fi

    # Fetch the activations
    activations=$(wsk activation list --limit $limit --skip $skip)

    # Print the activations
    echo "$activations"

    # Count the number of fetched activations
    count=$(echo "$activations" | wc -l)

    # Reduce the remaining activations to fetch
    remaining=$((remaining - count))

    # Break if we received fewer activations than the limit (no more activations left)
    if [ "$count" -lt "$limit" ]; then
        break
    fi

    # Increase the skip value for the next batch
    skip=$((skip + count))
done

