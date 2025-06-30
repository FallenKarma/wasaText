import apiClient from '../client'

export const messagesApi = {
  send(messageData) {
    return apiClient.post('/messages', messageData)
  },

  forward(messageId, targetConversationId) {
    return apiClient.post('/messages/forward', {
      messageId,
      targetConversationId,
    })
  },

  addReaction(messageId, reactionData) {
    return apiClient.post(`/messages/${messageId}/reaction`, reactionData)
  },

  removeReaction(messageId) {
    return apiClient.delete(`/messages/${messageId}/reaction`)
  },

  delete(messageId) {
    return apiClient.delete(`/messages/${messageId}`)
  },
}
