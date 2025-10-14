# Text2N4L Editor - Code Structure

## 📁 Project Structure

```
src/
├── main.js                 # Main application entry point (~400 lines)
├── style.css              # Global styles
├── app.html               # Editor layout
├── landing.html           # Landing page
├── lib/                   # Core libraries
│   ├── arrows.js          # Arrow validation & list management
│   ├── highlighter.js     # N4L syntax highlighting engine
│   ├── session.js         # Session persistence (localStorage)
│   └── fileUtils.js       # File reading, type detection, markdown
└── components/            # Reusable UI components
    ├── LineNumbers.js     # Line numbers display & sync
    └── ArrowMenu.js       # Arrow validation/edit modal
```

## 🧩 Module Responsibilities

### **main.js** (Entry Point)

- Application initialization
- Landing page → Editor routing
- DOM event management
- Orchestrates all modules

### **lib/arrows.js**

- `getValidArrowsList()` - Returns 300+ valid N4L arrows from SSTconfig
- `isValidArrow(text)` - Validates if text is a registered arrow

### **lib/highlighter.js**

- `highlightArrows(text)` - Main N4L syntax highlighting
- Context-aware: Distinguishes structural syntax from content blocks
- Highlights: comments, titles, sequences, references, arrows, etc.

### **lib/session.js**

- `saveSession(data)` - Persist editor state to localStorage
- `loadSession()` - Restore previous session (7-day expiry)
- `clearSession()` - Clear saved session data

### **lib/fileUtils.js**

- `detectFileType(filename, content)` - Auto-detect file format
- `markdownToHtml(markdown)` - Simple markdown → HTML converter
- `readFileAsText(file)` - Promise-based file reader

### **components/LineNumbers.js**

- `LineNumbers` class - Manages line number display
- Auto-updates on content changes (MutationObserver)
- Synchronized scrolling with output area

### **components/ArrowMenu.js**

- `createArrowMenu()` - Interactive arrow edit menu
- `showArrowValidationGuide()` - Modal with usage instructions
- Arrow suggestions based on current text
- Replace/delete functionality

## 🔄 Data Flow

```
Landing Page
    ↓ (select mode)
Editor Init
    ↓
Load Session → Restore State
    ↓
User Actions:
  • Upload File → fileUtils → Render
  • Convert Text → API → highlighter → Display
  • Edit Arrow → ArrowMenu → Update → Save Session
  • Edit Text → highlighter → LineNumbers Update
    ↓
Auto-save (every 2s) → session.js
```

## 🎯 Key Features

### Modular Benefits

- ✅ **Separation of Concerns** - Each module has single responsibility
- ✅ **Reusability** - Components can be used independently
- ✅ **Testability** - Pure functions, easy to unit test
- ✅ **Maintainability** - Changes isolated to specific modules
- ✅ **Readability** - Small focused files vs 1300+ line monolith

### Code Reduction

- **Before**: ~1371 lines in main.js
- **After**: ~400 lines in main.js + organized modules
- **70% reduction** in main file complexity

## 🛠️ Development

### Adding New Features

**1. New Arrow Type:**
Edit `lib/arrows.js` → Add to `getValidArrowsList()`

**2. New Syntax Highlighting:**
Edit `lib/highlighter.js` → Add pattern in `highlightArrows()`

**3. New UI Component:**
Create new file in `components/` → Import in `main.js`

**4. New File Format:**
Edit `lib/fileUtils.js` → Extend `detectFileType()` & add converter

### Testing Individual Modules

```javascript
// Test arrows
import { isValidArrow } from "./lib/arrows.js";
console.log(isValidArrow("(leads to)")); // true

// Test highlighter
import { highlightArrows } from "./lib/highlighter.js";
const html = highlightArrows("test (leads to) result");

// Test session
import { saveSession, loadSession } from "./lib/session.js";
saveSession({ test: "data" });
console.log(loadSession());
```

## 📝 Code Style

- **ES6 Modules** - `import`/`export` for clean dependencies
- **Pure Functions** - Where possible (highlighter, validators)
- **Class Components** - For stateful UI (LineNumbers)
- **Factory Functions** - For dynamic UI (ArrowMenu)
- **No frameworks** - Vanilla JS for minimal dependencies

## 🚀 Future Improvements

- [ ] Add JSDoc comments for better IDE support
- [ ] Unit tests for each module
- [ ] Async loading of arrow list (reduce initial bundle)
- [ ] Web Worker for heavy highlighting operations
- [ ] Virtual scrolling for large files
- [ ] Plugin system for custom highlighters

## 🔗 Dependencies

- **Zero npm dependencies** for core logic
- Vite for dev server & bundling
- Tailwind CSS v4 for styling
- Open Props for design tokens
