# N4L syntax highlighting for Emacs
# Place this in your .emacs or as a separate file and load it

(defvar n4l-mode-syntax-table
  (let ((table (make-syntax-table)))
    ;; Comments
    (modify-syntax-entry ?# "<" table)
    (modify-syntax-entry ?\n ">" table)
    ;; Strings
    (modify-syntax-entry ?\" "\"" table)
    (modify-syntax-entry ?' "\"" table)
    ;; Parentheses for relationships
    (modify-syntax-entry ?\( "()" table)
    (modify-syntax-entry ?\) ")(" table)
    table)
  "Syntax table for N4L mode.")

(defvar n4l-font-lock-keywords
  '(
    ;; Comments
    ("^[[:space:]]*\\(#\\|//\\).*$" . font-lock-comment-face)
    
    ;; Section headers
    ("^[[:space:]]*-[[:space:]]*\\(.+\\)$" 
     (1 font-lock-function-name-face))
    
    ;; Context declarations
    ("^[[:space:]]*\\([+\\-]?\\)\\(:{1,}\\)[[:space:]]*\\(.+?\\)[[:space:]]*\\(:{1,}\\)"
     (1 font-lock-keyword-face)
     (2 font-lock-delimiter-face)
     (3 font-lock-variable-name-face)
     (4 font-lock-delimiter-face))
    
    ;; Relationships
    ("(\\([^)]+\\))" 
     (1 font-lock-builtin-face))
    
    ;; Aliases
    ("@\\([a-zA-Z_][a-zA-Z0-9_]*\\)"
     (1 font-lock-constant-face))
    
    ;; References
    ("\\$\\([a-zA-Z_][a-zA-Z0-9_]*\\(?:\\.[0-9]+\\)?\\|[0-9]+\\|PREV\\.[0-9]+\\)"
     (1 font-lock-variable-name-face))
    
    ;; Special markers
    ("%\\([a-zA-Z_][a-zA-Z0-9_\"[:space:]]*\\)"
     (1 font-lock-type-face))
    ("=\\([a-zA-Z_][a-zA-Z0-9_]*\\)"
     (1 font-lock-warning-face))
    ("\\*\\([a-zA-Z_][a-zA-Z0-9_]*\\)"
     (1 font-lock-keyword-face))
    
    ;; Strings
    ("\"[^\"]*\"" . font-lock-string-face)
    ("'[^']*'" . font-lock-string-face)
    
    ;; URLs
    ("https?://[^[:space:]\"')}]+" . font-lock-doc-face)
    
    ;; TODO items
    ("^[[:space:]]*[A-Z][A-Z0-9[:space:]]+[A-Z0-9][[:space:]]*$" . font-lock-warning-face)
    
    ;; Continuation
    ("^[[:space:]]*\"[[:space:]]*" . font-lock-preprocessor-face))
  "Font lock keywords for N4L mode.")

(defvar n4l-mode-map
  (let ((map (make-sparse-keymap)))
    ;; Add key bindings here if needed
    map)
  "Keymap for N4L mode.")

(define-derived-mode n4l-mode text-mode "N4L"
  "Major mode for editing N4L (Notes for Learning) files."
  :syntax-table n4l-mode-syntax-table
  (setq-local font-lock-defaults '(n4l-font-lock-keywords))
  (setq-local comment-start "# ")
  (setq-local comment-end "")
  (setq-local comment-start-skip "\\(#\\|//\\)\\s-*")
  
  ;; Indentation
  (setq-local indent-tabs-mode nil)
  (setq-local tab-width 2)
  
  ;; Auto-completion setup (if available)
  (when (featurep 'company)
    (add-to-list 'company-backends 'company-n4l-relations)))

;; Auto-mode association
(add-to-list 'auto-mode-alist '("\\.n4l\\'" . n4l-mode))

;; Optional: Define company completion backend for relations
(defun company-n4l-relations (command &optional arg &rest ignored)
  "Company backend for N4L relation completion."
  (interactive (list 'interactive))
  (cl-case command
    (interactive (company-begin-backend 'company-n4l-relations))
    (prefix (and (eq major-mode 'n4l-mode)
                 (company-grab-symbol)))
    (candidates 
     (when (string-match-p "^(" (company-grab-line-to-point))
       '("note" "e.g." "ex" "prop" "rule" "then" "url" "img" 
         "he" "ph" "hp" "eh" "pe" "about" "expr" "source"
         "cond" "attr" "descr" "version" "title")))))

(provide 'n4l-mode)