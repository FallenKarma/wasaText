<!-- src/components/conversations/ConversationItem.vue -->
<template>
  <div class="conversation-item" :class="{ active: isActive }" @click="$emit('click')">
    <div class="conversation-details">
      <div class="conversation-top">
        <h3 class="conversation-name">{{ conversationName }}</h3>
        <span class="conversation-time">{{ formatTime(conversation.lastMessage?.timestamp) }}</span>
      </div>

      <div class="conversation-bottom">
        <p class="last-message" :class="{ unread: hasUnreadMessages }">
          <span v-if="conversation.lastMessage" class="sender-name"> {{ getSenderName() }}: </span>
          {{ getLastMessagePreview() }}
        </p>

        <div class="conversation-indicators">
          <div v-if="hasUnreadMessages" class="unread-badge">
            {{ unreadCount }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useAuthStore } from '@/store/auth'

export default {
  name: 'ConversationItem',
  props: {
    conversation: {
      type: Object,
      required: true,
    },
    isActive: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const authStore = useAuthStore()

    // Computed Properties
    const currentUserId = authStore.user?.id

    const conversationName = computed(() => {
      if (props.conversation.type === 'group') {
        return props.conversation.name
      } else {
        // For direct conversations, show the other person's name
        return otherUser.value?.name || 'Unknown User'
      }
    })

    const otherUser = computed(() => {
      if (props.conversation.type !== 'direct') {
        return null
      }
      return props.conversation.participants.find((member) => member.id !== currentUserId)
    })

    const hasUnreadMessages = computed(() => {
      return props.conversation.unreadCount > 0
    })

    const unreadCount = computed(() => {
      return props.conversation.unreadCount > 99 ? '99+' : props.conversation.unreadCount
    })

    // Methods
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

    const getSenderName = () => {
      if (!props.conversation.lastMessage) return ''

      const senderId = props.conversation.lastMessage.sender.id

      if (senderId === currentUserId) {
        return 'You'
      }

      return props.conversation.lastMessage.sender.name || 'Unknown'
    }

    const getLastMessagePreview = () => {
      if (!props.conversation.lastMessage) return 'No messages yet'

      const message = props.conversation.lastMessage

      // For media messages
      if (message.type === 'photo') {
        return '🖼️ Photo'
      }

      // For text messages
      if (message.content) {
        // Truncate long messages
        if (message.deletedAt) {
          return 'Message deleted'
        }

        return message.content.length > 40
          ? message.content.substring(0, 40) + '...'
          : message.content
      }

      return 'Empty message'
    }

    const formatTime = (timestamp) => {
      if (!timestamp) return ''

      const date = new Date(timestamp)
      const now = new Date()

      // If the message is from today, show only the time
      if (date.toDateString() === now.toDateString()) {
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
      }

      // If the message is from the last 7 days, show the day name
      const weekAgo = new Date()
      weekAgo.setDate(weekAgo.getDate() - 7)

      if (date > weekAgo) {
        return date.toLocaleDateString([], { weekday: 'short' })
      }

      // Otherwise show the date
      return date.toLocaleDateString([], { month: 'short', day: 'numeric' })
    }

    return {
      currentUserId,
      conversationName,
      otherUser,
      hasUnreadMessages,
      unreadCount,
      getInitials,
      getGroupInitial,
      getSenderName,
      getLastMessagePreview,
      formatTime,
    }
  },
}
</script>

<style scoped>
.conversation-item {
  display: flex;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #f3f4f6;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.conversation-item:hover {
  background-color: #f9fafb;
}

.conversation-item.active {
  background-color: #eff6ff;
}

.conversation-details {
  flex: 1;
  min-width: 0; /* Ensure text truncates properly inside flexbox */
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.conversation-top,
.conversation-bottom {
  display: flex;
  justify-content: space-between;
}

.conversation-name {
  font-weight: 600;
  font-size: 0.9375rem;
  color: #111827;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 70%;
}

.conversation-time {
  font-size: 0.75rem;
  color: #6b7280;
  white-space: nowrap;
}

.last-message {
  font-size: 0.8125rem;
  color: #6b7280;
  margin: 0.125rem 0 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 80%;
}

.last-message.unread {
  color: #111827;
  font-weight: 500;
}

.sender-name {
  font-weight: 500;
  margin-right: 0.125rem;
}

.conversation-indicators {
  display: flex;
  align-items: center;
}

.unread-badge {
  min-width: 20px;
  height: 20px;
  background-color: #3b82f6;
  color: white;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 6px;
}
</style>
