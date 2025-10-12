// --- Auto-refresh: poll reload.txt every 2s ---
let lastReload = null;
setInterval(function ()
{
  fetch('/static/../tmp/reload.txt', { cache: 'no-store' })
    .then(r => r.text())
    .then(txt =>
    {
      if (lastReload && txt.trim() !== lastReload)
      {
        window.location.reload();
      }
      lastReload = txt.trim();
    });
}, 2000);

// --- Session persistence using localStorage ---
function saveSession(data)
{
  localStorage.setItem('n4l_session', JSON.stringify(data));
}
function loadSession()
{
  let s = localStorage.getItem('n4l_session');
  if (!s) return null;
  try { return JSON.parse(s); } catch { return null; }
}
function clearSession()
{
  localStorage.removeItem('n4l_session');
}

// On file upload, save session info
document.addEventListener('htmx:afterOnLoad', function (evt)
{
  if (evt.detail && evt.detail.target && evt.detail.target.id === 'upload-result')
  {
    // Parse filename and batch count from upload-result
    let html = evt.detail.target.innerHTML;
    let m = html.match(/File uploaded: ([^<]+) \((\d+) batches?\)/);
    if (m)
    {
      saveSession({ filename: m[1], batch: 0 });
    }
  }
});

// On batch navigation, update session
document.addEventListener('htmx:afterOnLoad', function (evt)
{
  if (evt.detail && evt.detail.target && evt.detail.target.id === 'editor-container')
  {
    // Try to extract batch index from header
    let html = evt.detail.target.innerHTML;
    let m = html.match(/Batch (\d+) of (\d+)/);
    let session = loadSession();
    if (m && session)
    {
      session.batch = parseInt(m[1], 10);
      saveSession(session);
    }
  }
});

// Auto-submit upload form when file is chosen
document.addEventListener('DOMContentLoaded', function ()
{
  var uploadForm = document.getElementById('upload-form');
  var fileInput = document.getElementById('file-input');
  if (uploadForm && fileInput)
  {
    fileInput.addEventListener('change', function ()
    {
      if (fileInput.files && fileInput.files.length > 0)
      {
        uploadForm.requestSubmit();
      }
    });
  }
  // Clear button logic: warn, suggest download, clear session, reset UI via HTMX
  var clearBtn = document.getElementById('clear-session-btn');
  if (clearBtn)
  {
    clearBtn.onclick = function ()
    {
      if (confirm('Are you sure you want to clear your work?\n\nAll unsaved changes will be lost.\n\nTip: Download your N4L output first if you want to keep your work.'))
      {
        clearSession();
        // Use HTMX to reload the upload form area only (SPA style)
        htmx.ajax('GET', '/', { target: 'body', swap: 'innerHTML' });
      }
    };
  }
});
