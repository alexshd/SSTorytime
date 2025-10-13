import './style.css';
import { saveSession, loadSession, clearSession } from './lib/session.js';
import { highlightArrows } from './lib/highlighter.js';
import { detectFileType, markdownToHtml, readFileAsText } from './lib/fileUtils.js';
import { N4LEditor } from './components/N4LEditor.js';
import { createArrowMenu, showArrowValidationGuide } from './components/ArrowMenu.js';

// Global state
let editorMode = 'landing';
let n4lEditor = null; // CodeMirror editor instance
let state = {
  currentFileName: '',
  currentFileType: 'text',
  isScrollSyncing: false,
  showingPreview: false
};

// Check if we have a saved session with mode - restore editor directly
const savedSession = loadSession();
if (savedSession?.mode)
{
  editorMode = savedSession.mode;
  loadEditor();
}
else
{
  // No saved session - show landing page
  initLanding();
}

function initLanding()
{
  fetch('/src/landing.html')
    .then(res => res.text())
    .then(html =>
    {
      document.querySelector('#app').innerHTML = html;

      document.querySelector('#option-convert').addEventListener('click', () =>
      {
        editorMode = 'convert';
        loadEditor();
      });

      document.querySelector('#option-edit').addEventListener('click', () =>
      {
        editorMode = 'edit';
        loadEditor();
      });
    });
}

function loadEditor()
{
  fetch('/src/app.html')
    .then(res => res.text())
    .then(html =>
    {
      document.querySelector('#app').innerHTML = html;
      initEditor();
    })
    .catch(error =>
    {
      console.error('Error loading app:', error);
      document.getElementById('app').innerHTML = '<p>Error loading application</p>';
    });
}

function initEditor()
{
  // DOM elements
  const elements = {
    inputText: document.querySelector('#input-text'),
    inputPreview: document.querySelector('#input-preview'),
    outputEditorContainer: document.querySelector('#output-editor-container'),
    convertBtn: document.querySelector('#convert-btn'),
    copyBtn: document.querySelector('#copy-btn'),
    uploadBtn: document.querySelector('#upload-btn'),
    fileInput: document.querySelector('#file-input'),
    arrowsBtn: document.querySelector('#arrows-btn'),
    saveBtn: document.querySelector('#save-btn'),
    clearSessionBtn: document.querySelector('#clear-session-btn'),
    sessionIndicator: document.querySelector('#session-indicator'),
    statusMessage: document.querySelector('#status-message'),
    titleSection: document.querySelector('#title-section'),
  };

  // Initialize CodeMirror editor
  n4lEditor = new N4LEditor(elements.outputEditorContainer);

  // Set up onChange handler for auto-save
  n4lEditor.setOnChange(() =>
  {
    saveCurrentSession(elements);
  });

  // Set up mode-specific UI
  setupModeUI(elements);

  // Event listeners
  setupEventListeners(elements);

  // Global arrow menu functions
  setupGlobalArrowFunctions(elements);

  // Restore session
  restoreSession(elements);

  // Auto-save
  setInterval(() => saveCurrentSession(elements), 2000);

  // Initial focus
  elements.inputText.focus();
}

function setupModeUI(elements)
{
  if (editorMode === 'edit')
  {
    elements.titleSection.innerHTML = '<h1 class="text-2xl font-bold text-gray-800">✏️ N4L Editor</h1><p class="text-sm text-gray-600">Upload a .n4l file to edit and validate</p>';
    elements.uploadBtn.innerHTML = '📁 Upload N4L';
    elements.convertBtn.classList.add('hidden');
  } else
  {
    elements.titleSection.innerHTML = '<h1 class="text-2xl font-bold text-gray-800">▶️ Text to N4L Converter</h1><p class="text-sm text-gray-600">Upload text, HTML, or Markdown to convert</p>';
    elements.uploadBtn.innerHTML = '📁 Upload File';
  }
}

function setupEventListeners(elements)
{
  elements.convertBtn.addEventListener('click', () => convertText(elements));
  elements.copyBtn.addEventListener('click', () => copyToClipboard(elements));
  elements.uploadBtn.addEventListener('click', () => elements.fileInput.click());
  elements.fileInput.addEventListener('change', (e) => handleFileSelect(e, elements));
  elements.arrowsBtn.addEventListener('click', showArrowValidationGuide);
  elements.saveBtn.addEventListener('click', () => saveAsN4L(elements));
  elements.clearSessionBtn.addEventListener('click', () => clearSessionData(elements));

  // Scroll sync
  elements.outputArea.addEventListener('scroll', () => syncScroll(elements.outputArea, elements));
  elements.inputText.addEventListener('scroll', () => syncScroll(elements.inputText, elements));
  elements.inputPreview.addEventListener('scroll', () => syncScroll(elements.inputPreview, elements));

  // Keyboard shortcuts
  elements.inputText.addEventListener('keydown', (e) =>
  {
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') convertText(elements);
  });

  elements.outputArea.addEventListener('keydown', (e) =>
  {
    if ((e.ctrlKey || e.metaKey) && e.key === 's')
    {
      e.preventDefault();
      saveAsN4L(elements);
    }
  });

  // Update highlighting on edit
  elements.outputArea.addEventListener('input', () =>
  {
    const selection = window.getSelection();
    const range = selection.rangeCount > 0 ? selection.getRangeAt(0) : null;
    const offset = range ? range.startOffset : 0;

    const plainText = elements.outputArea.innerText;
    elements.outputArea.innerHTML = highlightArrows(plainText);

    // Update line numbers with source text
    if (lineNumbers)
    {
      lineNumbers.updateFromSource(plainText);
    }

    // Restore cursor position (simplified)
    if (range)
    {
      try
      {
        const newRange = document.createRange();
        const textNode = findTextNode(elements.outputArea, offset);
        if (textNode)
        {
          newRange.setStart(textNode, Math.min(offset, textNode.length));
          newRange.collapse(true);
          selection.removeAllRanges();
          selection.addRange(newRange);
        }
      } catch (e)
      {
        console.warn('Could not restore cursor position');
      }
    }
  });
}

function setupGlobalArrowFunctions(elements)
{
  window.showArrowMenu = function (event, arrowSpan)
  {
    event.stopPropagation();

    // Remove existing menu
    const existingMenu = document.querySelector('.arrow-menu');
    if (existingMenu) existingMenu.remove();

    const menu = createArrowMenu(
      arrowSpan,
      (newArrow, menuEl) => replaceArrow(arrowSpan, newArrow, menuEl, elements),
      (arrow, menuEl) => deleteArrow(arrowSpan, menuEl, elements)
    );

    // Position menu
    const rect = arrowSpan.getBoundingClientRect();
    menu.style.left = Math.min(rect.left, window.innerWidth - 320) + 'px';
    menu.style.top = (rect.bottom + 5) + 'px';

    document.body.appendChild(menu);

    // Close on outside click
    setTimeout(() =>
    {
      document.addEventListener('click', function closeMenu()
      {
        menu.remove();
        document.removeEventListener('click', closeMenu);
      });
    }, 100);
  };
}

async function convertText(elements)
{
  const inputEl = state.showingPreview ? elements.inputPreview : elements.inputText;
  const text = inputEl.value || inputEl.innerText || inputEl.textContent;

  if (!text.trim())
  {
    showStatus('Please enter some text to convert', true, elements);
    return;
  }

  elements.convertBtn.disabled = true;
  showStatus('Converting...', false, elements);

  try
  {
    // Send as form data, not JSON
    const formData = new URLSearchParams();
    formData.append('text', text);

    const response = await fetch('http://localhost:8080/api/convert', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: formData
    });

    if (!response.ok) throw new Error('Conversion failed');

    const n4lText = await response.text();
    console.log('Received N4L text length:', n4lText.length);
    console.log('First 200 chars:', n4lText.substring(0, 200));

    const highlighted = highlightArrows(n4lText);
    console.log('Highlighted HTML length:', highlighted.length);
    console.log('Output area element:', elements.outputArea);

    elements.outputArea.innerHTML = highlighted;
    console.log('innerHTML set, content length:', elements.outputArea.innerHTML.length);

    // Update line numbers with source text
    if (lineNumbers)
    {
      lineNumbers.updateFromSource(n4lText);
    }

    elements.copyBtn.disabled = false;
    elements.saveBtn.disabled = false;

    showStatus('Conversion complete!', false, elements);
    saveCurrentSession(elements);
  } catch (error)
  {
    showStatus('Error: ' + error.message, true, elements);
  } finally
  {
    elements.convertBtn.disabled = false;
  }
}

async function copyToClipboard(elements)
{
  const text = elements.outputArea.innerText || elements.outputArea.textContent;
  try
  {
    await navigator.clipboard.writeText(text);
    showStatus('Copied to clipboard!', false, elements);
  } catch (error)
  {
    showStatus('Failed to copy', true, elements);
  }
}

async function handleFileSelect(event, elements)
{
  const file = event.target.files[0];
  if (!file) return;

  if (file.size > 10 * 1024 * 1024)
  {
    showStatus('File too large. Max 10MB.', true, elements);
    return;
  }

  if (editorMode === 'edit' && !file.name.endsWith('.n4l'))
  {
    showStatus('Please select a .n4l file', true, elements);
    event.target.value = '';
    return;
  }

  try
  {
    const content = await readFileAsText(file);
    state.currentFileName = file.name;
    state.currentFileType = detectFileType(file.name, content);

    if (editorMode === 'edit')
    {
      // Load directly into output
      elements.outputArea.innerHTML = highlightArrows(content);

      // Update line numbers with source text
      if (lineNumbers)
      {
        lineNumbers.updateFromSource(content);
      }

      elements.inputText.value = '';
      elements.copyBtn.disabled = false;
      elements.saveBtn.disabled = false;
      elements.titleSection.innerHTML = `<h1 class="text-2xl font-bold text-gray-800">✏️ Editing: ${file.name}</h1>`;
      showStatus(`N4L file "${file.name}" loaded!`, false, elements);
      elements.outputArea.focus();
    } else
    {
      // Load for conversion
      renderFileContent(content, state.currentFileType, elements);
      elements.titleSection.classList.add('hidden');
      const typeLabel = state.currentFileType === 'html' ? 'HTML' : state.currentFileType === 'markdown' ? 'Markdown' : 'Text';
      showStatus(`File "${file.name}" loaded (${typeLabel})`, false, elements);
    }

    saveCurrentSession(elements);
  } catch (error)
  {
    showStatus('Error reading file', true, elements);
  } finally
  {
    event.target.value = '';
  }
}

function renderFileContent(content, type, elements)
{
  if (type === 'html')
  {
    elements.inputPreview.innerHTML = content;
    elements.inputPreview.classList.remove('hidden');
    elements.inputText.classList.add('hidden');
    state.showingPreview = true;
  } else if (type === 'markdown')
  {
    elements.inputPreview.innerHTML = markdownToHtml(content);
    elements.inputPreview.classList.remove('hidden');
    elements.inputText.classList.add('hidden');
    state.showingPreview = true;
  } else
  {
    elements.inputText.value = content;
    elements.inputText.classList.remove('hidden');
    elements.inputPreview.classList.add('hidden');
    state.showingPreview = false;
  }
}

function replaceArrow(arrowSpan, newArrow, menu, elements)
{
  const oldArrow = arrowSpan.getAttribute('data-arrow');
  const text = elements.outputArea.innerText;
  const updatedText = text.replace(oldArrow, newArrow);
  elements.outputArea.innerHTML = highlightArrows(updatedText);

  // Update line numbers with source text
  if (lineNumbers)
  {
    lineNumbers.updateFromSource(updatedText);
  }

  menu.remove();
  saveCurrentSession(elements);
  showStatus(`Replaced "${oldArrow}" with "${newArrow}"`, false, elements);
}

function deleteArrow(arrowSpan, menu, elements)
{
  const arrow = arrowSpan.getAttribute('data-arrow');
  const text = elements.outputArea.innerText;
  const updatedText = text.replace(new RegExp(escapeRegex(arrow), 'g'), '');
  elements.outputArea.innerHTML = highlightArrows(updatedText);

  // Update line numbers with source text
  if (lineNumbers)
  {
    lineNumbers.updateFromSource(updatedText);
  }

  menu.remove();
  saveCurrentSession(elements);
  showStatus(`Deleted "${arrow}"`, false, elements);
}

function saveAsN4L(elements)
{
  const content = elements.outputArea.innerText || elements.outputArea.textContent;
  const filename = state.currentFileName || 'output.n4l';
  const blob = new Blob([content], { type: 'text/plain' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename.endsWith('.n4l') ? filename : filename + '.n4l';
  a.click();
  URL.revokeObjectURL(url);
  showStatus('File saved!', false, elements);
}

function saveCurrentSession(elements)
{
  const sessionData = {
    inputText: elements.inputText.value,
    outputHTML: elements.outputArea.innerHTML,
    fileName: state.currentFileName,
    fileType: state.currentFileType,
    mode: editorMode
  };

  if (saveSession(sessionData))
  {
    elements.sessionIndicator.classList.remove('hidden');
    setTimeout(() => elements.sessionIndicator.classList.add('hidden'), 2000);
  }
}

function restoreSession(elements)
{
  const session = loadSession();
  if (!session) return;

  // Restore mode if saved
  if (session.mode)
  {
    editorMode = session.mode;
  }

  elements.inputText.value = session.inputText || '';
  elements.outputArea.innerHTML = session.outputHTML || '';

  // Update line numbers from restored output
  if (lineNumbers && session.outputHTML)
  {
    // Extract text content from the HTML to count lines
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = session.outputHTML;
    const plainText = tempDiv.textContent || tempDiv.innerText || '';
    lineNumbers.updateFromSource(plainText);
  }

  state.currentFileName = session.fileName || '';
  state.currentFileType = session.fileType || 'text';

  if (session.outputHTML?.trim())
  {
    elements.copyBtn.disabled = false;
    elements.saveBtn.disabled = false;
  }

  if (session.fileName)
  {
    elements.titleSection.querySelector('h1').textContent = session.fileName;
  }

  showStatus('Session restored (auto-save enabled)', false, elements);
}

function clearSessionData(elements)
{
  if (confirm('Clear session and return to start page?'))
  {
    clearSession();
    editorMode = 'landing';

    // Reload the page to show landing page
    window.location.reload();
  }
}

function syncScroll(source, elements)
{
  if (state.isScrollSyncing) return;
  state.isScrollSyncing = true;

  const isOutput = source === elements.outputArea;
  const inputEl = state.showingPreview ? elements.inputPreview : elements.inputText;
  const sourceEl = isOutput ? elements.outputArea : inputEl;
  const targetEl = isOutput ? inputEl : elements.outputArea;

  const sourceScrollHeight = sourceEl.scrollHeight - sourceEl.clientHeight;
  const scrollRatio = sourceScrollHeight > 0 ? sourceEl.scrollTop / sourceScrollHeight : 0;
  const targetScrollHeight = targetEl.scrollHeight - targetEl.clientHeight;
  targetEl.scrollTop = scrollRatio * targetScrollHeight;

  setTimeout(() => { state.isScrollSyncing = false; }, 50);
}

function showStatus(message, isError, elements)
{
  const p = elements.statusMessage.querySelector('p');
  p.textContent = message;
  p.style.color = isError ? '#dc2626' : '#059669';
  elements.statusMessage.classList.remove('hidden');
  setTimeout(() => elements.statusMessage.classList.add('hidden'), 3000);
}

function findTextNode(node, offset)
{
  if (node.nodeType === Node.TEXT_NODE)
  {
    return offset <= node.length ? node : null;
  }

  for (let child of node.childNodes)
  {
    const found = findTextNode(child, offset);
    if (found) return found;
    if (child.nodeType === Node.TEXT_NODE)
    {
      offset -= child.length;
    }
  }

  return null;
}

function escapeRegex(str)
{
  return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
