# Mazeexample Documentation Quick Reference

## Viewing Documentation

### Live Interactive Docs (Recommended)

```bash
task docs-godoc
# Then open: http://localhost:6060/pkg/main/mazeexample/maze/
```

### Generated API Documentation

```bash
task docs-pkg
cat docs/godoc/API.txt
```

### Command Line

```bash
go doc ./maze                  # Package overview
go doc ./maze.Vertex          # Specific function
go doc -all ./maze            # Everything
```

## Documentation Files

### Generated (Auto-created)

- `docs/godoc/API.txt` - Full package API documentation

### Project Documentation

- `README.md` - Main project overview
- `docs/TEST_README.md` - Testing guide
- `docs/JSON_VIEWER_README.md` - Visualization guide
- `docs/BENCHMARK_COMPARISON.md` - Performance analysis
- `docs/TASKFILE_UPDATES.md` - Taskfile enhancements

## Quick Commands

```bash
# Generate docs
task docs-pkg

# View docs
task docs-godoc

# Clean generated docs
task docs-clean

# Show doc info
task docs

# Full help
task help
```

## Documentation Tasks

| Task         | Description                    |
| ------------ | ------------------------------ |
| `docs`       | Show documentation access info |
| `docs-godoc` | Start godoc server on :6060    |
| `docs-pkg`   | Generate API.txt documentation |
| `docs-clean` | Remove generated docs          |

## Writing Documentation

Add comments above exported items:

```go
// Vertex creates or retrieves a node in the graph.
//
// If a node with the given name already exists, it returns
// the existing node. Otherwise, it creates a new node with
// the specified context.
//
// Parameters:
//   - graph: The SST graph context
//   - name: Unique identifier for the node
//   - context: Additional metadata (e.g., "chapter1")
//
// Returns the node pointer.
func Vertex(graph *LinkedSST, name, context string) *Node {
    // implementation
}
```

Then regenerate:

```bash
task docs-clean
task docs-pkg
```

---

**Tip**: Keep godoc server running in a separate terminal for live updates!
