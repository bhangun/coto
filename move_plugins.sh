#!/bin/bash
# Move plugin files from cmd/plugins to cmd/main

cd "$(dirname "$0")"

echo "Moving plugin files..."

# Move all .go files from cmd/plugins to cmd/main
for file in cmd/plugins/*.go; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        echo "Moving $filename to cmd/main/"
        mv "$file" "cmd/main/$filename"
    fi
done

echo "Done! Plugin files moved to cmd/main/"
echo ""
echo "Now rebuild with:"
echo "  go build -o bin/coto ./cmd/main"
