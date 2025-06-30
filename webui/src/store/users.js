import { defineStore } from 'pinia'
import usersApi from '@/api/endpoints/users'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    isLoading: false,
    error: null,
    users: [],
  }),

  getters: {
    currentUser: (state) => state.user,
    isAuthenticated: (state) => !!state.user,
    allUsers: (state) => state.users,
    userDisplayName: (state) => {
      if (!state.user) return ''
      return state.user.name || 'Unknown User'
    },
    userProfilePhoto: (state) => {
      if (!state.user || !state.user.photo) return null
      return state.user.photo
    },
  },

  actions: {
    // Set user information (typically called after login)
    setUser(userData) {
      this.user = userData
    },

    // Clear user data (typically called during logout)
    clearUser() {
      this.user = null
    },

    // Fetch current user details (from token)
    async fetchCurrentUser() {
      this.isLoading = true
      this.error = null

      try {
        const response = await usersApi.fetchCurrentUser()
        this.user = response.data
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to fetch user profile'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Fetch all users
    async fetchUsers() {
      this.isLoading = true
      this.error = null

      try {
        const response = await usersApi.fetchUsers()
        this.users = response.data
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to fetch users'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Update user's username
    async updateUsername(name) {
      this.isLoading = true
      this.error = null

      try {
        const response = await usersApi.updateUsername(name)

        // Update local user data with new name
        if (this.user) {
          this.user = { ...this.user, name: response.data.name }
        }

        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to update username'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Upload user profile photo
    async uploadProfilePhoto(photoFile) {
      this.isLoading = true
      this.error = null

      try {
        const response = await usersApi.uploadPhoto(photoFile)

        // Update local user data with new photo URL
        if (this.user) {
          this.user = { ...this.user, photo: response.data.photo }
        }

        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to upload profile photo'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Update user status (online, offline, etc.)
    async updateStatus(status) {
      this.isLoading = true
      this.error = null

      try {
        // This assumes you have an endpoint for updating status
        // You may need to add this to your usersApi if it doesn't exist
        const response = await usersApi.updateStatus(status)

        if (this.user) {
          this.user = { ...this.user, status: response.data.status }
        }

        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to update status'
        throw error
      } finally {
        this.isLoading = false
      }
    },
  },
})
