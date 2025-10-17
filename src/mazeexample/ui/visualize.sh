#!/bin/bash
# Maze Solver Visualization Helper Script

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧩 Maze Solver Visualization${NC}"
echo ""

# Check if binary exists
if [ ! -f "./mazeexample" ]; then
	echo -e "${YELLOW}Building mazeexample...${NC}"
	go build -o mazeexample
	echo -e "${GREEN}✓ Build complete${NC}"
fi

# Generate JSON output
echo -e "${YELLOW}Generating JSON output...${NC}"
./mazeexample --json >results.json

# Show statistics
echo -e "${GREEN}✓ Results generated${NC}"
echo ""
echo -e "${BLUE}Statistics:${NC}"
if command -v jq &>/dev/null; then
	cat results.json | jq '.statistics'
else
	echo "  (Install 'jq' to see formatted statistics)"
fi

echo ""
echo -e "${BLUE}Solution Summary:${NC}"
if command -v jq &>/dev/null; then
	cat results.json | jq '.solutions[] | {id, type, total_length}'
else
	echo "  (Install 'jq' to see formatted solution summary)"
fi

# Get the absolute path
VIEWER_PATH="$(pwd)/viewer.html"
RESULTS_PATH="$(pwd)/results.json"

echo ""
echo -e "${GREEN}✓ All files generated successfully!${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "  1. Open ${YELLOW}viewer.html${NC} in your web browser"
echo "  2. Click '📂 Load JSON File' button"
echo "  3. Select ${YELLOW}results.json${NC}"
echo ""
echo -e "${BLUE}File locations:${NC}"
echo "  Viewer:  ${VIEWER_PATH}"
echo "  Results: ${RESULTS_PATH}"
echo ""
echo -e "${BLUE}Quick commands:${NC}"
echo "  # Open with default browser (Linux)"
echo "  xdg-open viewer.html"
echo ""
echo "  # View raw JSON"
echo "  cat results.json | jq ."
echo ""
echo "  # Extract specific data"
echo "  cat results.json | jq '.statistics'"
echo "  cat results.json | jq '.solutions[0].path'"
