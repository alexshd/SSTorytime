// Line numbers component

export class LineNumbers
{
  constructor(outputArea, lineNumbersEl)
  {
    this.outputArea = outputArea;
    this.lineNumbersEl = lineNumbersEl;
    this.observer = null;
    this.sourceText = ''; // Store the original source text

    this.init();
  }

  init()
  {
    // Update line numbers when content changes
    const observerConfig = { childList: true, subtree: true, characterData: true };
    this.observer = new MutationObserver(() => this.update());
    this.observer.observe(this.outputArea, observerConfig);

    // Sync scroll between output area and line numbers
    this.outputArea.addEventListener('scroll', () => this.syncScroll());

    // Update line numbers on window resize (affects wrapping)
    window.addEventListener('resize', () => this.updateFromSource(this.sourceText));

    // Initial line numbers
    this.update();
  }

  // Update with the source text (called from main.js when content changes)
  updateFromSource(sourceText)
  {
    this.sourceText = sourceText;
    this.updateLineNumbers(sourceText);
  }

  update()
  {
    // If we have source text, use it; otherwise fall back to textContent
    if (this.sourceText)
    {
      this.updateLineNumbers(this.sourceText);
    } else
    {
      // Fallback: get text from the output area (may not be accurate with HTML)
      const text = this.outputArea.textContent || '';
      this.updateLineNumbers(text);
    }
  }

  updateLineNumbers(text)
  {
    if (!text)
    {
      this.lineNumbersEl.innerHTML = '<div>1</div>';
      return;
    }

    // Split by actual newlines to get source lines
    const lines = text.split('\n');
    const lineCount = lines.length;

    // Create a temporary div to measure each line's rendered height
    const tempDiv = document.createElement('div');
    tempDiv.style.cssText = window.getComputedStyle(this.outputArea).cssText;
    tempDiv.style.position = 'absolute';
    tempDiv.style.visibility = 'hidden';
    tempDiv.style.width = this.outputArea.offsetWidth + 'px';
    tempDiv.style.height = 'auto';
    tempDiv.style.whiteSpace = 'pre-wrap';
    tempDiv.style.wordWrap = 'break-word';
    document.body.appendChild(tempDiv);

    let lineNumbersHTML = '';

    for (let i = 0; i < lineCount; i++)
    {
      // Measure the height of this line when rendered
      tempDiv.innerHTML = this.escapeHtml(lines[i] || '\n');
      const height = tempDiv.offsetHeight;

      // Create a line number div with the exact same height
      lineNumbersHTML += `<div style="height: ${height}px; line-height: ${height}px;">${i + 1}</div>`;
    }

    document.body.removeChild(tempDiv);
    this.lineNumbersEl.innerHTML = lineNumbersHTML;
  }

  escapeHtml(text)
  {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  syncScroll()
  {
    this.lineNumbersEl.scrollTop = this.outputArea.scrollTop;
  }

  destroy()
  {
    if (this.observer)
    {
      this.observer.disconnect();
    }
  }
}

if (this.observer)
{
  this.observer.disconnect();
}
  }
}
