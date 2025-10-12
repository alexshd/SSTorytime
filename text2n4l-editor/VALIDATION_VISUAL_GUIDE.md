# Quick Reference: N4L Arrow Validation Visual Guide

## What You'll See in the Editor

### ✅ Valid Arrow (Blue)

```
Promise theory (similar to) agent theory
                ^^^^^^^^^^^
                Blue background
                Clickable for alternatives
```

### ❌ Invalid Arrow (Red with Warning)

```
Promise theory ⚠️(appears close to) agent theory
                ^^^^^^^^^^^^^^^^^^
                Red background with wavy underline
                Click to see suggested fixes
```

## Arrow Menu Examples

### For Invalid Arrow

Click on `⚠️(appears close to)` shows:

```
┌─────────────────────────────────────────┐
│ ⚠️ Invalid N4L Arrow                    │
│                                         │
│ This arrow is not recognized in the     │
│ N4L configuration. Please select a      │
│ valid arrow from the list below.        │
├─────────────────────────────────────────┤
│ Replace with N4L Arrow:                 │
│                                         │
│ ✓ Suggested Matches:                    │
│ ┌─────────────────────────────────────┐ │
│ │ (similar to)              (sim)     │ │ ← Highlighted in green
│ │ Type: NR-0                          │ │
│ └─────────────────────────────────────┘ │
│ ┌─────────────────────────────────────┐ │
│ │ (near to)                 (nr)      │ │
│ │ Type: NR-0                          │ │
│ └─────────────────────────────────────┘ │
│                                         │
│ Browse All Arrows:                      │
│                                         │
│ 🔗 Similarity (Non-directional)         │
│   (associated with)         (ass)       │
│   (see also)               (see)        │
│   ...                                   │
│                                         │
│ ➡️ Causality (Leads To)                 │
│ 📦 Composition (Contains)               │
│ 🏷️ Properties (Expresses)               │
│ ⭐ Special Annotations                   │
├─────────────────────────────────────────┤
│ 🗑️ Delete entire line                   │
└─────────────────────────────────────────┘
```

### For Valid Arrow

Click on `(similar to)` shows:

```
┌─────────────────────────────────────────┐
│ Change N4L Arrow:                       │
│                                         │
│ 🔗 Similarity (NR-0)                    │
│   (associated with)         (ass)       │
│   (see also)               (see)        │
│   (alias)                  (alias)      │
│   ...                                   │
│                                         │
│ ➡️ Causality (LT-1)                     │
│   (leads to)               (fwd)        │
│   (causes)                 (cause)      │
│   ...                                   │
│                                         │
│ 📦 Composition (CN-2)                   │
│ 🏷️ Properties (EP-3)                   │
│ ⭐ Special                              │
├─────────────────────────────────────────┤
│ 🗑️ Delete entire line                   │
└─────────────────────────────────────────┘
```

## Color Legend

| Color             | Meaning                                | Action                                    |
| ----------------- | -------------------------------------- | ----------------------------------------- |
| 🔵 **Light Blue** | Valid arrow, recognized by parser      | Click to change or browse alternatives    |
| 🔴 **Light Red**  | Invalid arrow, will cause parser error | **MUST FIX** - Click to see valid options |
| 🟢 **Green**      | Suggested match based on keywords      | Recommended replacement                   |
| ⚪ **White**      | Other available arrows                 | Valid alternatives to consider            |

## Common Patterns

### Pattern 1: Natural Language → N4L Arrow

```
Input:  X (appears close to) Y
        ⚠️ RED - Invalid

Fix:    X (similar to) Y
        ✅ BLUE - Valid (NR-0)
```

### Pattern 2: Verb Phrase → Causality Arrow

```
Input:  A (is related to) B
        ⚠️ RED - Invalid

Fix:    A (associated with) B  [if non-causal]
        ✅ BLUE - Valid (NR-0)

    or: A (leads to) B  [if causal]
        ✅ BLUE - Valid (LT-1)
```

### Pattern 3: Contains/Part-of Relationships

```
Input:  System (has) component
        ⚠️ RED - Invalid (too vague)

Fix:    System (contains) component
        ✅ BLUE - Valid (CN-2)

    or: System (has component) component
        ✅ BLUE - Valid (CN-2)
```

### Pattern 4: Description/Properties

```
Input:  Theory (is about) promises
        ⚠️ RED - Invalid

Fix:    Theory (is about) promises
        ⚠️ Still invalid! Use:

        Theory (is about topic/them) promises
        ✅ BLUE - Valid (EP-3)

    or: Theory (discusses) promises
        ✅ BLUE - Valid (EP-3)
```

## Quick Fixes Table

| If You See           | Replace With                        | Type      |
| -------------------- | ----------------------------------- | --------- |
| `(appears close to)` | `(similar to)`                      | NR-0      |
| `(is similar to)`    | `(similar to)`                      | NR-0      |
| `(relates to)`       | `(associated with)`                 | NR-0      |
| `(is related to)`    | `(associated with)` or `(leads to)` | NR-0/LT-1 |
| `(talks about)`      | `(discusses)`                       | EP-3      |
| `(mentioned in)`     | `(is discussed in)`                 | EP-3      |
| `(is about)`         | `(is about topic/them)`             | EP-3      |
| `(has)`              | `(contains)` or `(has component)`   | CN-2      |
| `(makes)`            | `(creates)`                         | LT-1      |
| `(part of)`          | `(is a part of)`                    | CN-2      |

## Workflow

1. **Convert** text to N4L → generates file with arrows
2. **Upload** to editor
3. **Scan** for red (invalid) arrows
4. **Click** each red arrow
5. **Select** suggested or browse alternatives
6. **Save** when all arrows are blue
7. **Parse** with confidence! ✅

## Keyboard Tips

- **Click arrow**: Open menu
- **Click suggestion**: Replace arrow
- **Click outside**: Close menu
- **Hover arrow**: See preview

## Remember

- ✅ **Blue = Safe** → Parser will accept
- ⚠️ **Red = Error** → Parser will reject
- 🟢 **Green = Suggested** → Smart match
- 📋 **Always validate before parsing!**

---

**Pro Tip**: The menu shows the arrow's **type** (NR-0, LT-1, CN-2, EP-3) so you can learn the semantic categories as you edit!
