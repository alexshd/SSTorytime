# Text to N4L Converter Frontend

A modern web interface for converting text to N4L DSL format with **real-time arrow validation** using Vite and Tailwind CSS v4.

## ✨ Key Features

### Core Functionality

- Clean, responsive UI built with Tailwind CSS v4
- Real-time text conversion via API
- File upload support (any text-based file)
- Copy to clipboard functionality
- Save edited N4L files
- Keyboard shortcuts (Ctrl/Cmd + Enter to convert)

### 🎯 NEW: Arrow Validation (v1.0)

**Real-time N4L arrow validation** to prevent parser errors!

- **Visual Error Detection**: Invalid arrows are highlighted in red with ⚠️ warning
- **Valid Arrow Highlighting**: Recognized arrows shown in blue
- **Smart Suggestions**: Keyword-based matching suggests correct arrows
- **Categorized Arrow Menu**: Browse 300+ valid arrows organized by semantic type:
  - 🔗 NR-0: Similarity (78 arrows)
  - ➡️ LT-1: Causality (70+ arrows)
  - 📦 CN-2: Composition (90+ arrows)
  - 🏷️ Properties (80+ arrows)
  - ⭐ Special annotations
- **Interactive Editing**: Click any arrow to see alternatives or fix errors
- **Parser Error Prevention**: Fix issues before they cause compilation errors

#### Quick Example

```
❌ Before: X (appears close to) Y     [RED - Invalid]
✅ After:  X (similar to) Y           [BLUE - Valid]
```

See [ARROW_VALIDATION.md](./ARROW_VALIDATION.md) for complete documentation.

## 📚 Documentation

- **[N4L_EDITING_GUIDE.md](./N4L_EDITING_GUIDE.md)** - Complete guide to N4L syntax and editing workflow
- **[ARROW_VALIDATION.md](./ARROW_VALIDATION.md)** - Technical details on arrow validation
- **[VALIDATION_VISUAL_GUIDE.md](./VALIDATION_VISUAL_GUIDE.md)** - Visual reference with examples

## Development

### Prerequisites

- Node.js (v16 or higher)
- npm
- Go backend API server running on port 8080

### Setup

1. Install dependencies:

```bash
npm install
```

2. Start the development server:

```bash
npm run dev
```

3. Open your browser to `http://localhost:5173`

### API Integration

The frontend connects to a Go backend API at `/api/convert` endpoint. Make sure the backend server is running on port 8080 for the proxy to work correctly.

### Building for Production

```bash
npm run build
```

The built files will be in the `dist/` directory.

## Technology Stack

- **Vite** - Fast build tool and development server
- **Tailwind CSS v4** - Utility-first CSS framework
- **Vanilla JavaScript** - No framework dependencies
- **Proxy Configuration** - Routes `/api/*` to Go backend

## File Structure

```
src/
  main.js       # Main application logic and UI
  style.css     # Tailwind CSS imports
index.html      # HTML entry point
vite.config.js  # Vite configuration with proxy
```

## Features

- Responsive dual-pane layout
- Form validation and error handling
- Loading states and user feedback
- Accessibility features
- Copy to clipboard functionality
