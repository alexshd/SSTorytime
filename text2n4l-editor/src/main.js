import './style.css';

// Load main layout from app.html and initialize app
fetch('/src/app.html')
  .then(res => res.text())
  .then(html =>
  {
    document.querySelector('#app').innerHTML = html;

    // Select DOM elements
    const inputText = document.querySelector('#input-text');
    const outputArea = document.querySelector('#output-area');
    const convertBtn = document.querySelector('#convert-btn');
    const copyBtn = document.querySelector('#copy-btn');
    const uploadBtn = document.querySelector('#upload-btn');
    const fileInput = document.querySelector('#file-input');
    const arrowsBtn = document.querySelector('#arrows-btn');
    const saveBtn = document.querySelector('#save-btn');
    const statusMessage = document.querySelector('#status-message');
    const titleSection = document.querySelector('#title-section');
    const arrowsSection = document.querySelector('#arrows-section');
    const arrowsContent = document.querySelector('#arrows-content');

    let currentFileName = '';
    let isScrollSyncing = false;

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
        arrowsBtn.disabled = false;
        showStatus('Text converted successfully!');

      } catch (error)
      {
        console.error('Conversion error:', error);
        showStatus('Error: ' + error.message, true);
        outputArea.innerHTML = '';
        copyBtn.disabled = true;
        saveBtn.disabled = true;
        arrowsBtn.disabled = true;
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

    function handleFileSelect(event)
    {
      const file = event.target.files[0];

      if (!file)
      {
        return;
      }

      // Check file type
      const allowedTypes = ['text/plain', 'text/markdown'];
      const fileExtension = file.name.toLowerCase().split('.').pop();
      const allowedExtensions = ['txt', 'md', 'text'];

      if (!allowedTypes.includes(file.type) && !allowedExtensions.includes(fileExtension))
      {
        showStatus('Please select a text file (.txt, .md, .text)', true);
        return;
      }

      // Check file size (limit to 10MB)
      if (file.size > 10 * 1024 * 1024)
      {
        showStatus('File too large. Please select a file smaller than 10MB.', true);
        return;
      }

      const reader = new FileReader();

      reader.onload = function (e)
      {
        const content = e.target.result;
        inputText.value = content;
        currentFileName = file.name;

        // Hide title section when file is uploaded
        titleSection.classList.add('hidden');

        showStatus('File "' + file.name + '" loaded successfully!');

        // Clear any previous output
        outputArea.innerText = '';
        copyBtn.disabled = true;
        saveBtn.disabled = true;
        arrowsBtn.disabled = true;
        arrowsSection.classList.add('hidden');

        // Reset file input
        fileInput.value = '';
      }; reader.onerror = function ()
      {
        showStatus('Error reading file. Please try again.', true);
        fileInput.value = '';
      };

      reader.readAsText(file);
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

      // Only highlight parenthetical arrow descriptions like "(appears close to)"
      // This matches any text in parentheses that looks like a relationship description
      const parenArrowRegex = /\([a-z][a-z\s,;:-]*\)/gi;

      // Highlight parenthetical arrow descriptions
      const result = escaped.replace(parenArrowRegex, function (match)
      {
        // Only highlight if it's not just a single word (to avoid matching things like (x) or (a))
        // and if it contains at least one space or common relationship words
        if (match.length > 4 && (match.includes(' ') || /(?:to|by|from|in|on|at|as)/i.test(match)))
        {
          return '<span class="n4l-arrow-highlight" data-arrow="' + match + '" onclick="showArrowMenu(event, this)">' +
            match +
            '</span>';
        }
        return match; // Don't highlight short or non-relational parenthetical text
      });

      return result;
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
      menu.style.cssText = 'position: absolute; z-index: 100; background: white; border: 1px solid #38bdf8; border-radius: 0.5rem; box-shadow: 0 4px 12px rgba(0,0,0,0.15); padding: 0.5rem; min-width: 200px;';

      // All available arrows with relationship keywords for matching
      const allArrows = [
        { symbol: '->', name: 'Implies', keywords: ['implies', 'causes', 'leads to', 'results in'] },
        { symbol: '+>', name: 'Positive', keywords: ['positive', 'encourages', 'promotes'] },
        { symbol: '~>', name: 'Approximately', keywords: ['approximately', 'appears close to', 'similar to', 'like'] },
        { symbol: '<->', name: 'Bidirectional', keywords: ['bidirectional', 'mutual', 'reciprocal'] },
        { symbol: '?>', name: 'Conditional', keywords: ['conditional', 'if', 'maybe', 'possibly'] },
        { symbol: '>>', name: 'Sequence', keywords: ['precedes', 'before', 'follows', 'then', 'sequence'] },
        { symbol: '::>', name: 'Context', keywords: ['context', 'depends on', 'given', 'in context'] },
        { symbol: '!>', name: 'Negative', keywords: ['negative', 'prevents', 'blocks', 'inhibits'] },
        { symbol: '|>', name: 'Requires', keywords: ['requires', 'needs', 'demands', 'constrains'] },
        { symbol: '~?>', name: 'Maybe', keywords: ['uncertain', 'might', 'could', 'suggests'] }
      ];

      let menuHTML = '';

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

        menuHTML += '<div style="font-weight: 600; color: #0369a1; margin-bottom: 0.5rem; font-size: 0.875rem;">Replace with arrow:</div>';

        // Show suggested arrows first
        if (suggestedArrows.length > 0)
        {
          menuHTML += '<div style="color: #16a34a; font-size: 0.75rem; margin-bottom: 0.25rem; font-weight: 500;">✓ Suggested:</div>';
          suggestedArrows.forEach(function (arr)
          {
            menuHTML += '<div class="arrow-menu-item" style="padding: 0.4rem 0.5rem; cursor: pointer; border-radius: 0.25rem; font-size: 0.875rem; display: flex; justify-content: space-between; align-items: center; background: #f0fdf4;" ' +
              'onmouseover="this.style.background=\'#dcfce7\'" ' +
              'onmouseout="this.style.background=\'#f0fdf4\'" ' +
              'onclick="replaceArrow(\'' + arrow.replace(/'/g, "\\'") + '\', \'' + arr.symbol + '\', this.closest(\'.arrow-menu\'))">' +
              '<span style="font-family: monospace; color: #0369a1; font-weight: 600;">' + arr.symbol + '</span>' +
              '<span style="color: #64748b;">' + arr.name + '</span>' +
              '</div>';
          });

          menuHTML += '<div style="color: #64748b; font-size: 0.75rem; margin: 0.5rem 0 0.25rem; font-weight: 500;">Other options:</div>';
        }

        // Show other arrows
        const otherArrows = suggestedArrows.length > 0
          ? allArrows.filter(function (arr) { return !suggestedArrows.includes(arr); })
          : allArrows;

        otherArrows.forEach(function (arr)
        {
          menuHTML += '<div class="arrow-menu-item" style="padding: 0.4rem 0.5rem; cursor: pointer; border-radius: 0.25rem; font-size: 0.875rem; display: flex; justify-content: space-between; align-items: center;" ' +
            'onmouseover="this.style.background=\'#e0f2fe\'" ' +
            'onmouseout="this.style.background=\'white\'" ' +
            'onclick="replaceArrow(\'' + arrow.replace(/'/g, "\\'") + '\', \'' + arr.symbol + '\', this.closest(\'.arrow-menu\'))">' +
            '<span style="font-family: monospace; color: #0369a1; font-weight: 600;">' + arr.symbol + '</span>' +
            '<span style="color: #64748b;">' + arr.name + '</span>' +
            '</div>';
        });
      }
      else
      {
        // For arrow symbols, show options to change to other arrows
        menuHTML += '<div style="font-weight: 600; color: #0369a1; margin-bottom: 0.5rem; font-size: 0.875rem;">Change arrow:</div>';

        allArrows.forEach(function (arr)
        {
          if (arr.symbol !== arrow)
          {
            menuHTML += '<div class="arrow-menu-item" style="padding: 0.4rem 0.5rem; cursor: pointer; border-radius: 0.25rem; font-size: 0.875rem; display: flex; justify-content: space-between; align-items: center;" ' +
              'onmouseover="this.style.background=\'#e0f2fe\'" ' +
              'onmouseout="this.style.background=\'white\'" ' +
              'onclick="replaceArrow(\'' + arrow.replace(/'/g, "\\'") + '\', \'' + arr.symbol + '\', this.closest(\'.arrow-menu\'))">' +
              '<span style="font-family: monospace; color: #0369a1; font-weight: 600;">' + arr.symbol + '</span>' +
              '<span style="color: #64748b;">' + arr.name + '</span>' +
              '</div>';
          }
        });
      }

      menuHTML += '<hr style="margin: 0.5rem 0; border-color: #e2e8f0;">';
      menuHTML += '<div class="arrow-menu-item" style="padding: 0.4rem 0.5rem; cursor: pointer; border-radius: 0.25rem; color: #dc2626; font-size: 0.875rem; font-weight: 500;" ' +
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
      const text = outputArea.innerHTML;

      // Escape special characters for regex
      function escapeRegex(str)
      {
        return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
      }

      // Escape HTML entities in the old arrow
      const oldEscaped = oldArrow
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');

      // Escape HTML entities in the new arrow
      const newEscaped = newArrow
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');

      // Find and replace the span containing the old arrow
      const regex = new RegExp('<span class="n4l-arrow-highlight"[^>]*data-arrow="[^"]*"[^>]*>' + escapeRegex(oldEscaped) + '</span>');
      const newText = text.replace(regex, '<span class="n4l-arrow-highlight" data-arrow="' + newArrow + '" onclick="showArrowMenu(event, this)">' + newEscaped + '</span>');
      outputArea.innerHTML = newText;

      menu.remove();
      showStatus('Arrow changed to ' + newArrow);
    };

    // Delete arrow from the output
    window.deleteArrow = function (arrow, menu)
    {
      const text = outputArea.innerHTML;

      // Escape special characters for regex
      function escapeRegex(str)
      {
        return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
      }

      // Escape HTML entities
      const escaped = arrow
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');

      // Find the line containing the arrow and delete the entire line
      // Match from start of line (or after a <br> or newline) to end of line (or before <br> or newline)
      const lineRegex = new RegExp(
        '(?:^|(?:<br>)|(?:\\n))([^\\n<]*<span class="n4l-arrow-highlight"[^>]*data-arrow="[^"]*"[^>]*>' +
        escapeRegex(escaped) +
        '</span>[^\\n<]*)(?:(?:<br>)|(?:\\n)|$)',
        'g'
      );

      let newText = text.replace(lineRegex, '');

      // If that didn't work (no <br> tags), try a simpler approach - just remove the whole line
      if (newText === text)
      {
        // Split by line breaks, filter out lines containing the arrow, rejoin
        const lines = text.split(/\n/);
        const filteredLines = lines.filter(function (line)
        {
          return !line.includes('<span class="n4l-arrow-highlight"[^>]*data-arrow="[^"]*"[^>]*>' + escapeRegex(escaped) + '</span>');
        });
        newText = filteredLines.join('\n');
      }

      outputArea.innerHTML = newText;

      menu.remove();
      showStatus('Line with arrow deleted');
    };

    function extractArrows()
    {
      // Toggle arrows section visibility
      const isHidden = arrowsSection.classList.contains('hidden');

      if (isHidden)
      {
        const n4lText = outputArea.innerText.trim();
        if (!n4lText)
        {
          showStatus('No N4L output to show arrows for', true);
          return;
        }

        // Define available arrow types for editing
        const availableArrows = [
          {
            type: 'Implication',
            symbol: '->',
            description: 'Basic implication or causation',
            example: 'A -> B',
            category: 'Basic'
          },
          {
            type: 'Strong Implication',
            symbol: '+>',
            description: 'Strong or positive implication',
            example: 'A +> B',
            category: 'Basic'
          },
          {
            type: 'Weak Implication',
            symbol: '~>',
            description: 'Weak or uncertain implication',
            example: 'A ~> B',
            category: 'Basic'
          },
          {
            type: 'Bidirectional',
            symbol: '<->',
            description: 'Two-way relationship',
            example: 'A <-> B',
            category: 'Bidirectional'
          },
          {
            type: 'Conditional',
            symbol: '?>',
            description: 'Conditional relationship',
            example: 'A ?> B',
            category: 'Conditional'
          },
          {
            type: 'Temporal',
            symbol: '>>',
            description: 'Temporal sequence',
            example: 'A >> B',
            category: 'Temporal'
          },
          {
            type: 'Contextual',
            symbol: '::>',
            description: 'Context-dependent relationship',
            example: 'A ::> B',
            category: 'Contextual'
          },
          {
            type: 'Negation',
            symbol: '!>',
            description: 'Negative implication',
            example: 'A !> B',
            category: 'Negation'
          },
          {
            type: 'Constraint',
            symbol: '|>',
            description: 'Constraint or requirement',
            example: 'A |> B',
            category: 'Constraint'
          },
          {
            type: 'Optional',
            symbol: '~?>',
            description: 'Optional or possible relationship',
            example: 'A ~?> B',
            category: 'Optional'
          }
        ];

        displayAvailableArrows(availableArrows);
        arrowsSection.classList.remove('hidden');
        showStatus('Showing available arrow types');
      }
      else
      {
        arrowsSection.classList.add('hidden');
        showStatus('Arrow types hidden');
      }
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

    function syncScroll()
    {
      if (isScrollSyncing) return;

      isScrollSyncing = true;

      const outputLines = outputArea.innerText.split('\n');
      const inputLines = inputText.value.split('\n');

      // Calculate scroll position ratio
      const outputScrollTop = outputArea.scrollTop;
      const outputScrollHeight = outputArea.scrollHeight - outputArea.clientHeight;
      const scrollRatio = outputScrollHeight > 0 ? outputScrollTop / outputScrollHeight : 0;

      // Apply same ratio to input
      const inputScrollHeight = inputText.scrollHeight - inputText.clientHeight;
      inputText.scrollTop = scrollRatio * inputScrollHeight;

      setTimeout(() => { isScrollSyncing = false; }, 50);
    }

    convertBtn.addEventListener('click', convertText);
    copyBtn.addEventListener('click', copyToClipboard);
    uploadBtn.addEventListener('click', handleFileUpload);
    fileInput.addEventListener('change', handleFileSelect);
    arrowsBtn.addEventListener('click', extractArrows);
    saveBtn.addEventListener('click', saveAsN4L);
    outputArea.addEventListener('scroll', syncScroll);

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

    // Initial focus
    inputText.focus();
  })
  .catch(error =>
  {
    console.error('Error loading app:', error);
    document.getElementById('app').innerHTML = '<p>Error loading application</p>';
  });
