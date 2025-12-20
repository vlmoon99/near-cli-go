#!/usr/bin/env bash

OUTPUT="go_context.txt"

# Clear output file
> "$OUTPUT"

echo "### GO PROJECT CONTEXT (ROOT FILES ONLY)" >> "$OUTPUT"
echo "### Generated at: $(date)" >> "$OUTPUT"
echo "" >> "$OUTPUT"

for file in *.go; do
  # Skip if no .go files exist
  [ -e "$file" ] || continue

  echo "==================================================" >> "$OUTPUT"
  echo "FILE: $file" >> "$OUTPUT"
  echo "==================================================" >> "$OUTPUT"
  echo "" >> "$OUTPUT"

  cat "$file" >> "$OUTPUT"
  echo -e "\n\n" >> "$OUTPUT"
done

echo "Context written to $OUTPUT"
