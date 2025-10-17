# Go Workspace Configuration

This repository uses **Go Workspaces** (`go.work`) to manage multiple Go modules in a unified development environment.

## What is This?

The `go.work` file tells Go to treat all these modules as a single workspace:

```
SSTorytime/                          # Root module
├── editors/tree-sitter-n4l/bindings/go/  # Tree-sitter Go bindings
├── n4l/                            # N4L module
├── pkg/SSTorytime/                 # Core SSTorytime package
├── src/                            # Source module
├── src/mazeexample/                # Maze solver example
└── text2n4l-web/                   # Web interface module
```

## Benefits

✅ **No more `GOWORK=off`** - All modules work together seamlessly
✅ **Unified testing** - Run tests across all modules from any directory
✅ **Local development** - Changes in one module immediately visible to others
✅ **IDE support** - Better autocomplete and navigation across modules

## Usage

### Running Tests

From any module directory, tests now work without special flags:

```bash
# Before (needed GOWORK=off)
cd src/mazeexample
GOWORK=off go test ./maze

# Now (workspace-aware)
cd src/mazeexample
go test ./maze
```

### Building

```bash
# Build any module from anywhere
cd src/mazeexample
go build

cd text2n4l-web
go build ./cmd/web
```

### Benchmarking

```bash
# Run benchmarks in any module
cd src/mazeexample
go test -bench=. ./maze
```

## Workspace Commands

```bash
# Show workspace info
go work edit -print

# Add a new module to workspace
go work use ./path/to/new/module

# Remove a module from workspace
go work edit -dropuse ./path/to/module

# Sync dependencies across workspace
go work sync
```

## Important Notes

⚠️ **`go.work` is NOT committed to git** by default (it's in `.gitignore`)

- This is a local development tool
- Each developer can customize their workspace
- CI/CD builds each module independently

⚠️ **Module independence is preserved**

- Each module still has its own `go.mod`
- Modules can still be built/tested independently
- The workspace is just for convenience

## Verifying Workspace

Check that all modules are recognized:

```bash
cd /home/alex/SHDProj/SSTorytime
cat go.work
```

Should show all 7 modules listed under `use (...)`.

## Troubleshooting

**Problem**: Tests fail with "module not in workspace"
**Solution**: Add the module with `go work use ./path/to/module`

**Problem**: Changes in one module not visible in another
**Solution**: Run `go work sync` to synchronize dependencies

**Problem**: Want to build without workspace
**Solution**: Use `GOWORK=off go build` to temporarily disable workspace

---

Created: 2025-10-17
