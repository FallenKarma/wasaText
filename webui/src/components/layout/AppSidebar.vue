<!-- filepath: src/components/layout/AppSidebar.vue -->
<template>
  <div class="app-sidebar">
    <div class="sidebar-header">
      <h1 class="app-title">WasaText</h1>
    </div>

    <!-- Add the conversation list component here -->
    <div class="conversations-container">
      <ConversationList
        :active-conversation-id="activeConversationId"
        @conversation-selected="handleConversationSelected"
      />
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import ConversationList from '@/components/conversations/ConversationList.vue'

export default {
  name: 'AppSidebar',
  components: {
    ConversationList,
  },
  setup() {
    const router = useRouter()
    const route = useRoute()

    // Extract the conversation ID from the route if available
    const activeConversationId = computed(() => {
      const path = route.path
      if (path.startsWith('/conversations/')) {
        return path.split('/').pop()
      }
      return null
    })

    const handleConversationSelected = (conversation) => {
      // This function will receive the selected conversation from the ConversationList component
      router.push(`/conversations/${conversation.id}`)
    }

    return {
      activeConversationId,
      handleConversationSelected,
    }
  },
}
</script>

<style scoped>
.app-sidebar {
  width: 250px;
  background-color: #ffffff;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.sidebar-header {
  height: 64px;
  border-bottom: 1px solid #e5e7eb;
  text-align: center;
}

.app-title {
  padding: 1rem;
  font-size: 1.5rem;
  font-weight: 600;
  color: #4a6cf7;
}

/* Add new styles for the conversations container */
.conversations-container {
  flex: 1;
  overflow-y: auto;
  border-top: 1px solid #e5e7eb;
}
</style>
