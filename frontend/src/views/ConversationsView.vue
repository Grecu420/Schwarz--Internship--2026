<template>
  <div class="messages-container">
    <div 
      v-if="messageStore.isLoading && !messageStore.conversations.length" 
      class="loading-state"
    >
      <div class="spinner"></div>
      <span class="loading-text">Loading conversations...</span>
    </div>

    <template v-else>
      <aside class="sidebar-panel">
        <ConversationSidebar />
      </aside>

      <main class="chat-panel">
        <EmptyState v-if="!activeConversationId" />
        <ChatArea v-else :conversation-id="activeConversationId" />
      </main>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useMessageStore } from '@/stores/messages'

import ConversationSidebar from '@/components/ConversationSidebar.vue'
import ChatArea from '@/components/ChatArea.vue'
import EmptyState from '@/components/EmptyState.vue'

const messageStore = useMessageStore()

const activeConversationId = computed(() => messageStore.activeConversationId)

let pollingTimer: number | null = null

onMounted(async () => {
  await messageStore.fetchConversations()

  pollingTimer = window.setInterval(() => {
    messageStore.pollUpdates()
  }, 10000)
})

onUnmounted(() => {
  if (pollingTimer !== null) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
})
</script>

<style scoped>
.messages-container {
  max-width: 1280px;
  height: calc(100vh - 100px);
  margin: 0 auto;
  padding: 1.5rem;
  display: flex;
  gap: 1.5rem;
  box-sizing: border-box;
}

.loading-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  gap: 0.75rem;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 4px solid #1e40af;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.loading-text {
  font-size: 0.875rem;
  color: #6b7280;
  font-weight: 500;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.sidebar-panel {
  width: 320px;
  min-width: 300px;
  max-width: 380px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background-color: #ffffff;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.chat-panel {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background-color: #f9fafb;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}
</style>