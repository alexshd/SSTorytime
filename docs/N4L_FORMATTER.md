# N4L Formatter Specification

This document defines formatting rules for the N4L DSL that can be implemented in various tools.

## Formatting Rules

### Indentation

- **Base**: Use 2 spaces for indentation
- **Sections**: No indentation (column 1)
- **Context blocks**: No indentation (column 1)
- **Items**: Base indentation (2 spaces from margin)
- **Relationship continuations**: Additional 2 spaces per level
- **Examples and sub-items**: Additional 2 spaces per nesting level

### Line Spacing

- **Section headers**: 1 blank line before, 1 blank line after
- **Context blocks**: 1 blank line before and after
- **Relationship chains**: No blank lines within chains
- **Major topic changes**: 2 blank lines between unrelated sections

### Alignment

- **Relationship operators**: Align `(relation)` vertically in chains when practical
- **Continuation quotes**: Align `"` at consistent column in chains
- **Comments**: Right-align inline comments at column 40+ when space permits

### Line Length

- **Soft limit**: 80 characters
- **Hard limit**: 120 characters
- **Breaking**: Break long relationship chains at logical relation boundaries

### Whitespace

- **Around operators**: Single space around `(relation)`
- **Context delimiters**: Single space after opening `::` and before closing `::`
- **Trailing whitespace**: Remove all trailing whitespace
- **Multiple spaces**: Collapse to single space except for alignment

## Example Transformations

### Before Formatting

```n4l
-chinese notes
::basic greetings::
hello(eh)你好(hp)nǐhǎo
"(e.g.)nǐhǎo, wǒ jiào Mark(ph)你好，我叫马克(he)hello, my name is Mark
thank you(eh)谢谢(hp)xièxiè
"(note)very common greeting
::grammar rules::
@negation using 不/bù versus 没有/méiyǒu(rule)不/bù negates present and future
"(rule)没/méi negates past actions
"(e.g.)我没去学校(ph)wǒ méi qù xuéxiào(he)I didn't go to school
$negation.1(ex)我不吃肉(ph)wǒ bù chī ròu(he)I don't eat meat
```

### After Formatting

```n4l
-chinese notes

:: basic greetings ::

  hello (eh) 你好 (hp) nǐhǎo
    "   (e.g.) nǐhǎo, wǒ jiào Mark (ph) 你好，我叫马克 (he) hello, my name is Mark

  thank you (eh) 谢谢 (hp) xièxiè
        "   (note) very common greeting

:: grammar rules ::

  @negation using 不/bù versus 没有/méiyǒu (rule) 不/bù negates present and future
         "  (rule) 没/méi negates past actions
         "  (e.g.) 我没去学校 (ph) wǒ méi qù xuéxiào (he) I didn't go to school

  $negation.1 (ex) 我不吃肉 (ph) wǒ bù chī ròu (he) I don't eat meat
```

## Formatter Implementation Guide

### Parsing Strategy

1. **Tokenize** the input into logical elements
2. **Identify** structure: sections, contexts, relationships, continuations
3. **Group** related lines (relationship chains, context blocks)
4. **Apply** formatting rules per element type
5. **Preserve** semantic meaning while improving readability

### Core Algorithm

```pseudo
for each line in input:
  line_type = classify_line(line)

  switch line_type:
    case SECTION:
      output_with_spacing(format_section(line), before=1, after=1)

    case CONTEXT:
      output_with_spacing(format_context(line), before=1, after=1)

    case RELATIONSHIP_START:
      relationship_chain = collect_chain(line, next_lines)
      output(format_relationship_chain(relationship_chain))

    case CONTINUATION:
      # Handled as part of relationship chain
      continue

    case COMMENT:
      output(format_comment(line))

    case TODO:
      output(format_todo(line))

    case PLAIN_ITEM:
      output(format_item(line))
```

### Key Functions

#### `format_section(line)`

- Remove extra whitespace
- Ensure single `-` prefix
- Trim to reasonable length

#### `format_context(line)`

- Normalize `::` delimiters
- Add spaces around content
- Handle `+` and `-` modifiers

#### `format_relationship_chain(lines)`

- Align continuation quotes
- Standardize spacing around `(relation)`
- Maintain logical grouping
- Break long lines appropriately

#### `format_item(line)`

- Apply base indentation
- Handle special markers (`%`, `=`, `*`)
- Preserve quotes and escape sequences

### Configuration Options

```yaml
indent_size: 2
max_line_length: 80
align_relations: true
align_continuations: true
spaces_around_relations: true
normalize_quotes: true
preserve_blank_lines: false
sort_contexts: false # Keep original order
```

## Editor Integration

### VSCode

- Implement as language server extension
- Use `DocumentFormattingEditProvider`
- Support format-on-save and format-on-type

### Vim

- Create `formatprg` script
- Integrate with `:Format` command
- Support visual selection formatting

### Emacs

- Implement as `format-buffer` function
- Integrate with `format-region`
- Support auto-formatting hooks

### Command Line Tool

```bash
# Format single file
n4l-format input.n4l > output.n4l

# Format in place
n4l-format -i input.n4l

# Check formatting (exit 1 if not formatted)
n4l-format --check input.n4l

# Format with custom config
n4l-format --config .n4lrc input.n4l
```

## Validation Rules

### Syntax Validation

- Check for unmatched quotes
- Validate context block closure
- Ensure alias uniqueness
- Verify reference targets exist

### Style Validation

- Warn on overly long lines
- Suggest alias usage for repeated items
- Check for inconsistent indentation
- Flag potential TODO items not in caps

### Semantic Validation

- Verify relation names against config files
- Check for circular references
- Validate URL formats
- Ensure sequence mode consistency

This formatter specification provides a solid foundation for implementing consistent N4L code formatting across different editors and tools.
