# Session Persistence Feature

## Overview

The N4L editor now includes automatic session persistence to prevent data loss during editing. Your work is automatically saved to your browser's local storage and restored when you reload the page.

## Features

### Auto-Save

- **Automatic saving** after conversions, file uploads, and arrow edits
- **Debounced input saving** - saves 2 seconds after you stop typing
- **Visual indicator** - Shows "💾 Auto-saved" briefly when session is saved
- **No manual action required** - Everything saves automatically

### Session Restoration

- **Automatic restore on page load** - Previous session loads automatically
- **7-day retention** - Sessions older than 7 days are automatically cleared
- **Preserves everything**:
  - Input text
  - Converted N4L output (with all highlighting)
  - Uploaded filename
  - File type (HTML/Markdown/Text)

### Manual Controls

- **🗑️ Clear button** - Manually clear the current session and reset editor
- **Confirmation dialog** - Prevents accidental clearing
- **Complete reset** - Clears all content, disables buttons, resets state

## What Gets Saved

The following data is automatically saved to browser localStorage:

```javascript
{
  inputText: "Your input text...",
  outputHTML: "Converted N4L with highlighting...",
  fileName: "document.txt",
  fileType: "text" | "html" | "markdown",
  timestamp: "2025-10-12T10:30:00.000Z"
}
```

## Storage Location

Data is stored in browser's `localStorage` under key: `n4l-editor-session`

- **Storage limit**: ~5-10MB (browser dependent)
- **Privacy**: Data never leaves your browser
- **Persistence**: Survives browser restarts but not cache clearing
- **Per-domain**: Each domain has separate storage

## Large File Handling

For large files (approaching 10MB):

- Session may fail to save if localStorage is full
- Console warnings will appear if save fails
- Consider downloading work regularly with "💾 Save" button
- Browser shows warning if storage quota exceeded

## Usage Tips

1. **For long editing sessions**: Work is auto-saved continuously
2. **For large files**: Download periodically with Save button as backup
3. **Switching browsers**: Session is browser-specific, not synced
4. **Privacy mode**: Session may not persist in incognito/private mode
5. **Cache clearing**: Will delete saved sessions

## Technical Details

### Save Triggers

- After text conversion (immediate)
- After file upload (immediate)
- After arrow replacement (immediate)
- After arrow deletion (immediate)
- During input typing (2-second delay)
- During output editing (2-second delay)

### Session Validation

- Checks age on load (max 7 days)
- Validates JSON structure
- Handles corrupted data gracefully
- Falls back to empty state on error

### Browser Compatibility

- Chrome/Edge: ✅ Full support
- Firefox: ✅ Full support
- Safari: ✅ Full support
- Mobile browsers: ✅ Usually supported

## Troubleshooting

### Session Not Restoring

1. Check if more than 7 days have passed
2. Check browser's localStorage settings
3. Verify not in private/incognito mode
4. Check console for error messages

### Storage Full Error

1. Clear old sessions from other sites
2. Download current work with Save button
3. Clear session with 🗑️ Clear button
4. Reduce file size if possible

### Data Lost

- If browser cache was cleared, data cannot be recovered
- Always use "💾 Save" button for critical work
- Consider version control for important documents
