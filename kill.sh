#!/bin/bash

# Check if PID is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <pid>"
    exit 1
fi

# Store the PID
PID=$1

# Send SIGTERM
echo "Sending SIGTERM to process $PID"
kill -15 "$PID"

# Wait 5 seconds
echo "Waiting 5 seconds."
sleep 5

# Check if process is still running
if kill -0 "$PID" 2>/dev/null; then
    # If process is still alive, send SIGKILL
    echo "Sending SIGKILL."
    kill -9 "$PID"
fi