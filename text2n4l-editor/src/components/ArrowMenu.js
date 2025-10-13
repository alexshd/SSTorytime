// Arrow menu modal component

import { getValidArrowsList } from '../lib/arrows.js';

export function showArrowValidationGuide()
{
  const panel = createHelpPanel();
  const container = document.getElementById('help-panel-container');

  if (container)
  {
    // Toggle visibility
    if (container.style.display === 'none' || !container.style.display)
    {
      container.innerHTML = '';
      container.appendChild(panel);
      container.style.display = 'block';
    } else
    {
      container.style.display = 'none';
    }
  }
}

export function createArrowMenu(arrowSpan, replaceCallback, deleteCallback)
{
  const currentArrow = arrowSpan.getAttribute('data-arrow');
  const isValid = arrowSpan.getAttribute('data-valid') === 'true';

  const menu = document.createElement('div');
  menu.className = 'arrow-menu';
  menu.style.cssText = 'position: fixed; background: white; border: 2px solid #3b82f6; border-radius: 8px; padding: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.15); z-index: 1000; max-height: 400px; overflow-y: auto; min-width: 300px;';

  const title = document.createElement('div');
  title.style.cssText = 'font-weight: 600; margin-bottom: 8px; color: #1f2937; font-size: 14px;';
  title.innerHTML = (isValid ? '✓ Valid Arrow' : '⚠️ Invalid Arrow') + '<br><code style="color: #3b82f6;">' + currentArrow + '</code>';
  menu.appendChild(title);

  const actions = document.createElement('div');
  actions.style.cssText = 'display: flex; gap: 4px; margin-bottom: 12px; padding-bottom: 12px; border-bottom: 1px solid #e5e7eb;';

  const deleteBtn = createButton('🗑️ Delete', '#ef4444', () => deleteCallback(currentArrow, menu));
  actions.appendChild(deleteBtn);
  menu.appendChild(actions);

  if (!isValid)
  {
    const suggestions = createArrowSuggestions(currentArrow, replaceCallback, menu);
    menu.appendChild(suggestions);
  }

  const allArrowsSection = createAllArrowsSection(replaceCallback, menu);
  menu.appendChild(allArrowsSection);

  return menu;
}

function createButton(text, color, onClick)
{
  const btn = document.createElement('button');
  btn.textContent = text;
  btn.style.cssText = `flex: 1; padding: 6px 12px; background: ${color}; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 12px; font-weight: 500;`;
  btn.onclick = onClick;
  return btn;
}

function createArrowSuggestions(currentArrow, replaceCallback, menu)
{
  const container = document.createElement('div');
  container.style.cssText = 'margin-bottom: 12px;';

  const heading = document.createElement('div');
  heading.style.cssText = 'font-weight: 600; margin-bottom: 6px; color: #1f2937; font-size: 13px;';
  heading.textContent = '💡 Suggested Replacements:';
  container.appendChild(heading);

  const allArrows = getAllArrowsData();
  const cleanCurrent = currentArrow.replace(/[()]/g, '').toLowerCase();
  const words = cleanCurrent.split(/\s+/);

  const suggestedArrows = allArrows.filter(arr =>
    arr.keywords.some(keyword => words.some(word => keyword.toLowerCase().includes(word) || word.includes(keyword.toLowerCase())))
  );

  if (suggestedArrows.length > 0)
  {
    suggestedArrows.slice(0, 5).forEach(arr =>
    {
      container.appendChild(createArrowButton(arr, replaceCallback, menu, true));
    });
  } else
  {
    const noSuggestions = document.createElement('div');
    noSuggestions.style.cssText = 'color: #6b7280; font-size: 12px; font-style: italic; padding: 4px;';
    noSuggestions.textContent = 'No similar arrows found';
    container.appendChild(noSuggestions);
  }

  return container;
}

function createAllArrowsSection(replaceCallback, menu)
{
  const container = document.createElement('div');

  const heading = document.createElement('div');
  heading.style.cssText = 'font-weight: 600; margin-bottom: 6px; color: #1f2937; font-size: 13px;';
  heading.textContent = '📋 All Valid Arrows:';
  container.appendChild(heading);

  const allArrows = getAllArrowsData();
  const byType = groupArrowsByType(allArrows);

  ['NR-0', 'LT-1', 'CN-2', 'EP-3'].forEach(type =>
  {
    if (byType[type])
    {
      const typeHeader = document.createElement('div');
      typeHeader.style.cssText = 'font-weight: 500; margin: 8px 0 4px; color: #4b5563; font-size: 11px; text-transform: uppercase;';
      typeHeader.textContent = getTypeLabel(type);
      container.appendChild(typeHeader);

      byType[type].slice(0, 10).forEach(arr =>
      {
        container.appendChild(createArrowButton(arr, replaceCallback, menu, false));
      });
    }
  });

  return container;
}

function createArrowButton(arrow, replaceCallback, menu, isSuggestion)
{
  const btn = document.createElement('button');
  btn.style.cssText = `display: block; width: 100%; text-align: left; padding: 6px 8px; margin: 2px 0; background: ${isSuggestion ? '#dbeafe' : '#f3f4f6'}; border: 1px solid ${isSuggestion ? '#93c5fd' : '#d1d5db'}; border-radius: 4px; cursor: pointer; font-size: 12px; transition: all 0.2s;`;
  btn.innerHTML = `<strong>${arrow.symbol}</strong> <span style="color: #6b7280;">${arrow.example}</span>`;
  btn.onclick = () => replaceCallback(arrow.symbol, menu);
  btn.onmouseenter = () => btn.style.background = isSuggestion ? '#bfdbfe' : '#e5e7eb';
  btn.onmouseleave = () => btn.style.background = isSuggestion ? '#dbeafe' : '#f3f4f6';
  return btn;
}

function createModal()
{
  const modal = document.createElement('div');
  modal.style.cssText = 'position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 9999; opacity: 0; transition: opacity 0.3s;';

  const content = document.createElement('div');
  content.style.cssText = 'background: white; border-radius: 12px; padding: 24px; max-width: 800px; max-height: 80vh; overflow-y: auto; box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1);';

  content.innerHTML = `
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
      <h2 style="font-size: 24px; font-weight: 700; color: #1f2937;">🏹 N4L Arrow Validation Guide</h2>
      <button onclick="this.closest('div[style*=fixed]').remove()" style="background: #ef4444; color: white; border: none; border-radius: 6px; padding: 8px 16px; cursor: pointer; font-weight: 500;">Close</button>
    </div>
    <p style="color: #6b7280; margin-bottom: 16px;">Click on any highlighted arrow in the output to validate, replace, or delete it. Valid arrows are highlighted in <span style="background: #dbeafe; color: #1e40af; padding: 2px 6px; border-radius: 4px;">blue</span>, invalid arrows in <span style="background: #fee2e2; color: #991b1b; padding: 2px 6px; border-radius: 4px;">red</span>.</p>
    <div style="background: #f9fafb; border-radius: 8px; padding: 16px;">
      <h3 style="font-weight: 600; margin-bottom: 8px; color: #374151;">📝 300+ Valid N4L Arrows Available</h3>
      <p style="color: #6b7280; font-size: 14px;">All arrows are loaded from SSTconfig files and validated in real-time.</p>
    </div>
  `;

  modal.appendChild(content);
  modal.onclick = (e) => { if (e.target === modal) modal.remove(); };

  return modal;
}

function createHelpPanel()
{
  const panel = document.createElement('div');
  panel.style.cssText = 'background: white; border: 1px solid #e5e7eb; border-radius: 8px; padding: 20px; margin-bottom: 16px; max-height: 400px; overflow-y: auto;';

  const allArrows = getAllArrowsData();
  const byType = groupArrowsByType(allArrows);

  panel.innerHTML = `
    <div style="margin-bottom: 16px;">
      <h3 style="font-size: 18px; font-weight: 700; color: #1f2937; margin-bottom: 8px;">🏹 N4L Editing Reference Guide</h3>
      <p style="color: #6b7280; font-size: 14px; line-height: 1.6;">
        Complete reference for editing N4L files with <strong>${allArrows.length} valid arrows</strong> organized by semantic type.
      </p>
    </div>

    <div style="background: #eff6ff; border-left: 4px solid #3b82f6; padding: 12px; margin-bottom: 16px; border-radius: 4px;">
      <h4 style="font-weight: 600; color: #1e40af; margin-bottom: 4px;">📝 Editing Workflow</h4>
      <ol style="color: #1e3a8a; font-size: 13px; margin-left: 20px; line-height: 1.8;">
        <li><strong>Convert:</strong> Upload text/HTML/Markdown to generate initial N4L</li>
        <li><strong>Edit:</strong> Click invalid arrows (red) to see suggestions and replace</li>
        <li><strong>Enhance:</strong> Add semantic relationships using valid arrows below</li>
        <li><strong>Validate:</strong> All arrows are checked against SSTconfig files</li>
        <li><strong>Save:</strong> Download the validated .n4l file ready for upload</li>
      </ol>
    </div>

    <div style="background: #fef3c7; border-left: 4px solid #f59e0b; padding: 12px; margin-bottom: 16px; border-radius: 4px;">
      <h4 style="font-weight: 600; color: #92400e; margin-bottom: 4px;">⚠️ Common Invalid Arrows</h4>
      <p style="color: #78350f; font-size: 12px; margin-bottom: 8px;">The text2n4l converter generates natural language that needs correction:</p>
      <table style="width: 100%; font-size: 12px; color: #78350f;">
        <tr><td style="padding: 2px;"><code>(appears close to)</code></td><td>→</td><td><strong>(similar to)</strong></td></tr>
        <tr><td style="padding: 2px;"><code>(is similar to)</code></td><td>→</td><td><strong>(similar to)</strong></td></tr>
        <tr><td style="padding: 2px;"><code>(relates to)</code></td><td>→</td><td><strong>(associated with)</strong></td></tr>
        <tr><td style="padding: 2px;"><code>(mentioned in)</code></td><td>→</td><td><strong>(is discussed in)</strong></td></tr>
      </table>
    </div>

    <div style="margin-bottom: 16px;">
      <h4 style="font-weight: 600; color: #374151; margin-bottom: 12px; font-size: 16px;">📚 Valid Arrow Categories</h4>
      
      <details style="margin-bottom: 8px; border: 1px solid #e5e7eb; border-radius: 6px; padding: 8px;">
        <summary style="cursor: pointer; font-weight: 600; color: #dc2626; font-size: 14px;">
          🔴 NR-0: Similarity & Equivalence (${byType['NR-0'].length} arrows)
        </summary>
        <div style="margin-top: 8px; padding-left: 8px; max-height: 200px; overflow-y: auto;">
          ${byType['NR-0'].slice(0, 20).map(arr =>
    `<div style="padding: 4px; font-size: 12px; font-family: monospace; color: #4b5563;">${arr.symbol}</div>`
  ).join('')}
          ${byType['NR-0'].length > 20 ? `<div style="padding: 4px; font-size: 11px; color: #6b7280; font-style: italic;">...and ${byType['NR-0'].length - 20} more</div>` : ''}
        </div>
      </details>

      <details style="margin-bottom: 8px; border: 1px solid #e5e7eb; border-radius: 6px; padding: 8px;">
        <summary style="cursor: pointer; font-weight: 600; color: #ea580c; font-size: 14px;">
          🟠 LT-1: Causality & Temporal (${byType['LT-1'].length} arrows)
        </summary>
        <div style="margin-top: 8px; padding-left: 8px; max-height: 200px; overflow-y: auto;">
          ${byType['LT-1'].slice(0, 20).map(arr =>
    `<div style="padding: 4px; font-size: 12px; font-family: monospace; color: #4b5563;">${arr.symbol}</div>`
  ).join('')}
          ${byType['LT-1'].length > 20 ? `<div style="padding: 4px; font-size: 11px; color: #6b7280; font-style: italic;">...and ${byType['LT-1'].length - 20} more</div>` : ''}
        </div>
      </details>

      <details style="margin-bottom: 8px; border: 1px solid #e5e7eb; border-radius: 6px; padding: 8px;">
        <summary style="cursor: pointer; font-weight: 600; color: #0891b2; font-size: 14px;">
          🔵 CN-2: Containment & Structure (${byType['CN-2'].length} arrows)
        </summary>
        <div style="margin-top: 8px; padding-left: 8px; max-height: 200px; overflow-y: auto;">
          ${byType['CN-2'].slice(0, 20).map(arr =>
    `<div style="padding: 4px; font-size: 12px; font-family: monospace; color: #4b5563;">${arr.symbol}</div>`
  ).join('')}
          ${byType['CN-2'].length > 20 ? `<div style="padding: 4px; font-size: 11px; color: #6b7280; font-style: italic;">...and ${byType['CN-2'].length - 20} more</div>` : ''}
        </div>
      </details>

      <details style="margin-bottom: 8px; border: 1px solid #e5e7eb; border-radius: 6px; padding: 8px;">
        <summary style="cursor: pointer; font-weight: 600; color: #7c3aed; font-size: 14px;">
          🟣 EP-3: Expression & Properties (${byType['EP-3'].length} arrows)
        </summary>
        <div style="margin-top: 8px; padding-left: 8px; max-height: 200px; overflow-y: auto;">
          ${byType['EP-3'].slice(0, 20).map(arr =>
    `<div style="padding: 4px; font-size: 12px; font-family: monospace; color: #4b5563;">${arr.symbol}</div>`
  ).join('')}
          ${byType['EP-3'].length > 20 ? `<div style="padding: 4px; font-size: 11px; color: #6b7280; font-style: italic;">...and ${byType['EP-3'].length - 20} more</div>` : ''}
        </div>
      </details>
    </div>

    <div style="background: #f0fdf4; border-left: 4px solid #10b981; padding: 12px; border-radius: 4px;">
      <h4 style="font-weight: 600; color: #065f46; margin-bottom: 4px;">💡 Pro Tips</h4>
      <ul style="color: #047857; font-size: 12px; margin-left: 20px; line-height: 1.8;">
        <li><strong>Click any arrow</strong> in the output to validate and replace</li>
        <li><strong>Red underlined arrows</strong> are invalid - replace them first</li>
        <li><strong>Use suggestions</strong> - the menu shows similar valid arrows</li>
        <li><strong>Line numbers</strong> match your source - wrapped lines keep same number</li>
        <li><strong>Scroll sync</strong> - input and output scroll together for easy reference</li>
      </ul>
    </div>
  `;

  return panel;
}


function getAllArrowsData()
{
  const validArrows = getValidArrowsList();
  return validArrows.map(arrow => ({
    symbol: `(${arrow})`,
    example: arrow,
    keywords: arrow.split(/\s+/)
  }));
}

function groupArrowsByType(arrows)
{
  // Simple heuristic grouping - you can improve this
  const groups = { 'NR-0': [], 'LT-1': [], 'CN-2': [], 'EP-3': [] };

  arrows.forEach(arr =>
  {
    const text = arr.example.toLowerCase();
    if (text.includes('similar') || text.includes('same') || text.includes('equals'))
    {
      groups['NR-0'].push(arr);
    } else if (text.includes('leads') || text.includes('causes') || text.includes('result'))
    {
      groups['LT-1'].push(arr);
    } else if (text.includes('contains') || text.includes('part of') || text.includes('member'))
    {
      groups['CN-2'].push(arr);
    } else
    {
      groups['EP-3'].push(arr);
    }
  });

  return groups;
}

function getTypeLabel(type)
{
  const labels = {
    'NR-0': 'Similarity/Equivalence',
    'LT-1': 'Causality/Temporal',
    'CN-2': 'Containment/Structure',
    'EP-3': 'Expression/Properties'
  };
  return labels[type] || type;
}
