<template>
  <div class="message-list-wrapper">
    <div class="message-list" ref="messageListRef">
      <div v-if="isLoadingMessages && !messages.length" class="loading-state">
        <div class="loading-spinner"></div>
        <p>Loading messages...</p>
      </div>

      <div v-else-if="!messages" class="empty-state">
        <div class="empty-icon">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="40"
            height="40"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
        </div>
        <p>No messages yet</p>
        <p class="help-text">Start the conversation by sending a message below.</p>
      </div>

      <template v-else>
        <div class="messages-container">
          <template v-for="(message, index) in sortedMessages" :key="message.id">
            <!-- Message item -->
            <MessageItem
              :message="message"
              :isOwn="isOwnMessage(message)"
              :showAvatar="shouldShowAvatar(message, index)"
              @reaction="handleReaction"
              @reply="handleReply"
            />
          </template>
        </div>

        <!-- Bottom scroll button -->
        <button v-if="showScrollToBottom" @click="scrollToBottom" class="scroll-bottom-button">
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
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>
      </template>
    </div>

    <!-- Reply box - moved outside message-list for proper positioning -->
    <div v-if="replyingTo" class="reply-box">
      <div class="reply-content">
        <div class="reply-indicator">
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
            <polyline points="9 14 4 9 9 4"></polyline>
            <path d="M20 20v-7a4 4 0 0 0-4-4H4"></path>
          </svg>
        </div>
        <div class="reply-text">
          <div class="reply-author">{{ replyingTo.sender.name }}</div>
          <span v-if="!replyingTo.content.startsWith('/')" class="reply-message">{{
            truncateText(replyingTo.content, 40)
          }}</span>
          <span v-else class="replied-message">🖼️ Photo</span>
        </div>
      </div>
      <button @click="cancelReply" class="cancel-reply-button">
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
</template>

<script>
import { ref, computed, onMounted, onUpdated, watch, nextTick, onBeforeUnmount } from 'vue'
import MessageItem from './MessageItem.vue'
import { useAuthStore } from '@/store/auth'
import { useMessageStore } from '@/store/messages'

export default {
  name: 'MessageList',
  components: {
    MessageItem,
  },
  props: {
    conversationId: {
      type: [String, Number],
      required: true,
    },
    isLoadingMessages: {
      type: Boolean,
      default: false,
    },
    replyingTo: {
      type: Object,
      default: null, // Object with message data or null if not replying
    },
  },
  emits: ['reply-message', 'cancel-reply'],
  setup(props, { emit }) {
    const messageListRef = ref(null)
    const showScrollToBottom = ref(false)
    const autoScrollToBottom = ref(true)
    const isUserScrolling = ref(false)
    const scrollTimeout = ref(null)
    const authStore = useAuthStore()
    const messageStore = useMessageStore()
    const pollingInterval = ref(null)
    const isPolling = ref(false)
    const POLLING_INTERVAL_MS = 3000 // Poll every 3 seconds

    // Computed
    const currentUserId = computed(() => authStore.user?.id)
    const messages = computed(() => messageStore.allMessages)

    // Function to "hydrate" the replied-to message data
    const hydrateReplyToMessages = (messages) => {
      if (!messages || messages.length === 0) return []

      // Create a map for quick lookup of messages by ID
      const messageMap = new Map(messages.map((msg) => [msg.id, msg]))

      return messages.map((message) => {
        // Check if the message has a replyTo object with a valid String ID
        if (message.replyTo) {
          const repliedMessageId = message.replyTo
          const repliedToMessage = messageMap.get(repliedMessageId)

          // If the replied-to message is found in the current list,
          // attach its relevant data.
          if (repliedToMessage) {
            return {
              ...message,
              repliedToMessageData: {
                // New property to hold the hydrated data
                id: repliedToMessage.id,
                sender: repliedToMessage.sender.name,
                content: repliedToMessage.content,
              },
            }
          }
        }
        return message // Return the message as is if no reply or invalid replyTo
      })
    }

    // Sort messages to ensure newest messages are at the bottom
    const sortedMessages = computed(() => {
      if (!messages.value.length) return []

      // First, hydrate the replyTo data
      const hydratedMessages = hydrateReplyToMessages(messages.value)

      // Then, sort the hydrated messages
      return [...hydratedMessages].sort((a, b) => {
        const dateA = new Date(a.createdAt)
        const dateB = new Date(b.createdAt)
        return dateA - dateB // Ascending order (oldest first, newest at bottom)
      })
    })

    // Methods
    const isOwnMessage = (message) => {
      return message.sender.id === currentUserId.value
    }

    const shouldShowAvatar = (message, index) => {
      // Show avatar if it's the first message or if the previous message is from a different sender
      if (index === 0) return true

      const prevMessage = sortedMessages.value[index - 1]
      return prevMessage.sender.id !== message.sender.id
    }

    const formatMessageDate = (dateString) => {
      const date = new Date(dateString)
      const today = new Date()
      const yesterday = new Date(today)
      yesterday.setDate(yesterday.getDate() - 1)

      if (date.toDateString() === today.toDateString()) {
        return 'Today'
      } else if (date.toDateString() === yesterday.toDateString()) {
        return 'Yesterday'
      } else {
        return date.toLocaleDateString(undefined, {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
        })
      }
    }

    const isScrolledToBottom = () => {
      if (!messageListRef.value) return false

      const { scrollTop, scrollHeight, clientHeight } = messageListRef.value
      return scrollHeight - scrollTop - clientHeight <= 20 // Increased threshold for better reliability
    }

    const handleScroll = () => {
      if (!messageListRef.value) return

      if (scrollTimeout.value) {
        clearTimeout(scrollTimeout.value)
      }

      isUserScrolling.value = true

      const scrolledToBottom = isScrolledToBottom()
      showScrollToBottom.value = !scrolledToBottom

      if (scrolledToBottom) {
        autoScrollToBottom.value = true
      } else {
        const { scrollTop, scrollHeight, clientHeight } = messageListRef.value
        if (scrollTop < scrollHeight - clientHeight - 50) {
          autoScrollToBottom.value = false
        }
      }

      // Reset user scrolling flag after a delay
      scrollTimeout.value = setTimeout(() => {
        isUserScrolling.value = false
      }, 150)
    }

    const scrollToBottom = (force = false) => {
      if (!messageListRef.value) return

      // Use requestAnimationFrame for smoother scrolling
      requestAnimationFrame(() => {
        if (messageListRef.value && (autoScrollToBottom.value || force)) {
          const element = messageListRef.value
          element.scrollTop = element.scrollHeight
          autoScrollToBottom.value = true
          showScrollToBottom.value = false
        }
      })
    }

    const handleReaction = ({ messageId, reaction }) => {
      // TODO: Implement reaction handling
      console.log('Reaction:', messageId, reaction)
    }

    const handleReply = (message) => {
      emit('reply-message', message)
    }

    const cancelReply = () => {
      emit('cancel-reply', null)
    }

    const truncateText = (text, maxLength) => {
      if (!text) return ''
      return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
    }

    const startPolling = () => {
      if (pollingInterval.value) {
        clearInterval(pollingInterval.value)
      }

      isPolling.value = true
      pollingInterval.value = setInterval(async () => {
        try {
          // Only poll if we're not already loading and the component is still mounted
          if (!props.isLoadingMessages && messageListRef.value) {
            await fetchMessages()
          }
        } catch (error) {
          console.error('Error during polling:', error)
          // Don't stop polling on error, just log it
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

    // Modify your existing fetchMessages function to be smarter about new messages
    const fetchMessages = async () => {
      try {
        const previousMessageIds = new Set(messages.value.map((msg) => msg.id))
        await messageStore.fetchMessages(props.conversationId)

        const hasNewMessages = messages.value.some((msg) => !previousMessageIds.has(msg.id))

        if (hasNewMessages && autoScrollToBottom.value) {
          nextTick(() => {
            scrollToBottom()
          })
        }
      } catch (error) {
        console.error('Failed to fetch messages:', error)
      }
    }

    // Add visibility change handler to pause/resume polling
    const handleVisibilityChange = () => {
      if (document.hidden) {
        stopPolling()
      } else {
        startPolling()
      }
    }

    // Lifecycle hooks
    onMounted(() => {
      fetchMessages()
      startPolling()

      // Add visibility change listener
      document.addEventListener('visibilitychange', handleVisibilityChange)

      if (messageListRef.value) {
        messageListRef.value.addEventListener('scroll', handleScroll, { passive: true })

        // Initial scroll to bottom after ensuring content is rendered
        nextTick(() => {
          setTimeout(() => {
            scrollToBottom(true)
          }, 100)
        })
      }
    })

    onBeforeUnmount(() => {
      stopPolling() // Stop polling when component unmounts
      document.removeEventListener('visibilitychange', handleVisibilityChange)

      if (messageListRef.value) {
        messageListRef.value.removeEventListener('scroll', handleScroll)
      }
      if (scrollTimeout.value) {
        clearTimeout(scrollTimeout.value)
      }
    })

    onUpdated(() => {
      if (props.isLoadingMessages || isUserScrolling.value) return

      // Scroll to bottom for new messages if auto-scroll is enabled
      if (autoScrollToBottom.value) {
        nextTick(() => {
          scrollToBottom()
        })
      }
    })

    // Watch for new messages with improved logic
    watch(
      () => props.conversationId,
      (newConversationId, oldConversationId) => {
        if (newConversationId !== oldConversationId) {
          stopPolling()
          if (newConversationId) {
            fetchMessages()
            startPolling()
          }
        }
      },
    )

    // Watch for changes in loading state
    watch(
      () => props.isLoadingMessages,
      (isLoading, wasLoading) => {
        // When loading finishes and we have messages, scroll to bottom if needed
        if (wasLoading && !isLoading && sortedMessages.value.length > 0) {
          if (autoScrollToBottom.value) {
            nextTick(() => {
              scrollToBottom()
            })
          }
        }
      },
    )

    return {
      messages,
      messageListRef,
      showScrollToBottom,
      sortedMessages,
      isOwnMessage,
      shouldShowAvatar,
      formatMessageDate,
      scrollToBottom,
      handleReaction,
      handleReply,
      cancelReply,
      truncateText,
    }
  },
}
</script>

<style scoped>
.message-list-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  min-height: 0; /* Important for flex shrinking in parent container */
  overflow: hidden; /* Prevent wrapper from overflowing */
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 1rem 0.5rem;
  display: flex;
  flex-direction: column;
  scroll-behavior: smooth;
  min-height: 0; /* Important for flex shrinking */
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #6b7280;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top: 3px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #6b7280;
  text-align: center;
}

.empty-icon {
  color: #9ca3af;
  margin-bottom: 1rem;
}

.empty-state p {
  margin: 0;
  font-size: 0.875rem;
}

.empty-state .help-text {
  color: #9ca3af;
  margin-top: 0.5rem;
}

.messages-container {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.date-separator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 1rem 0;
}

.date-separator span {
  background-color: #f3f4f6;
  color: #6b7280;
  font-size: 0.75rem;
  padding: 0.25rem 0.75rem;
  border-radius: 0.375rem;
}

.scroll-bottom-button {
  position: absolute;
  bottom: 1rem;
  right: 1rem;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border: none;
  transition: background-color 0.2s ease;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.scroll-bottom-button:hover {
  background-color: #2563eb;
}

/* Fixed reply box positioning - stays at bottom of MessageList component */
.reply-box {
  background-color: #f9fafb;
  border-top: 1px solid #e5e7eb;
  padding: 0.75rem 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
  z-index: 5;
  margin-top: auto; /* Push to bottom of wrapper */
}

.reply-content {
  display: flex;
  align-items: center;
  flex: 1;
}

.reply-indicator {
  color: #6b7280;
  margin-right: 0.5rem;
  flex-shrink: 0;
}

.reply-text {
  flex: 1;
  min-width: 0; /* Allow text to shrink */
}

.reply-author {
  font-size: 0.75rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 0.125rem;
}

.reply-message,
.replied-message {
  font-size: 0.75rem;
  color: #6b7280;
  display: block;
  word-break: break-word;
}

.cancel-reply-button {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 0.25rem;
  transition: background-color 0.2s ease;
  flex-shrink: 0;
}

.cancel-reply-button:hover {
  background-color: #e5e7eb;
}
</style>
