import apiClient from '../client'

export const groupsApi = {
  addMember(groupId, userId) {
    return apiClient.post(`/groups/${groupId}/members`, { userId })
  },

  leave(groupId) {
    return apiClient.post(`/groups/${groupId}/leave`)
  },

  setName(groupId, name) {
    return apiClient.put(`/groups/${groupId}/name`, { name })
  },

  setPhoto(groupId, photoFile) {
    console.log('Setting group photo:', groupId, photoFile)
    const formData = new FormData()
    formData.append('photo', photoFile)

    return apiClient.put(`/groups/${groupId}/photo`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
  },
}

export default groupsApi
