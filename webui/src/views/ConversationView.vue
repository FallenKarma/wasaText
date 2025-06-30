<!-- src/views/ConversationView.vue -->
<template>
  <AppLayout>
    <div class="conversation-container">
      <div v-if="isLoading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>Loading conversation...</p>
      </div>

      <template v-else-if="conversation">
        <ConversationHeader :conversation="conversation" @info="showConversationInfo" />

        <MessageList
          :conversationId="conversationId"
          :isLoadingMessages="isLoadingMessages"
          :replyingTo="replyingTo"
          @load-more="loadMoreMessages"
          @reply-message="handleReplyMessage"
          @cancel-reply="handleCancelReply"
        />

        <MessageInput
          :conversationId="conversationId"
          :replyingTo="replyingTo"
          @message-sent="handleMessageSent"
          @cancel-reply="handleCancelReply"
        />
      </template>

      <div v-else class="empty-state">
        <div class="empty-icon">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="48"
            height="48"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
          </svg>
        </div>
        <h3>Conversation not found</h3>
        <p>The conversation you're looking for doesn't exist or you don't have access to it.</p>
        <button class="back-button" @click="goBack">Go Back</button>
      </div>

      <!-- Conversation Info Sidebar -->
      <div v-if="showInfo" class="conversation-info-sidebar" :class="{ active: showInfo }">
        <div class="info-header">
          <h3>Conversation Info</h3>
          <button class="close-button" @click="showInfo = false">
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
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="info-content">
          <div class="conversation-avatar-large">
            <img
              v-if="conversation?.avatarUrl"
              :src="conversation.avatarUrl"
              alt="Conversation avatar"
            />
            <div v-else class="avatar-placeholder-large">
              {{ getInitials(conversation?.name) }}
            </div>
          </div>

          <h2 class="conversation-name-large">{{ conversation?.name }}</h2>

          <div class="info-section">
            <h4>Members</h4>
            <div class="member-list">
              <div v-for="member in conversation?.members" :key="member.id" class="member-item">
                <div class="member-avatar">
                  <img v-if="member.avatarUrl" :src="member.avatarUrl" alt="Member avatar" />
                  <div v-else class="avatar-placeholder">
                    {{ getInitials(member.name) }}
                  </div>
                </div>
                <div class="member-name">{{ member.name }}</div>
                <div v-if="member.isOnline" class="member-status online"></div>
              </div>
            </div>
          </div>

          <div class="info-actions">
            <button class="danger-button" @click="leaveConversation">Leave Conversation</button>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import ConversationHeader from '@/components/conversations/ConversationHeader.vue'
import MessageList from '@/components/messages/MessageList.vue'
import MessageInput from '@/components/messages/MessageInput.vue'
import { useConversationStore } from '@/store/conversations'
import { useMessageStore } from '@/store/messages'
import { useUserStore } from '@/store/users'

export default {
  name: 'ConversationView',
  components: {
    AppLayout,
    ConversationHeader,
    MessageList,
    MessageInput,
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const conversationStore = useConversationStore()
    const messageStore = useMessageStore()

    // State
    const isLoading = ref(false)
    const isLoadingMessages = ref(false)
    const showInfo = ref(false)
    const replyingTo = ref(null)

    // Computed
    const conversationId = computed(() => route.params.conversationId)
    const conversation = computed(() => conversationStore.currentConversation)
    const messages = computed(() => conversation.value?.messages || [])

    // Methods
    const fetchConversation = async () => {
      if (!conversationId.value) return

      try {
        isLoading.value = true
        await conversationStore.fetchConversation(conversationId.value)
      } catch (error) {
        console.error('Failed to fetch conversation:', error)
      } finally {
        isLoading.value = false
      }
    }

    const loadMoreMessages = async () => {
      if (!conversationId.value || messageStore.isLoading) return

      try {
        isLoadingMessages.value = true
        await conversationStore.loadMoreMessages(conversationId.value)
      } catch (error) {
        console.error('Failed to load more messages:', error)
      } finally {
        isLoadingMessages.value = false
      }
    }

    const handleMessageSent = (message) => {
      // Optimistically add message to the list
      conversationStore.addMessageToConversation(conversationId.value, message)
      replyingTo.value = null

      // Scroll to the bottom of the message list
      setTimeout(() => {
        const messageList = document.querySelector('.message-list')
        if (messageList) {
          messageList.scrollTop = messageList.scrollHeight
        }
      }, 100)
    }

    const handleReplyMessage = (message) => {
      replyingTo.value = message
    }

    const handleCancelReply = () => {
      replyingTo.value = null
    }

    const showConversationInfo = () => {
      showInfo.value = true
      fetchMediaFiles()
    }

    const leaveConversation = async () => {
      if (!conversationId.value) return

      if (confirm('Are you sure you want to leave this conversation?')) {
        try {
          await conversationStore.leaveConversation(conversationId.value)
          router.push('/')
        } catch (error) {
          console.error('Failed to leave conversation:', error)
        }
      }
    }

    const goBack = () => {
      router.push('/')
    }

    const getInitials = (name) => {
      if (!name) return ''
      return name.charAt(0).toUpperCase()
    }

    // Lifecycle
    onMounted(() => {
      fetchConversation()
    })

    // Watch for route changes
    watch(
      () => route.params.conversationId,
      (newId, oldId) => {
        if (newId !== oldId) {
          fetchConversation()
          showInfo.value = false
        }
      },
    )

    return {
      isLoading,
      isLoadingMessages,
      conversationId,
      conversation,
      messages,
      showInfo,
      replyingTo,
      loadMoreMessages,
      handleReplyMessage,
      handleCancelReply,
      handleMessageSent,
      showConversationInfo,
      leaveConversation,
      goBack,
      getInitials,
    }
  },
}
</script>

<style scoped>
.conversation-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  position: relative;
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
  padding: 2rem;
  text-align: center;
}

.empty-icon {
  color: #9ca3af;
  margin-bottom: 1rem;
}

.empty-state h3 {
  margin: 0 0 0.5rem;
  font-size: 1.25rem;
  font-weight: 500;
  color: #4b5563;
}

.empty-state p {
  margin: 0 0 1.5rem;
  max-width: 24rem;
}

.back-button {
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

.back-button:hover {
  background-color: #2563eb;
}

.conversation-info-sidebar {
  position: absolute;
  top: 0;
  right: 0;
  width: 320px;
  height: 100%;
  background-color: white;
  border-left: 1px solid #e5e7eb;
  transform: translateX(100%);
  transition: transform 0.3s ease;
  z-index: 10;
  display: flex;
  flex-direction: column;
}

.conversation-info-sidebar.active {
  transform: translateX(0);
}

.info-header {
  padding: 1rem;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-header h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 500;
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

.info-content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.conversation-avatar-large {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  margin-bottom: 1rem;
}

.conversation-avatar-large img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder-large {
  width: 100%;
  height: 100%;
  background-color: #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 2rem;
}

.conversation-name-large {
  margin: 0 0 1.5rem;
  font-size: 1.25rem;
  font-weight: 600;
  text-align: center;
}

.info-section {
  width: 100%;
  margin-bottom: 1.5rem;
}

.info-section h4 {
  margin: 0 0 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.member-list {
  display: flex;
  flex-direction: column;
}

.member-item {
  display: flex;
  align-items: center;
  padding: 0.5rem 0;
}

.member-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 0.75rem;
}

.member-avatar img {
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

.member-name {
  flex: 1;
  font-size: 0.875rem;
}

.member-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.member-status.online {
  background-color: #10b981;
}

.info-actions {
  margin-top: auto;
  width: 100%;
  padding-top: 1rem;
}

.danger-button {
  width: 100%;
  padding: 0.625rem;
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.danger-button:hover {
  background-color: #dc2626;
}
</style>
