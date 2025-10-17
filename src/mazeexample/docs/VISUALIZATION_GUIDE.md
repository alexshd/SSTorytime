# Visualization Guide

## Quick Reference

### Best Way (HTTP Server)

```bash
task serve
```

This will:

1. Build the binary (if needed)
2. Generate `results.json`
3. Show statistics
4. Start HTTP server on port 8000
5. Display the URL: `http://localhost:8000/viewer.html`

**Advantages:**

- ✅ Auto-loads `results.json` (no file picker)
- ✅ Instant refresh when regenerating results
- ✅ No CORS issues
- ✅ Professional workflow

### Alternative Ways

```bash
# Generate JSON only
task visualize

# Open in browser (file picker mode)
task view

# Just run and get JSON output
task run-json
```

## Troubleshooting

### Port 8000 Already in Use

If you see `Address already in use`, either:

1. **Stop the existing server:**

   ```bash
   # Find the process
   lsof -i :8000
   # Kill it
   kill <PID>
   ```

2. **Use a different port:**
   ```bash
   # Edit Taskfile.yml and change 8000 to another port (e.g., 8080)
   # Or run Python directly:
   python3 -m http.server 8080
   ```

### File Picker Doesn't Work

When opening `viewer.html` directly with `file://` protocol:

**Symptoms:**

- Button click does nothing
- No file picker appears
- Console shows errors

**Solution:**
Use the HTTP server method instead:

```bash
task serve
```

### Auto-load Not Working

If the viewer doesn't auto-load `results.json` when using HTTP server:

**Check:**

1. Is `results.json` in the same directory as `viewer.html`?
2. Are you accessing via `http://` (not `file://`)?
3. Open browser console (F12) to see error messages

**Generate results:**

```bash
task visualize  # Generates fresh results.json
```

### Browser CORS Errors

```
Access to fetch at 'file:///...' from origin 'null' has been blocked by CORS policy
```

**Solution:** You MUST use HTTP server for auto-load feature:

```bash
task serve
```

This is a browser security restriction - the `file://` protocol doesn't allow JavaScript to load other local files.

## Workflow Examples

### Development Workflow

```bash
# Terminal 1: Keep server running
cd /home/alex/SHDProj/SSTorytime/src/mazeexample
task serve

# Terminal 2: Make changes and regenerate
cd /home/alex/SHDProj/SSTorytime/src/mazeexample
# Edit code...
go build -o mazeexample
./mazeexample --json > results.json
# Refresh browser to see new results
```

### One-time View

```bash
task visualize
# Then manually open viewer.html and load results.json
```

### CI/CD Pipeline

```bash
# Generate JSON in headless environment
./mazeexample --json > results.json

# Save as artifact
cp results.json viewer.html artifacts/

# View later by serving artifacts directory
cd artifacts && python3 -m http.server 8000
```

## Advanced Usage

### Custom Port

Edit `Taskfile.yml`:

```yaml
serve:
  cmds:
    - python3 -m http.server 9000 # Change port here
```

### Using npx serve

```bash
task visualize
npx serve  # More features than Python http.server
```

### Docker

```dockerfile
FROM python:3-alpine
WORKDIR /app
COPY viewer.html results.json ./
CMD ["python3", "-m", "http.server", "8000"]
EXPOSE 8000
```

```bash
docker build -t maze-viewer .
docker run -p 8000:8000 maze-viewer
```

## Browser Compatibility

### Supported Browsers

| Browser | Version | Status          |
| ------- | ------- | --------------- |
| Chrome  | 90+     | ✅ Full support |
| Firefox | 88+     | ✅ Full support |
| Safari  | 14+     | ✅ Full support |
| Edge    | 90+     | ✅ Full support |

### Features Used

- **Fetch API** - For auto-loading results.json
- **CSS Grid** - For responsive layout
- **Arrow Functions** - ES6 JavaScript
- **Template Literals** - ES6 JavaScript

All modern browsers support these features.

## Mobile Viewing

The viewer is fully responsive and works on mobile devices:

1. Generate results on your development machine
2. Start HTTP server: `task serve`
3. Find your local IP: `ip addr show` or `ifconfig`
4. Access from mobile: `http://192.168.x.x:8000/viewer.html`

**Note:** Make sure both devices are on the same network.

## Performance

### Large Result Files

If your maze generates large JSON files (>1MB):

**Options:**

1. Use the compact format (if implemented)
2. Limit search depth
3. Filter solutions in the visualization

**Current file size:**

```bash
ls -lh results.json
# Example: 11K (very small, no issues)
```

### Many Solutions

The viewer can handle hundreds of solutions, but browser performance may degrade with thousands. Consider:

1. Paginating solutions
2. Lazy loading
3. Virtual scrolling
4. Filtering by type

## Tips

### Quick Iteration

Keep the HTTP server running and just regenerate:

```bash
# Terminal 1: Keep running
task serve

# Terminal 2: Quick regeneration
./mazeexample --json > results.json && echo "Refresh browser"
```

### Comparing Results

```bash
# Save different runs
./mazeexample --json > results-v1.json
# Make changes
./mazeexample --json > results-v2.json

# Compare
diff <(jq '.statistics' results-v1.json) <(jq '.statistics' results-v2.json)
```

### Pretty Print JSON

```bash
cat results.json | jq . | less
cat results.json | jq '.solutions[0]' | less
cat results.json | jq '.statistics'
```

### Extract Specific Data

```bash
# Solution count
jq '.statistics.total_solutions' results.json

# Path lengths
jq '.solutions[].total_length' results.json

# Frontier sizes at each step
jq '.search_steps[].left_frontier | length' results.json
```

## Common Issues

### Issue: "Python not found"

**Solution:**

```bash
# Install Python
sudo apt install python3  # Debian/Ubuntu
sudo yum install python3  # RHEL/CentOS
brew install python3      # macOS

# Or use Node.js instead
npx serve
```

### Issue: Browser doesn't open

**Solution:**

```bash
# Linux
xdg-open http://localhost:8000/viewer.html

# macOS
open http://localhost:8000/viewer.html

# Or just paste URL in browser
```

### Issue: Results not updating

**Solution:**

1. Hard refresh browser: `Ctrl+Shift+R` (or `Cmd+Shift+R` on macOS)
2. Check timestamp: `ls -l results.json`
3. Verify regeneration: `task visualize`

## See Also

- [README.md](README.md) - Main documentation
- [JSON_VIEWER_README.md](JSON_VIEWER_README.md) - JSON schema details
- [JSON_OUTPUT_SUMMARY.md](JSON_OUTPUT_SUMMARY.md) - Implementation details
- [TEST_README.md](TEST_README.md) - Testing documentation
