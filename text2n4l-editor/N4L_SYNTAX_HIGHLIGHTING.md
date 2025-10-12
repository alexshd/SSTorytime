# N4L Syntax Highlighting Guide

## Overview

The N4L editor now provides comprehensive syntax highlighting for all N4L DSL elements, making it easier to read, write, and validate your knowledge graphs.

## Highlighted Syntax Elements

### 1. 📝 Comments

**Syntax**: `# comment text`  
**Color**: Gray, italic  
**Example**: `# This is a comment`

Comments appear in gray italics to distinguish them from actual content.

---

### 2. 📚 Title Markers

**Syntax**: `- Title text` (dash at start of line)  
**Color**: Purple, bold  
**Example**: `- My Notes About AI`

The dash marker appears normal, but the title text is highlighted in purple to make document structure clear.

---

### 3. 🔄 Sequence Markers

**Syntax**: `+::` (start sequence), `-::` (end sequence)  
**Color**: Orange, bold  
**Example**:

```n4l
+:: _sequence_ , Mary had a little lamb, poem ::
Some content here
-:: _sequence_ ::
```

Sequence mode markers are highlighted in bold orange to show document structure.

---

### 4. 🏷️ Tag References

**Syntax**: `@identifier` (at-symbol followed by name)  
**Color**: Cyan, bold  
**Examples**:

- `@title`
- `@senX`
- `@chapter.1`

Tag references (like `@title`) are highlighted in cyan to show they reference named concepts.

---

### 5. 💲 Value References

**Syntax**: `$identifier` (dollar-symbol followed by name)  
**Color**: Cyan, bold  
**Examples**:

- `$title.1`
- `$concept`

Value references use the same cyan color as tags to show they're referencing something.

---

### 6. 🔁 Ditto Symbol

**Syntax**: `"` (quote at start of line)  
**Color**: Green, bold, slightly larger  
**Example**:

```n4l
Mary had a little lamb
  " (written by) Mary's mum
```

The ditto symbol (") means "repeat the previous concept" and is highlighted in green.

---

### 7. 🎯 Special Annotations

**Syntax**: Single-character annotations  
**Color**: Red, bold  
**Symbols**:

- `%` - discusses
- `=` - involves
- `**` - is a special case of
- `>>` - is an example of
- `>` - has actor/subject role
- `<` - has affected object role

**Example**: `>> (is an example of)`

These special short-form annotations are highlighted in red for visibility.

---

### 8. ⚠️ ALL CAPS Reminders

**Syntax**: 3+ consecutive uppercase words  
**Color**: Yellow background with brown text, bold  
**Example**: `IF YOU WRITE IN ALL CAPS, YOU WILL BE REMINDED OF THE NOTE LATER!`

All-caps text (at least 3 words) gets a yellow highlight to indicate it's a reminder/important note.

---

### 9. 🏹 Semantic Arrows (Parenthetical)

**Syntax**: `(relationship description)`  
**Colors**:

- **Valid arrows**: Blue background with darker blue text
- **Invalid arrows**: Red background with error icon

**Examples**:

- `(leads to)` - Blue (valid)
- `(contains)` - Blue (valid)
- `(leadsto)` - Red (invalid - space missing)

**Interactive**: Click any arrow to:

- See validation status
- Replace with correct arrow
- Delete the line
- View all valid alternatives

---

## Complete Example with All Highlighting

```n4l
- My Notes About Mary          # Purple title, gray comment

@title Mary had a little lamb   # Cyan @title
         (note) Had means possessed   # Blue arrow

+:: _sequence_ , poem ::        # Orange sequence marker

  "  (written by) Mary's mum    # Green ditto, blue arrow

$title.1 (is an example of) Nursery rhyme   # Cyan $ref, blue arrow

>> (has example) Another verse   # Red >>, blue arrow

IF YOU WRITE IN CAPS IT STANDS OUT   # Yellow highlight

-:: _sequence_ ::               # Orange end marker
```

## Color Palette Reference

| Element              | Color                        | Purpose                       |
| -------------------- | ---------------------------- | ----------------------------- |
| Comments `#`         | `#6b7280` Gray               | De-emphasize meta-information |
| Titles `-`           | `#7c3aed` Purple             | Highlight document structure  |
| Sequences `+::/-::`  | `#ea580c` Orange             | Show mode changes             |
| Tags/Refs `@/$`      | `#0891b2` Cyan               | Identify references           |
| Ditto `"`            | `#059669` Green              | Show repetition               |
| Annotations `%/=/>>` | `#dc2626` Red                | Highlight special syntax      |
| Reminders (CAPS)     | `#fef3c7` bg, `#92400e` text | Draw attention                |
| Valid Arrows         | `#dbeafe` bg, `#0369a1` text | Show valid semantics          |
| Invalid Arrows       | `#fee2e2` bg, `#dc2626` text | Flag errors                   |

## Benefits

1. **Improved Readability**: Each syntax element has distinct visual treatment
2. **Error Detection**: Invalid arrows are immediately visible
3. **Structure Clarity**: Document organization (titles, sequences) stands out
4. **Reference Tracking**: Easy to spot where concepts are defined and referenced
5. **Validation Feedback**: Real-time indication of N4L compliance
6. **Learning Aid**: Visual cues help users learn N4L syntax faster

## Technical Implementation

The highlighting is applied in stages:

1. **Escape HTML** to preserve text
2. **Line-by-line processing** for line-based syntax
3. **Regex replacements** for each syntax element
4. **Arrow validation** against 300+ valid arrows from SSTconfig
5. **Interactive wrapping** for clickable arrow editing

All highlighting is done client-side with zero impact on the actual N4L text content.

## Future Enhancements

Potential additions:

- Syntax highlighting themes (light/dark/colorblind-friendly)
- Customizable colors via settings
- Highlight matching references (click `@title` to highlight `$title` uses)
- Fold/unfold sequences
- Syntax error underlining with tooltips

---

**Enjoy your colorful N4L editing experience!** 🎨
