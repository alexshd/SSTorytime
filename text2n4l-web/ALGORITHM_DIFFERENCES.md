# Algorithm Differences Between CLI and Web text2N4L

## Issue: Web and CLI versions produce different outputs for same input

**Test Case:** `obama.dat` converted at 100% using both tools

### Key Differences Found

#### 1. **Text Fractionation (Sentence Splitting)**

**CLI (`src/text2N4L.go` + `pkg/SSTorytime/SSTorytime.go`):**

- Uses `CleanText()` which:
  - Replaces HTML tags with `:\n`
  - Handles abbreviations (Mr., Dr., etc.)
  - Marks sentence boundaries with `#` after `?!.。` followed by space/newline
  - Converts double newlines to `>>` for paragraph markers
  - Preserves parenthetical content with special handling
- Uses `SplitIntoParaSentences()` which:
  - Splits on `>>` for paragraphs
  - Splits on `#` for sentences
  - Handles parentheses specially (expands and keeps unexpanded versions)
  - Calls `SplitPunctuationText()` for further fractionation
- Preserves original sentence boundaries better

**Web (`text2n4l-web/internal/analyzer/converter.go`):**

- Uses simpler `FractionateTextFile()` which:
  - Splits on `. `, `! `, `? `, `\n`, `\r\n`
  - No special handling for abbreviations
  - No paragraph marking
  - No parenthetical expansion
- Result: **Merges multiple sentences together** (see sen1 example below)

**Example:**

```
CLI sen1:   OBAMA, My fellow citizens:
CLI sen2:   I stand here today humbled by the task before us...

Web sen1:   OBAMA: My fellow citizens: I stand here today humbled by the task before us...
            (merged two sentences into one!)
```

#### 2. **Source Filename Handling**

**CLI:**

- Uses actual filename: `../examples/example_data/obama.dat`
- Preserves file path in output

**Web:**

- Hardcoded to: `uploaded.txt`
- Should use actual uploaded filename

**Location:** `text2n4l-web/internal/analyzer/converter.go` line ~11

#### 3. **Topic/Context Extraction**

**CLI:**

- Header: `:: _sequence_ , obama ::`
- Context per section: `:: new age ::`
- Uses `ExtractIntentionalTokens()` from SSTorytime package
- Sophisticated n-gram frequency tracking with location history
- Uses `STM_NGRAM_FREQ` and `STM_NGRAM_LOCA` global maps

**Web:**

- Header: `:: and we will, men and women ::`
- Different n-gram extraction algorithm
- Simpler frequency counting without location tracking
- Different ranking/selection of top phrases

**Example:**

```
CLI context keywords:
  "this transition", "our nation", "our planet", "this day", "last year"

Web context keywords:
  "and we will", "men and women", "we the people", "is not complete"
```

#### 4. **N-gram Constants**

**CLI (`pkg/SSTorytime/SSTorytime.go`):**

```go
const N_GRAM_MIN = 1
const N_GRAM_MAX = 4  // 1,2,3-grams
```

**Web (`text2n4l-web/internal/analyzer/types.go`):**

```go
const N1GRAM = 1
const N2GRAM = 2
const N3GRAM = 3
const N_GRAM_MAX = 4
```

These appear aligned but used differently.

## Required Fixes

### Priority 1: Text Fractionation

**Goal:** Make web version split sentences identically to CLI

**Changes needed in `text2n4l-web/internal/analyzer/converter.go`:**

1. Port `CleanText()` from CLI:

   - Handle abbreviations (Mr., Dr., Mrs., Ms., St.)
   - Use `#` markers for sentence boundaries
   - Use `>>` for paragraph markers
   - Handle ellipsis and em-dashes
   - Strip HTML tags properly

2. Port `SplitIntoParaSentences()` logic:

   - Split on `>>` first (paragraphs)
   - Split on `#` second (sentences)
   - Handle parenthetical content

3. Port `SplitPunctuationText()` if needed for finer fractionation

### Priority 2: Filename Handling

**Goal:** Use actual uploaded filename instead of "uploaded.txt"

**Changes needed:**

- Modify web handler to pass actual filename through to converter
- Update `N4LSkeletonOutput()` to use provided filename

### Priority 3: N-gram Extraction

**Goal:** Use same intentionality token extraction as CLI

**Changes needed:**

- Port the CLI's `Fractionate()` function logic
- Port the global n-gram frequency and location tracking
- Match the CLI's n-gram selection algorithm for context

### Priority 4: Context Selection

**Goal:** Generate same "appears close to" keywords as CLI

**Changes needed:**

- Port CLI's context building logic from `src/text2N4L.go`
- Match the selection criteria for ambient vs. anomalous phrases

## Testing

After fixes, both tools should produce **identical output** for same input:

- Same sentence boundaries
- Same topic extraction
- Same context keywords
- Same "appears close to" keywords

**Test command:**

```bash
# CLI
./text2N4L -% 100 examples/example_data/obama.dat

# Web
curl -X POST http://localhost:8080/convert \
  -F "file=@examples/example_data/obama.dat" \
  -F "percentage=100"

# Compare
diff obama.dat_edit_me.n4l obama.dat.n4l
```

## Root Cause

The web version was reimplemented from scratch with simplified algorithms rather than calling the existing CLI codebase. To ensure consistency, the web version should either:

1. **Option A:** Port all the exact algorithms from CLI (current path)
2. **Option B:** Refactor CLI code into a shared library and call it from web (better long-term)
3. **Option C:** Make CLI a subprocess call from web (simplest but less efficient)
