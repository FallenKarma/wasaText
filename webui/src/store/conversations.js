import { defineStore } from 'pinia'
import conversationsApi from '@/api/endpoints/conversations'
import groupsApi from '@/api/endpoints/groups'

export const useConversationStore = defineStore('conversations', {
  state: () => ({
    conversations: [],
    currentConversation: null,
    isLoading: false,
    error: null,
  }),

  getters: {
    allConversations: (state) => state.conversations,
  },

  actions: {
    // Fetch all conversations for the current user
    async fetchConversations() {
      this.isLoading = true
      this.error = null

      try {
        const response = await conversationsApi.getAll()
        this.conversations = response.data
        if (this.conversations == null) {
          this.conversations = []
        }
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to fetch conversations'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Fetch a single conversation by ID
    async fetchConversation(conversationId) {
      this.isLoading = true
      this.error = null

      try {
        const response = await conversationsApi.getById(conversationId)
        this.currentConversation = response.data
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to fetch conversation'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Create a new conversation
    async createConversation(conversationData) {
      this.isLoading = true
      this.error = null

      try {
        const response = await conversationsApi.create(conversationData)
        this.conversations.unshift(response.data)
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to create conversation'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Update an existing conversation
    async updateConversation({ conversationId, data }) {
      this.isLoading = true
      this.error = null

      try {
        const response = await conversationsApi.update(conversationId, data)
        const updatedConversation = response.data

        // Update in the array
        const index = this.conversations.findIndex((c) => c.id === updatedConversation.id)
        if (index !== -1) {
          this.conversations.splice(index, 1, updatedConversation)
        }

        // Also update current conversation if it's the same one
        if (this.currentConversation && this.currentConversation.id === updatedConversation.id) {
          this.currentConversation = updatedConversation
        }

        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to update conversation'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Delete a conversation
    async deleteConversation(conversationId) {
      this.isLoading = true
      this.error = null

      try {
        await conversationsApi.delete(conversationId)
        this.conversations = this.conversations.filter((c) => c.id !== conversationId)

        // Clear current conversation if it was the one removed
        if (this.currentConversation && this.currentConversation.id === conversationId) {
          this.currentConversation = null
        }

        return true
      } catch (error) {
        this.error = error.message || 'Failed to delete conversation'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async addMessageToConversation(conversationId, message) {
      const conversation = this.conversations.find((c) => c.id === conversationId)
      if (conversation) {
        conversation.messages = conversation.messages || []
        conversation.messages.unshift(message)
      }
    },

    async leaveGroupConversation(conversationId) {
      this.isLoading = true
      this.error = null

      try {
        await groupsApi.leave(conversationId)
        this.conversations = this.conversations.filter((c) => c.id !== conversationId)

        // Clear current conversation if it was the one removed
        if (this.currentConversation && this.currentConversation.id === conversationId) {
          this.currentConversation = null
        }

        return true
      } catch (error) {
        this.error = error.message || 'Failed to leave group conversation'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async setGroupPhoto(conversationId, photo) {
      this.isLoading = true
      this.error = null

      try {
        const response = await groupsApi.setPhoto(conversationId, photo)
        const updatedConversation = response.data

        // Update in the array
        const index = this.conversations.findIndex((c) => c.id === updatedConversation.id)
        if (index !== -1) {
          this.conversations.splice(index, 1, updatedConversation)
        }

        // Also update current conversation if it's the same one
        if (this.currentConversation && this.currentConversation.id === updatedConversation.id) {
          this.currentConversation = updatedConversation
        }

        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to set group photo'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async addMembersToGroupConversation(conversationId, userIds) {
      this.isLoading = true
      this.error = null

      let finalUpdatedConversation = null // To store the latest updated conversation object

      try {
        for (const userId of userIds) {
          try {
            // Call the API for each individual user
            const response = await groupsApi.addMember(conversationId, userId) // Assuming userId is passed directly
            finalUpdatedConversation = response.data // Keep track of the most recent update

            // Update in the array for each successful addition
            const index = this.conversations.findIndex((c) => c.id === finalUpdatedConversation.id)
            if (index !== -1) {
              this.conversations.splice(index, 1, finalUpdatedConversation)
            }

            // Also update current conversation if it's the same one
            if (
              this.currentConversation &&
              this.currentConversation.id === finalUpdatedConversation.id
            ) {
              this.currentConversation = finalUpdatedConversation
            }
          } catch (individualError) {
            // Handle errors for individual user additions
            console.error(
              `Failed to add user ${userId} to group ${conversationId}:`,
              individualError,
            )
            // You might want to collect these errors or display a partial success message
            // For now, we'll continue to try adding other users.
          }
        }

        // Return the final updated conversation object after all attempts
        // This will be the state of the conversation after the last successful addition
        return finalUpdatedConversation
      } catch (overallError) {
        // This outer catch block would only be hit if there's an issue before the loop starts,
        // or if you re-throw individual errors. Given the current structure, individual errors
        // are caught inside the loop.
        this.error =
          overallError.message || 'An unexpected error occurred during member addition process'
        throw overallError // Re-throw if there's an unexpected outer error
      } finally {
        this.isLoading = false
      }
    },
  },
})
