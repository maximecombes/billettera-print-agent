#!/bin/bash
# Billettera Print Agent - macOS installer
# Removes quarantine attribute and launches the agent

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
AGENT="$SCRIPT_DIR/BilletteraPrintAgent"

if [ ! -f "$AGENT" ]; then
    echo "Error: BilletteraPrintAgent not found in $SCRIPT_DIR"
    exit 1
fi

echo "Removing macOS quarantine..."
xattr -cr "$AGENT"
chmod +x "$AGENT"
echo "Done. Launching BilletteraPrintAgent..."
open "$AGENT"
