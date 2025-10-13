import './style.css';

// Global mode for the editor
let editorMode = 'landing'; // 'landing', 'convert', or 'edit'

// Load landing page first
fetch('/src/landing.html')
  .then(res => res.text())
  .then(html =>
  {
    document.querySelector('#app').innerHTML = html;

    // Set up landing page option handlers
    const optionConvert = document.querySelector('#option-convert');
    const optionEdit = document.querySelector('#option-edit');

    optionConvert.addEventListener('click', () =>
    {
      editorMode = 'convert';
      loadEditor();
    });

    optionEdit.addEventListener('click', () =>
    {
      editorMode = 'edit';
      loadEditor();
    });
  });

// Load the main editor
function loadEditor()
{
  fetch('/src/app.html')
    .then(res => res.text())
    .then(html =>
    {
      document.querySelector('#app').innerHTML = html;

      // Select DOM elements
      const inputText = document.querySelector('#input-text');
      const inputPreview = document.querySelector('#input-preview');
      const outputArea = document.querySelector('#output-area');
      const lineNumbers = document.querySelector('#line-numbers');
      const convertBtn = document.querySelector('#convert-btn');
      const copyBtn = document.querySelector('#copy-btn');
      const uploadBtn = document.querySelector('#upload-btn');
      const fileInput = document.querySelector('#file-input');
      const arrowsBtn = document.querySelector('#arrows-btn');
      const saveBtn = document.querySelector('#save-btn');
      const clearSessionBtn = document.querySelector('#clear-session-btn');
      const sessionIndicator = document.querySelector('#session-indicator');
      const statusMessage = document.querySelector('#status-message');
      const titleSection = document.querySelector('#title-section');
      const arrowsSection = document.querySelector('#arrows-section');
      const arrowsContent = document.querySelector('#arrows-content');

      let currentFileName = '';
      let currentFileType = 'text'; // 'text', 'html', or 'markdown'
      let isScrollSyncing = false;
      let showingPreview = false;

      // === LINE NUMBERS ===

      function updateLineNumbers()
      {
        const text = outputArea.innerText || outputArea.textContent || '';
        const lines = text.split('\n');
        const lineCount = lines.length;

        let lineNumbersHTML = '';
        for (let i = 1; i <= lineCount; i++)
        {
          lineNumbersHTML += '<div>' + i + '</div>';
        }
        lineNumbers.innerHTML = lineNumbersHTML;
      }

      function syncLineNumbersScroll()
      {
        lineNumbers.scrollTop = outputArea.scrollTop;
      }

      // Update line numbers when content changes
      const observerConfig = { childList: true, subtree: true, characterData: true };
      const lineNumberObserver = new MutationObserver(() =>
      {
        updateLineNumbers();
      });
      lineNumberObserver.observe(outputArea, observerConfig);

      // Sync scroll between output area and line numbers
      outputArea.addEventListener('scroll', syncLineNumbersScroll);

      // Initial line numbers
      updateLineNumbers();

      // === SESSION PERSISTENCE ===

      function saveSession()
      {
        try
        {
          const sessionData = {
            inputText: inputText.value,
            outputHTML: outputArea.innerHTML,
            fileName: currentFileName,
            fileType: currentFileType,
            timestamp: new Date().toISOString()
          };
          localStorage.setItem('n4l-editor-session', JSON.stringify(sessionData));

          // Show indicator briefly
          sessionIndicator.classList.remove('hidden');
          setTimeout(() =>
          {
            sessionIndicator.classList.add('hidden');
          }, 2000);
        }
        catch (error)
        {
          console.warn('Could not save session:', error);
        }
      }

      function loadSession()
      {
        try
        {
          const saved = localStorage.getItem('n4l-editor-session');
          if (!saved) return false;

          const sessionData = JSON.parse(saved);

          // Check if session is less than 7 days old
          const sessionAge = Date.now() - new Date(sessionData.timestamp).getTime();
          const sevenDays = 7 * 24 * 60 * 60 * 1000;

          if (sessionAge > sevenDays)
          {
            localStorage.removeItem('n4l-editor-session');
            return false;
          }

          // Restore session
          inputText.value = sessionData.inputText || '';
          outputArea.innerHTML = sessionData.outputHTML || '';
          currentFileName = sessionData.fileName || '';
          currentFileType = sessionData.fileType || 'text';

          // Enable buttons if there's output
          if (outputArea.innerHTML.trim())
          {
            copyBtn.disabled = false;
            saveBtn.disabled = false;
          }

          // Update title if we have a filename
          if (currentFileName)
          {
            const titleText = titleSection.querySelector('h1');
            if (titleText) titleText.textContent = currentFileName;
          }

          showStatus('Previous session restored (auto-save enabled)');
          return true;
        }
        catch (error)
        {
          console.warn('Could not load session:', error);
          localStorage.removeItem('n4l-editor-session');
          return false;
        }
      }

      function clearSessionData()
      {
        if (confirm('Clear the current session? This will reset the editor and remove saved work.'))
        {
          localStorage.removeItem('n4l-editor-session');

          // Clear the editor
          inputText.value = '';
          outputArea.innerHTML = '';
          currentFileName = '';
          currentFileType = 'text';

          // Disable buttons
          copyBtn.disabled = true;
          saveBtn.disabled = true;

          // Hide sections
          arrowsSection.classList.add('hidden');
          titleSection.classList.remove('hidden');

          showStatus('Session cleared');
        }
      }

      // === ALL FUNCTION DEFINITIONS MUST BE HERE ===

      function showStatus(message, isError = false)
      {
        const p = statusMessage.querySelector('p');
        p.textContent = message;
        p.className = 'text-sm ' + (isError ? 'text-red-600' : 'text-green-600');
        statusMessage.classList.remove('hidden');
        setTimeout(() => { statusMessage.classList.add('hidden'); }, 3000);
      }

      async function convertText()
      {
        const text = inputText.value.trim();

        if (!text)
        {
          showStatus('Please enter some text to convert', true);
          return;
        }

        convertBtn.disabled = true;
        convertBtn.textContent = 'Converting...';

        try
        {
          const formData = new FormData();
          formData.append('text', text);

          const response = await fetch('/api/convert', {
            method: 'POST',
            body: formData
          });

          if (!response.ok)
          {
            throw new Error('API Error: ' + response.status);
          }

          const n4lOutput = await response.text();
          outputArea.innerHTML = highlightArrows(n4lOutput);
          copyBtn.disabled = false;
          saveBtn.disabled = false;
          showStatus('Text converted successfully!');

          // Save session after successful conversion
          saveSession();

        } catch (error)
        {
          console.error('Conversion error:', error);
          showStatus('Error: ' + error.message, true);
          outputArea.innerHTML = '';
          copyBtn.disabled = true;
          saveBtn.disabled = true;
        } finally
        {
          convertBtn.disabled = false;
          convertBtn.textContent = 'Convert to N4L';
        }
      }

      async function copyToClipboard()
      {
        try
        {
          await navigator.clipboard.writeText(outputArea.innerText);
          showStatus('Copied to clipboard!');
        } catch (error)
        {
          console.error('Copy error:', error);
          showStatus('Failed to copy to clipboard', true);
        }
      }

      function handleFileUpload()
      {
        fileInput.click();
      }

      // Detect file type based on content and extension
      function detectFileType(filename, content)
      {
        const ext = filename.split('.').pop().toLowerCase();

        // Check by extension
        if (['html', 'htm'].includes(ext)) return 'html';
        if (['md', 'markdown'].includes(ext)) return 'markdown';

        // Check by content
        const sample = content.substring(0, 500).toLowerCase();
        if (sample.includes('<!doctype html') || sample.includes('<html') ||
          (sample.includes('<body') && sample.includes('<head')))
        {
          return 'html';
        }
        if (sample.match(/^#{1,6}\s+/m) || sample.match(/\[.*\]\(.*\)/) ||
          sample.includes('```'))
        {
          return 'markdown';
        }

        return 'text';
      }

      // Simple markdown to HTML converter
      function markdownToHtml(markdown)
      {
        let html = markdown;

        // Headers
        html = html.replace(/^### (.*$)/gim, '<h3 class="text-lg font-bold mt-4 mb-2">$1</h3>');
        html = html.replace(/^## (.*$)/gim, '<h2 class="text-xl font-bold mt-5 mb-3">$1</h2>');
        html = html.replace(/^# (.*$)/gim, '<h1 class="text-2xl font-bold mt-6 mb-4">$1</h1>');

        // Bold and italic
        html = html.replace(/\*\*\*(.+?)\*\*\*/g, '<strong><em>$1</em></strong>');
        html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
        html = html.replace(/\*(.+?)\*/g, '<em>$1</em>');
        html = html.replace(/___(.+?)___/g, '<strong><em>$1</em></strong>');
        html = html.replace(/__(.+?)__/g, '<strong>$1</strong>');
        html = html.replace(/_(.+?)_/g, '<em>$1</em>');

        // Links
        html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" class="text-blue-600 hover:underline">$1</a>');

        // Code blocks
        html = html.replace(/```([^`]+)```/g, '<pre class="bg-gray-100 p-3 rounded my-2 overflow-x-auto"><code>$1</code></pre>');
        html = html.replace(/`([^`]+)`/g, '<code class="bg-gray-100 px-1 rounded">$1</code>');

        // Lists
        html = html.replace(/^\* (.+)$/gim, '<li class="ml-4">• $1</li>');
        html = html.replace(/^\- (.+)$/gim, '<li class="ml-4">• $1</li>');
        html = html.replace(/^\d+\. (.+)$/gim, '<li class="ml-4">$1</li>');

        // Line breaks
        html = html.replace(/\n\n/g, '</p><p class="mb-2">');
        html = '<p class="mb-2">' + html + '</p>';

        return html;
      }

      // Render file content based on type
      function renderFileContent(content, type)
      {
        inputPreview.innerHTML = '';
        inputText.classList.remove('hidden');
        inputPreview.classList.add('hidden');
        showingPreview = false;

        if (type === 'html')
        {
          // Render HTML in preview
          inputPreview.innerHTML = content;
          inputPreview.classList.remove('hidden');
          inputText.classList.add('hidden');
          showingPreview = true;
          // Keep raw content in textarea for conversion
          inputText.value = content;
        }
        else if (type === 'markdown')
        {
          // Render markdown as HTML in preview
          const html = markdownToHtml(content);
          inputPreview.innerHTML = '<div class="prose max-w-none">' + html + '</div>';
          inputPreview.classList.remove('hidden');
          inputText.classList.add('hidden');
          showingPreview = true;
          // Keep raw content in textarea for conversion
          inputText.value = content;
        }
        else
        {
          // Plain text - just show in textarea
          inputText.value = content;
          inputText.classList.remove('hidden');
          inputPreview.classList.add('hidden');
        }
      }

      function handleFileSelect(event)
      {
        const file = event.target.files[0];

        if (!file)
        {
          return;
        }

        // Check file size (limit to 10MB)
        if (file.size > 10 * 1024 * 1024)
        {
          showStatus('File too large. Please select a file smaller than 10MB.', true);
          return;
        }

        // In 'edit' mode, only allow .n4l files
        if (editorMode === 'edit' && !file.name.endsWith('.n4l'))
        {
          showStatus('Please select a .n4l file for editing.', true);
          fileInput.value = '';
          return;
        }

        // We'll attempt to read any file as text
        // The browser will handle it gracefully if it's not a text file
        const reader = new FileReader();

        reader.onload = function (e)
        {
          const content = e.target.result;
          currentFileName = file.name;

          if (editorMode === 'edit')
          {
            // Edit mode: Load N4L file directly into OUTPUT area (already N4L format)
            currentFileType = 'n4l';

            // Clear input area
            inputText.value = '';
            inputText.classList.remove('hidden');
            inputPreview.classList.add('hidden');
            showingPreview = false;

            // Load N4L content into OUTPUT area with highlighting
            outputArea.innerHTML = highlightArrows(content);
            outputArea.setAttribute('contenteditable', 'true');

            // Enable copy and save buttons
            copyBtn.disabled = false;
            saveBtn.disabled = false;
            arrowsSection.classList.add('hidden');

            // Update title
            titleSection.classList.remove('hidden');
            titleSection.innerHTML = '<h1 class="text-2xl font-bold text-gray-800">✏️ Editing: ' + file.name + '</h1>';

            showStatus('N4L file "' + file.name + '" loaded for editing!');

            // Focus on output area for immediate editing
            outputArea.focus();
          }
          else
          {
            // Convert mode: Load and detect file type for conversion
            currentFileType = detectFileType(file.name, content);

            // Render content based on file type
            renderFileContent(content, currentFileType);

            // Hide title section when file is uploaded
            titleSection.classList.add('hidden');

            const typeLabel = currentFileType === 'html' ? 'HTML' :
              currentFileType === 'markdown' ? 'Markdown' : 'Text';
            showStatus('File "' + file.name + '" loaded successfully! (' + typeLabel + ' format)');

            // Clear any previous output
            outputArea.innerText = '';
            copyBtn.disabled = true;
            saveBtn.disabled = true;
            arrowsSection.classList.add('hidden');
          }

          // Save session after file load
          saveSession();

          // Reset file input
          fileInput.value = '';
        }; reader.onerror = function ()
        {
          showStatus('Error reading file. Please try again.', true);
          fileInput.value = '';
        };

        reader.readAsText(file);
      }

      // Complete list of valid N4L arrows from SSTconfig files
      // This is used for validation - if an arrow is not in this list, it's invalid
      function getValidArrowsList()
      {
        return [
          // NR-0: Similarity arrows (symmetric)
          'similar to', 'sim', 'associated with', 'ass', 'see also', 'see', 'alias',
          'equals', '=', 'compare to', 'compare', 'is not', 'not the same as',
          'looks like', 'll', 'sounds like', 'sl', 'don\'t confuse with', 'confuse',
          'alternate spelling', 'sp', 'stands for', 'sfor', 'is sometimes mistaken for', 'mistaken',
          'has same supplier', 'sm-sup', 'has common origin', 'cm-org', 'is entangled with', 'entg',
          'has same destination', 'sm-dir', 'same class', 'sm_cls', 'same group', 'sm_grp',
          'same tribe', 'sm_trb', 'same hometown', 'sm_hmt', 'same modus operandi', 'sm_MO',
          'is a variant of', 'variant', 'near to', 'nr', 'terminates together with', 'termwith',

          // LT-1: Causality arrows (directional)
          'leads to', 'fwd', 'comes from', 'bwd', 'brings about', 'brings', 'was brought about by', 'brought-by',
          'from which we derive', 'derive', 'derives from', 'derive-fr', 'leads ahead to', 'fwd1',
          'leads back to', 'bwd1', 'succeeds', 'succ1', 'is succeeded by', 'succ-by',
          'precedes', 'prec', 'is preceded by', 'prec-by', 'succeeded by', 'succ2',
          'preceded by', 'pre-by', 'comes before', 'bfr', 'comes after', 'after',
          'supplies', 'suppl', 'has been supplied by', 'suppl-by', 'evolved to', 'evolve',
          'evolved from', 'evolve-fr', 'thereafter', '=>', 'whence', '<=',
          'next if yes', 'ifyes', 'is affirmative outcome of', 'bifyes',
          'next if no', 'if no', 'is negitive outcome of', 'bifno',
          'affects', 'aff', 'is affected by', 'aff-by', 'causes', 'cause',
          'is caused by', 'cause-by', 'creates', 'cr', 'is created by', 'crtd',
          'redirects', 'redir', 'is redirected by', 'redir-by',
          'is used by', 'used-by', 'makes use of', 'uses', 'binds to', 'bind-to',
          'is bound by', 'bound-by', 'results in', 'result', 'was a result of', 'result-of',
          'enables', 'depends on', 'dep', 'is depended upon by', 'dep-by',
          'invokes', 'invoke', 'is invoked by', 'invoke-by', 'determines', 'det',

          // CN-2: Composition arrows
          'contains', 'contain', 'belongs to', 'belong', 'consists of', 'consists',
          'makes up part of', 'mkpt', 'uses word', 'useword', 'is a word used in', 'word-in',
          'uses material', 'material', 'is a material for', 'mat-for',
          'may contain', 'm-cont', 'may be contained by', 'm-cont-by',
          'has component', 'has-cmpt', 'is component of', 'cmpt-of',
          'has a part', 'has-pt', 'is a part of', 'pt-of',
          'has no part', 'no-pt', 'is not part of', 'not-pt-of',
          'has ingredient', 'ingred', 'is an ingredient of', 'ingred-of',
          'does not contain', 'has no', 'is not from set', 'is not',
          'is a set of', 'setof', 'is part of the set', 'in-set',
          'is a group of', 'grp-of', 'is an element of', 'in-grp',
          'is a herd of', 'herd', 'is part of herd', 'in-herd',
          'is a tribe of', 'tribe', 'is part of tribe', 'in-tribe',
          'is a region of', 'region', 'is within region', 'within',
          'is an interval of', 'interval', 'is part of interval', 'in-intvl',
          'species has member', 'smember', 'belongs to species', 'sbelongs',
          'is an umbrella for', 'umbrella', 'under the aegis of', 'aegis',
          'is the company of', 'company', 'works for', 'works',
          'is the employer of', 'employer', 'is an employee of', 'employee',
          'has many copies of', 'has-many', 'is one of many in', 'one-of',
          'has member', 'has-memb', 'is a member of', 'is-memb',
          'subsumes', 'sub', 'is subsumed by', 'sub-by', 'encloses', 'encl',
          'is enclosed by', 'is-encl', 'has absorbed', 'sw', 'absorbed by', 'ab-by',
          'houses', 'house', 'is housed by', 'is-housed', 'situates', 'situate',
          'is situated in', 'is-sit', 'bullet point', 'bullet', 'is a bulletpoint for', 'isbullet',
          'occurs within', 'occurs', 'situates occurrences of', 'occurence-of',
          'is a shared resource for', 'shared-by', 'shares resource', 'shares',

          // EP-3: Properties/Expresses arrows
          'note', 'remark', 'is a note or remark about', 'isnotefor',
          'added remark', 'is a remark about', 'remabout', 'please note!', 'NB',
          'is an important remark concerning', 'important regarding',
          'has example', 'e.g.', 'is an example of', 'ex', 'has no example', 'noex',
          'is not an example of', 'not-ex', 'has abbreviation', 'abbrev',
          'is short for', 'short for', 'has version', 'version', 'is a version of', 'is-version',
          'has supplier', 'supply', 'is the supplier of', 'supply-by',
          'has source', 'source', 'the source for', 'source-of',
          'discusses', 'disc', 'is discussed in', 'is-disc',
          'described as', 'descr', 'is used as a description of', 'isdescr',
          'has age character', 'age', 'is the age characteristic of', 'isageof',
          'mentions study', 'ment', 'a kind of study mentioned in', 'isment',
          'has a pun about', 'pun on', 'was made into pun', 'has pun',
          'famously uses', 'fuses', 'famously used in', 'fused',
          'has title', 'title', 'is the title of', 'title of',
          'has description', 'description', 'is the description for', 'descrip-for',
          'is about topic/them', 'about', 'is the topic/theme of', 'theme-of',
          'is wrong about', 'wrongabt', 'is not accurately represented by', 'narb',
          'expresses', 'expresses property', 'expr', 'is/are expressed by', 'expr-by',
          'has value', 'hasX', 'is the value of', 'isXof',
          'may have value', 'maybeX', 'is the possible value for', 'mXof',
          'has frequency', 'freq', 'is the frequency of', 'freq-of',
          'is characterized by', 'state', 'is a/the state of', 'state-of',
          'is expressed in formula', 'formula', 'is a formula for', 'isformula',
          'has condition', 'cond', 'is the condition of', 'cond-of',
          'has an attribute', 'attr', 'is an attribute of', 'attr-of',
          'has feature', 'feat', 'is a feature of', 'feat-of',
          'has property', 'prop', 'is a property of', 'prop-of',

          // Special arrows
          'offers', 'accepts', 'observes', 'ovserves', 'may be influenced by',
          'has overlap', 'has promiser', 'has promisee', 'observed by',
          'has intended promisee', 'has intended promiser', 'imposes'
        ];
      }

      function isValidArrow(arrowText)
      {
        // Remove parentheses and normalize
        const normalized = arrowText.replace(/^\(|\)$/g, '').trim().toLowerCase();
        const validArrows = getValidArrowsList();

        // Check if it's a valid arrow (exact match or close match)
        return validArrows.some(function (valid)
        {
          return valid.toLowerCase() === normalized;
        });
      }

      function highlightArrows(text)
      {
        // First escape HTML entities to preserve text as-is
        let escaped = text
          .replace(/&/g, '&amp;')
          .replace(/</g, '&lt;')
          .replace(/>/g, '&gt;')
          .replace(/"/g, '&quot;')
          .replace(/'/g, '&#39;');

        // Split into lines to handle line-based syntax and content blocks
        const lines = escaped.split('\n');
        const highlightedLines = [];
        let inContentBlock = false;
        let contentIndentLevel = 0;

        for (let i = 0; i < lines.length; i++)
        {
          const line = lines[i];
          let result = line;

          // Check if this line starts a content block (has @ tag)
          const hasAtTag = /^\s*@[a-zA-Z_][a-zA-Z0-9_\.]*/.test(line);

          // Check if this line is a ditto line with arrow (ends content block)
          const isDittoWithArrow = /^\s*&quot;\s+\([a-z]/.test(line);

          // Check current line indentation
          const currentIndent = line.match(/^\s*/)[0].length;

          // If we hit a ditto line with arrow, we're exiting content block
          if (inContentBlock && isDittoWithArrow)
          {
            inContentBlock = false;
          }

          // If we're in a content block
          if (inContentBlock)
          {
            // Check if we've gone back to lower indentation (exiting content block)
            if (currentIndent <= contentIndentLevel && line.trim() !== '')
            {
              inContentBlock = false;
            }
            else
            {
              // In content block - no special highlighting except @ tag if present
              if (hasAtTag)
              {
                // Highlight the @ tag only
                result = result.replace(/(@[a-zA-Z_][a-zA-Z0-9_\.]*)/g, '<span style="color: #0891b2; font-weight: 600;">$1</span>');
              }
              // Leave rest as plain text
              highlightedLines.push(result);
              continue;
            }
          }

          // If this line starts with @ tag, enter content block mode
          if (hasAtTag)
          {
            inContentBlock = true;
            contentIndentLevel = currentIndent;

            // Highlight @ tag on this line, but leave rest as is
            result = result.replace(/(@[a-zA-Z_][a-zA-Z0-9_\.]*)/g, '<span style="color: #0891b2; font-weight: 600;">$1</span>');

            // Check if there's an arrow on the same line (inline annotation)
            const hasInlineArrow = /\([a-z][a-z\s,;:.\-'\/]{3,}\)/i.test(line);
            if (hasInlineArrow)
            {
              // Apply arrow highlighting only
              result = applyArrowHighlighting(result);
            }

            highlightedLines.push(result);
            continue;
          }

          // Not in content block - apply all highlighting rules

          // 1. Highlight comments (# anything after)
          result = result.replace(/(#.*)$/g, '<span style="color: #6b7280; font-style: italic;">$1</span>');

          // 2. Highlight title markers (- at start of line followed by text)
          result = result.replace(/^(\s*-\s+)([^#]+)/g, '$1<span style="color: #7c3aed; font-weight: 600;">$2</span>');

          // 3. Highlight sequence markers (+:: and -::)
          result = result.replace(/(\+::|−::)/g, '<span style="color: #ea580c; font-weight: 700;">$1</span>');
          result = result.replace(/(-::)/g, '<span style="color: #ea580c; font-weight: 700;">$1</span>');

          // 4. Highlight $ symbols (references like $title.1)
          result = result.replace(/(\$[a-zA-Z_][a-zA-Z0-9_\.]*)/g, '<span style="color: #0891b2; font-weight: 600;">$1</span>');

          // 5. Highlight ditto symbol (standalone " at beginning of line)
          result = result.replace(/^(\s*)(&quot;)(\s)/g, '$1<span style="color: #059669; font-weight: 700; font-size: 1.1em;">$2</span>$3');

          // 6. Highlight special single-char annotations (%, =, **, >>, >, <)
          result = result.replace(/(\*\*|&gt;&gt;|&gt;|&lt;|%|=)(?=\s)/g, '<span style="color: #dc2626; font-weight: 700;">$1</span>');

          // 7. Highlight ALL CAPS text (3+ consecutive uppercase words) - reminders
          result = result.replace(/\b([A-Z][A-Z]+(?:\s+[A-Z][A-Z]+){2,})\b/g, '<span style="background: #fef3c7; color: #92400e; font-weight: 600; padding: 0 0.2em;">$1</span>');

          highlightedLines.push(result);
        }

        // Join lines back together
        const joined = highlightedLines.join('\n');

        // 8. Apply arrow highlighting to parenthetical expressions
        return applyArrowHighlighting(joined);
      }

      // Helper function to apply arrow highlighting
      function applyArrowHighlighting(text)
      {
        // Highlight parenthetical arrow descriptions like "(appears close to)"
        const parenArrowRegex = /\([a-z][a-z\s,;:.\-'\/]*\)/gi;

        return text.replace(parenArrowRegex, function (match)
        {
          // Only highlight if it's not just a single word (to avoid matching things like (x) or (a))
          // and if it contains at least one space or common relationship words
          if (match.length > 4 && (match.includes(' ') || /(?:to|by|from|in|on|at|as|with|of)/i.test(match)))
          {
            // Check if this is a valid N4L arrow
            const isValid = isValidArrow(match);
            const cssClass = isValid ? 'n4l-arrow-highlight' : 'n4l-arrow-error';

            return '<span class="' + cssClass + '" data-arrow="' + match + '" data-valid="' + isValid + '" onclick="showArrowMenu(event, this)">' +
              match +
              '</span>';
          }
          return match; // Don't highlight short or non-relational parenthetical text
        });
      }

      // Global function for arrow menu (called from onclick in HTML)
      window.showArrowMenu = function (event, arrowSpan)
      {
        event.stopPropagation();

        // Remove any existing menu
        const existingMenu = document.querySelector('.arrow-menu');
        if (existingMenu) existingMenu.remove();

        const arrow = arrowSpan.dataset.arrow;
        const isParenthetical = arrow.startsWith('(');

        // Create menu
        const menu = document.createElement('div');
        menu.className = 'arrow-menu';
        menu.style.cssText = 'position: absolute; z-index: 100; background: white; border: 1px solid #38bdf8; border-radius: 0.5rem; box-shadow: 0 4px 12px rgba(0,0,0,0.15); padding: 0.75rem; min-width: 280px; max-height: 400px; overflow-y: auto;';

        // N4L Semantic Arrow Types based on SSTconfig
        const allArrows = [
          // NR-0: Similarity (symmetric, non-directional)
          { type: 'NR-0', symbol: '(similar to)', abbrev: '(sim)', keywords: ['similar', 'like', 'resembles', 'appears close to'] },
          { type: 'NR-0', symbol: '(associated with)', abbrev: '(ass)', keywords: ['associated', 'related', 'connected'] },
          { type: 'NR-0', symbol: '(see also)', abbrev: '(see)', keywords: ['see', 'refer', 'also', 'related'] },
          { type: 'NR-0', symbol: '(alias)', abbrev: '(alias)', keywords: ['also called', 'known as', 'alias'] },
          { type: 'NR-0', symbol: '(equals)', abbrev: '(=)', keywords: ['equals', 'same as', 'equivalent', 'is'] },
          { type: 'NR-0', symbol: '(compare to)', abbrev: '(compare)', keywords: ['compare', 'contrast', 'versus'] },
          { type: 'NR-0', symbol: '(is not)', abbrev: '(!eq)', keywords: ['not', 'different', 'unlike', 'isn\'t'] },

          // LT-1: Leads To (causality, processes)
          { type: 'LT-1', symbol: '(leads to)', abbrev: '(fwd)', keywords: ['leads', 'causes', 'results', 'produces', 'brings about'] },
          { type: 'LT-1', symbol: '(causes)', abbrev: '(cause)', keywords: ['causes', 'triggers', 'initiates', 'starts'] },
          { type: 'LT-1', symbol: '(creates)', abbrev: '(cr)', keywords: ['creates', 'generates', 'makes', 'builds'] },
          { type: 'LT-1', symbol: '(results in)', abbrev: '(result)', keywords: ['results', 'ends', 'concludes', 'becomes'] },
          { type: 'LT-1', symbol: '(enables)', abbrev: '(enables)', keywords: ['enables', 'allows', 'permits', 'facilitates'] },
          { type: 'LT-1', symbol: '(affects)', abbrev: '(aff)', keywords: ['affects', 'impacts', 'influences', 'changes'] },
          { type: 'LT-1', symbol: '(precedes)', abbrev: '(prec)', keywords: ['precedes', 'before', 'prior', 'earlier'] },
          { type: 'LT-1', symbol: '(comes from)', abbrev: '(bwd)', keywords: ['from', 'derives', 'originates', 'stems'] },
          { type: 'LT-1', symbol: '(determines)', abbrev: '(det)', keywords: ['determines', 'decides', 'controls', 'defines'] },

          // CN-2: Contains (composition, part-whole)
          { type: 'CN-2', symbol: '(contains)', abbrev: '(contain)', keywords: ['contains', 'includes', 'has', 'comprises'] },
          { type: 'CN-2', symbol: '(is a part of)', abbrev: '(pt-of)', keywords: ['part of', 'within', 'inside', 'component of'] },
          { type: 'CN-2', symbol: '(consists of)', abbrev: '(consists)', keywords: ['consists', 'made up', 'composed'] },
          { type: 'CN-2', symbol: '(belongs to)', abbrev: '(belong)', keywords: ['belongs', 'member', 'element'] },
          { type: 'CN-2', symbol: '(has component)', abbrev: '(has-cmpt)', keywords: ['component', 'piece', 'section'] },
          { type: 'CN-2', symbol: '(is a set of)', abbrev: '(setof)', keywords: ['set', 'collection', 'group'] },

          // EP-3: Expresses/Properties (attributes, descriptions)
          { type: 'EP-3', symbol: '(expresses)', abbrev: '(expr)', keywords: ['expresses', 'shows', 'demonstrates', 'reveals'] },
          { type: 'EP-3', symbol: '(note)', abbrev: '(note)', keywords: ['note', 'remark', 'comment', 'annotation'] },
          { type: 'EP-3', symbol: '(defined as)', abbrev: '(def)', keywords: ['defined', 'means', 'definition', 'is defined'] },
          { type: 'EP-3', symbol: '(has example)', abbrev: '(e.g.)', keywords: ['example', 'instance', 'such as', 'for example'] },
          { type: 'EP-3', symbol: '(is about)', abbrev: '(about)', keywords: ['about', 'concerning', 'regarding', 'topic'] },
          { type: 'EP-3', symbol: '(has attribute)', abbrev: '(attr)', keywords: ['attribute', 'property', 'characteristic', 'feature'] },
          { type: 'EP-3', symbol: '(has feature)', abbrev: '(feat)', keywords: ['feature', 'aspect', 'quality'] },
          { type: 'EP-3', symbol: '(describes)', abbrev: '(descr)', keywords: ['describes', 'explains', 'details'] },

          // Special annotations
          { type: 'Special', symbol: '(is a special case of)', abbrev: '(**)', keywords: ['special case', 'exception', 'specific'] },
          { type: 'Special', symbol: '(is an example of)', abbrev: '(>>)', keywords: ['example of', 'instance of', 'exemplifies'] },
          { type: 'Special', symbol: '(discusses)', abbrev: '(%)', keywords: ['discusses', 'talks about', 'addresses'] },
          { type: 'Special', symbol: '(involves)', abbrev: '(=)', keywords: ['involves', 'includes', 'entails'] }
        ];

        let menuHTML = '';

        // Check if this arrow is invalid and show error message
        const isValid = arrowSpan.dataset.valid === 'true';
        if (!isValid && isParenthetical)
        {
          menuHTML += '<div style="background: #fee2e2; color: #991b1b; padding: 0.75rem; border-radius: 0.375rem; margin-bottom: 0.75rem; border-left: 4px solid #dc2626;">' +
            '<div style="font-weight: 700; font-size: 0.875rem; margin-bottom: 0.25rem;">⚠️ Invalid N4L Arrow</div>' +
            '<div style="font-size: 0.75rem; line-height: 1.4;">This arrow is not recognized in the N4L configuration. ' +
            'Please select a valid arrow from the list below or the N4L parser will report an error.</div>' +
            '</div>';
        }

        if (isParenthetical)
        {
          // For parenthetical text, suggest matching arrows
          const lowerArrow = arrow.toLowerCase();
          const suggestedArrows = allArrows.filter(function (arr)
          {
            return arr.keywords.some(function (keyword)
            {
              return lowerArrow.includes(keyword);
            });
          });

          menuHTML += '<div style="font-weight: 700; color: #0369a1; margin-bottom: 0.75rem; font-size: 0.875rem; border-bottom: 2px solid #38bdf8; padding-bottom: 0.5rem;">Replace with N4L Arrow:</div>';

          // Show suggested arrows first
          if (suggestedArrows.length > 0)
          {
            menuHTML += '<div style="color: #16a34a; font-size: 0.75rem; margin-bottom: 0.5rem; font-weight: 600; text-transform: uppercase;">✓ Suggested Matches:</div>';
            suggestedArrows.forEach(function (arr)
            {
              menuHTML += '<div class="arrow-menu-item" style="padding: 0.5rem; cursor: pointer; border-radius: 0.375rem; font-size: 0.875rem; margin-bottom: 0.25rem; background: #f0fdf4; border-left: 3px solid #22c55e;" ' +
                'onmouseover="this.style.background=\'#dcfce7\'" ' +
                'onmouseout="this.style.background=\'#f0fdf4\'" ' +
                'onclick="replaceArrow(\'' + arrow.replace(/'/g, "\\'") + '\', \'' + arr.symbol.replace(/'/g, "\\'") + '\', this.closest(\'.arrow-menu\'))">' +
                '<div style="display: flex; justify-content: space-between; align-items: center;">' +
                '<span style="font-family: monospace; color: #0369a1; font-weight: 600; flex: 1;">' + arr.symbol + '</span>' +
                '<span style="font-size: 0.75rem; color: #64748b; margin-left: 0.5rem;">' + arr.abbrev + '</span>' +
                '</div>' +
                '<div style="font-size: 0.7rem; color: #059669; margin-top: 0.25rem;">Type: ' + arr.type + '</div>' +
                '</div>';
            });

            menuHTML += '<div style="color: #64748b; font-size: 0.75rem; margin: 0.75rem 0 0.5rem; font-weight: 600; text-transform: uppercase;">Browse All Arrows:</div>';
          }

          // Group other arrows by type
          const arrowsByType = {};
          const otherArrows = suggestedArrows.length > 0
            ? allArrows.filter(function (arr) { return !suggestedArrows.includes(arr); })
            : allArrows;

          otherArrows.forEach(function (arr)
          {
            if (!arrowsByType[arr.type]) arrowsByType[arr.type] = [];
            arrowsByType[arr.type].push(arr);
          });

          // Display arrows grouped by type
          const typeOrder = ['NR-0', 'LT-1', 'CN-2', 'EP-3', 'Special'];
          const typeNames = {
            'NR-0': '🔗 Similarity (Non-directional)',
            'LT-1': '➡️ Causality (Leads To)',
            'CN-2': '📦 Composition (Contains)',
            'EP-3': '🏷️ Properties (Expresses)',
            'Special': '⭐ Special Annotations'
          };

          typeOrder.forEach(function (type)
          {
            if (arrowsByType[type] && arrowsByType[type].length > 0)
            {
              menuHTML += '<div style="font-size: 0.7rem; color: #475569; margin: 0.75rem 0 0.25rem; font-weight: 600;">' + typeNames[type] + '</div>';
              arrowsByType[type].forEach(function (arr)
              {
                menuHTML += '<div class="arrow-menu-item" style="padding: 0.4rem 0.5rem; cursor: pointer; border-radius: 0.25rem; font-size: 0.8rem; margin-bottom: 0.15rem;" ' +
                  'onmouseover="this.style.background=\'#e0f2fe\'" ' +
                  'onmouseout="this.style.background=\'white\'" ' +
                  'onclick="replaceArrow(\'' + arrow.replace(/'/g, "\\'") + '\', \'' + arr.symbol.replace(/'/g, "\\'") + '\', this.closest(\'.arrow-menu\'))">' +
                  '<div style="display: flex; justify-content: space-between; align-items: center;">' +
                  '<span style="font-family: monospace; color: #0369a1; font-weight: 500; font-size: 0.75rem;">' + arr.symbol + '</span>' +
                  '<span style="font-size: 0.65rem; color: #94a3b8;">' + arr.abbrev + '</span>' +
                  '</div>' +
                  '</div>';
              });
            }
          });
        }
        else
        {
          // For arrow symbols, show options to change to other arrows grouped by type
          menuHTML += '<div style="font-weight: 700; color: #0369a1; margin-bottom: 0.75rem; font-size: 0.875rem; border-bottom: 2px solid #38bdf8; padding-bottom: 0.5rem;">Change N4L Arrow:</div>';

          // Group arrows by type
          const arrowsByType = {};
          allArrows.forEach(function (arr)
          {
            if (arr.symbol !== arrow)
            {
              if (!arrowsByType[arr.type]) arrowsByType[arr.type] = [];
              arrowsByType[arr.type].push(arr);
            }
          });

          const typeOrder = ['NR-0', 'LT-1', 'CN-2', 'EP-3', 'Special'];
          const typeNames = {
            'NR-0': '🔗 Similarity (NR-0)',
            'LT-1': '➡️ Causality (LT-1)',
            'CN-2': '📦 Composition (CN-2)',
            'EP-3': '🏷️ Properties (EP-3)',
            'Special': '⭐ Special'
          };

          typeOrder.forEach(function (type)
          {
            if (arrowsByType[type] && arrowsByType[type].length > 0)
            {
              menuHTML += '<div style="font-size: 0.7rem; color: #475569; margin: 0.75rem 0 0.25rem; font-weight: 600;">' + typeNames[type] + '</div>';
              arrowsByType[type].forEach(function (arr)
              {
                menuHTML += '<div class="arrow-menu-item" style="padding: 0.4rem 0.5rem; cursor: pointer; border-radius: 0.25rem; font-size: 0.8rem; margin-bottom: 0.15rem;" ' +
                  'onmouseover="this.style.background=\'#e0f2fe\'" ' +
                  'onmouseout="this.style.background=\'white\'" ' +
                  'onclick="replaceArrow(\'' + arrow.replace(/'/g, "\\'") + '\', \'' + arr.symbol.replace(/'/g, "\\'") + '\', this.closest(\'.arrow-menu\'))">' +
                  '<div style="display: flex; justify-content: space-between; align-items: center;">' +
                  '<span style="font-family: monospace; color: #0369a1; font-weight: 500; font-size: 0.75rem;">' + arr.symbol + '</span>' +
                  '<span style="font-size: 0.65rem; color: #94a3b8;">' + arr.abbrev + '</span>' +
                  '</div>' +
                  '</div>';
              });
            }
          });
        }

        menuHTML += '<hr style="margin: 0.75rem 0; border-color: #e2e8f0;">';
        menuHTML += '<div class="arrow-menu-item" style="padding: 0.5rem; cursor: pointer; border-radius: 0.375rem; color: #dc2626; font-size: 0.875rem; font-weight: 600; text-align: center;" ' +
          'onmouseover="this.style.background=\'#fee2e2\'" ' +
          'onmouseout="this.style.background=\'white\'" ' +
          'onclick="deleteArrow(\'' + arrowSpan.dataset.arrow.replace(/'/g, "\\'") + '\', this.closest(\'.arrow-menu\'))">🗑️ Delete entire line</div>';

        menu.innerHTML = menuHTML;


        // Position menu near the arrow
        const rect = arrowSpan.getBoundingClientRect();
        menu.style.left = rect.left + 'px';
        menu.style.top = (rect.bottom + 5) + 'px';

        document.body.appendChild(menu);

        // Close menu when clicking outside
        setTimeout(function ()
        {
          document.addEventListener('click', function closeMenu()
          {
            const m = document.querySelector('.arrow-menu');
            if (m) m.remove();
            document.removeEventListener('click', closeMenu);
          });
        }, 100);
      };

      // Replace arrow in the output
      window.replaceArrow = function (oldArrow, newArrow, menu)
      {
        // Get plain text content
        const plainText = outputArea.innerText;

        // Check if arrow exists
        if (!plainText.includes(oldArrow))
        {
          menu.remove();
          showStatus('Could not find arrow: ' + oldArrow, true);
          return;
        }

        // Replace the arrow (only first occurrence)
        const newText = plainText.replace(oldArrow, newArrow);

        // Re-highlight and update
        outputArea.innerHTML = highlightArrows(newText);

        menu.remove();
        showStatus('Arrow changed to ' + newArrow);

        // Save session after change
        saveSession();
      };

      // Delete arrow from the output
      window.deleteArrow = function (arrow, menu)
      {
        const text = outputArea.innerText;
        const htmlContent = outputArea.innerHTML;

        // Escape special characters for regex
        function escapeRegex(str)
        {
          return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
        }

        // Find the line in plain text that contains this arrow
        const lines = text.split('\n');
        let lineToDelete = -1;

        for (let i = 0; i < lines.length; i++)
        {
          if (lines[i].includes(arrow))
          {
            lineToDelete = i;
            break;
          }
        }

        if (lineToDelete !== -1)
        {
          // Remove the line from the array
          lines.splice(lineToDelete, 1);

          // Rejoin the lines
          const newText = lines.join('\n');

          // Re-highlight the arrows in the new text
          outputArea.innerHTML = highlightArrows(newText);

          menu.remove();
          showStatus('Line with arrow "' + arrow + '" deleted');

          // Save session after deletion
          saveSession();
        }
        else
        {
          menu.remove();
          showStatus('Could not find line to delete', true);
        }
      };

      function extractArrows()
      {
        // Show a popup modal with the N4L Arrow Validation Guide
        showArrowValidationGuide();
      }

      function showArrowValidationGuide()
      {
        // Create modal overlay
        const modal = document.createElement('div');
        modal.className = 'fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4';
        modal.onclick = (e) =>
        {
          if (e.target === modal) modal.remove();
        };

        // Create modal content
        const modalContent = document.createElement('div');
        modalContent.className = 'bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-hidden flex flex-col';

        modalContent.innerHTML = `
        <div class="flex items-center justify-between p-4 border-b">
          <h2 class="text-xl font-bold text-gray-800">🏹 N4L Arrow Validation Guide</h2>
          <button onclick="this.closest('.fixed').remove()" class="text-gray-500 hover:text-gray-700 text-2xl leading-none">&times;</button>
        </div>
        <div class="overflow-y-auto p-6 flex-1" style="font-family: system-ui, -apple-system, sans-serif;">
          <div class="prose max-w-none">
            <p class="text-gray-600 mb-4">
              N4L uses semantic arrows to represent different types of relationships between concepts. 
              Click on any <span class="text-blue-600 font-semibold">blue highlighted arrow</span> in the output to change it, 
              or see <span class="text-red-600 font-semibold">red warnings</span> for invalid arrows.
            </p>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
              <div class="border-l-4 border-blue-500 pl-4 py-2 bg-blue-50">
                <h3 class="text-sm font-bold text-blue-900 mb-1">✓ Valid Arrow</h3>
                <code class="text-sm bg-white px-2 py-1 rounded border border-blue-200">(leads to)</code>
                <p class="text-xs text-blue-700 mt-1">Blue highlight = recognized by N4L parser</p>
              </div>
              <div class="border-l-4 border-red-500 pl-4 py-2 bg-red-50">
                <h3 class="text-sm font-bold text-red-900 mb-1">✗ Invalid Arrow</h3>
                <code class="text-sm bg-white px-2 py-1 rounded border border-red-200">(leadsto)</code>
                <p class="text-xs text-red-700 mt-1">Red highlight = will cause parser error</p>
              </div>
            </div>

            <h3 class="text-lg font-bold text-gray-800 mb-3 border-b pb-2">Four Semantic Arrow Types</h3>

            <!-- NR-0 -->
            <div class="mb-4 border rounded-lg p-4 bg-gray-50">
              <h4 class="font-bold text-gray-800 mb-2">🔗 NR-0: Similarity (Non-directional)</h4>
              <p class="text-sm text-gray-600 mb-2">Symmetric relationships, order doesn't matter</p>
              <div class="bg-white p-3 rounded border">
                <div class="grid grid-cols-2 gap-2 text-xs">
                  <code class="text-blue-600">(similar to)</code>
                  <code class="text-blue-600">(associated with)</code>
                  <code class="text-blue-600">(see also)</code>
                  <code class="text-blue-600">(alias)</code>
                  <code class="text-blue-600">(equals)</code>
                  <code class="text-blue-600">(compare to)</code>
                  <code class="text-blue-600">(is not)</code>
                  <code class="text-blue-600">(looks like)</code>
                </div>
              </div>
              <p class="text-xs text-gray-500 mt-2 italic">Example: Alice (alias) A | Person (is not) Thing</p>
            </div>

            <!-- LT-1 -->
            <div class="mb-4 border rounded-lg p-4 bg-gray-50">
              <h4 class="font-bold text-gray-800 mb-2">➡️ LT-1: Causality (Leads To)</h4>
              <p class="text-sm text-gray-600 mb-2">Directional, cause-effect, temporal ordering</p>
              <div class="bg-white p-3 rounded border">
                <div class="grid grid-cols-2 gap-2 text-xs">
                  <code class="text-blue-600">(leads to)</code>
                  <code class="text-blue-600">(causes)</code>
                  <code class="text-blue-600">(creates)</code>
                  <code class="text-blue-600">(results in)</code>
                  <code class="text-blue-600">(enables)</code>
                  <code class="text-blue-600">(affects)</code>
                  <code class="text-blue-600">(precedes)</code>
                  <code class="text-blue-600">(comes from)</code>
                </div>
              </div>
              <p class="text-xs text-gray-500 mt-2 italic">Example: Rain (leads to) Wet ground | Study (enables) Success</p>
            </div>

            <!-- CN-2 -->
            <div class="mb-4 border rounded-lg p-4 bg-gray-50">
              <h4 class="font-bold text-gray-800 mb-2">📦 CN-2: Composition (Contains)</h4>
              <p class="text-sm text-gray-600 mb-2">Part-whole, membership, hierarchical</p>
              <div class="bg-white p-3 rounded border">
                <div class="grid grid-cols-2 gap-2 text-xs">
                  <code class="text-blue-600">(contains)</code>
                  <code class="text-blue-600">(is a part of)</code>
                  <code class="text-blue-600">(consists of)</code>
                  <code class="text-blue-600">(belongs to)</code>
                  <code class="text-blue-600">(has component)</code>
                  <code class="text-blue-600">(is a set of)</code>
                </div>
              </div>
              <p class="text-xs text-gray-500 mt-2 italic">Example: Book (contains) Chapter | Wheel (is a part of) Car</p>
            </div>

            <!-- EP-3 -->
            <div class="mb-4 border rounded-lg p-4 bg-gray-50">
              <h4 class="font-bold text-gray-800 mb-2">🏷️ EP-3: Properties (Expresses)</h4>
              <p class="text-sm text-gray-600 mb-2">Attributes, descriptions, annotations</p>
              <div class="bg-white p-3 rounded border">
                <div class="grid grid-cols-2 gap-2 text-xs">
                  <code class="text-blue-600">(expresses)</code>
                  <code class="text-blue-600">(note)</code>
                  <code class="text-blue-600">(defined as)</code>
                  <code class="text-blue-600">(has example)</code>
                  <code class="text-blue-600">(is about)</code>
                  <code class="text-blue-600">(has attribute)</code>
                  <code class="text-blue-600">(describes)</code>
                </div>
              </div>
              <p class="text-xs text-gray-500 mt-2 italic">Example: Sky (has attribute) Blue | API (defined as) Interface</p>
            </div>

            <div class="bg-yellow-50 border-l-4 border-yellow-400 p-4 mt-4">
              <p class="text-sm text-yellow-800">
                <strong>💡 Pro Tip:</strong> Click any arrow in the output area to see suggestions and change it. 
                The editor validates against 300+ official N4L arrows from SSTconfig.
              </p>
            </div>
          </div>
        </div>
        <div class="border-t p-4 bg-gray-50 flex justify-end">
          <button onclick="this.closest('.fixed').remove()" 
                  class="px-6 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors">
            Got it!
          </button>
        </div>
      `;

        modal.appendChild(modalContent);
        document.body.appendChild(modal);
      }

      function displayAvailableArrows(arrows)
      {
        const categories = [...new Set(arrows.map(arrow => arrow.category))];

        let content = '<div class="space-y-4">';

        categories.forEach(category =>
        {
          const categoryArrows = arrows.filter(arrow => arrow.category === category);


          content += '<div class="category-section">' +
            '<h4 class="font-semibold text-blue-700 mb-2 text-sm">' + category + ' Arrows</h4>' +
            '<div class="grid grid-cols-1 md:grid-cols-2 gap-2">';

          categoryArrows.forEach(function (arrow)
          {
            content += '<div class="arrow-item p-3 bg-white rounded border border-gray-200 hover:border-blue-300 cursor-pointer transition-colors" ' +
              'onclick="insertArrowTemplate(\'' + arrow.symbol + '\', \' ' + arrow.example + '\')">' +
              '<div class="flex items-center justify-between mb-1">' +
              '<span class="font-mono text-lg text-blue-600">' + arrow.symbol + '</span>' +
              '<span class="text-xs text-gray-500">' + arrow.type + '</span>' +
              '</div>' +
              '<div class="text-sm text-gray-700 mb-1">' + arrow.description + '</div>' +
              '<div class="text-xs font-mono text-gray-500">' + arrow.example + '</div>' +
              '</div>';
          });

          content += '</div></div>';
        });

        content += '</div>';
        arrowsContent.innerHTML = content;
      }

      function insertArrowTemplate(symbol, example)
      {
        const template = '\n' + example + '\n';
        insertAtCaret(outputArea, template);
        outputArea.focus();
        showStatus('Inserted ' + symbol + ' arrow template');
      }

      function insertAtCaret(editableDiv, text)
      {
        const sel = window.getSelection();
        if (!sel.rangeCount) return;
        const range = sel.getRangeAt(0);
        range.deleteContents();
        const textNode = document.createTextNode(text);
        range.insertNode(textNode);
        range.setStartAfter(textNode);
        range.setEndAfter(textNode);
        sel.removeAllRanges();
        sel.addRange(range);
      }

      window.insertArrowTemplate = insertArrowTemplate;

      function saveAsN4L()
      {
        const content = outputArea.innerText.trim();

        if (!content)
        {
          showStatus('No content to save', true);
          return;
        }

        const fileName = currentFileName
          ? currentFileName.replace(/\.[^/.]+$/, '.n4l')
          : 'output.n4l';

        const blob = new Blob([content], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = fileName;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);

        showStatus('File saved as ' + fileName);
      }

      function syncScroll(source)
      {
        if (isScrollSyncing) return;

        isScrollSyncing = true;

        // Determine which element triggered the scroll
        const isOutput = source === outputArea;

        // Get the visible input element (either textarea or preview)
        const inputEl = inputText.classList.contains('hidden') ? inputPreview : inputText;

        const sourceEl = isOutput ? outputArea : inputEl;
        const targetEl = isOutput ? inputEl : outputArea;

        // Calculate scroll position ratio
        const sourceScrollTop = sourceEl.scrollTop;
        const sourceScrollHeight = sourceEl.scrollHeight - sourceEl.clientHeight;
        const scrollRatio = sourceScrollHeight > 0 ? sourceScrollTop / sourceScrollHeight : 0;

        // Apply same ratio to target
        const targetScrollHeight = targetEl.scrollHeight - targetEl.clientHeight;
        targetEl.scrollTop = scrollRatio * targetScrollHeight;

        setTimeout(() => { isScrollSyncing = false; }, 50);
      }

      convertBtn.addEventListener('click', convertText);
      copyBtn.addEventListener('click', copyToClipboard);
      uploadBtn.addEventListener('click', handleFileUpload);
      fileInput.addEventListener('change', handleFileSelect);
      arrowsBtn.addEventListener('click', extractArrows);
      saveBtn.addEventListener('click', saveAsN4L);
      clearSessionBtn.addEventListener('click', clearSessionData);

      // Bidirectional scroll sync
      outputArea.addEventListener('scroll', function () { syncScroll(outputArea); });
      inputText.addEventListener('scroll', function () { syncScroll(inputText); });
      inputPreview.addEventListener('scroll', function () { syncScroll(inputPreview); });

      // Allow Ctrl/Cmd + Enter to convert
      inputText.addEventListener('keydown', (e) =>
      {
        if ((e.ctrlKey || e.metaKey) && e.key === 'Enter')
        {
          convertText();
        }
      });

      // Allow Ctrl/Cmd + S to save
      outputArea.addEventListener('keydown', (e) =>
      {
        if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's')
        {
          e.preventDefault();
          saveAsN4L();
        }
      });

      // Auto-convert on paste
      inputText.addEventListener('paste', () =>
      {
        setTimeout(() =>
        {
          convertText();
        }, 100);
      });

      // Auto-save on input changes (debounced)
      let autoSaveTimeout;
      inputText.addEventListener('input', () =>
      {
        clearTimeout(autoSaveTimeout);
        autoSaveTimeout = setTimeout(() =>
        {
          saveSession();
        }, 2000); // Save 2 seconds after user stops typing
      });

      // Auto-save when output is manually edited
      outputArea.addEventListener('input', () =>
      {
        clearTimeout(autoSaveTimeout);
        autoSaveTimeout = setTimeout(() =>
        {
          saveSession();
        }, 2000);
      });

      // Try to restore previous session on load
      loadSession();

      // Set mode-specific title and behavior
      if (editorMode === 'edit')
      {
        titleSection.innerHTML = '<h1 class="text-2xl font-bold text-gray-800">✏️ N4L Editor</h1><p class="text-sm text-gray-600">Upload a .n4l file to edit and validate</p>';
        uploadBtn.innerHTML = '📁 Upload N4L';
        convertBtn.classList.add('hidden'); // Hide convert button in edit mode
      }
      else if (editorMode === 'convert')
      {
        titleSection.innerHTML = '<h1 class="text-2xl font-bold text-gray-800">▶️ Text to N4L Converter</h1><p class="text-sm text-gray-600">Upload text, HTML, or Markdown to convert</p>';
        uploadBtn.innerHTML = '📁 Upload File';
      }

      // Initial focus
      inputText.focus();
    })
    .catch(error =>
    {
      console.error('Error loading app:', error);
      document.getElementById('app').innerHTML = '<p>Error loading application</p>';
    });
} // End of loadEditor function
