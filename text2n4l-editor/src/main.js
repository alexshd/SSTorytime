import './style.css';
import { saveSession, loadSession, clearSession } from './lib/session.js';
import { detectFileType, markdownToHtml, readFileAsText } from './lib/fileUtils.js';
import { N4LEditor } from './components/N4LEditor.js';
import { showArrowValidationGuide } from './components/ArrowMenu.js';

// Global state
let editorMode = 'landing';
let n4lEditor = null;
let state = {
  currentFileName: '',
  currentFileType: 'text',
  showingPreview: false,
  isScrollSyncing: false
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

  // Initialize CodeMirror editor with proper line numbers
  n4lEditor = new N4LEditor(elements.outputEditorContainer);

  // Set up onChange handler for auto-save
  n4lEditor.setOnChange(() =>
  {
    saveCurrentSession(elements);
  });

  // Set up arrow action callbacks
  n4lEditor.setArrowCallbacks(
    (oldArrow, newArrow) =>
    {
      showStatus(`Replaced "${oldArrow}" with "${newArrow}"`, false, elements);
      saveCurrentSession(elements);
    },
    (arrow) =>
    {
      showStatus(`Deleted "${arrow}"`, false, elements);
      saveCurrentSession(elements);
    }
  );

  // Set up scroll sync from output editor to input
  n4lEditor.onScroll(() => syncScrollFromOutput(elements));

  // Set up mode-specific UI
  setupModeUI(elements);

  // Event listeners
  setupEventListeners(elements);

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
    elements.titleSection.innerHTML = '<h1 class="text-2xl font-bold text-gray-800">✏️ N4L Editor</h1><p class="text-sm text-gray-600">Upload a .n4l file to edit and validate - Line numbers work correctly!</p>';
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

  // Scroll sync between input and output
  elements.inputText.addEventListener('scroll', () => syncScrollFromInput(elements.inputText, elements));
  elements.inputPreview.addEventListener('scroll', () => syncScrollFromInput(elements.inputPreview, elements));

  // Keyboard shortcuts
  elements.inputText.addEventListener('keydown', (e) =>
  {
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') convertText(elements);
  });
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
    const formData = new URLSearchParams();
    formData.append('text', text);

    const response = await fetch('http://localhost:8080/api/convert', {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: formData
    });

    if (!response.ok) throw new Error('Conversion failed');

    const n4lText = await response.text();

    // Set content in CodeMirror editor
    n4lEditor.setContent(n4lText);

    elements.copyBtn.disabled = false;
    elements.saveBtn.disabled = false;

    showStatus('Conversion complete! Line numbers show actual lines.', false, elements);
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
  const text = n4lEditor.getContent();
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
      // Load directly into CodeMirror editor
      n4lEditor.setContent(content);

      elements.inputText.value = '';
      elements.copyBtn.disabled = false;
      elements.saveBtn.disabled = false;
      elements.titleSection.innerHTML = `<h1 class="text-2xl font-bold text-gray-800">✏️ Editing: ${file.name}</h1><p class="text-sm text-gray-600">Line numbers correspond to actual lines in file</p>`;
      showStatus(`N4L file "${file.name}" loaded!`, false, elements);
      n4lEditor.focus();
    } else
    {
      // Load for conversion
      renderFileContent(content, state.currentFileType, elements);
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

function saveAsN4L(elements)
{
  const content = n4lEditor.getContent();
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
    outputText: n4lEditor ? n4lEditor.getContent() : '',
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

  // Restore output content in CodeMirror
  if (session.outputText)
  {
    n4lEditor.setContent(session.outputText);
  }

  state.currentFileName = session.fileName || '';
  state.currentFileType = session.fileType || 'text';

  if (session.outputText?.trim())
  {
    elements.copyBtn.disabled = false;
    elements.saveBtn.disabled = false;
  }

  if (session.fileName)
  {
    elements.titleSection.innerHTML = `<h1 class="text-2xl font-bold text-gray-800">✏️ Editing: ${session.fileName}</h1><p class="text-sm text-gray-600">Line numbers correspond to actual lines in file</p>`;
  }

  showStatus('Session restored (auto-save enabled)', false, elements);
}

function clearSessionData(elements)
{
  if (confirm('Clear session and return to start page?'))
  {
    clearSession();
    editorMode = 'landing';
    window.location.reload();
  }
}

function showStatus(message, isError, elements)
{
  const p = elements.statusMessage.querySelector('p');
  p.textContent = message;
  p.style.color = isError ? '#dc2626' : '#059669';
  elements.statusMessage.classList.remove('hidden');
  setTimeout(() => elements.statusMessage.classList.add('hidden'), 3000);
}

// Scroll sync: from input to output
function syncScrollFromInput(inputElement, elements)
{
  if (state.isScrollSyncing) return;
  state.isScrollSyncing = true;

  const sourceScrollHeight = inputElement.scrollHeight - inputElement.clientHeight;
  const scrollRatio = sourceScrollHeight > 0 ? inputElement.scrollTop / sourceScrollHeight : 0;

  // Sync to output editor
  if (n4lEditor)
  {
    n4lEditor.scrollToRatio(scrollRatio);
  }

  setTimeout(() => { state.isScrollSyncing = false; }, 50);
}

// Scroll sync: from output to input
function syncScrollFromOutput(elements)
{
  if (state.isScrollSyncing) return;
  state.isScrollSyncing = true;

  const inputEl = state.showingPreview ? elements.inputPreview : elements.inputText;
  const scrollRatio = n4lEditor.getScrollRatio();

  const targetScrollHeight = inputEl.scrollHeight - inputEl.clientHeight;
  inputEl.scrollTop = scrollRatio * targetScrollHeight;

  setTimeout(() => { state.isScrollSyncing = false; }, 50);
}
