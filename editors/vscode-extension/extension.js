const vscode = require('vscode');
const fs = require('fs');
const path = require('path');

// Dynamic arrow data (loaded from SSTconfig/arrows-*.sst)
let arrowTokens = new Set();            // Set of all tokens in parentheses
let arrowSynonyms = new Map();          // token -> Set(tokens from same line)
let arrowCategories = new Map();        // category -> Set(tokens)
let arrowLineMap = new Map();           // token -> original line (for hover)

const ARROW_FILE_PREFIX = 'arrows-';

function loadArrowDefinitions(workspaceRoot)
{
  try
  {
    arrowTokens.clear();
    arrowSynonyms.clear();
    arrowCategories.clear();
    arrowLineMap.clear();

    const sstDir = path.join(workspaceRoot, 'SSTconfig');
    if (!fs.existsSync(sstDir)) return;
    const files = fs.readdirSync(sstDir).filter(f => f.startsWith(ARROW_FILE_PREFIX) && f.endsWith('.sst'));
    for (const file of files)
    {
      const categoryMatch = file.match(/arrows-([A-Z]+-[0-9]+)/);
      const category = categoryMatch ? categoryMatch[1] : 'misc';
      const content = fs.readFileSync(path.join(sstDir, file), 'utf8');
      const lines = content.split(/\r?\n/);
      for (const raw of lines)
      {
        const line = raw.trim();
        if (!line) continue;
        if (line.startsWith('#') || line.startsWith('//')) continue;
        // find groups like (token)
        const matches = [...line.matchAll(/\(([a-zA-Z0-9_!<>\-=\/\s\.']{1,40})\)/g)].map(m => m[1].trim());
        if (matches.length === 0) continue;
        // treat tokens in same line as synonyms / alternatives
        for (const tok of matches)
        {
          arrowTokens.add(tok);
          if (!arrowSynonyms.has(tok)) arrowSynonyms.set(tok, new Set());
          if (!arrowCategories.has(category)) arrowCategories.set(category, new Set());
          arrowCategories.get(category).add(tok);
          if (!arrowLineMap.has(tok)) arrowLineMap.set(tok, raw.trim());
        }
        for (const a of matches)
        {
          for (const b of matches)
          {
            if (a !== b) arrowSynonyms.get(a).add(b);
          }
        }
      }
    }
    console.log(`Loaded ${arrowTokens.size} N4L arrow tokens from SSTconfig.`);
  } catch (e)
  {
    console.error('Failed to load arrow definitions', e);
  }
}

function isValidArrow(arrowText)
{
  const cleanArrow = arrowText.replace(/^\(|\)$/g, '').trim();
  return arrowTokens.has(cleanArrow);
}

function findArrowSuggestions(arrowText)
{
  const base = arrowText.replace(/^\(|\)$/g, '').trim();
  if (arrowSynonyms.has(base) && arrowSynonyms.get(base).size > 0)
  {
    return [...arrowSynonyms.get(base)].slice(0, 10);
  }
  const lc = base.toLowerCase();
  return [...arrowTokens].filter(t => t.toLowerCase().includes(lc) && t !== base).slice(0, 10);
}

class N4LDiagnosticsProvider
{
  constructor()
  {
    this.diagnostics = vscode.languages.createDiagnosticCollection('n4l');
  }

  provideDiagnostics(document)
  {
    if (document.languageId !== 'n4l' && document.languageId !== 'sst')
    {
      return;
    }

    const diagnostics = [];
    const text = document.getText();
    const lines = text.split('\n');

    // Check for invalid arrows (capture broader token charset)
    const arrowRegex = /\(([a-zA-Z0-9_!<>\-=\/\s\.']{1,40})\)/g;

    for (let i = 0; i < lines.length; i++)
    {
      const line = lines[i];
      let match;

      while ((match = arrowRegex.exec(line)) !== null)
      {
        const arrowText = match[0];

        if (!isValidArrow(arrowText))
        {
          const startPos = new vscode.Position(i, match.index);
          const endPos = new vscode.Position(i, match.index + arrowText.length);
          const range = new vscode.Range(startPos, endPos);
          const diagnostic = new vscode.Diagnostic(
            range,
            `Unrecognized N4L arrow token: ${arrowText}. Run 'Reload N4L Arrows' if recently added.`,
            vscode.DiagnosticSeverity.Warning
          );
          diagnostic.code = 'invalid-arrow';
          diagnostic.source = 'n4l';
          diagnostics.push(diagnostic);
        }
      }
    }

    this.diagnostics.set(document.uri, diagnostics);
  }

  dispose()
  {
    this.diagnostics.dispose();
  }
}

class N4LCodeActionProvider
{
  provideCodeActions(document, range, context)
  {
    const actions = [];
    // Provide refactors for all arrows on the line (even valid) to switch synonyms
    const lineText = document.lineAt(range.start.line).text;
    const arrowRegex = /\(([a-zA-Z0-9_!<>\-=\/\s\.']{1,40})\)/g;
    let m;
    while ((m = arrowRegex.exec(lineText)) !== null)
    {
      const token = m[1];
      const start = new vscode.Position(range.start.line, m.index);
      const end = new vscode.Position(range.start.line, m.index + m[0].length);
      const tokenRange = new vscode.Range(start, end);
      const suggestions = findArrowSuggestions(token).slice(0, 8);
      for (const s of suggestions)
      {
        const action = new vscode.CodeAction(`Replace (${token}) → (${s})`, vscode.CodeActionKind.RefactorRewrite);
        action.edit = new vscode.WorkspaceEdit();
        action.edit.replace(document.uri, tokenRange, `(${s})`);
        actions.push(action);
      }
    }
    // Quick fixes for invalid tokens
    for (const diagnostic of context.diagnostics)
    {
      if (diagnostic.code === 'invalid-arrow')
      {
        const arrowText = document.getText(diagnostic.range);
        const suggestions = findArrowSuggestions(arrowText);
        for (const s of suggestions)
        {
          const fix = new vscode.CodeAction(`Change to (${s})`, vscode.CodeActionKind.QuickFix);
          fix.edit = new vscode.WorkspaceEdit();
          fix.edit.replace(document.uri, diagnostic.range, `(${s})`);
          fix.diagnostics = [diagnostic];
          actions.push(fix);
        }
        const del = new vscode.CodeAction('Remove unknown arrow', vscode.CodeActionKind.QuickFix);
        del.edit = new vscode.WorkspaceEdit();
        del.edit.replace(document.uri, diagnostic.range, '');
        del.diagnostics = [diagnostic];
        actions.push(del);
      }
    }
    return actions;
  }
}

class N4LHoverProvider
{
  provideHover(document, position)
  {
    const range = document.getWordRangeAtPosition(position, /\(([a-zA-Z0-9_!<>\-=\/\s\.']{1,40})\)/);
    if (!range) return;
    const arrowText = document.getText(range);
    const token = arrowText.replace(/^\(|\)$/g, '').trim();
    const valid = isValidArrow(arrowText);
    const synonyms = arrowSynonyms.get(token) ? [...arrowSynonyms.get(token)].sort() : [];
    const catEntries = [...arrowCategories.entries()].filter(([c, set]) => set.has(token));
    const categories = catEntries.map(([c]) => c).join(', ');
    let md = valid ? `✅ **N4L Arrow** \`(${token})\`` : `⚠️ **Unknown N4L Arrow** \`(${token})\``;
    if (categories) md += `\n\nCategory: ${categories}`;
    const origin = arrowLineMap.get(token);
    if (origin)
    {
      md += `\n\nSource: \`${origin.replace(/`/g, '')}\``;
    }
    if (synonyms.length > 0)
    {
      md += `\n\n**Alternatives:**\n` + synonyms.map(s => `- \`(${s})\``).join('\n');
    }
    if (!valid)
    {
      const suggestions = findArrowSuggestions(token);
      if (suggestions.length > 0)
      {
        md += `\n\n**Suggestions:**\n` + suggestions.map(s => `- \`(${s})\``).join('\n');
      }
    }
    return new vscode.Hover(new vscode.MarkdownString(md));
  }
}

function activate(context)
{
  console.log('N4L Language Support extension activated');
  const root = vscode.workspace.workspaceFolders && vscode.workspace.workspaceFolders[0] ? vscode.workspace.workspaceFolders[0].uri.fsPath : undefined;
  if (root) loadArrowDefinitions(root);

  const diagnosticsProvider = new N4LDiagnosticsProvider();
  const onDidChangeTextDocument = vscode.workspace.onDidChangeTextDocument(event =>
  {
    if (event.document.languageId === 'n4l' || event.document.languageId === 'sst')
    {
      diagnosticsProvider.provideDiagnostics(event.document);
    }
  });
  const onDidOpenTextDocument = vscode.workspace.onDidOpenTextDocument(doc => diagnosticsProvider.provideDiagnostics(doc));

  const codeActionProvider = vscode.languages.registerCodeActionsProvider(['n4l', 'sst'], new N4LCodeActionProvider());
  const hoverProvider = vscode.languages.registerHoverProvider(['n4l', 'sst'], new N4LHoverProvider());

  const validateArrowsCommand = vscode.commands.registerCommand('n4l.validateArrows', () =>
  {
    const editor = vscode.window.activeTextEditor;
    if (editor) diagnosticsProvider.provideDiagnostics(editor.document);
    vscode.window.showInformationMessage('N4L arrow validation complete');
  });

  const showArrowListCommand = vscode.commands.registerCommand('n4l.showArrowList', () =>
  {
    const panel = vscode.window.createWebviewPanel('n4lArrows', 'N4L Arrows', vscode.ViewColumn.Beside, { enableScripts: false });
    let html = '<h1>N4L Arrow Tokens</h1>';
    for (const [cat, set] of arrowCategories.entries())
    {
      html += `<h2>${cat}</h2><div style="line-height:1.6">${[...set].sort().map(t => `<code>(${t})</code>`).join(' ')}</div>`;
    }
    panel.webview.html = html;
  });

  const reloadArrowsCommand = vscode.commands.registerCommand('n4l.reloadArrows', () =>
  {
    if (root)
    {
      loadArrowDefinitions(root);
      const editor = vscode.window.activeTextEditor;
      if (editor) diagnosticsProvider.provideDiagnostics(editor.document);
      vscode.window.showInformationMessage('Reloaded N4L arrow definitions.');
    }
  });

  context.subscriptions.push(
    diagnosticsProvider,
    onDidChangeTextDocument,
    onDidOpenTextDocument,
    codeActionProvider,
    hoverProvider,
    validateArrowsCommand,
    showArrowListCommand,
    reloadArrowsCommand
  );
}

function deactivate()
{
  console.log('N4L Language Support extension deactivated');
}

module.exports = {
  activate,
  deactivate
};