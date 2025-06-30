import { defineStore } from 'pinia'
import apiClient from '@/api/client'

const USER_TOKEN_KEY = 'user_token'
const USER_KEY = 'user_data'
const STORAGE_TYPE_KEY = 'auth_storage_type'

// Helper functions for token and user management
const getStorage = () => {
  const storageType = localStorage.getItem(STORAGE_TYPE_KEY)
  return storageType === 'local' ? localStorage : sessionStorage
}

const getStoredToken = () => {
  return getStorage().getItem(USER_TOKEN_KEY)
}
const storeToken = (token, useLocalStorage) => {
  const storage = useLocalStorage ? localStorage : sessionStorage

  // Remember storage type preference
  localStorage.setItem(STORAGE_TYPE_KEY, useLocalStorage ? 'local' : 'session')

  // Store token in selected storage
  storage.setItem(USER_TOKEN_KEY, token)
}
const removeToken = () => {
  localStorage.removeItem(USER_TOKEN_KEY)
  sessionStorage.removeItem(USER_TOKEN_KEY)
}

const getStoredUser = () => {
  const userData = getStorage().getItem(USER_KEY)
  return userData ? JSON.parse(userData) : null
}
const storeUser = (user, useLocalStorage) => {
  const storage = useLocalStorage ? localStorage : sessionStorage
  storage.setItem(USER_KEY, JSON.stringify(user))
}
const removeUser = () => {
  localStorage.removeItem(USER_KEY)
  sessionStorage.removeItem(USER_KEY)
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: getStoredToken() || null,
    user: getStoredUser() || null,
    isAuthenticated: !!getStoredToken(),
    isLoading: false,
  }),

  getters: {
    currentUser: (state) => state.user,
    authToken: (state) => state.token,
    name: (state) => state.user?.name || '',
    userId: (state) => state.user?.id || state.token,
  },

  actions: {
    setAuthToken(token, useLocalStorage = true) {
      this.token = token
      this.isAuthenticated = !!token

      if (token) {
        storeToken(token, useLocalStorage)
      } else {
        removeToken()
      }
    },

    setUser(user, useLocalStorage = true) {
      this.user = user

      if (user) {
        storeUser(user, useLocalStorage)
      } else {
        removeUser()
      }
    },

    setLoading(isLoading) {
      this.isLoading = isLoading
    },

    // Log in user or create a new account with just a username
    async login({ username, rememberMe = false }) {
      this.setLoading(true)

      try {
        const response = await apiClient.post('/session', { name: username })

        const { id } = response.data

        // Create user object that includes all data from response
        const user = {
          id: id,
          name: username,
        }

        // The user ID serves as the authentication token in this system
        const token = id

        // Set token in API client for future requests
        apiClient.defaults.headers.common['Authorization'] = `Bearer ${token}`

        this.setAuthToken(token, rememberMe)
        this.setUser(user, rememberMe)

        return {
          id,
          username,
        }
      } catch (error) {
        this.setAuthToken(null)
        this.setUser(null)
        throw error
      } finally {
        this.setLoading(false)
      }
    },

    // Log out user
    async logout() {
      // Remove token from API client
      delete apiClient.defaults.headers.common['Authorization']

      // Clear authentication data
      this.setAuthToken(null)
      this.setUser(null)
    },

    // Check if user is already authenticated
    checkAuthStatus() {
      if (!this.token) {
        return false
      }

      return true
    },

    // Update user profile
    async updateProfile(userData) {
      this.setLoading(true)

      try {
        const response = await apiClient.put('/auth/profile', userData)

        // Update only the user data, keep the same token
        this.setUser(response.data)

        return response.data
      } catch (error) {
        throw error
      } finally {
        this.setLoading(false)
      }
    },
  },
})
