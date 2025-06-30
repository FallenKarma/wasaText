// src/api/endpoints/conversations.js
import apiClient from '../client'

const conversationsApi = {
  getAll() {
    return apiClient.get('/conversations')
  },

  getById(id) {
    return apiClient.get(`/conversations/${id}`)
  },

  create(conversationData) {
    return apiClient.post('/conversations', conversationData)
  },

  update(id, conversationData) {
    return apiClient.put(`/conversations/${id}`, conversationData)
  },
}

export default conversationsApi
