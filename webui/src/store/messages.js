import { defineStore } from 'pinia'
import apiClient from '@/api/client'
import { messagesApi } from '@/api/endpoints/messages'
import { useAuthStore } from '@/store/auth'

export const useMessageStore = defineStore('messages', {
  state: () => ({
    messages: [],
    isLoading: false,
    error: null,
  }),

  getters: {
    allMessages: (state) => state.messages,
    isLoadingMessages: (state) => state.isLoading,
    hasError: (state) => !!state.error,
  },

  actions: {
    // Fetch messages for a conversation
    async fetchMessages(conversationId) {
      this.isLoading = true
      this.error = null

      try {
        const response = await apiClient.get(`/conversations/${conversationId}`)
        this.messages = response.data.messages || []
        return this.messages
      } catch (error) {
        this.error = error.message || 'Failed to fetch messages'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Send a new message
    async sendTextMessage(messageData) {
      this.isLoading = true
      this.error = null
      try {
        const response = await messagesApi.send(messageData)

        this.messages.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to send message'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Action to send a photo message
    async sendPhotoMessage(conversationId, photoFile, replyToId = '') {
      this.isLoading = true
      try {
        const formData = new FormData()
        formData.append('conversationId', conversationId)
        formData.append('photo', photoFile) // 'photo' matches the backend's expected field name
        if (replyToId) {
          formData.append('replyTo', replyToId)
        }

        const response = await apiClient.post('/messages', formData, {
          headers: {
            'Content-Type': 'multipart/form-data', // Important for file uploads
          },
        })
        const newMessage = response.data
        // Optionally add the new message to your state
        this.messages.push(newMessage)
        return newMessage
      } catch (error) {
        console.error('Error sending photo message:', error)
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Add a message locally (for optimistic updates)
    addMessage(message) {
      this.messages = [message, ...this.messages]
    },

    // Update a message (for editing)
    async updateMessage({ messageId, content }) {
      this.isLoading = true
      this.error = null

      try {
        const response = await apiClient.put(`/messages/${messageId}`, { content })

        const messageToUpdate = this.messages.find((m) => m.id === messageId)
        if (messageToUpdate) {
          messageToUpdate.content = content
        }

        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to update message'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Delete a message
    async deleteMessage(messageId) {
      this.isLoading = true
      this.error = null

      try {
        await messagesApi.delete(messageId)
        this.messages = this.messages.filter((m) => m.id !== messageId)
        return true
      } catch (error) {
        this.error = error.message || 'Failed to delete message'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Add a reaction to a message
    async addReaction(messageId, reaction) {
      try {
        await messagesApi.addReaction(messageId, reaction)

        const messageToUpdate = this.messages.find((m) => m.id === messageId)
        if (messageToUpdate) {
          if (!messageToUpdate.reactions) {
            messageToUpdate.reactions = []
          }
          messageToUpdate.reactions.push(reaction)
        }
      } catch (error) {
        console.error('Failed to add reaction:', error)
        throw error
      }
    },

    // Remove a reaction from a message
    async removeReaction(messageId, emoji) {
      const authStore = useAuthStore()
      const currentUserId = authStore.user?.id
      try {
        const response = await messagesApi.removeReaction(messageId)

        const messageToUpdate = this.messages.find((m) => m.id === messageId)
        const reactionIndex = messageToUpdate.reactions.findIndex(
          (r) => r.emoji === emoji && r.userId === currentUserId,
        )

        if (reactionIndex !== -1) {
          messageToUpdate.reactions.splice(reactionIndex, 1)
        }

        return response.data
      } catch (error) {
        console.error('Failed to remove reaction:', error)
        throw error
      }
    },

    // Reset pagination
    resetPagination() {
      this.currentPage = 1
      this.hasMoreMessages = true
      this.lastMessageTimestamp = null
    },

    // Clear messages when leaving a conversation
    clearMessages() {
      this.messages = []
      this.resetPagination()
    },
  },
})
