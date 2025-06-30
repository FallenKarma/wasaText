<template>
  <div class="conversation-header">
    <div class="conversation-info">
      <button v-if="showBackButton" class="back-button" @click="goBack">
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
          <line x1="19" y1="12" x2="5" y2="12"></line>
          <polyline points="12 19 5 12 12 5"></polyline>
        </svg>
      </button>

      <div class="avatar-container">
        <div v-if="conversation.type == 'group'" class="user-avatar">
          <img
            v-if="conversation.photo"
            :src="getFullPhotoUrl(conversation.photo)"
            alt="Group photo"
          />
          <span v-else>{{ getGroupInitial() }}</span>
        </div>
        <div v-else class="user-avatar">
          <img
            v-if="otherUser && otherUser.photo"
            :src="getFullPhotoUrl(otherUser.photo)"
            alt="User Avatar"
          />
          <div v-else class="avatar-placeholder">{{ getInitials() }}</div>
        </div>
      </div>

      <div class="conversation-details">
        <h2 class="conversation-name">{{ conversationName }}</h2>
        <div class="conversation-status">
          <span v-if="conversation.type == 'group'">
            {{ conversation.participants.length }} members
          </span>
        </div>
      </div>
    </div>

    <div class="conversation-actions">
      <button
        v-if="conversation.type == 'group'"
        class="action-button"
        @click="openGroupInfo"
        title="Group info"
      >
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
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="16" x2="12" y2="12"></line>
          <line x1="12" y1="8" x2="12.01" y2="8"></line>
        </svg>
      </button>

      <button
        v-if="conversation.type === 'group'"
        class="action-button"
        @click.stop="toggleMenu"
        title="More options"
      >
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
          <circle cx="12" cy="12" r="1"></circle>
          <circle cx="12" cy="5" r="1"></circle>
          <circle cx="12" cy="19" r="1"></circle>
        </svg>
      </button>
    </div>

    <div v-if="showMenu" class="dropdown-menu" ref="menuRef">
      <div v-if="conversation.type === 'group'" class="menu-item" @click="leaveGroup">
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
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
        <span>Leave group</span>
      </div>

      <div v-if="conversation.type === 'group'" class="menu-item" @click="openAddUserDialog">
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
          <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
          <circle cx="8.5" cy="7" r="4"></circle>
          <polyline points="17 11 19 13 23 9"></polyline>
        </svg>
        <span>Add member</span>
      </div>

      <div v-if="conversation.type === 'group'" class="menu-item" @click="triggerPhotoUpload">
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
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
          <circle cx="8.5" cy="8.5" r="1.5"></circle>
          <polyline points="21 15 16 10 5 21"></polyline>
        </svg>
        <span>Set group photo</span>
      </div>
      <input
        type="file"
        ref="fileInput"
        @change="handlePhotoUpload"
        accept="image/*"
        class="hidden-file-input"
      />
    </div>

    <div v-if="showAddUserDialog" class="modal-backdrop" @click="closeAddUserDialog">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Add Members to {{ conversationName }}</h3>
          <button class="close-button" @click="closeAddUserDialog">
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
            <label for="user-search-add">Search Users</label>
            <div class="search-container">
              <input
                id="user-search-add"
                type="text"
                placeholder="Search users to add..."
                v-model="addUserSearchQuery"
                class="search-input"
              />
            </div>

            <div class="selected-members" v-if="selectedUsersToAdd.length > 0">
              <div v-for="user in selectedUsersToAdd" :key="user.id" class="selected-member">
                <span>{{ user.name }}</span>
                <button class="remove-member" @click="removeSelectedUserToAdd(user)">
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
                v-for="user in filteredUsersToAdd"
                :key="user.id"
                class="user-item"
                @click="toggleUserSelectionToAdd(user)"
                :class="{ selected: isUserSelectedToAdd(user) }"
              >
                <div class="user-avatar">
                  <img v-if="user.avatarUrl" :src="user.avatarUrl" alt="User avatar" />
                  <div v-else class="avatar-placeholder">{{ getInitialsForUser(user.name) }}</div>
                </div>
                <div class="user-info">
                  <div class="user-name">{{ user.name }}</div>
                </div>
                <div v-if="isUserSelectedToAdd(user)" class="selected-indicator">
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
          <button class="cancel-button" @click="closeAddUserDialog">Cancel</button>
          <button
            class="create-button"
            @click="addUsersToGroup"
            :disabled="selectedUsersToAdd.length === 0"
          >
            Add
          </button>
        </div>
      </div>
    </div>

    <div v-if="showGroupInfoDialog" class="modal-backdrop" @click="closeGroupInfoDialog">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Group Members ({{ conversationName }})</h3>
          <button class="close-button" @click="closeGroupInfoDialog">
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
          <div class="user-list">
            <div v-for="member in conversation.participants" :key="member.id" class="user-item">
              <div class="user-avatar">
                <img v-if="member.photo" :src="getFullPhotoUrl(member.photo)" alt="Member avatar" />
                <div v-else class="avatar-placeholder">{{ getInitialsForUser(member.name) }}</div>
              </div>
              <div class="user-info">
                <div class="user-name">{{ member.name }}</div>
                <div class="user-status" v-if="member.id === currentUserId">(You)</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { useConversationStore } from '@/store/conversations'
import { useUserStore } from '@/store/users'
import { getFullPhotoUrl } from '@/utilities/helpers'

export default {
  name: 'ConversationHeader',
  props: {
    conversation: {
      type: Object,
      required: true,
    },
    showBackButton: {
      type: Boolean,
      default: true,
    },
  },
  setup(props) {
    const authStore = useAuthStore()
    const router = useRouter()
    const conversationsStore = useConversationStore()
    const userStore = useUserStore()

    // State
    const showMenu = ref(false)
    const menuRef = ref(null)
    const showAddUserDialog = ref(false)
    const addUserSearchQuery = ref('')
    const selectedUsersToAdd = ref([])
    const fileInput = ref(null) // Reference to the hidden file input
    const showGroupInfoDialog = ref(false) // New state for group info modal

    // Computed
    const currentUserId = authStore.user?.id
    const allUsers = computed(() => userStore.allUsers)

    const conversationName = computed(() => {
      if (props.conversation.type === 'group') {
        return props.conversation.name
      } else {
        return otherUser.value?.name || 'Unknown User'
      }
    })

    const otherUser = computed(() => {
      if (props.conversation.type === 'group') {
        return null
      }
      return props.conversation.participants.find((member) => member.id !== currentUserId)
    })

    const filteredUsersToAdd = computed(() => {
      const currentParticipantIds = new Set(props.conversation.participants.map((p) => p.id))

      let usersToFilter = allUsers.value.filter(
        (user) => user.id !== currentUserId && !currentParticipantIds.has(user.id),
      )

      if (addUserSearchQuery.value) {
        usersToFilter = usersToFilter.filter((user) =>
          user.name.toLowerCase().includes(addUserSearchQuery.value.toLowerCase()),
        )
      }
      return usersToFilter
    })

    // Methods
    const goBack = () => {
      router.push('/conversations')
    }

    const getInitials = () => {
      if (!otherUser.value || !otherUser.value.name) return '?'

      return otherUser.value.name
        .split(' ')
        .map((word) => word.charAt(0).toUpperCase())
        .join('')
        .substring(0, 2)
    }

    const getGroupInitial = () => {
      if (!props.conversation.name) return '#'
      return props.conversation.name.charAt(0).toUpperCase()
    }

    const getInitialsForUser = (name) => {
      if (!name) return ''
      return name.charAt(0).toUpperCase()
    }

    const toggleMenu = () => {
      showMenu.value = !showMenu.value
    }

    // Modified openGroupInfo to show modal
    const openGroupInfo = () => {
      console.log('Opening group info for:', props.conversation)
      showGroupInfoDialog.value = true
      showMenu.value = false // Close dropdown if it was open
    }

    // New method to close the group info modal
    const closeGroupInfoDialog = () => {
      showGroupInfoDialog.value = false
    }

    const leaveGroup = async () => {
      if (confirm(`Are you sure you want to leave ${props.conversation.name}?`)) {
        try {
          await conversationsStore.leaveGroupConversation(props.conversation.id)
          showMenu.value = false
          router.push('/conversations')
        } catch (error) {
          console.error('Failed to leave group:', error)
        }
      }
    }

    const openAddUserDialog = () => {
      showAddUserDialog.value = true
      addUserSearchQuery.value = ''
      selectedUsersToAdd.value = []
      showMenu.value = false
      userStore.fetchUsers()
    }

    const closeAddUserDialog = () => {
      showAddUserDialog.value = false
    }

    const toggleUserSelectionToAdd = (user) => {
      const index = selectedUsersToAdd.value.findIndex((selected) => selected.id === user.id)
      if (index > -1) {
        selectedUsersToAdd.value.splice(index, 1)
      } else {
        selectedUsersToAdd.value.push(user)
      }
    }

    const isUserSelectedToAdd = (user) => {
      return selectedUsersToAdd.value.some((selected) => selected.id === user.id)
    }

    const removeSelectedUserToAdd = (user) => {
      selectedUsersToAdd.value = selectedUsersToAdd.value.filter(
        (selected) => selected.id !== user.id,
      )
    }

    const addUsersToGroup = async () => {
      try {
        const userIdsToAdd = selectedUsersToAdd.value.map((user) => user.id)
        await conversationsStore.addMembersToGroupConversation(props.conversation.id, userIdsToAdd)
        closeAddUserDialog()
        // No need to explicitly refetch conversation here,
        // the store's action should handle updating `currentConversation`
      } catch (error) {
        console.error('Failed to add users to group:', error)
        // Optionally show an error message to the user
      }
    }

    const triggerPhotoUpload = () => {
      fileInput.value.click()
      showMenu.value = false
    }

    const handlePhotoUpload = async (event) => {
      const file = event.target.files[0]
      if (!file) {
        return
      }

      if (!file.type.startsWith('image/')) {
        alert('Please select an image file.')
        event.target.value = ''
        return
      }

      try {
        await conversationsStore.setGroupPhoto(props.conversation.id, file)
        alert('Group photo updated successfully!')
      } catch (error) {
        console.error('Failed to set group photo:', error)
        alert('Failed to set group photo. Please try again.')
      } finally {
        event.target.value = ''
      }
    }

    const handleClickOutside = (event) => {
      if (menuRef.value && !menuRef.value.contains(event.target) && showMenu.value) {
        showMenu.value = false
      }
    }

    // Lifecycle hooks
    onMounted(() => {
      document.addEventListener('click', handleClickOutside)
    })

    onBeforeUnmount(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    watch(
      () => props.conversation,
      () => {
        closeAddUserDialog()
        // Also close the group info dialog if the conversation changes
        closeGroupInfoDialog()
      },
      { deep: true },
    )

    return {
      showMenu,
      menuRef,
      conversationName,
      otherUser,
      getInitials,
      getGroupInitial,
      goBack,
      toggleMenu,
      openGroupInfo,
      leaveGroup,
      getFullPhotoUrl,

      showAddUserDialog,
      addUserSearchQuery,
      selectedUsersToAdd,
      filteredUsersToAdd,
      openAddUserDialog,
      closeAddUserDialog,
      toggleUserSelectionToAdd,
      isUserSelectedToAdd,
      removeSelectedUserToAdd,
      addUsersToGroup,
      getInitialsForUser,

      fileInput,
      triggerPhotoUpload,
      handlePhotoUpload,

      // New state and methods for group info modal
      showGroupInfoDialog,
      closeGroupInfoDialog,
      currentUserId, // Ensure currentUserId is returned for the modal
    }
  },
}
</script>

<style scoped>
/* Existing styles for conversation header and dropdown menu remain the same */
.conversation-header {
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #e5e7eb;
  background-color: white;
  height: 64px;
}

.conversation-info {
  display: flex;
  align-items: center;
  overflow: hidden;
}

.back-button {
  background: none;
  border: none;
  padding: 0.5rem;
  margin-right: 0.5rem;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.375rem;
  transition: background-color 0.2s ease;
}

.back-button:hover {
  background-color: #f3f4f6;
}

.avatar-container {
  margin-right: 0.75rem;
}

.user-avatar,
.group-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
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

.conversation-details {
  overflow: hidden;
}

.conversation-name {
  font-size: 1rem;
  font-weight: 600;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #111827;
}

.conversation-status {
  font-size: 0.75rem;
  color: #6b7280;
}

.online-status {
  color: #10b981;
}

.offline-status {
  color: #6b7280;
}

.member-count {
  color: #6b7280;
}

.conversation-actions {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.action-button {
  background: none;
  border: none;
  padding: 0.5rem;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.375rem;
  transition: background-color 0.2s ease;
}

.action-button:hover {
  background-color: #f3f4f6;
  color: #4b5563;
}

/* Dropdown menu */
.dropdown-menu {
  position: absolute;
  top: 60px;
  right: 12px;
  width: 220px;
  background-color: white;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  z-index: 10;
  overflow: hidden;
}

.menu-item {
  color: #111827;
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.menu-item:hover {
  background-color: #f3f4f6;
}

.menu-item svg {
  margin-right: 0.75rem;
  color: #6b7280;
}

.menu-item.danger {
  color: #ef4444;
}

.menu-item.danger svg {
  color: #ef4444;
}

/* Modal styles (shared for both Add User and Group Info) */
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

.search-container {
  position: relative;
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
  left: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  color: #9ca3af;
}

.user-list {
  margin-top: 0.5rem;
  max-height: 250px; /* Limit height for scrollability */
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
  flex-shrink: 0; /* Prevent avatar from shrinking */
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

/* Hidden file input for styling */
.hidden-file-input {
  display: none;
}
</style>
