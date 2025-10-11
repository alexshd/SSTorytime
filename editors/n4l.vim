# N4L syntax highlighting for Vim
# Place this file in ~/.vim/syntax/n4l.vim or create a plugin

if exists("b:current_syntax")
  finish
endif

" Comments
syntax match n4lComment "^[[:space:]]*#.*$"
syntax match n4lComment "^[[:space:]]*//.*$"

" Section headers
syntax match n4lSection "^[[:space:]]*-[[:space:]]*.*$"

" Context declarations
syntax match n4lContext "^[[:space:]]*[+\-]\?:{1,}[[:space:]]*.*[[:space:]]*:{1,}[[:space:]]*$"
syntax match n4lContextOp "^[[:space:]]*[+\-]" contained

" Relationships in parentheses
syntax region n4lRelation start="(" end=")" oneline

" Aliases and references
syntax match n4lAlias "@[a-zA-Z_][a-zA-Z0-9_]*"
syntax match n4lReference "\$[a-zA-Z_][a-zA-Z0-9_]*\(\.[0-9]\+\)\?"
syntax match n4lReference "\$[0-9]\+"
syntax match n4lReference "\$PREV\.[0-9]\+"

" Special markers
syntax match n4lConcept "%[a-zA-Z_][a-zA-Z0-9_\"[:space:]]*"
syntax match n4lSpecial "=[a-zA-Z_][a-zA-Z0-9_]*"
syntax match n4lEmphasis "\*[a-zA-Z_][a-zA-Z0-9_]*"

" Strings
syntax region n4lString start='"' end='"' oneline
syntax region n4lString start="'" end="'" oneline

" URLs
syntax match n4lURL "https\?://[^[:space:]\"')}]\+"

" TODO items (all caps)
syntax match n4lTodo "^[[:space:]]*[A-Z][A-Z0-9[:space:]]\+[A-Z0-9][[:space:]]*$"

" Continuation markers
syntax match n4lContinuation "^[[:space:]]*\"[[:space:]]*"

" Define highlighting
highlight default link n4lComment Comment
highlight default link n4lSection Title
highlight default link n4lContext Tag
highlight default link n4lContextOp Operator
highlight default link n4lRelation Function
highlight default link n4lAlias Constant
highlight default link n4lReference Identifier
highlight default link n4lConcept Type
highlight default link n4lSpecial Special
highlight default link n4lEmphasis Underlined
highlight default link n4lString String
highlight default link n4lURL Underlined
highlight default link n4lTodo Todo
highlight default link n4lContinuation Operator

let b:current_syntax = "n4l"