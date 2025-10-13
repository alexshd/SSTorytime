// Session management for editor state persistence

const SESSION_KEY = 'n4l-editor-session';
const SESSION_MAX_AGE = 7 * 24 * 60 * 60 * 1000; // 7 days

export function saveSession(data)
{
  try
  {
    const sessionData = {
      ...data,
      timestamp: new Date().toISOString()
    };
    localStorage.setItem(SESSION_KEY, JSON.stringify(sessionData));
    return true;
  } catch (error)
  {
    console.warn('Could not save session:', error);
    return false;
  }
}

export function loadSession()
{
  try
  {
    const saved = localStorage.getItem(SESSION_KEY);
    if (!saved) return null;

    const sessionData = JSON.parse(saved);

    // Check if session is less than 7 days old
    const sessionAge = Date.now() - new Date(sessionData.timestamp).getTime();

    if (sessionAge > SESSION_MAX_AGE)
    {
      localStorage.removeItem(SESSION_KEY);
      return null;
    }

    return sessionData;
  } catch (error)
  {
    console.warn('Could not load session:', error);
    return null;
  }
}

export function clearSession()
{
  try
  {
    localStorage.removeItem(SESSION_KEY);
    return true;
  } catch (error)
  {
    console.warn('Could not clear session:', error);
    return false;
  }
}
