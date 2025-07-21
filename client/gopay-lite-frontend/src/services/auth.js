const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';

/**
 * Universal API fetch handler with timeout and error management
 * @param {string} endpoint - API endpoint (e.g., '/auth/login')
 * @param {RequestInit} options - Fetch options
 * @returns {Promise<any>}
 */
async function fetchAPI(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  
  // Default headers with override capability
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  // Timeout configuration (8 seconds)
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 8000);

  try {
    const response = await fetch(url, {
      ...options,
      headers,
      signal: controller.signal,
      credentials: 'include', // Required for cookies/sessions
    });
    clearTimeout(timeout);

    // Handle empty responses (e.g., 204 No Content)
    if (response.status === 204) return null;

    const data = await response.json();

    if (!response.ok) {
      const error = new Error(data?.message || `HTTP ${response.status}`);
      error.status = response.status;
      error.data = data; // Attach full error payload
      throw error;
    }

    return data;
  } catch (error) {
    clearTimeout(timeout);
    if (error.name === 'AbortError') {
      throw new Error('Request timed out (8s)');
    }
    throw error; // Re-throw for service-specific handling
  }
}

// ==================== AUTH SERVICES ====================

/**
 * Register a new user
 * @param {{
 *   name: string,
 *   email: string,
 *   password: string
 * }} userData - User registration data
 * @returns {Promise<{
 *   message: string,
 *   token?: string,
 *   user?: { id: string, email: string }
 * }>}
 */
export async function register(userData) {
  return fetchAPI('/auth/register', {
    method: 'POST',
    body: JSON.stringify(userData),
  });
}

/**
 * Authenticate user
 * @param {{
 *   email: string,
 *   password: string
 * }} credentials - Login credentials
 * @returns {Promise<{
 *   message: string,
 *   token: string,
 *   expiresIn?: number
 * }>}
 */
export async function login(credentials) {
  return fetchAPI('/auth/login', {
    method: 'POST',
    body: JSON.stringify(credentials),
  });
}

/**
 * Fetch current user profile (protected)
 * @param {string} token - JWT access token
 * @returns {Promise<{
 *   user_id: string,
 *   email: string,
 *   name: string,
 *   roles?: string[]
 * }>}
 */
export async function getMe(token) {
  return fetchAPI('/auth/me', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

/**
 * Invalidate user session
 * @returns {Promise<{ message: string }>}
 */
export async function logout() {
  return fetchAPI('/auth/logout', {
    method: 'POST',
  });
}

/**
 * Refresh access token
 * @param {string} refreshToken 
 * @returns {Promise<{ token: string, expiresIn: number }>}
 */
export async function refreshToken(refreshToken) {
  return fetchAPI('/auth/refresh', {
    method: 'POST',
    body: JSON.stringify({ refreshToken }),
  });
}

/**
 * Request password reset
 * @param {string} email 
 * @returns {Promise<{ message: string }>}
 */
export async function forgotPassword(email) {
  return fetchAPI('/auth/forgot-password', {
    method: 'POST',
    body: JSON.stringify({ email }),
  });
}

/**
 * Verify email with token
 * @param {string} token 
 * @returns {Promise<{ message: string }>}
 */
export async function verifyEmail(token) {
  return fetchAPI(`/auth/verify-email?token=${token}`);
}