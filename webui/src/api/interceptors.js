// src/api/interceptors.js

const USER_TOKEN_KEY = 'user_token'
const STORAGE_TYPE_KEY = 'auth_storage_type'

// Helper function to get the correct storage and token
function getAuthToken() {
  // First check which storage type to use
  const storageType = localStorage.getItem(STORAGE_TYPE_KEY)
  const storage = storageType === 'local' ? localStorage : sessionStorage

  // Then get the token from the appropriate storage
  return storage.getItem(USER_TOKEN_KEY)
}

export function setupInterceptors(apiClient) {
  // Request interceptor
  apiClient.interceptors.request.use(
    (config) => {
      // Get token using the same logic as the auth store
      const token = getAuthToken()

      // If token exists, add to headers
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }

      return config
    },
    (error) => Promise.reject(error),
  )

  // Response interceptor
  apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
      // Handle 401 Unauthorized errors
      if (error.response && error.response.status === 401) {
        // Clear token from both storages to be safe
        localStorage.removeItem(USER_TOKEN_KEY)
        sessionStorage.removeItem(USER_TOKEN_KEY)

        // Only redirect if we're not already on the login page
        // This helps prevent redirect loops
        if (!window.location.pathname.includes('/login')) {
          // Save the current location to redirect back after login
          const currentPath = window.location.pathname
          if (currentPath !== '/' && !currentPath.includes('/login')) {
            sessionStorage.setItem('redirect_after_login', currentPath)
          }

          window.location.href = '/login'
        }
      }
      return Promise.reject(error)
    },
  )
}
