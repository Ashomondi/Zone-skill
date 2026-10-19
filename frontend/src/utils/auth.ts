/**
 * Utility for handling authentication headers in API requests.
 * When VITE_DISABLE_AUTH is set to 'true', authentication is bypassed.
 */

export function getAuthHeaders(): Record<string, string> {
  const disableAuth = import.meta.env.VITE_DISABLE_AUTH === 'true';
  
  if (disableAuth) {
    return {
      'X-Disable-Auth': 'true',
    };
  }
  
  const token = localStorage.getItem('token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}
