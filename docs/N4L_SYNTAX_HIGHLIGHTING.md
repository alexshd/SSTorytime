# N4L DSL Syntax Highlighting and Formatting Rules

This document defines syntax highlighting and formatting rules for the N4L (Notes for Learning) domain-specific language used in the SSTorytime project.

## Language Overview

N4L is a simple knowledge management language for creating structured notes that can be converted into semantic spacetime graphs. It focuses on relationships between concepts using a minimal syntax.

## File Extension

- `.n4l` - N4L source files

## Syntax Elements

### Comments

```n4l
# Single line comment
// Alternative comment style
```

### Section/Chapter Declarations

```n4l
-section name
-chapter title
```

### Context Sets

```n4l
: list, context, words :         # Basic context
:: list, context, words ::       # Extended context (any number of colons)
+:: extend-list, context ::      # Extend existing context
-:: delete, words ::             # Remove from context
```

### Items and Text

```n4l
Simple item
"Quoted item with spaces"
'Single quoted item with "internal quotes"'
```

### Relationships

```n4l
A (relation) B                   # Basic relationship
A (relation) B (relation) C      # Chain relationship
" (relation) D                   # Continuation from previous item
$1 (relation) D                  # Reference to first previous item
$2 (relation) E                  # Reference to second previous item
```

### Aliases and References

```n4l
@myalias                         # Define alias for this line
$myalias.1                       # Reference to aliased line
$PREV.2                          # Reference to previous items
```

### Special Markers

```n4l
NOTE TO SELF ALLCAPS             # TODO items (all caps)
=specialword                     # Special word marker
*emphasized                      # Emphasis marker
%concept                         # Concept marker
```

### URLs and Media

```n4l
item (url) "https://example.com"
item (img) "https://example.com/image.jpg"
```

### Sequence Mode

```n4l
+:: _sequence_ , context ::      # Start sequence mode
-:: _sequence_ ::                # End sequence mode
```

## Syntax Highlighting Rules

### Colors and Styles

#### Comments

- **Pattern**: `^[[:space:]]*#.*$`, `^[[:space:]]*//.*$`
- **Style**: Italic, muted color (gray)
- **Scope**: `comment.line.hash.n4l`, `comment.line.double-slash.n4l`

#### Section/Chapter Headers

- **Pattern**: `^[[:space:]]*-[[:space:]]*(.+)$`
- **Style**: Bold, prominent color (blue/purple)
- **Scope**: `entity.name.section.n4l`

#### Context Declarations

- **Pattern**: `^[[:space:]]*[+\-]?:{1,}[[:space:]]*(.+?)[[:space:]]*:{1,}[[:space:]]*$`
- **Style**: Bold, distinctive color (orange/yellow)
- **Scope**: `entity.name.tag.context.n4l`
- **Captures**:
  - `+` prefix: `keyword.operator.context.extend.n4l`
  - `-` prefix: `keyword.operator.context.remove.n4l`
  - Content: `string.unquoted.context.n4l`

#### Quoted Strings

- **Pattern**: `"([^"]*)"`, `'([^']*)'`
- **Style**: Standard string color (green)
- **Scope**: `string.quoted.double.n4l`, `string.quoted.single.n4l`

#### Relationships (Parentheses)

- **Pattern**: `\(([^)]+)\)`
- **Style**: Bold, accent color (cyan/teal)
- **Scope**: `entity.name.function.relation.n4l`
- **Captures**:
  - Parentheses: `punctuation.definition.relation.n4l`
  - Relation name: `entity.name.function.relation.name.n4l`

#### Aliases and References

- **Pattern**: `@([a-zA-Z_][a-zA-Z0-9_]*)`
- **Style**: Bold, identifier color (purple)
- **Scope**: `entity.name.constant.alias.n4l`

- **Pattern**: `\$([a-zA-Z_][a-zA-Z0-9_]*(?:\.[0-9]+)?|\d+|PREV\.[0-9]+)`
- **Style**: Italic, reference color (magenta)
- **Scope**: `variable.other.reference.n4l`

#### Special Markers

- **Pattern**: `%([a-zA-Z_][a-zA-Z0-9_]*)`
- **Style**: Bold italic, concept color (dark blue)
- **Scope**: `entity.name.type.concept.n4l`

- **Pattern**: `=([a-zA-Z_][a-zA-Z0-9_]*)`
- **Style**: Bold, marker color (red)
- **Scope**: `markup.bold.special.n4l`

- **Pattern**: `\*([a-zA-Z_][a-zA-Z0-9_]*)`
- **Style**: Italic, emphasis color (orange)
- **Scope**: `markup.italic.emphasis.n4l`

#### ALL CAPS (TODO items)

- **Pattern**: `^[[:space:]]*[A-Z][A-Z0-9[:space:]]+[A-Z0-9]$`
- **Style**: Bold, warning color (yellow/orange background)
- **Scope**: `markup.bold.todo.n4l`

#### Continuation Markers

- **Pattern**: `^[[:space:]]*"[[:space:]]*`
- **Style**: Subtle, operator color (gray)
- **Scope**: `keyword.operator.continuation.n4l`

#### URLs

- **Pattern**: `https?://[^\s"']+`
- **Style**: Underlined, link color (blue)
- **Scope**: `markup.underline.link.n4l`

### Token Classification

#### Keywords

```
Operators: +, -, @, $, ", '
Reserved Relations: then, url, img, note, e.g., ex, prop, rule
Context Keywords: _sequence_
```

#### Identifiers

```
Aliases: @identifier
References: $identifier, $identifier.number, $PREV.number, $number
Concepts: %identifier
Special: =identifier, *identifier
```

#### Literals

```
Strings: "text", 'text'
URLs: http://..., https://...
```

#### Comments

```
Hash: # comment
Slash: // comment
```

## Formatting Rules

### Indentation

- Use 2 spaces for indentation
- Context declarations align at column 1
- Relationship continuations indent +2 from parent
- Nested items indent +2 per level

### Line Spacing

- Single blank line after section headers
- Single blank line between major context groups
- No blank lines within relationship chains
- Double blank line between major sections

### Alignment

- Align relationship operators `(relation)` vertically when in chains
- Align continuation quotes `"` at consistent column
- Right-align comments when on same line as content

### Line Length

- Soft limit: 80 characters
- Hard limit: 120 characters
- Break long relationship chains across lines

### Example Formatted Code

```n4l
-chinese language notes

:: basic greetings ::

hello (eh) 你好 (hp) nǐhǎo
  "   (e.g.) nǐhǎo, wǒ jiào Mark  (ph) 你好，我叫马克  (he) hello, my name is Mark

thank you (eh) 谢谢 (hp) xièxiè
      "   (note) very common greeting

:: grammar rules ::

@negation using 不/bù versus 没有/méiyǒu (rule) 不/bù negates present and future
        " (rule) 没/méi negates past actions
        " (e.g.) 我没去学校  (ph) wǒ méi qù xuéxiào  (he) I didn't go to school

$negation.1 (ex) 我不吃肉  (ph) wǒ bù chī ròu  (he) I don't eat meat
```

## TextMate/VSCode Grammar Structure

```json
{
  "name": "N4L",
  "scopeName": "source.n4l",
  "fileTypes": ["n4l"],
  "patterns": [
    {
      "include": "#comments"
    },
    {
      "include": "#sections"
    },
    {
      "include": "#contexts"
    },
    {
      "include": "#relationships"
    },
    {
      "include": "#aliases"
    },
    {
      "include": "#references"
    },
    {
      "include": "#special-markers"
    },
    {
      "include": "#strings"
    },
    {
      "include": "#todos"
    },
    {
      "include": "#urls"
    }
  ],
  "repository": {
    "comments": {
      "patterns": [
        {
          "name": "comment.line.hash.n4l",
          "match": "^\\s*#.*$"
        },
        {
          "name": "comment.line.double-slash.n4l",
          "match": "^\\s*//.*$"
        }
      ]
    },
    "sections": {
      "name": "entity.name.section.n4l",
      "match": "^\\s*-\\s*(.+)$",
      "captures": {
        "1": {
          "name": "string.unquoted.section.n4l"
        }
      }
    },
    "contexts": {
      "name": "entity.name.tag.context.n4l",
      "match": "^\\s*([+\\-]?)(:{1,})\\s*(.+?)\\s*(:{1,})\\s*$",
      "captures": {
        "1": {
          "name": "keyword.operator.context.n4l"
        },
        "2": {
          "name": "punctuation.definition.context.begin.n4l"
        },
        "3": {
          "name": "string.unquoted.context.n4l"
        },
        "4": {
          "name": "punctuation.definition.context.end.n4l"
        }
      }
    },
    "relationships": {
      "name": "entity.name.function.relation.n4l",
      "match": "\\(([^)]+)\\)",
      "captures": {
        "1": {
          "name": "entity.name.function.relation.name.n4l"
        }
      }
    },
    "aliases": {
      "name": "entity.name.constant.alias.n4l",
      "match": "@([a-zA-Z_][a-zA-Z0-9_]*)"
    },
    "references": {
      "name": "variable.other.reference.n4l",
      "match": "\\$([a-zA-Z_][a-zA-Z0-9_]*(?:\\.[0-9]+)?|\\d+|PREV\\.[0-9]+)"
    },
    "special-markers": {
      "patterns": [
        {
          "name": "entity.name.type.concept.n4l",
          "match": "%([a-zA-Z_][a-zA-Z0-9_]*)"
        },
        {
          "name": "markup.bold.special.n4l",
          "match": "=([a-zA-Z_][a-zA-Z0-9_]*)"
        },
        {
          "name": "markup.italic.emphasis.n4l",
          "match": "\\*([a-zA-Z_][a-zA-Z0-9_]*)"
        }
      ]
    },
    "strings": {
      "patterns": [
        {
          "name": "string.quoted.double.n4l",
          "match": "\"([^\"]*)\""
        },
        {
          "name": "string.quoted.single.n4l",
          "match": "'([^']*)'"
        }
      ]
    },
    "todos": {
      "name": "markup.bold.todo.n4l",
      "match": "^\\s*[A-Z][A-Z0-9\\s]+[A-Z0-9]$"
    },
    "urls": {
      "name": "markup.underline.link.n4l",
      "match": "https?://[^\\s\"']+"
    }
  }
}
```

## Language Server Features

### Semantic Tokens

- Relations: Function calls
- Aliases: Constants
- References: Variables
- Concepts: Types
- Special markers: Keywords
- Sections: Namespaces
- Contexts: Structs

### Validation Rules

1. Alias definitions must be unique within file
2. References must point to existing aliases or valid patterns
3. Context sets must be properly closed
4. Relationship chains should be well-formed
5. Section names should be unique

### Auto-completion

- Relation names from SSTconfig files
- Common patterns: (e.g.), (note), (rule), (then)
- Alias references: $alias.number patterns
- Context keywords: _sequence_

### Folding

- Section blocks (from header to next header)
- Context blocks (between :: markers)
- Relationship chains (multi-line)
- Comment blocks

This comprehensive syntax definition provides full highlighting support for the N4L DSL while maintaining readability and consistency with the language's minimalist philosophy.
