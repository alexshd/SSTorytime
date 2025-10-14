# tree-sitter-n4l

Experimental Tree-sitter grammar for the N4L (Notes for Learning) language.

## Status

Early draft: captures sections (- heading), context markers (:: ... ::), aliases (@name), numeric and alias references ($1, $foo), relation blocks ((relation stuff)), TODO blocks, and generic statements. Comments (#, //) are recognized.

## Development

```bash
npm install
npx tree-sitter generate
npx tree-sitter test
```

## NeoVim Usage (manual)

1. Build the parser:
   ```bash
   npm install
   npx tree-sitter generate
   cc -I./src -shared -fPIC src/parser.c -o parser.so
   ```
2. Copy `parser.so` to: `~/.config/nvim/parsers/n4l.so`
3. Add to your Neovim `init.lua` (using nvim-treesitter):

   ```lua
   require('nvim-treesitter.parsers').get_parser_configs().n4l = {
     install_info = {
       url = '/absolute/path/to/SSTorytime/editors/tree-sitter-n4l',
       files = { 'src/parser.c' },
       generate_requires_npm = true,
       requires_generate_from_grammar = true,
     },
     filetype = 'n4l',
   }

   vim.filetype.add({ extension = { n4l = 'n4l' } })
   ```

4. Open an `.n4l` file and run `:TSInstall n4l` (if using custom local config install_info).

## Next Steps

- Add captures for emphasis (%concept, =special, \*emphasis)
- Multi-line continuation and string literal improvements
- Query files (highlights.scm, folds.scm)
