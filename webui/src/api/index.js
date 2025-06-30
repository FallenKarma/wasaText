// src/api/index.js
import apiClient from './client'
import { authApi } from './endpoints/auth'
import { usersApi } from './endpoints/users'
import { conversationsApi } from './endpoints/conversations'
import { messagesApi } from './endpoints/messages'
import { groupsApi } from './endpoints/groups'

export { apiClient, authApi, usersApi, conversationsApi, messagesApi, groupsApi }
