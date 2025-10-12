# Text to N4L Converter Frontend

A modern web interface for converting text to N4L DSL format using Vite and Tailwind CSS v4.

## Features

- Clean, responsive UI built with Tailwind CSS v4
- Real-time text conversion via API
- Copy to clipboard functionality
- Keyboard shortcuts (Ctrl/Cmd + Enter to convert)
- Modern development stack with Vite

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
