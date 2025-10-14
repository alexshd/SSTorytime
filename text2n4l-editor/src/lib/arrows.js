// Dynamic arrow loading from SSTconfig on the server side.
// Assumes SSTconfig is served statically at /SSTconfig/ (adjust path if needed)

let cachedArrows = [];
let loaded = false;

async function fetchArrowFile(file)
{
  const resp = await fetch(`/SSTconfig/${file}`);
  if (!resp.ok) return '';
  return await resp.text();
}

function parseArrowTokens(text)
{
  const tokens = [];
  const lines = text.split(/\r?\n/);
  for (const raw of lines)
  {
    const line = raw.trim();
    if (!line) continue;
    if (line.startsWith('#') || line.startsWith('//')) continue;
    const matches = [...line.matchAll(/\(([a-zA-Z0-9_!<>\-=\/\s\.']{1,40})\)/g)].map(m => m[1].trim());
    for (const t of matches) tokens.push(t);
  }
  return tokens;
}

export async function loadArrows()
{
  if (loaded && cachedArrows.length) return cachedArrows;
  const files = [
    'arrows-NR-0.sst',
    'arrows-LT-1.sst',
    'arrows-CN-2.sst',
    'arrows-EP-3.sst'
  ];
  const all = [];
  for (const f of files)
  {
    try
    {
      const txt = await fetchArrowFile(f);
      all.push(...parseArrowTokens(txt));
    } catch (e)
    {
      console.warn('Failed loading', f, e);
    }
  }
  cachedArrows = Array.from(new Set(all));
  loaded = true;
  return cachedArrows;
}

export function getValidArrowsList()
{
  return cachedArrows;
}

export function isValidArrow(arrowText)
{
  const clean = arrowText.replace(/^\(|\)$/g, '').trim();
  return cachedArrows.includes(clean);
}

// Utility: find synonym suggestions (tokens appearing on same line) - simple heuristic can be added later
export function findSuggestions(token)
{
  const lc = token.toLowerCase();
  return cachedArrows.filter(t => t.toLowerCase().includes(lc) && t !== token).slice(0, 12);
}
