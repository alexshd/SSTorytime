# Maze Solver JSON Visualization

## Quick Start

### Option 1: Using Task (Recommended)

```bash
# Generate JSON and get instructions
task visualize

# Or serve with HTTP server (best experience)
task serve
# Then open http://localhost:8000/viewer.html in your browser
```

### Option 2: Manual

1. **Generate JSON output:**

   ```bash
   ./mazeexample --json > results.json
   ```

2. **Serve with HTTP (recommended):**

   ```bash
   # Python 3
   python3 -m http.server 8000
   # or Python 2
   python -m SimpleHTTPServer 8000
   # or Node.js
   npx serve

   # Then open http://localhost:8000/viewer.html
   ```

   The viewer will automatically load `results.json` when served via HTTP.

3. **Or open directly (requires file picker):**
   - Open `viewer.html` in your web browser (double-click or drag into browser)
   - Click "Load JSON File" button
   - Select the `results.json` file you just generated

## Why Use HTTP Server?

When served via HTTP, the viewer automatically loads `results.json` without requiring manual file selection. This provides a smoother experience, especially when regenerating results multiple times.

## Features

The visualization displays:

### 📊 Statistics Dashboard

- **Total Solutions:** Number of valid paths found
- **Total Loops:** Number of loop corrections detected
- **Max Left/Right Depth:** Maximum search depth reached from each direction
- **Search Steps:** Total iterations of the bidirectional search

### 🎯 Solutions Found

Each solution shows:

- **Solution ID:** Unique identifier
- **Type:** tree/dag/graph classification
- **Depths:** Left and right traversal depths (L:X R:Y)
- **Total Length:** Number of steps in the path
- **Visual Path:** Node-by-node visualization with arrow types

### 🔄 Loop Corrections (if any)

Shows any loops detected during the search with full path details

### 📈 Search Timeline

Step-by-step visualization of the bidirectional search:

- **Turn number:** Current iteration
- **Depths:** Current search depth from left and right
- **Frontier nodes:** Nodes being explored in this iteration
- **Solutions found:** Highlighted when paths connect

## Example Usage

```bash
# Build the example
go build -o mazeexample

# Generate text output (default)
./mazeexample

# Generate JSON output
./mazeexample --json

# Save JSON to file for visualization
./mazeexample --json > my_results.json

# Or pipe to another tool
./mazeexample --json | jq '.statistics'
```

## JSON Structure

The output JSON contains:

```json
{
  "start_node": "maze_a7",
  "end_node": "maze_i6",
  "max_depth": 16,
  "solutions": [...],
  "loops": [...],
  "search_steps": [...],
  "statistics": {
    "total_solutions": 1,
    "total_loops": 0,
    "max_left_depth": 15,
    "max_right_depth": 14,
    "total_search_steps": 29
  }
}
```

### Path Link Format

Each step in a path is represented as:

```json
{
  "from": "maze_b7",
  "to": "maze_b6",
  "arrow": "fwd",
  "weight": 1
}
```

Arrow types:

- `fwd` - Forward arrow
- `bkw` - Backward arrow
- `neu` - Neutral/bidirectional arrow

### Search Step Format

Each search iteration includes:

```json
{
  "turn": 25,
  "left_depth": 13,
  "right_depth": 12,
  "left_frontier": ["maze_e3", "maze_d4"],
  "right_frontier": ["maze_e3", "maze_f5"],
  "solutions_found": 1,
  "loops_found": 0
}
```

## Integration with Your UI

The JSON output is designed to be consumed by web applications:

1. **Fetch the JSON:**

   ```javascript
   fetch("results.json")
     .then((response) => response.json())
     .then((data) => {
       console.log("Solutions:", data.solutions);
       console.log("Statistics:", data.statistics);
     });
   ```

2. **Render solutions:**

   ```javascript
   data.solutions.forEach((solution) => {
     console.log(`Solution ${solution.id}:`, solution.path);
   });
   ```

3. **Animate search progress:**
   ```javascript
   data.search_steps.forEach((step, index) => {
     setTimeout(() => {
       updateVisualization(step);
     }, index * 100);
   });
   ```

## Compatibility

The viewer works in all modern browsers:

- ✅ Chrome/Edge 90+
- ✅ Firefox 88+
- ✅ Safari 14+

No external dependencies required - just open the HTML file!

## Tips

- The search timeline is scrollable - use it to see the complete search progression
- Solutions are color-coded: green for solutions, yellow for metadata
- Hover over solution cards to see the highlight effect
- The viewer is fully responsive and works on mobile devices
