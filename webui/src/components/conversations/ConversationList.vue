<!-- src/components/conversations/ConversationList.vue -->
<template>
  <div class="conversation-list-container">
    <div class="conversation-list-header">
      <h2>Conversations</h2>
      <button class="new-chat-button" @click="openNewConversation">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        <span>New</span>
      </button>
    </div>

    <div class="search-container">
      <input
        type="text"
        placeholder="Search conversations..."
        v-model="searchQuery"
        class="search-input"
      />
      <svg
        class="search-icon"
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <circle cx="11" cy="11" r="8"></circle>
        <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
      </svg>
    </div>

    <!--     <div v-if="conversationStore.isLoading" class="loading-container">
      <div class="loading-spinner"></div>
    </div> -->

    <div v-if="filteredConversations.length === 0" class="empty-list">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="32"
        height="32"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
      </svg>
      <p>{{ searchQuery ? 'No conversations found' : 'No conversations yet' }}</p>
      <button v-if="searchQuery" class="clear-search" @click="clearSearch">Clear search</button>
      <button v-else class="start-chat" @click="openNewConversation">
        Start a new conversation
      </button>
    </div>

    <div v-else class="conversation-list">
      <ConversationItem
        v-for="conversation in filteredConversations"
        :key="conversation.id"
        :conversation="conversation"
        :is-active="activeConversationId === conversation.id"
        @click="selectConversation(conversation)"
      />
    </div>

    <!-- New Conversation Dialog -->
    <div v-if="showNewConversationDialog" class="modal-backdrop" @click="closeNewConversation">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>New Conversation</h3>
          <button class="close-button" @click="closeNewConversation">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label for="conversation-type">Conversation Type</label>
            <div class="conversation-type-options">
              <label class="type-option">
                <input type="radio" v-model="newConversationType" value="direct" />
                <span>Direct Message</span>
              </label>
              <label class="type-option">
                <input type="radio" v-model="newConversationType" value="group" />
                <span>Group Chat</span>
              </label>
            </div>
          </div>

          <div v-if="newConversationType === 'direct'" class="form-group">
            <label for="user-search">Select User</label>
            <div class="search-container">
              <input
                id="user-search"
                type="text"
                placeholder="Search users..."
                v-model="userSearchQuery"
                class="search-input"
              />
            </div>
            <div class="user-list">
              <div
                v-for="user in filteredUsers"
                :key="user.id"
                class="user-item"
                @click="selectUser(user)"
              >
                <div class="user-avatar">
                  <img v-if="user.avatarUrl" :src="user.avatarUrl" alt="User avatar" />
                  <div v-else class="avatar-placeholder">{{ getInitials(user.name) }}</div>
                </div>
                <div class="user-info">
                  <div class="user-name">{{ user.name }}</div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="form-group">
            <label for="group-name">Group Name</label>
            <input
              id="group-name"
              type="text"
              placeholder="Enter group name"
              v-model="groupName"
              class="search-input"
            />

            <label class="mt-4">Add Members</label>
            <div class="search-container">
              <input
                type="text"
                placeholder="Search users to add..."
                v-model="userSearchQuery"
                class="search-input"
              />
            </div>

            <div class="selected-members" v-if="selectedUsers.length > 0">
              <div v-for="user in selectedUsers" :key="user.name" class="selected-member">
                <span>{{ user.name }}</span>
                <button class="remove-member" @click="removeSelectedUser(user)">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              </div>
            </div>

            <div class="user-list">
              <div
                v-for="user in filteredUsers"
                :key="user.name"
                class="user-item"
                @click="toggleUserSelection(user)"
                :class="{ selected: isUserSelected(user) }"
              >
                <div class="user-avatar">
                  <img v-if="user.avatarUrl" :src="user.avatarUrl" alt="User avatar" />
                  <div v-else class="avatar-placeholder">{{ getInitials(user.name) }}</div>
                </div>
                <div class="user-info">
                  <div class="user-name">{{ user.name }}</div>
                </div>
                <div v-if="isUserSelected(user)" class="selected-indicator">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="cancel-button" @click="closeNewConversation">Cancel</button>
          <button
            class="create-button"
            @click="createConversation"
            :disabled="!canCreateConversation"
          >
            Create
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useConversationStore } from '@/store/conversations'
import { useAuthStore } from '@/store/auth'
import { useUserStore } from '@/store/users'
import ConversationItem from './ConversationItem.vue'

export default {
  name: 'ConversationList',
  components: {
    ConversationItem,
  },
  props: {
    activeConversationId: {
      type: String,
      default: null,
    },
  },
  emits: ['conversation-selected'],
  setup(props, { emit }) {
    const router = useRouter()

    const conversationStore = useConversationStore()
    const authStore = useAuthStore()
    const userStore = useUserStore()

    const isLoading = ref(false)
    const searchQuery = ref('')
    const showNewConversationDialog = ref(false)
    const newConversationType = ref('direct')
    const userSearchQuery = ref('')
    const selectedUser = ref(null)
    const selectedUsers = ref([])
    const groupName = ref('')
    const pollingInterval = ref(null)
    const isPolling = ref(false)
    const POLLING_INTERVAL_MS = 3000 // Poll every 3 seconds

    const conversations = computed(() => conversationStore.allConversations)
    const users = computed(() => userStore.allUsers)

    const filteredConversations = computed(() => {
      let conversationsToFilter = conversations.value || []

      if (searchQuery.value) {
        conversationsToFilter = conversationsToFilter.filter((conversation) => {
          return conversation.name.toLowerCase().includes(searchQuery.value.toLowerCase())
        })
      }

      // Sort the filtered conversations by lastMessage.timestamp
      return conversationsToFilter.sort((a, b) => {
        // Handle cases where a conversation might not have a lastMessage (e.g., brand new conversation)
        const timestampA = a.lastMessage ? new Date(a.lastMessage.timestamp).getTime() : 0
        const timestampB = b.lastMessage ? new Date(b.lastMessage.timestamp).getTime() : 0

        // Sort in descending order (latest on top)
        return timestampB - timestampA
      })
    })

    const filteredUsers = computed(() => {
      if (!userSearchQuery.value) return users.value.filter((user) => user.id !== authStore.user.id)

      return users.value.filter(
        (user) =>
          user.name.toLowerCase().includes(userSearchQuery.value.toLowerCase()) &&
          user.id !== authStore.user.id,
      )
    })

    const canCreateConversation = computed(() => {
      if (newConversationType.value === 'direct') {
        return !!selectedUser.value
      } else {
        return groupName.value.trim() !== '' && selectedUsers.value.length > 0
      }
    })

    const fetchConversations = async () => {
      try {
        await conversationStore.fetchConversations()
      } catch (error) {
        console.error('Failed to fetch conversations:', error)
      }
    }

    const fetchUsers = async () => {
      try {
        await userStore.fetchUsers()
      } catch (error) {
        console.error('Failed to fetch users:', error)
      }
    }

    const selectConversation = (conversation) => {
      emit('conversation-selected', conversation)
      router.push(`/conversations/${conversation.id}`)
    }

    const startPolling = () => {
      if (pollingInterval.value) {
        clearInterval(pollingInterval.value)
      }

      isPolling.value = true
      pollingInterval.value = setInterval(async () => {
        try {
          await fetchConversations()
        } catch (error) {
          console.error('Error during polling:', error)
        }
      }, POLLING_INTERVAL_MS)
    }

    const stopPolling = () => {
      if (pollingInterval.value) {
        clearInterval(pollingInterval.value)
        pollingInterval.value = null
      }
      isPolling.value = false
    }

    const handleVisibilityChange = () => {
      if (document.hidden) {
        stopPolling()
      } else {
        startPolling()
      }
    }

    const clearSearch = () => {
      searchQuery.value = ''
    }

    const openNewConversation = () => {
      showNewConversationDialog.value = true
      newConversationType.value = 'direct'
      userSearchQuery.value = ''
      selectedUser.value = null
      selectedUsers.value = []
      groupName.value = ''
    }

    const closeNewConversation = () => {
      showNewConversationDialog.value = false
    }

    const selectUser = (user) => {
      selectedUser.value = user.id
      console.log('Selected user:', selectedUser.value)
      createConversation()
    }

    const toggleUserSelection = (user) => {
      if (isUserSelected(user)) {
        removeSelectedUser(user)
      } else {
        selectedUsers.value.push(user)
      }
    }

    const isUserSelected = (user) => {
      return selectedUsers.value.some((selectedUser) => selectedUser.id === user.id)
    }

    const removeSelectedUser = (user) => {
      selectedUsers.value = selectedUsers.value.filter(
        (selectedUser) => selectedUser.id !== user.id,
      )
    }

    const getInitials = (name) => {
      if (!name) return ''
      return name.charAt(0).toUpperCase()
    }

    const createConversation = async () => {
      try {
        let conversationData = {}

        if (newConversationType.value === 'direct') {
          conversationData = {
            type: 'direct',
            participants: [selectedUser.value],
          }
        } else {
          conversationData = {
            type: 'group',
            name: groupName.value,
            participants: selectedUsers.value.map((user) => user.id),
          }
        }

        const newConversation = await conversationStore.createConversation(conversationData)
        closeNewConversation()
        selectConversation(newConversation)
      } catch (error) {
        console.error('Failed to create conversation:', error)
      }
    }

    // Lifecycle
    onMounted(() => {
      document.addEventListener('visibilitychange', handleVisibilityChange)
      fetchConversations()
      startPolling()
      fetchUsers()
    })

    onBeforeUnmount(() => {
      stopPolling() // Stop polling when component unmounts
      document.removeEventListener('visibilitychange', handleVisibilityChange)
    })

    return {
      conversationStore, // Expose the store to the template
      isLoading,
      searchQuery,
      conversations,
      filteredConversations,
      selectConversation,
      clearSearch,
      openNewConversation,
      closeNewConversation,
      showNewConversationDialog,
      newConversationType,
      userSearchQuery,
      selectedUser,
      selectedUsers,
      groupName,
      users,
      filteredUsers,
      selectUser,
      toggleUserSelection,
      isUserSelected,
      removeSelectedUser,
      canCreateConversation,
      createConversation,
      getInitials,
    }
  },
}
</script>

<style scoped>
.conversation-list-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  border-right: 1px solid #e5e7eb;
  background-color: #f9fafb;
}

.conversation-list-header {
  padding: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e5e7eb;
}

.conversation-list-header h2 {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
  color: #111827;
}

.new-chat-button {
  display: flex;
  align-items: center;
  padding: 0.5rem 0.75rem;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.new-chat-button:hover {
  background-color: #2563eb;
}

.new-chat-button svg {
  margin-right: 0.25rem;
}

.search-container {
  position: relative;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #e5e7eb;
}

.search-input {
  width: 100%;
  padding: 0.5rem 1rem 0.5rem 2rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  background-color: white;
  transition: border-color 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.search-icon {
  position: absolute;
  left: 1.5rem;
  top: 50%;
  transform: translateY(-50%);
  color: #9ca3af;
}

.loading-container {
  display: flex;
  justify-content: center;
  padding: 2rem 0;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid #e5e7eb;
  border-top: 2px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.empty-list {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 1.5rem;
  color: #6b7280;
  text-align: center;
}

.empty-list svg {
  color: #9ca3af;
  margin-bottom: 1rem;
}

.empty-list p {
  margin: 0 0 1rem;
  font-size: 0.875rem;
}

.clear-search,
.start-chat {
  padding: 0.5rem 0.75rem;
  background-color: #f3f4f6;
  color: #4b5563;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.clear-search:hover,
.start-chat:hover {
  background-color: #e5e7eb;
}

.conversation-list {
  flex: 1;
  overflow-y: auto;
}

/* Modal styles */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 50;
}

.modal-content {
  color: #111827;
  width: 90%;
  max-width: 480px;
  max-height: 90vh;
  background-color: white;
  border-radius: 0.5rem;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 1rem;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
}

.close-button {
  background: none;
  border: none;
  padding: 0.25rem;
  cursor: pointer;
  color: #6b7280;
  transition: color 0.2s ease;
}

.close-button:hover {
  color: #4b5563;
}

.modal-body {
  padding: 1rem;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 0.5rem;
  color: #4b5563;
}

.conversation-type-options {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.type-option {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.type-option input {
  margin-right: 0.5rem;
}

.user-list {
  margin-top: 0.5rem;
  max-height: 250px;
  overflow-y: auto;
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
}

.user-item {
  color: #111827;
  display: flex;
  align-items: center;
  padding: 0.75rem;
  border-bottom: 1px solid #f3f4f6;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.user-item:last-child {
  border-bottom: none;
}

.user-item:hover {
  background-color: #f9fafb;
}

.user-item.selected {
  background-color: #eff6ff;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 0.75rem;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background-color: #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1rem;
}

.user-info {
  flex: 1;
}

.user-name {
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 0.125rem;
}

.user-status {
  font-size: 0.75rem;
  color: #6b7280;
}

.user-status.online {
  color: #10b981;
}

.selected-indicator {
  color: #3b82f6;
}

.selected-members {
  color: #111827;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin: 0.5rem 0;
}

.selected-member {
  display: flex;
  align-items: center;
  background-color: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 1rem;
  padding: 0.25rem 0.5rem 0.25rem 0.75rem;
  font-size: 0.75rem;
}

.remove-member {
  background: none;
  border: none;
  padding: 0;
  margin-left: 0.25rem;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
}

.modal-footer {
  padding: 1rem;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.cancel-button {
  padding: 0.5rem 1rem;
  background-color: white;
  color: #4b5563;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.cancel-button:hover {
  background-color: #f9fafb;
}

.create-button {
  padding: 0.5rem 1rem;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.create-button:hover:not(:disabled) {
  background-color: #2563eb;
}

.create-button:disabled {
  background-color: #93c5fd;
  cursor: not-allowed;
}

.mt-4 {
  margin-top: 1rem;
}
</style>
