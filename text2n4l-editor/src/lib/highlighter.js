// N4L Syntax highlighting

import { isValidArrow } from './arrows.js';

export function highlightArrows(text)
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
      } else
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
