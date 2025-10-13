# UI State Transitions - Working Indicator

## Visual Guide

### State 1: Ready to Convert (Idle)

```
┌─────────────────────────────────────────────────────────┐
│  ▶️ Text to N4L Converter                               │
│  Upload text, HTML, or Markdown to convert              │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  [📁 Upload] [▶ Convert] [🏹 Help & Arrows]            │
│                    ↑                                     │
│                 enabled                                  │
│                                                          │
│  ┌────────────────────────┐                            │
│  │ Type or paste text...  │                            │
│  │                        │                            │
│  └────────────────────────┘                            │
└─────────────────────────────────────────────────────────┘
```

### State 2: Converting (Processing)

```
┌─────────────────────────────────────────────────────────┐
│  ▶️ Text to N4L Converter                               │
│  Upload text, HTML, or Markdown to convert              │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  [📁 Upload] [✕ Cancel] [🏹 Help & Arrows] 🔄 Processing...│
│                   ↑ RED                     ↑ ANIMATED  │
│                clickable                   spinner+dots  │
│                                                          │
│  Status: Converting...                                   │
│                                                          │
│  ┌────────────────────────┐                            │
│  │ Your text is here...   │                            │
│  │                        │                            │
│  └────────────────────────┘                            │
└─────────────────────────────────────────────────────────┘
```

### State 3: Cancel Confirmation

```
┌─────────────────────────────────────────────────────────┐
│                                                          │
│         ┌──────────────────────────────────────┐       │
│         │  ⚠️  Confirmation                     │       │
│         │                                       │       │
│         │  Are you sure you want to cancel     │       │
│         │  the conversion?                     │       │
│         │                                       │       │
│         │         [ Cancel ]  [  OK  ]         │       │
│         └──────────────────────────────────────┘       │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### State 4: Cancelled

```
┌─────────────────────────────────────────────────────────┐
│  ▶️ Text to N4L Converter                               │
│  Upload text, HTML, or Markdown to convert              │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  [📁 Upload] [▶ Convert] [🏹 Help & Arrows]            │
│                    ↑                                     │
│                 enabled                                  │
│                                                          │
│  Status: Conversion cancelled by user ⚠️                │
│                                                          │
│  ┌────────────────────────┐                            │
│  │ Your text is here...   │                            │
│  │                        │                            │
│  └────────────────────────┘                            │
└─────────────────────────────────────────────────────────┘
```

### State 5: Conversion Complete

```
┌─────────────────────────────────────────────────────────┐
│  ▶️ Text to N4L Converter                               │
│  Upload text, HTML, or Markdown to convert              │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  [📁 Upload] [▶ Convert] [🏹 Help & Arrows]            │
│                                                          │
│  Status: Conversion complete! ✅                         │
│                                                          │
│  ┌────────────────────────┐  ┌────────────────────────┐│
│  │ Your text...           │  │ @sen0 Result...        ││
│  │                        │  │                        ││
│  └────────────────────────┘  └────────────────────────┘│
│                               [📋 Copy] [💾 Save]       │
└─────────────────────────────────────────────────────────┘
```

## Animation Details

### Spinner SVG

```
    🔄  ← Rotating continuously (CSS animate-spin)
```

### Animated Dots Pattern

```
Frame 1:  Processing...
Frame 2:  Processing.
Frame 3:  Processing..
Frame 4:  Processing...
(repeats every 2 seconds)
```

### Button Color Transitions

**Convert Button (Purple)**

```css
Normal:  Purple gradient (#8B5CF6 → #7C3AED)
Hover:   Brighter purple with focus ring
Disabled: Opacity 50%, cursor not-allowed
Hidden:  display: none
```

**Cancel Button (Red)**

```css
Normal:  Red gradient (#EF4444 → #DC2626)
Hover:   Brighter red with focus ring
Hidden:  display: none (default)
Visible: Only during conversion
```

## Responsive Behavior

### Desktop (Wide Screen)

```
[Upload] [Convert] [Cancel] [Help] 🔄 Processing...
  ↑        ↑         ↑        ↑          ↑
 60px     80px      70px     120px    150px
                                    (visible when active)
```

### Mobile (Narrow Screen)

```
[📁] [▶] [✕] [🏹]
         🔄 Processing...
```

Buttons stack more compactly, text may wrap.

## Color Scheme

| Element         | Color  | Hex     | Purpose               |
| --------------- | ------ | ------- | --------------------- |
| Convert Button  | Purple | #8B5CF6 | Primary action        |
| Cancel Button   | Red    | #EF4444 | Destructive action    |
| Spinner         | Purple | #8B5CF6 | Matches convert theme |
| Success Status  | Green  | #10B981 | Positive feedback     |
| Error Status    | Red    | #EF4444 | Negative feedback     |
| Processing Text | Purple | #8B5CF6 | In-progress state     |

## Keyboard Shortcuts

| Key          | Action            | State                     |
| ------------ | ----------------- | ------------------------- |
| `Ctrl+Enter` | Convert           | Enabled when text present |
| `ESC`        | _(Future)_ Cancel | During conversion         |

## Accessibility

### ARIA Labels (Future Enhancement)

```html
<button aria-label="Cancel conversion in progress">✕ Cancel</button>
<div role="status" aria-live="polite">Processing...</div>
```

### Screen Reader Announcements

1. "Converting text to N4L format"
2. "Conversion complete"
3. "Conversion cancelled"
4. "Error: [message]"

## Performance Considerations

### Spinner Animation

- Uses CSS `animate-spin` (GPU-accelerated)
- No JavaScript involved in rotation
- Minimal CPU/battery impact

### Dots Animation

- JavaScript interval: 500ms
- Updates single text node
- Cleared immediately on completion
- No memory leaks

### Fetch Cancellation

- AbortController cleanup in `finally` block
- Prevents multiple active requests
- Releases resources promptly

## Browser Compatibility

| Feature            | Support                               |
| ------------------ | ------------------------------------- |
| `AbortController`  | Chrome 66+, Firefox 57+, Safari 12.1+ |
| CSS `animate-spin` | All modern browsers                   |
| `confirm()` dialog | Universal support                     |
| Fetch API          | All modern browsers                   |

## Testing Checklist

- [ ] Spinner appears when converting
- [ ] Dots animate continuously (500ms cycle)
- [ ] Convert button hides, Cancel button shows
- [ ] Cancel button triggers confirmation
- [ ] Confirming cancel aborts request
- [ ] Denying cancel continues conversion
- [ ] UI resets properly after completion
- [ ] UI resets properly after cancellation
- [ ] UI resets properly after error
- [ ] Multiple conversions work correctly
- [ ] Keyboard shortcuts still work
- [ ] No console errors
- [ ] No memory leaks (check dev tools)

## Edge Cases Handled

1. **Clicking Convert while already converting**: Disabled, can't click
2. **Clicking Cancel when not converting**: Function returns early, no effect
3. **Server timeout during conversion**: Fetch error caught, UI resets
4. **Network error during conversion**: Fetch error caught, UI resets
5. **Rapid clicking Cancel**: Confirmation prevents multiple aborts
6. **Converting empty text**: Validation prevents, shows error message
7. **Browser tab backgrounded**: Animation continues, no issues
8. **Page refresh during conversion**: Session lost, normal behavior

## Future Improvements

1. **ESC key to cancel**: Add keyboard listener
2. **Progress percentage**: Needs backend support
3. **Time estimate**: Based on text length
4. **Toast notifications**: Less intrusive than status bar
5. **Retry button**: After errors
6. **Queue system**: Multiple conversions
7. **Partial results**: Show N4L as it's generated
8. **Dark mode**: Adjust colors for dark theme
