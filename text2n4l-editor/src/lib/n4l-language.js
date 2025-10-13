// N4L Language support for CodeMirror

import { StreamLanguage } from '@codemirror/language';
import { isValidArrow } from './arrows.js';

// Define N4L syntax highlighting rules
const n4lLanguage = StreamLanguage.define({
  startState: () => ({
    inContentBlock: false,
    contentIndentLevel: 0
  }),

  token(stream, state)
  {
    // Get the line content
    const line = stream.string;
    const currentIndent = line.match(/^\s*/)[0].length;

    // Check for @ tag (starts content block)
    if (stream.match(/^@[a-zA-Z_][a-zA-Z0-9_\.]*/))
    {
      state.inContentBlock = true;
      state.contentIndentLevel = currentIndent;
      return 'annotation';
    }

    // Check for ditto line with arrow (exits content block)
    if (state.inContentBlock && /^\s*"\s+\([a-z]/.test(line))
    {
      state.inContentBlock = false;
    }

    // If in content block
    if (state.inContentBlock)
    {
      // Check if we've gone back to lower indentation
      if (currentIndent <= state.contentIndentLevel && line.trim() !== '')
      {
        state.inContentBlock = false;
      } else
      {
        // In content block - plain text
        stream.skipToEnd();
        return 'content';
      }
    }

    // Comments (# anything after)
    if (stream.match(/^#.*/))
    {
      return 'comment';
    }

    // Title markers (- at start)
    if (stream.sol() && stream.match(/^-\s+/))
    {
      return 'heading';
    }

    // Sequence markers (+:: and -::)
    if (stream.match(/^(\+::|−::|−::)/))
    {
      return 'keyword';
    }

    // $ references
    if (stream.match(/^\$[a-zA-Z_][a-zA-Z0-9_\.]*/))
    {
      return 'variableName';
    }

    // Ditto symbol (standalone " at beginning)
    if (stream.sol() && stream.match(/^"\s/))
    {
      return 'string';
    }

    // Arrows - (text inside parentheses)
    if (stream.match(/^\([a-z][a-z\s,;:.\-'\/]{3,}\)/i))
    {
      const arrow = stream.current();
      // Check if valid arrow
      const isValid = isValidArrow(arrow);
      return isValid ? 'arrow-valid' : 'arrow-invalid';
    }

    // % lines (context/meta)
    if (stream.sol() && stream.match(/^%.*/))
    {
      return 'meta';
    }

    // = lines (relations)
    if (stream.sol() && stream.match(/^=.*/))
    {
      return 'operator';
    }

    // >> lines (outcomes)
    if (stream.sol() && stream.match(/^>>.*/))
    {
      return 'keyword';
    }

    // Default: advance one character
    stream.next();
    return null;
  }
});

export { n4lLanguage };
