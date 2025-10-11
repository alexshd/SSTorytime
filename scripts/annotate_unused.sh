#!/usr/bin/env bash

# annotate_unused.sh
#
# Automatically add "// UNUSED:" comments to symbols identified in UNUSED_REPORT.md
# Based on the CSV-like markdown table format, this script parses each row and adds
# inline comments to the source files at the specified line numbers.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")"/.. && pwd)"
REPORT_FILE="$REPO_ROOT/UNUSED_REPORT.md"

if [[ ! -f "$REPORT_FILE" ]]; then
	echo "Error: $REPORT_FILE not found" >&2
	exit 1
fi

echo "Annotating unused symbols from $REPORT_FILE..." >&2

# Extract the table rows from the markdown (skip header lines)
awk '/^\| \.\/[^|]*\| [0-9]/ {
    # Remove leading/trailing spaces and pipes
    gsub(/^\| */, "")
    gsub(/ *\|$/, "")
    # Split by pipe and extract fields
    gsub(/ *\| */, "|")
    print $0
}' "$REPORT_FILE" | while IFS='|' read -r file line kind name package in_pkg cross_pkg; do
	# Clean up field values
	file=$(echo "$file" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
	line=$(echo "$line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
	kind=$(echo "$kind" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
	name=$(echo "$name" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')

	# Convert relative path to absolute
	abs_file="$REPO_ROOT/$file"

	if [[ ! -f "$abs_file" ]]; then
		echo "Warning: File $abs_file not found, skipping" >&2
		continue
	fi

	echo "Annotating $name ($kind) at line $line in $file" >&2

	# Create a temporary file for the modified content
	temp_file=$(mktemp)

	# Read the file line by line and add annotation
	line_num=1
	while IFS= read -r file_line || [[ -n "$file_line" ]]; do
		if [[ $line_num -eq $line ]]; then
			# Check if already annotated
			if [[ "$file_line" =~ "// UNUSED:" ]]; then
				echo "  Already annotated, skipping" >&2
				echo "$file_line" >>"$temp_file"
			else
				# Add the annotation at the end of the line
				echo "${file_line} // UNUSED: $kind $name (0 refs)" >>"$temp_file"
			fi
		else
			echo "$file_line" >>"$temp_file"
		fi
		((line_num++))
	done <"$abs_file"

	# Replace the original file with the annotated version
	mv "$temp_file" "$abs_file"
done

echo "Annotation complete!" >&2
