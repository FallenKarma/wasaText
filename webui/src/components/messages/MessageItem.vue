<template>
  <div class="message-wrapper" :class="{ 'own-message': isOwn }">
    <div class="message-item" :class="{ 'with-avatar': showAvatar }">
      <div v-if="showAvatar && !isOwn" class="avatar-container">
        <img
          v-if="message.sender?.photo_url"
          :src="getFullPhotoUrl(message.sender.photo_url)"
          :alt="message.sender.name"
          class="avatar"
        />
        <div v-else class="default-avatar">
          {{ getInitials(message.sender?.name) }}
        </div>
      </div>
      <div v-else-if="!isOwn" class="avatar-spacer"></div>

      <div class="message-content-wrapper">
        <div v-if="showAvatar && !isOwn" class="sender-name">
          {{ message.sender?.name }}
        </div>

        <div v-if="message.repliedToMessageData && !message.deletedAt" class="reply-indicator">
          <div class="reply-line"></div>
          <div class="replied-content">
            <span class="replied-user">{{ message.repliedToMessageData.sender.name }}</span>
            <span
              v-if="!message.repliedToMessageData.content.startsWith('/')"
              class="replied-text"
              >{{ truncateText(message.repliedToMessageData.content, 40) }}</span
            >
            <span v-else class="replied-text">🖼️ Photo</span>
          </div>
        </div>

        <div
          class="message-content"
          :class="{ 'is-deleted': message.deletedAt, 'is-editing': isEditing }"
        >
          <div v-if="isEditing" class="edit-mode">
            <textarea
              ref="editTextarea"
              v-model="editContent"
              class="edit-textarea"
              @keydown="handleEditKeydown"
              @blur="handleEditBlur"
              rows="1"
            ></textarea>
            <div class="edit-actions">
              <button class="edit-action-button save" @click="saveEdit">
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
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>
              </button>
              <button class="edit-action-button cancel" @click="cancelEdit">
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

          <div v-else>
            <div v-if="message.type === 'photo' && !message.deletedAt" class="message-photo">
              <img :src="getFullPhotoUrl(message.content)" :alt="'Photo message'" />
            </div>

            <div v-else-if="message.type === 'text' && !message.deletedAt" class="message-text">
              {{ message.content }}
            </div>

            <div v-else-if="message.deletedAt" class="deleted-message">
              <span class="deleted-icon">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="12"
                  height="12"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <polyline points="3 6 5 6 21 6"></polyline>
                  <path
                    d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                  ></path>
                </svg>
              </span>
              Message deleted
            </div>

            <MessageReactions
              v-if="message.reactions && message.reactions.length > 0"
              :reactions="message.reactions"
              :messageId="message.id"
              @add-reaction="$emit('reaction', { messageId: message.id, reaction: $event })"
              @remove-reaction="$emit('reaction', { messageId: message.id, reaction: null })"
            />

            <div class="message-time">
              {{ formatMessageTime(message.timestamp) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="message-actions" v-if="!message.deletedAt && !isEditing">
      <button class="action-button emoji-button" @click="toggleEmojiPicker">
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
          <circle cx="12" cy="12" r="10"></circle>
          <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
          <line x1="9" y1="9" x2="9.01" y2="9"></line>
          <line x1="15" y1="9" x2="15.01" y2="9"></line>
        </svg>
      </button>
      <button class="action-button reply-button" @click="$emit('reply', message)">
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
          <polyline points="9 17 4 12 9 7"></polyline>
          <path d="M20 18v-2a4 4 0 0 0-4-4H4"></path>
        </svg>
      </button>
      <button v-if="isOwn" class="action-button more-button" @click.stop="toggleMoreActions">
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
          <circle cx="12" cy="12" r="1"></circle>
          <circle cx="19" cy="12" r="1"></circle>
          <circle cx="5" cy="12" r="1"></circle>
        </svg>
      </button>

      <div v-if="showMoreActions" class="more-actions-dropdown">
        <button class="dropdown-item" @click="startEdit">
          <span class="dropdown-icon">
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
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
            </svg>
          </span>
          Edit
        </button>
        <button class="dropdown-item delete" @click="confirmDelete">
          <span class="dropdown-icon">
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
              <polyline points="3 6 5 6 21 6"></polyline>
              <path
                d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
              ></path>
            </svg>
          </span>
          Delete
        </button>
      </div>
    </div>

    <div v-if="showEmojiPicker" class="emoji-picker">
      <div class="emoji-list">
        <button
          v-for="emoji in commonEmojis"
          :key="emoji"
          class="emoji-item"
          @click="addReaction(emoji)"
        >
          {{ emoji }}
        </button>
      </div>
    </div>

    <div v-if="showDeleteConfirmation" class="delete-confirmation">
      <div class="confirmation-dialog">
        <div class="confirmation-title">Delete message?</div>
        <div class="confirmation-text">This cannot be undone.</div>
        <div class="confirmation-actions">
          <button class="cancel-button" @click="showDeleteConfirmation = false">Cancel</button>
          <button class="delete-button" @click="deleteMessage">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import MessageReactions from './MessageReactions.vue'
import { useMessageStore } from '@/store/messages'
import { useAuthStore } from '@/store/auth'
import { getFullPhotoUrl } from '@/utilities/helpers'

export default {
  name: 'MessageItem',
  components: {
    MessageReactions,
  },
  props: {
    message: {
      type: Object,
      required: true,
    },
    isOwn: {
      type: Boolean,
      default: false,
    },
    showAvatar: {
      type: Boolean,
      default: true,
    },
  },
  emits: ['reaction', 'reply'],
  setup(props) {
    const showActions = ref(false)
    const showMoreActions = ref(false)
    const showEmojiPicker = ref(false)
    const showDeleteConfirmation = ref(false)
    const isEditing = ref(false)
    const editContent = ref('')
    const originalContent = ref('')
    const editTextarea = ref(null)
    const messageStore = useMessageStore()
    const authStore = useAuthStore()

    // Common emojis
    const commonEmojis = ['👍', '❤️', '😂', '😮', '😢', '👏', '🎉', '🤔']

    // Methods
    const getInitials = (username) => {
      if (!username) return '?'
      return username
        .split(' ')
        .map((word) => word.charAt(0).toUpperCase())
        .slice(0, 2)
        .join('')
    }

    const formatMessageTime = (timestamp) => {
      if (!timestamp) return ''

      const date = new Date(timestamp)
      const now = new Date()

      // Calculate difference in milliseconds
      const diffMs = now.getTime() - date.getTime()
      const diffMinutes = Math.floor(diffMs / (1000 * 60))
      const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24)) // Difference in full days

      // Get current date, yesterday, and 7 days ago (for comparison)
      const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
      const messageDateDay = new Date(date.getFullYear(), date.getMonth(), date.getDate())

      const yesterday = new Date(today)
      yesterday.setDate(today.getDate() - 1)

      const sevenDaysAgo = new Date(today)
      sevenDaysAgo.setDate(today.getDate() - 7)

      // Option for time formatting (e.g., "09:30 AM" or "09:30")
      // You can adjust 'en-US' and 'hour12' based on your locale preference
      const timeFormatOptions = {
        hour: '2-digit',
        minute: '2-digit',
        hour12: true, // Set to false for 24-hour format
      }
      const formattedTime = date.toLocaleTimeString('it-IT', timeFormatOptions) // Or your preferred locale

      // Option for day name formatting (e.g., "Monday")
      const dayNameFormatOptions = { weekday: 'long' }
      const formattedDayName = date.toLocaleDateString('en-US', dayNameFormatOptions) // Or your preferred locale

      if (diffMinutes < 1) {
        return 'Just now'
      } else if (diffMinutes < 60) {
        return `${diffMinutes}m ago`
      } else if (diffHours < 24) {
        return `Today ${formattedTime}` // e.g., "Today 09:30 AM"
      } else if (messageDateDay.getTime() === yesterday.getTime()) {
        return 'Yesterday ' + formattedTime
      } else if (messageDateDay > sevenDaysAgo && messageDateDay < today) {
        // If it's within the last 7 days (not including today and yesterday)
        return formattedDayName + ' ' + formattedTime // e.g., "Monday 09:30 AM"
      } else {
        return date.toLocaleDateString(undefined, {
          month: 'short',
          day: 'numeric',
          hour: '2-digit',
          minute: '2-digit',
        })
      }
    }

    const truncateText = (text, maxLength) => {
      if (!text) return ''
      return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
    }

    const isImage = (attachment) => {
      if (!attachment || !attachment.mimeType) return false
      return attachment.mimeType.startsWith('image/')
    }

    const formatFileSize = (bytes) => {
      if (!bytes) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
    }

    const toggleMoreActions = () => {
      showMoreActions.value = !showMoreActions.value
      if (showMoreActions.value) {
        showEmojiPicker.value = false
      }
    }

    const toggleEmojiPicker = () => {
      showEmojiPicker.value = !showEmojiPicker.value
      if (showEmojiPicker.value) {
        showMoreActions.value = false
      }
    }

    const addReaction = (emoji) => {
      showEmojiPicker.value = false
      const reaction = {
        emoji,
        userId: authStore.user.id,
        messageId: props.message.id,
      }
      messageStore.addReaction(props.message.id, reaction)
    }

    // Edit functionality
    const startEdit = () => {
      isEditing.value = true
      editContent.value = props.message.content
      originalContent.value = props.message.content
      showMoreActions.value = false

      // Focus textarea and adjust height after DOM update
      nextTick(() => {
        if (editTextarea.value) {
          editTextarea.value.focus()
          adjustTextareaHeight()
          // Select all text for easy editing
          editTextarea.value.select()
        }
      })
    }

    const cancelEdit = () => {
      isEditing.value = false
      editContent.value = originalContent.value
    }

    const saveEdit = async () => {
      if (!editContent.value.trim()) {
        // Don't save empty messages
        return
      }

      if (editContent.value.trim() === originalContent.value.trim()) {
        // No changes made
        cancelEdit()
        return
      }

      try {
        await messageStore.updateMessage({
          messageId: props.message.id,
          content: editContent.value.trim(),
        })
        isEditing.value = false
        // Update the message content in the prop directly for immediate UI reflection
        // This assumes message is reactive and this mutation is acceptable.
        // For stricter immutability, you might emit an event and let the parent handle the update.
        props.message.content = editContent.value.trim()
      } catch (error) {
        console.error('Failed to update message:', error)
        // Optionally show an error message to the user
      }
    }

    const handleEditKeydown = (event) => {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault()
        saveEdit()
      } else if (event.key === 'Escape') {
        event.preventDefault()
        cancelEdit()
      }

      // Adjust textarea height as user types
      adjustTextareaHeight()
    }

    const handleEditBlur = () => {
      // Small delay to allow save button click to register
      setTimeout(() => {
        if (isEditing.value) {
          saveEdit()
        }
      }, 150)
    }

    const adjustTextareaHeight = () => {
      nextTick(() => {
        if (editTextarea.value) {
          editTextarea.value.style.height = 'auto'
          editTextarea.value.style.height = editTextarea.value.scrollHeight + 'px'
        }
      })
    }

    const confirmDelete = () => {
      showMoreActions.value = false
      showDeleteConfirmation.value = true
    }

    const deleteMessage = () => {
      messageStore.deleteMessage(props.message.id)
      showDeleteConfirmation.value = false
    }

    // Handle click outside to close dropdowns
    const handleClickOutside = (event) => {
      const isClickInsideMoreActions = event.target.closest('.more-button')
      const isClickInsideEmojiButton = event.target.closest('.emoji-button')

      if (!isClickInsideMoreActions && showMoreActions.value) {
        showMoreActions.value = false
      }

      if (!isClickInsideEmojiButton && showEmojiPicker.value) {
        showEmojiPicker.value = false
      }
    }

    // Mouse events
    const handleMouseEnter = () => {
      showActions.value = true
    }

    const handleMouseLeave = () => {
      showActions.value = false

      // Don't hide dropdowns when mouse leaves if they're active
      if (!showMoreActions.value && !showEmojiPicker.value) {
        showActions.value = false
      }
    }

    // Lifecycle hooks
    onMounted(() => {
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      showActions,
      showMoreActions,
      showEmojiPicker,
      showDeleteConfirmation,
      isEditing,
      editContent,
      editTextarea,
      commonEmojis,
      getInitials,
      formatMessageTime,
      truncateText,
      isImage,
      getFullPhotoUrl,
      formatFileSize,
      toggleMoreActions,
      toggleEmojiPicker,
      addReaction,
      startEdit,
      cancelEdit,
      saveEdit,
      handleEditKeydown,
      handleEditBlur,
      confirmDelete,
      deleteMessage,
      handleMouseEnter,
      handleMouseLeave,
    }
  },
}
</script>

<style scoped>
.message-photo img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 0.375rem;
  object-fit: contain;
  display: block;
  margin-bottom: 0.25rem;
}

.message-content {
  padding: 0.5rem 0.75rem;
}

.message-content .message-photo {
  padding: 0;
  margin: 0;
  border-radius: 0.375rem;
  overflow: hidden;
}

.message-content .message-photo + .message-time,
.message-content .message-photo + .message-reactions {
  margin-top: 0.5rem;
}

.message-wrapper {
  position: relative;
  margin-bottom: 0.5rem;
  padding: 0.25rem 0.5rem;
  border-radius: 0.5rem;
  transition: background-color 0.2s;
}

.message-wrapper:hover {
  background-color: #f9fafb;
}

.message-item {
  display: flex;
  position: relative;
}

.avatar-container {
  margin-right: 0.75rem;
  flex-shrink: 0;
}

.avatar-spacer {
  width: 2.25rem;
  margin-right: 0.75rem;
  flex-shrink: 0;
}

.avatar {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  object-fit: cover;
}

.default-avatar {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  background-color: #e5e7eb;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  font-weight: 500;
}

.message-content-wrapper {
  flex: 1;
  min-width: 0;
}

.sender-name {
  font-weight: 500;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
  color: #374151;
}

.message-content {
  position: relative;
  background-color: #f3f4f6;
  padding: 0.5rem 0.75rem;
  border-radius: 0.375rem;
  max-width: 85%;
  word-break: break-word;
}

.message-content.is-editing {
  background-color: #ffffff;
  border: 2px solid #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.own-message .message-content {
  background-color: #dbeafe;
  margin-left: auto;
}

.own-message .message-content.is-editing {
  background-color: #ffffff;
}

.message-text {
  color: #111827;
  white-space: pre-wrap;
  font-size: 0.9375rem;
  line-height: 1.5;
}

.message-time {
  font-size: 0.75rem;
  color: #6b7280;
  margin-top: 0.25rem;
  text-align: right;
}

.edited-indicator {
  font-size: 0.75rem;
  color: #9ca3af;
  margin-left: 0.25rem;
}

.deleted-message {
  color: #9ca3af;
  font-style: italic;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
}

.deleted-icon {
  margin-right: 0.375rem;
  display: flex;
  align-items: center;
}

/* Edit mode styles */
.edit-mode {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.edit-textarea {
  width: 100%;
  border: none;
  outline: none;
  background: none;
  resize: none;
  font-family: inherit;
  font-size: 0.9375rem;
  line-height: 1.5;
  color: #111827;
  padding: 0;
  min-height: 1.5rem;
  overflow: hidden;
}

.edit-actions {
  display: flex;
  gap: 0.25rem;
  justify-content: flex-end;
}

.edit-action-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.75rem;
  height: 1.75rem;
  border: 1px solid #d1d5db;
  background-color: #ffffff;
  border-radius: 0.25rem;
  cursor: pointer;
  transition: all 0.2s;
}

.edit-action-button.save {
  color: #059669;
  border-color: #059669;
}

.edit-action-button.save:hover {
  background-color: #ecfdf5;
}

.edit-action-button.cancel {
  color: #dc2626;
  border-color: #dc2626;
}

.edit-action-button.cancel:hover {
  background-color: #fef2f2;
}

/* Reply styles */
.reply-indicator {
  display: flex;
  margin-bottom: 0.25rem;
  padding-left: 0.5rem;
  border-left: 2px solid #d1d5db;
}

.replied-content {
  font-size: 0.75rem;
  color: #6b7280;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.replied-user {
  font-weight: 500;
  color: #4b5563;
  margin-right: 0.25rem;
}

/* Attachments */
.attachments {
  margin-top: 0.5rem;
}

.attachment {
  margin-top: 0.5rem;
  cursor: pointer;
}

.image-attachment img {
  max-width: 100%;
  max-height: 200px;
  border-radius: 0.25rem;
  object-fit: contain;
}

.file-attachment {
  display: flex;
  align-items: center;
  padding: 0.5rem;
  background-color: rgba(255, 255, 255, 0.5);
  border-radius: 0.25rem;
  border: 1px solid #e5e7eb;
}

.file-icon {
  margin-right: 0.5rem;
  color: #6b7280;
}

.file-details {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 0.875rem;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 0.75rem;
  color: #6b7280;
}

/* Message actions */
.message-actions {
  position: absolute;
  top: -0.75rem;
  right: 0.5rem;
  display: flex;
  background-color: white;
  border-radius: 0.375rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
  opacity: 0;
  transition: opacity 0.2s;
}

.message-wrapper:hover .message-actions {
  opacity: 1;
}

.action-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: none;
  background: none;
  color: #6b7280;
  cursor: pointer;
  border-radius: 0.25rem;
  transition: all 0.2s;
}

.action-button:hover {
  background-color: #f3f4f6;
  color: #374151;
}

/* Emoji picker */
.emoji-picker {
  position: absolute;
  top: -2.5rem;
  right: 0.5rem;
  background-color: white;
  border-radius: 0.375rem;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border: 1px solid #e5e7eb;
  z-index: 10;
  padding: 0.5rem;
}

.emoji-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.emoji-item {
  font-size: 1.125rem;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 0.25rem;
  transition: all 0.2s;
}

.emoji-item:hover {
  background-color: #f3f4f6;
}

/* More actions dropdown */
.more-actions-dropdown {
  position: absolute;
  top: 2rem;
  right: 0;
  background-color: white;
  border-radius: 0.375rem;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border: 1px solid #e5e7eb;
  z-index: 10;
  width: 8rem;
}

.dropdown-item {
  display: flex;
  align-items: center;
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  color: #374151;
  cursor: pointer;
  border: none;
  background: none;
  width: 100%;
  text-align: left;
  transition: all 0.2s;
}

.dropdown-item:hover {
  background-color: #f3f4f6;
}

.dropdown-item.delete {
  color: #ef4444;
}

.dropdown-item.delete:hover {
  background-color: #fef2f2;
}

.dropdown-icon {
  margin-right: 0.5rem;
  display: flex;
  align-items: center;
}

/* Delete confirmation */
.delete-confirmation {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 50;
}

.confirmation-dialog {
  background-color: white;
  border-radius: 0.5rem;
  box-shadow:
    0 10px 15px -3px rgba(0, 0, 0, 0.1),
    0 4px 6px -2px rgba(0, 0, 0, 0.05);
  padding: 1.25rem;
  width: 20rem;
  max-width: 90%;
}

.confirmation-title {
  font-size: 1rem;
  font-weight: 600;
  color: #111827;
  margin-bottom: 0.5rem;
}

.confirmation-text {
  font-size: 0.875rem;
  color: #6b7280;
  margin-bottom: 1rem;
}

.confirmation-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.cancel-button {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  border-radius: 0.375rem;
  border: 1px solid #d1d5db;
  background-color: white;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.cancel-button:hover {
  background-color: #f9fafb;
}

.delete-button {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  border-radius: 0.375rem;
  border: none;
  background-color: #ef4444;
  color: white;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.delete-button:hover {
  background-color: #dc2626;
}
</style>
