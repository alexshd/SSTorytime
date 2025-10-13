// Arrow menu modal component

import { getValidArrowsList } from '../lib/arrows.js';

export function showArrowValidationGuide()
{
  const modal = createModal();
  document.body.appendChild(modal);

  setTimeout(() => modal.classList.add('opacity-100'), 10);
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
