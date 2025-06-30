<template>
  <AppLayout>
    <div class="home-container">
      <h1>Welcome to WasaText</h1>
      <p>Select a conversation from the sidebar or start a new one.</p>

      <div class="home-stats">
        <div class="stat-card">
          <h3>{{ conversationsCount }}</h3>
          <p>Conversations</p>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { useConversationStore } from '@/store/conversations'
import AppLayout from '@/components/layout/AppLayout.vue'

const conversationStore = useConversationStore()

const conversations = computed(() => conversationStore.allConversations)
const conversationsCount = computed(() => conversations.value?.length || 0)
const unreadCount = computed(() => {
  if (!conversations.value || !Array.isArray(conversations.value)) return 0
  return conversations.value.reduce((count, conv) => count + (conv.unread_count || 0), 0)
})

// Fetch conversations when the component is mounted
onMounted(() => {
  conversationStore.fetchConversations()
})
</script>

<style scoped>
.home-container {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

h1 {
  font-size: 1.75rem;
  margin-bottom: 1rem;
  color: #333;
}

p {
  color: #666;
  margin-bottom: 2rem;
}

.home-stats {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background-color: white;
  border-radius: 8px;
  padding: 1.5rem;
  flex: 1;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  text-align: center;
}

.stat-card h3 {
  font-size: 2rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: #4a6cf7;
}

.stat-card p {
  margin-bottom: 0;
  font-size: 0.9rem;
  color: #666;
}

.action-buttons {
  display: flex;
  gap: 1rem;
}

.btn-primary,
.btn-secondary {
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  font-size: 0.9rem;
  transition:
    background-color 0.3s,
    transform 0.2s;
}

.btn-primary {
  background-color: #4a6cf7;
  color: white;
}

.btn-primary:hover {
  background-color: #3a5cdc;
  transform: translateY(-2px);
}

.btn-secondary {
  background-color: #f0f0f0;
  color: #333;
}

.btn-secondary:hover {
  background-color: #e0e0e0;
  transform: translateY(-2px);
}
</style>
