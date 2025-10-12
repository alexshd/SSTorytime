# text2n4l-web - API-Only Backend

## Summary

The text2n4l-web directory has been cleaned and refactored into a pure API-only Go backend. All UI components, templates, and frontend assets have been removed.

## What Was Removed

- `/static/` - All static assets (CSS, JS)
- `/templates/` - All Go templates
- `/web/` - Web-specific static directories
- `.air.toml` - Hot reload configuration (not needed)
- `fix-templates.sh` - Template fixing script
- `/scripts/` - Build/deployment scripts
- Old test files related to templates and handlers
- Coverage reports and temporary files

## Current Structure

```
text2n4l-web/
├── cmd/
│   ├── cli/          # CLI tool (existing)
│   └── web/          # API server ✓
├── internal/
│   ├── analyzer/     # N4L conversion engine ✓
│   └── web/          # API handlers ✓
├── testdata/         # Test data
├── tests/            # Integration tests
├── bin/              # Build output
├── Makefile          # Build targets ✓
├── go.mod            # Go module
├── go.sum            # Dependencies
└── README.md         # API documentation ✓
```

## API Server Features

### Endpoints

1. **POST /api/convert** - Convert text to N4L format
2. **GET /debug/pprof/** - Profiling index
3. **GET /debug/pprof/profile** - CPU profiling
4. **GET /debug/pprof/heap** - Memory profiling
5. **GET /debug/pprof/goroutine** - Goroutine info

### Makefile Targets

- `make build` - Build the API server
- `make dev-server` - Run the API server
- `make test-api` - Run API tests
- `make test` - Run all tests
- `make fmt` - Format code
- `make lint` - Run linters
- `make profile-cpu` - Capture CPU profile
- `make profile-heap` - Capture heap profile
- `make clean` - Clean build artifacts
- `make help` - Show help

## Usage

### Start Server

```bash
cd text2n4l-web
make dev-server
```

### Test API

```bash
curl -X POST http://localhost:8080/api/convert \
  -d "text=Your text here."
```

### Profile Performance

```bash
# In terminal 1
make dev-server

# In terminal 2
make profile-cpu
go tool pprof cpu.prof
```

## Next Steps

The API backend is now clean and ready for:

1. **Frontend Development** - Build a separate Vite + Tailwind CSS v4 frontend
2. **Go Refactoring** - Further optimize the Go codebase
3. **Deployment** - Deploy as a containerized service

The backend has no dependencies on UI frameworks and can be consumed by any frontend.
