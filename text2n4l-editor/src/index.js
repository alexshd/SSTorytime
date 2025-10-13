// Main exports index for easy imports

// Libraries
export { getValidArrowsList, isValidArrow } from './lib/arrows.js';
export { highlightArrows } from './lib/highlighter.js';
export { saveSession, loadSession, clearSession } from './lib/session.js';
export { detectFileType, markdownToHtml, readFileAsText } from './lib/fileUtils.js';

// Components
export { LineNumbers } from './components/LineNumbers.js';
export { createArrowMenu, showArrowValidationGuide } from './components/ArrowMenu.js';

// Usage example:
// import { highlightArrows, getValidArrowsList } from './index.js';
