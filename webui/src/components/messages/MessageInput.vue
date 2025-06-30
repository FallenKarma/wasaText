<!-- src/components/messages/MessageInput.vue -->
<template>
  <div class="message-input-container">
    <!-- Attachment preview area -->
    <div v-if="attachments.length > 0" class="attachments-preview">
      <div v-for="(file, index) in attachments" :key="index" class="attachment-preview">
        <div v-if="isImage(file)" class="image-preview">
          <img :src="getPreviewUrl(file)" :alt="file.name" />
        </div>
        <div v-else class="file-preview">
          <div class="file-icon">
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
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <polyline points="10 9 9 9 8 9"></polyline>
            </svg>
          </div>
          <div class="file-name">{{ file.name }}</div>
        </div>
        <button class="remove-attachment" @click="removeAttachment(index)">
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
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>

    <!-- Message input area -->
    <div class="input-wrapper">
      <!-- Attachment button -->
      <label class="attachment-button" role="button">
        <input
          type="file"
          ref="fileInput"
          @change="handleFileSelect"
          multiple
          class="hidden-file-input"
          :disabled="isAttachmentButtonDisabled"
        />
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
          <path d="M21 10l-8-4-8 4m16 0l-8 4m8-4v10m0-10l-8 4m0-8v10m0 0l-8-4V6l8 4"></path>
        </svg>
      </label>

      <!-- Emoji button -->
      <button class="emoji-button" @click="toggleEmojiPicker">
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
          <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
          <line x1="9" y1="9" x2="9.01" y2="9"></line>
          <line x1="15" y1="9" x2="15.01" y2="9"></line>
        </svg>
      </button>

      <!-- Text input -->
      <textarea
        ref="inputRef"
        v-model="messageText"
        placeholder="Type a message..."
        class="message-textarea"
        :disabled="isTextareaDisabled"
        @keydown.enter.prevent="handleEnterPress"
        @input="autoGrow"
      ></textarea>

      <!-- Send button -->
      <button class="send-button" @click="sendMessage" :disabled="!canSend">
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
          class="send-icon"
        >
          <line x1="22" y1="2" x2="11" y2="13"></line>
          <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
        </svg>
      </button>
    </div>

    <!-- Emoji picker -->
    <div v-if="showEmojiPicker" class="emoji-picker">
      <div class="emoji-categories">
        <button
          v-for="(category, index) in emojiCategories"
          :key="index"
          class="category-button"
          :class="{ active: currentCategory === index }"
          @click="currentCategory = index"
        >
          {{ category.icon }}
        </button>
      </div>
      <div class="emoji-list">
        <button
          v-for="emoji in currentCategoryEmojis"
          :key="emoji"
          class="emoji-item"
          @click="insertEmoji(emoji)"
        >
          {{ emoji }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useMessageStore } from '@/store/messages'
import { useAuthStore } from '@/store/auth'

export default {
  name: 'MessageInput',
  props: {
    conversationId: {
      type: [String, Number],
      required: true,
    },
    replyingTo: {
      type: Object,
      default: null,
    },
    isLoading: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['message-sent', 'cancel-reply'],
  setup(props, { emit }) {
    const messageStore = useMessageStore()
    const authStore = useAuthStore()
    const messageText = ref('')
    const attachments = ref([]) // This will hold File objects for photos
    const fileInput = ref(null)
    const inputRef = ref(null)
    const showEmojiPicker = ref(false)
    const currentCategory = ref(0)

    // Emoji categories (rest of your emojiCategories array)
    const emojiCategories = [
      {
        name: 'Frequently Used',
        icon: '🕒',
        emojis: ['😊', '👍', '❤️', '😂', '🙏', '😍', '🥰', '😎', '🤔', '👏'],
      },
      {
        name: 'Smileys',
        icon: '😀',
        emojis: [
          '😀',
          '😃',
          '😄',
          '😁',
          '😆',
          '😅',
          '😂',
          '🤣',
          '😊',
          '😇',
          '🙂',
          '🙃',
          '😉',
          '😌',
          '😍',
          '🥰',
          '😘',
          '😗',
          '😙',
          '😚',
          '😋',
          '😛',
          '😝',
          '😜',
        ],
      },
      {
        name: 'Gestures',
        icon: '👋',
        emojis: [
          '👋',
          '🤚',
          '🖐️',
          '✋',
          '🖖',
          '👌',
          '🤌',
          '🤏',
          '✌️',
          '🤞',
          '🤟',
          '🤘',
          '🤙',
          '👈',
          '👉',
          '👆',
          '🖕',
          '👇',
          '👍',
          '👎',
          '✊',
          '👊',
          '🤛',
          '🤜',
        ],
      },
      {
        name: 'Animals',
        icon: '🐶',
        emojis: [
          '🐶',
          '🐱',
          '🐭',
          '🐹',
          '🐰',
          '🦊',
          '🐻',
          '🐼',
          '🐨',
          '🐯',
          '🦁',
          '🐮',
          '🐷',
          '🐸',
          '🐵',
          '🐔',
          '🐧',
          '🐦',
          '🐤',
          '🦆',
          '🦉',
          '🦇',
          '🐺',
          '🐗',
        ],
      },
      {
        name: 'Food',
        icon: '🍔',
        emojis: [
          '🍏',
          '🍎',
          '🍐',
          '🍊',
          '🍋',
          '🍌',
          '🍉',
          '🍇',
          '🍓',
          '🍈',
          '🍒',
          '🍑',
          '🥭',
          '🍍',
          '🥥',
          '🥝',
          '🍅',
          '🍆',
          '🥑',
          '🥦',
          '🥬',
          '🥒',
          '🌶️',
          '🫑',
        ],
      },
    ]

    const currentCategoryEmojis = computed(() => {
      return emojiCategories[currentCategory.value].emojis
    })

    // Computed
    const canSend = computed(() => {
      // Can send if there's text OR attachments, and not currently loading
      return (messageText.value.trim() !== '' || attachments.value.length > 0) && !props.isLoading
    })

    // Methods
    const autoGrow = () => {
      if (!inputRef.value) return

      // Reset height to auto to get the correct scrollHeight
      inputRef.value.style.height = 'auto'

      // Set to scrollHeight, but limit to max height
      const maxHeight = 150 // Max height in pixels
      const newHeight = Math.min(inputRef.value.scrollHeight, maxHeight)
      inputRef.value.style.height = `${newHeight}px`
    }

    const handleEnterPress = (event) => {
      // Send message on Enter, but allow Shift+Enter for new line
      if (!event.shiftKey && canSend.value) {
        sendMessage()
      }
    }

    const sendMessage = async () => {
      if (!canSend.value) return

      const content = messageText.value.trim()
      const files = [...attachments.value] // Clone attachments to avoid mutation issues
      const replyToId = props.replyingTo?.id || '' // Ensure it's an empty string if null/undefined
      const user = authStore.user

      try {
        let sentMessageData
        if (files.length > 0) {
          const photoFile = files[0]

          sentMessageData = await messageStore.sendPhotoMessage(
            props.conversationId,
            photoFile,
            replyToId,
          )
        } else {
          // Sending a text message
          const messageData = {
            conversationId: props.conversationId,
            sender: user,
            content: content,
            type: 'text',
            ...(replyToId && { replyTo: replyToId }),
          }
          sentMessageData = await messageStore.sendTextMessage(messageData)
        }

        // Clear input and attachments after sending
        messageText.value = ''
        attachments.value = []
        if (fileInput.value) {
          fileInput.value.value = ''
        }

        if (inputRef.value) {
          inputRef.value.style.height = 'auto'
        }

        // Emit events
        emit('message-sent', sentMessageData)
        if (props.replyingTo) {
          emit('cancel-reply')
        }

        // Close emoji picker if open
        showEmojiPicker.value = false
      } catch (error) {
        console.error('Failed to send message:', error)
        // Could show an error toast here (e.g., using a notification library)
      }
    }

    const handleFileSelect = (event) => {
      const selectedFiles = Array.from(event.target.files)
      if (selectedFiles.length === 0) {
        attachments.value = []
        return
      }

      // Filter for images only, or extend logic to handle other file types
      const imageFiles = selectedFiles.filter((file) => file.type.startsWith('image/'))

      if (imageFiles.length > 0) {
        attachments.value = [imageFiles[0]] // Replace existing attachments with the new image
        messageText.value = '' // Clear text if sending only an image
      } else {
        console.warn('Only image files are currently supported for attachments.')
      }
    }
    const isTextareaDisabled = computed(() => attachments.value.length >= 1)
    const isAttachmentButtonDisabled = computed(() => messageText.value.trim() !== '')

    const removeAttachment = (index) => {
      // Revoke URL if it's an image preview to prevent memory leaks
      if (isImage(attachments.value[index])) {
        URL.revokeObjectURL(getPreviewUrl(attachments.value[index]))
      }
      attachments.value.splice(index, 1)
      // Re-enable send button if there's text, otherwise it will be disabled
      // This handles the case where attachment was the only thing allowing send
      if (messageText.value.trim() !== '') {
        canSend.value = true
      }
    }

    const isImage = (file) => {
      return file.type.startsWith('image/')
    }

    const getPreviewUrl = (file) => {
      if (isImage(file)) {
        return URL.createObjectURL(file)
      }
      // For non-image files, you might return a default icon URL or handle differently
      return null
    }

    const toggleEmojiPicker = () => {
      showEmojiPicker.value = !showEmojiPicker.value
    }

    const insertEmoji = (emoji) => {
      const startPos = inputRef.value.selectionStart
      const endPos = inputRef.value.selectionEnd

      messageText.value =
        messageText.value.substring(0, startPos) + emoji + messageText.value.substring(endPos)

      // Update cursor position
      nextTick(() => {
        const newCursorPos = startPos + emoji.length
        inputRef.value.focus()
        inputRef.value.setSelectionRange(newCursorPos, newCursorPos)
        autoGrow()
      })
    }

    // Handle click outside to close emoji picker
    const handleClickOutside = (event) => {
      const isClickInsideEmojiButton = event.target.closest('.emoji-button')
      const isClickInsideEmojiPicker = event.target.closest('.emoji-picker')

      if (!isClickInsideEmojiButton && !isClickInsideEmojiPicker && showEmojiPicker.value) {
        showEmojiPicker.value = false
      }
    }

    // Set focus to input when component mounts
    onMounted(() => {
      if (inputRef.value) {
        inputRef.value.focus()
      }
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)

      // Clean up any object URLs to prevent memory leaks
      attachments.value.forEach((file) => {
        if (isImage(file)) {
          URL.revokeObjectURL(getPreviewUrl(file))
        }
      })
    })

    // Watch for replyingTo changes to adjust the input focus
    watch(
      () => props.replyingTo,
      (newValue) => {
        if (newValue && inputRef.value) {
          nextTick(() => {
            inputRef.value.focus()
          })
        }
      },
    )

    return {
      messageText,
      attachments,
      fileInput,
      inputRef,
      showEmojiPicker,
      emojiCategories,
      currentCategory,
      currentCategoryEmojis,
      canSend,
      autoGrow,
      handleEnterPress,
      sendMessage,
      handleFileSelect,
      removeAttachment,
      isImage,
      getPreviewUrl,
      toggleEmojiPicker,
      insertEmoji,
      isTextareaDisabled,
      isAttachmentButtonDisabled,
    }
  },
}
</script>

<style scoped>
.message-input-container {
  padding: 0.75rem;
  border-top: 1px solid #e5e7eb;
  background-color: white;
  position: relative;
}

.attachments-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
  max-height: 150px;
  overflow-y: auto;
  padding: 0.5rem;
  background-color: #f9fafb;
  border-radius: 0.375rem;
}

.attachment-preview {
  position: relative;
  border-radius: 0.25rem;
  overflow: hidden;
  border: 1px solid #e5e7eb;
  background-color: white;
}

.image-preview {
  width: 80px;
  height: 80px;
}

.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.file-preview {
  display: flex;
  align-items: center;
  padding: 0.5rem;
  width: 150px;
}

.file-icon {
  margin-right: 0.5rem;
  color: #6b7280;
}

.file-name {
  font-size: 0.75rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.remove-attachment {
  position: absolute;
  top: 0.25rem;
  right: 0.25rem;
  width: 1.25rem;
  height: 1.25rem;
  border-radius: 50%;
  background-color: rgba(0, 0, 0, 0.5);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s;
}

.remove-attachment:hover {
  background-color: rgba(0, 0, 0, 0.7);
}

.input-wrapper {
  display: flex;
  align-items: flex-end;
  background-color: #f3f4f6;
  border-radius: 1.5rem;
  padding: 0.5rem 0.75rem;
}

.attachment-button,
.emoji-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  margin-right: 0.5rem;
  border: none;
  background: none;
  color: #6b7280;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.2s;
  padding: 0;
}

.attachment-button:hover,
.emoji-button:hover {
  background-color: #e5e7eb;
  color: #374151;
}

.hidden-file-input {
  display: none;
}

.message-textarea {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  resize: none;
  min-height: 1.5rem;
  max-height: 150px;
  font-size: 0.9375rem;
  line-height: 1.5;
  padding: 0.25rem 0;
  font-family: inherit;
}

.message-textarea::placeholder {
  color: #9ca3af;
}

.send-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  margin-left: 0.5rem;
  border: none;
  background-color: #3b82f6;
  color: white;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.2s;
  padding: 0;
}

.send-button:hover {
  background-color: #2563eb;
}

.send-button:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}

.emoji-picker {
  position: absolute;
  bottom: calc(100% + 0.5rem);
  right: 1rem;
  width: 320px;
  height: 250px;
  background-color: white;
  border-radius: 0.5rem;
  box-shadow:
    0 10px 15px -3px rgba(0, 0, 0, 0.1),
    0 4px 6px -2px rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
  padding: 0.75rem;
  z-index: 10;
  display: flex;
  flex-direction: column;
}

.emoji-categories {
  display: flex;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 0.5rem;
  margin-bottom: 0.5rem;
}

.category-button {
  font-size: 1.25rem;
  width: 2.5rem;
  height: 2.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 0.25rem;
  transition: all 0.2s;
}

.category-button.active {
  background-color: #f3f4f6;
}

.category-button:hover {
  background-color: #f3f4f6;
}

.emoji-list {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 0.25rem;
  overflow-y: auto;
}

.emoji-item {
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 2.25rem;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 0.25rem;
  transition: all 0.2s;
}

.emoji-item:hover {
  background-color: #f3f4f6;
}
</style>
