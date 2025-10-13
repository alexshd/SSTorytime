// CodeMirror-based N4L editor with proper line numbers and syntax highlighting

import { EditorView, lineNumbers, highlightActiveLine, keymap, Decoration, ViewPlugin, WidgetType } from '@codemirror/view';
import { EditorState, RangeSetBuilder } from '@codemirror/state';
import { defaultKeymap } from '@codemirror/commands';
import { syntaxHighlighting, HighlightStyle, syntaxTree } from '@codemirror/language';
import { tags } from '@lezer/highlight';
import { n4lLanguage } from '../lib/n4l-language.js';
import { isValidArrow } from '../lib/arrows.js';
import { createArrowMenu } from './ArrowMenu.js';

export class N4LEditor
{
  constructor(container)
  {
    this.container = container;
    this.view = null;
    this.content = '';
    this.onChange = null;
    this.onArrowReplace = null;
    this.onArrowDelete = null;
    this.onScrollCallback = null;

    this.init();
  }

  init()
  {
    // N4L syntax highlighting theme
    const n4lHighlightStyle = HighlightStyle.define([
      { tag: tags.comment, color: '#6b7280', fontStyle: 'italic' },
      { tag: tags.heading, color: '#7c3aed', fontWeight: 'bold' },
      { tag: tags.keyword, color: '#ea580c', fontWeight: 'bold' },
      { tag: tags.variableName, color: '#0891b2', fontWeight: '600' },
      { tag: tags.string, color: '#059669', fontWeight: 'bold', fontSize: '1.1em' },
      { tag: tags.annotation, color: '#0891b2', fontWeight: '600' },
      { tag: tags.meta, color: '#8b5cf6', fontStyle: 'italic' },
      { tag: tags.operator, color: '#059669', fontWeight: '500' },
    ]);

    // Arrow decoration plugin
    const self = this; // Capture editor instance for arrow callbacks
    const arrowPlugin = ViewPlugin.fromClass(class
    {
      constructor(view)
      {
        this.decorations = this.buildDecorations(view);
      }

      update(update)
      {
        if (update.docChanged || update.viewportChanged)
        {
          this.decorations = this.buildDecorations(update.view);
        }
      }

      buildDecorations(view)
      {
        const builder = new RangeSetBuilder();
        const text = view.state.doc.toString();

        // Find all arrows in the text
        const arrowRegex = /\([a-z][a-z\s,;:.\-'\/]{3,}\)/gi;
        let match;

        while ((match = arrowRegex.exec(text)) !== null)
        {
          const arrow = match[0];
          const from = match.index;
          const to = from + arrow.length;
          const valid = isValidArrow(arrow);

          // Create a clickable decoration
          builder.add(
            from,
            to,
            Decoration.mark({
              class: valid ? 'cm-arrow-valid' : 'cm-arrow-invalid',
              attributes: {
                'data-arrow': arrow,
                'data-valid': valid,
                'style': `cursor: pointer; ${valid ? 'color: #059669; font-weight: 600;' : 'color: #dc2626; font-weight: 600; text-decoration: underline wavy;'}`
              }
            })
          );
        }

        return builder.finish();
      }
    }, {
      decorations: v => v.decorations,
      eventHandlers: {
        click: (e, view) =>
        {
          const target = e.target;
          if (target.hasAttribute('data-arrow'))
          {
            const arrow = target.getAttribute('data-arrow');
            const valid = target.getAttribute('data-valid') === 'true';

            // Show arrow menu
            const menu = createArrowMenu(
              target,
              (newArrow, menuEl) =>
              {
                // Replace arrow in editor
                const text = view.state.doc.toString();
                const newText = text.replace(arrow, newArrow);
                view.dispatch({
                  changes: { from: 0, to: text.length, insert: newText }
                });
                menuEl.remove();
                if (self.onArrowReplace)
                {
                  self.onArrowReplace(arrow, newArrow);
                }
              },
              (arrowToDelete, menuEl) =>
              {
                // Delete arrow from editor
                const text = view.state.doc.toString();
                const newText = text.replace(new RegExp(arrowToDelete.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g'), '');
                view.dispatch({
                  changes: { from: 0, to: text.length, insert: newText }
                });
                menuEl.remove();
                if (self.onArrowDelete)
                {
                  self.onArrowDelete(arrowToDelete);
                }
              }
            );

            // Position menu
            const rect = target.getBoundingClientRect();
            menu.style.left = Math.min(rect.left, window.innerWidth - 320) + 'px';
            menu.style.top = (rect.bottom + 5) + 'px';

            document.body.appendChild(menu);

            // Close on outside click
            setTimeout(() =>
            {
              document.addEventListener('click', function closeMenu(event)
              {
                if (!menu.contains(event.target))
                {
                  menu.remove();
                  document.removeEventListener('click', closeMenu);
                }
              });
            }, 100);

            e.preventDefault();
          }
        }
      }
    });

    // Create the editor state
    const startState = EditorState.create({
      doc: '',
      extensions: [
        lineNumbers(),
        highlightActiveLine(),
        EditorView.lineWrapping,
        n4lLanguage,
        syntaxHighlighting(n4lHighlightStyle),
        arrowPlugin,
        keymap.of(defaultKeymap),
        EditorView.updateListener.of((update) =>
        {
          if (update.docChanged)
          {
            this.content = update.state.doc.toString();
            if (this.onChange)
            {
              this.onChange(this.content);
            }
          }
        }),
        // Custom theme for N4L
        EditorView.theme({
          '&': {
            fontSize: '1rem',
            fontFamily: 'monospace',
            height: '100%'
          },
          '.cm-scroller': {
            overflow: 'auto'
          },
          '.cm-content': {
            fontFamily: 'monospace',
            padding: '0.75rem 0',
            caretColor: '#3b82f6',
            minHeight: '100%'
          },
          '.cm-line': {
            padding: '0 0.5rem',
            lineHeight: '1.5'
          },
          '.cm-gutters': {
            backgroundColor: '#f9fafb',
            color: '#6b7280',
            border: 'none',
            borderRight: '1px solid #e5e7eb',
            minWidth: '3rem',
            paddingRight: '0.5rem'
          },
          '.cm-lineNumbers .cm-gutterElement': {
            padding: '0 0.5rem',
            minWidth: '2rem',
            textAlign: 'right',
            fontSize: '0.875rem'
          },
          '.cm-activeLine': {
            backgroundColor: '#eff6ff'
          },
          '.cm-activeLineGutter': {
            backgroundColor: '#dbeafe'
          },
          '&.cm-focused': {
            outline: 'none'
          },
          '.cm-arrow-valid': {
            color: '#059669',
            fontWeight: '600',
            cursor: 'pointer'
          },
          '.cm-arrow-invalid': {
            color: '#dc2626',
            fontWeight: '600',
            textDecoration: 'underline wavy',
            cursor: 'pointer'
          }
        })
      ]
    });

    // Create the editor view
    this.view = new EditorView({
      state: startState,
      parent: this.container
    });

    // Add scroll event listener
    this.view.scrollDOM.addEventListener('scroll', () =>
    {
      if (this.onScrollCallback)
      {
        this.onScrollCallback();
      }
    });
  }

  setContent(text)
  {
    this.content = text;
    const transaction = this.view.state.update({
      changes: {
        from: 0,
        to: this.view.state.doc.length,
        insert: text
      }
    });
    this.view.dispatch(transaction);
  }

  getContent()
  {
    return this.content;
  }

  focus()
  {
    this.view.focus();
  }

  destroy()
  {
    if (this.view)
    {
      this.view.destroy();
    }
  }

  // Set callback for content changes
  setOnChange(callback)
  {
    this.onChange = callback;
  }

  // Set callbacks for arrow actions
  setArrowCallbacks(onReplace, onDelete)
  {
    this.onArrowReplace = onReplace;
    this.onArrowDelete = onDelete;
  }

  // Get the scroll DOM element for scroll sync
  getScrollDOM()
  {
    return this.view.scrollDOM;
  }

  // Scroll to a specific position (0-1 ratio)
  scrollToRatio(ratio)
  {
    const scrollDOM = this.view.scrollDOM;
    const maxScroll = scrollDOM.scrollHeight - scrollDOM.clientHeight;
    scrollDOM.scrollTop = ratio * maxScroll;
  }

  // Get current scroll ratio (0-1)
  getScrollRatio()
  {
    const scrollDOM = this.view.scrollDOM;
    const maxScroll = scrollDOM.scrollHeight - scrollDOM.clientHeight;
    return maxScroll > 0 ? scrollDOM.scrollTop / maxScroll : 0;
  }

  // Add scroll event listener
  onScroll(callback)
  {
    this.onScrollCallback = callback;
  }
}
