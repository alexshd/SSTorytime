// File handling utilities

export function detectFileType(filename, content)
{
  const ext = filename.split('.').pop().toLowerCase();

  if (ext === 'html' || ext === 'htm')
  {
    return 'html';
  } else if (ext === 'md' || ext === 'markdown')
  {
    return 'markdown';
  } else if (ext === 'n4l')
  {
    return 'n4l';
  }

  // Try to detect by content
  if (content.includes('<html') || content.includes('<!DOCTYPE'))
  {
    return 'html';
  } else if (content.match(/^#{1,6}\s+/m) || content.includes('```'))
  {
    return 'markdown';
  }

  return 'text';
}

export function markdownToHtml(markdown)
{
  let html = markdown;

  // Headers
  html = html.replace(/^### (.*$)/gim, '<h3>$1</h3>');
  html = html.replace(/^## (.*$)/gim, '<h2>$1</h2>');
  html = html.replace(/^# (.*$)/gim, '<h1>$1</h1>');

  // Bold
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
  html = html.replace(/__(.*?)__/g, '<strong>$1</strong>');

  // Italic
  html = html.replace(/\*(.*?)\*/g, '<em>$1</em>');
  html = html.replace(/_(.*?)_/g, '<em>$1</em>');

  // Links
  html = html.replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2">$1</a>');

  // Line breaks
  html = html.replace(/\n\n/g, '</p><p>');
  html = html.replace(/\n/g, '<br>');

  // Code blocks
  html = html.replace(/```(.*?)```/gs, '<pre><code>$1</code></pre>');
  html = html.replace(/`(.*?)`/g, '<code>$1</code>');

  // Wrap in paragraphs if not already wrapped
  if (!html.startsWith('<'))
  {
    html = '<p>' + html + '</p>';
  }

  return html;
}

export function readFileAsText(file)
{
  return new Promise((resolve, reject) =>
  {
    const reader = new FileReader();
    reader.onload = (e) => resolve(e.target.result);
    reader.onerror = () => reject(new Error('Failed to read file'));
    reader.readAsText(file);
  });
}
