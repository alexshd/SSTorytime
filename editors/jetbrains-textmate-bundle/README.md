# JetBrains TextMate Bundle for N4L & SST

This directory packages the existing TextMate grammars so JetBrains IDEs (IntelliJ IDEA, GoLand, WebStorm, etc.) can highlight `.n4l` and `.sst` files.

## Contents

- `syntaxes/n4l.tmLanguage.json`
- `syntaxes/sst.tmLanguage.json`

## Setup Instructions

1. Copy the two grammar files into this directory structure:
   ```
   jetbrains-textmate-bundle/
     syntaxes/
       n4l.tmLanguage.json
       sst.tmLanguage.json
   ```
2. In a JetBrains IDE:
   - Go to: Settings > Editor > TextMate Bundles
   - Click `+` and select the `jetbrains-textmate-bundle` directory
   - Apply & restart if needed
3. Associate file types:
   - Settings > Editor > File Types
   - Add `*.n4l` and `*.sst` patterns mapped to `TextMate` > `N4L` / `SST`

## Notes

- JetBrains uses the `scopeName` (e.g., `source.n4l`) for color theme mapping.
- You can customize colors by editing your IDE color scheme and mapping scopes.
- This is a lightweight integration; advanced features (structure view, inspections) would need a full plugin.
