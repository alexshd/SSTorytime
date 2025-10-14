# N4L Arrow Validation Feature

## Overview

The text2n4l editor now includes **real-time arrow validation** to help you identify and fix invalid N4L arrows before running the parser. This prevents parser errors and makes editing much more efficient.

## Visual Indicators

### Valid Arrows (Blue)

- **Background**: Light blue (`#e0f2fe`)
- **Text**: Dark blue (`#0369a1`)
- **Border**: Transparent, blue on hover
- **Behavior**: Clickable to show arrow menu with alternatives

### Invalid Arrows (Red)

- **Background**: Light red (`#fee2e2`)
- **Text**: Dark red (`#dc2626`)
- **Border**: Pink (`#fca5a5`)
- **Icon**: ⚠️ warning symbol
- **Underline**: Wavy red underline
- **Behavior**: Clickable to show error message and valid alternatives

## How It Works

### 1. Validation Database

The editor maintains a comprehensive list of **all valid N4L arrows** from the SSTconfig files:

```javascript
getValidArrowsList() {
  return [
    // NR-0: Similarity (78 arrows)
    'similar to', 'sim', 'associated with', 'ass', ...

    // LT-1: Causality (70+ arrows)
    'leads to', 'fwd', 'causes', 'creates', ...

    // CN-2: Composition (90+ arrows)
    'contains', 'belongs to', 'consists of', ...

    // EP-3: Properties (80+ arrows)
    'expresses', 'note', 'has example', 'about', ...

    // Special (10+ arrows)
    'offers', 'accepts', 'observes', ...
  ];
}
```

**Total**: 300+ valid arrows including both full names and abbreviations

### 2. Real-time Validation

When text is converted to N4L:

```javascript
function highlightArrows(text) {
  // For each parenthetical text like "(appears close to)"
  const isValid = isValidArrow(match);
  const cssClass = isValid ? "n4l-arrow-highlight" : "n4l-arrow-error";
  // Apply appropriate styling
}
```

### 3. Contextual Error Messages

When clicking an invalid arrow, the menu shows:

```
⚠️ Invalid N4L Arrow

This arrow is not recognized in the N4L configuration.
Please select a valid arrow from the list below or
the N4L parser will report an error.
```

Followed by **smart suggestions** based on keyword matching.

## Example Validation Scenarios

### Scenario 1: Common Invalid Arrows from text2n4l

| Generated Text       | Validation     | Suggested Fix                       |
| -------------------- | -------------- | ----------------------------------- |
| `(appears close to)` | ❌ **INVALID** | `(similar to)` or `(near to)`       |
| `(is related to)`    | ❌ **INVALID** | `(associated with)` or `(see also)` |
| `(talks about)`      | ❌ **INVALID** | `(discusses)` or `(is about)`       |
| `(mentioned in)`     | ❌ **INVALID** | `(is discussed in)` or `(mentions)` |

### Scenario 2: Valid Arrows

| Text           | Validation   | Type             |
| -------------- | ------------ | ---------------- |
| `(similar to)` | ✅ **VALID** | NR-0 Similarity  |
| `(leads to)`   | ✅ **VALID** | LT-1 Causality   |
| `(contains)`   | ✅ **VALID** | CN-2 Composition |
| `(expresses)`  | ✅ **VALID** | EP-3 Properties  |

### Scenario 3: Abbreviations

| Text     | Validation   | Full Form      |
| -------- | ------------ | -------------- |
| `(sim)`  | ✅ **VALID** | `(similar to)` |
| `(fwd)`  | ✅ **VALID** | `(leads to)`   |
| `(expr)` | ✅ **VALID** | `(expresses)`  |

## Parser Error Prevention

### Without Validation

```bash
$ N4L promisetheory1_edit_me.n4l
Error at line 45: ERR_NO_SUCH_ARROW: No such arrow has been declared
in the configuration: (appears close to)
```

### With Validation

1. Editor shows **red highlighting** on `(appears close to)`
2. Click arrow → see error message
3. Select suggested arrow `(similar to)` from menu
4. Save file
5. ✅ Parser runs successfully!

## Benefits

1. **Instant Feedback**: See errors immediately without running the parser
2. **Smart Suggestions**: Keyword-based matching suggests relevant arrows
3. **Learn as You Edit**: Browse all 300+ valid arrows organized by type
4. **Prevent Parser Failures**: Fix issues before they cause compilation errors
5. **Faster Workflow**: No more trial-and-error with the parser

## Technical Details

### Validation Sources

The validation list is compiled from:

- `/SSTconfig/arrows-NR-0.sst` (Similarity - 78 arrows)
- `/SSTconfig/arrows-LT-1.sst` (Causality - 70 arrows)
- `/SSTconfig/arrows-CN-2.sst` (Composition - 90 arrows)
- `/SSTconfig/arrows-EP-3.sst` (Properties - 80 arrows)
- Special arrows from example files

### Arrow Format Rules

Valid arrows must:

- Be enclosed in parentheses: `(arrow name)`
- Match exactly (case-insensitive) a configured arrow
- Use only defined abbreviations from SSTconfig files

### Matching Algorithm

```javascript
function isValidArrow(arrowText) {
  const normalized = arrowText
    .replace(/^\(|\)$/g, "") // Remove parens
    .trim()
    .toLowerCase();

  return validArrows.some((valid) => valid.toLowerCase() === normalized);
}
```

## Future Enhancements

Potential improvements:

1. **Fuzzy Matching**: Suggest close matches for typos
2. **Custom Arrows**: Allow user-defined arrows with validation
3. **Batch Validation**: Validate entire file and show summary
4. **Export Report**: Generate list of all invalid arrows in file
5. **Auto-fix**: Automatically replace common invalid patterns

## Configuration

The validation is automatically enabled. No configuration needed!

To disable validation (not recommended):

```javascript
// In main.js, always use 'n4l-arrow-highlight' class
const cssClass = "n4l-arrow-highlight"; // Remove validation check
```

## Testing

Test with the example file:

```bash
# Convert example file
cd text2n4l-web
./text2n4l ../examples/example_data/promisetheory1.dat

# Open in editor
cd ../text2n4l-editor
npm run dev

# Upload the generated *_edit_me.n4l file
# Click on arrows to see validation
```

## References

- N4L Parser: `/src/N4L.go`
- Arrow Configs: `/SSTconfig/arrows-*.sst`
- Editor Guide: `/text2n4l-editor/N4L_EDITING_GUIDE.md`
- Valid Example: `/examples/PromiseTheory.n4l`

---

**Status**: ✅ Implemented and tested
**Version**: 1.0
**Date**: October 12, 2025
