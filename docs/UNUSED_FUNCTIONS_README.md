# SSTorytime.go Unused Functions Analysis

**Generated:** October 11, 2025  
**Total Functions Analyzed:** 224  
**Used Functions:** 210  
**Unused Functions:** 14

## Unused Functions Summary

This document lists all unused functions in the `pkg/SSTorytime/SSTorytime.go` file. These functions have been marked with `// UNUSED:` comments in the source code.

### 1. GetNodeContext (Line 838)

```go
func GetNodeContext(ctx PoSST, node Node) []string
```

**Purpose:** Retrieves node context as a string slice  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** There is a similar function `GetNodeContextString` that is used instead

### 2. ScoreContext (Line 4027)

```go
func ScoreContext(i, j int) bool
```

**Purpose:** Intended to score context relevance (always returns true)  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** Appears to be a placeholder function that was never fully implemented

### 3. GetSparseOccupancy (Line 5300)

```go
func GetSparseOccupancy(m [][]float32, dim int) []int
```

**Purpose:** Calculates sparse matrix occupancy counts  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** Part of matrix analysis functionality that may have been experimental

### 4. TransposeMatrix (Line 5340)

```go
func TransposeMatrix(m [][]float32) [][]float32
```

**Purpose:** Transposes a 2D float32 matrix  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** Matrix operation that's not needed by current algorithms

### 5. NextLinkArrow (Line 5550)

```go
func NextLinkArrow(ctx PoSST, path []Link, arrows []ArrowPtr) string
```

**Purpose:** Finds the next arrow in a link path matching given arrow types  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** Path analysis functionality that may be superseded by other methods

### 6. ContextInterferometry (Line 6048)

```go
func ContextInterferometry(now_ctx string)
```

**Purpose:** Unknown - marked as "deleted" in comments  
**Status:** UNUSED - Explicitly marked as deleted, empty implementation  
**Reason:** Function was removed but declaration kept

### 7. GetUnixTimeKey (Line 7449)

```go
func GetUnixTimeKey(now int64) string
```

**Purpose:** Generates a database-suitable key from Unix timestamp  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** Time-based key generation that may be replaced by other time functions

### 8. NewNgramMap (Line 7536)

```go
func NewNgramMap() [N_GRAM_MAX]map[string]float64
```

**Purpose:** Creates a new n-gram frequency map array  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** N-gram processing functionality that may be handled differently

### 9. SplitCommandText (Line 7661)

```go
func SplitCommandText(s string) []string
```

**Purpose:** Splits command text using punctuation rules  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** Text processing utility that may be specialized for unused command parsing

### 10. Array2Str (Line 8480)

```go
func Array2Str(arr []string) string
```

**Purpose:** Converts string array to comma-separated string  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** String manipulation utility superseded by standard library functions

### 11. Str2Array (Line 8495)

```go
func Str2Array(s string) ([]string, int)
```

**Purpose:** Converts string to string array  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** String manipulation utility superseded by standard library functions

### 12. RunErr (Line 9009)

```go
func RunErr(message string)
```

**Purpose:** Prints colored error messages  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** Error reporting utility that may be replaced by standard error handling

### 13. ContextString (Line 9038)

```go
func ContextString(context []string) string
```

**Purpose:** Concatenates context strings with spaces  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** Simple string utility that may be replaced by standard library functions

### 14. Already (Line 9174)

```go
func Already(s string, cone map[int][]string) bool
```

**Purpose:** Checks if a string already exists in a cone data structure  
**Status:** UNUSED - Not called anywhere in the codebase  
**Reason:** Utility function for cone/graph algorithms that may be superseded

## Recommendations

### Immediate Actions

1. **Consider Removal**: Most unused functions can be safely removed to reduce code bloat
2. **Archive**: Some functions like matrix operations might be kept for future use
3. **Document**: Functions with potential future value should be better documented

### Functions to Keep (Potentially Useful)

- `TransposeMatrix` - Standard matrix operation that might be needed
- `GetSparseOccupancy` - Could be useful for matrix analysis
- `NewNgramMap` - Might be needed for text processing enhancements

### Functions to Remove (Safe to Delete)

- `ContextInterferometry` - Already marked as deleted
- `ScoreContext` - Placeholder with no implementation
- `Array2Str` / `Str2Array` - Superseded by standard library
- `RunErr` - Simple utility easily replaced
- `Already` - Simple search function easily replaced

## Usage Analysis Methodology

The analysis was performed by:

1. Extracting all function definitions from `SSTorytime.go`
2. Searching for function calls across 65 Go files in the project
3. Checking for both external usage (`SST.FunctionName`) and internal usage
4. Manually verifying results for accuracy

## File Statistics

- **Total Lines:** 9,308
- **Package Files Using SSTorytime:** 65
- **Function Density:** ~2.4% of lines contain function definitions
- **Unused Function Ratio:** 6.25% (14/224)

## Notes

- Some functions may be called dynamically or via reflection (not detected by this analysis)
- Functions may be intended for future use or experimental features
- Consider the maintenance cost vs. potential future value when deciding whether to remove unused functions
