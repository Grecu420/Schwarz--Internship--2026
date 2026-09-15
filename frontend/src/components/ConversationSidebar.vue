<template>
  <div class="sidebar-container">
    <div class="search-header">
      <div class="search-input-wrapper">
        <svg class="search-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Search messages..." 
          class="search-input"
        />
      </div>
    </div>

    <div class="conversation-list">
      <div v-if="filteredConversations.length === 0" class="no-conversations">
        No conversations found
      </div>

      <div 
        v-for="conv in filteredConversations" 
        :key="conv.id"
        @click="messageStore.setActiveConversation(conv.id)"
        class="conversation-item"
        :class="{ 'active': conv.id === messageStore.activeConversationId }"
      >
        <div class="avatar-wrapper">
          <OnyxAvatar 
            :src="conv.partnerAvatar" 
            :fullName="conv.partnerName || 'User'" 
          />
        </div>

        <div class="details-wrapper">
          <div class="top-row">
            <h4 class="partner-name" :class="{ 'unread-text': conv.hasUnread }">
              {{ conv.partnerName }}
            </h4>
          </div>

          <div class="bottom-row">
            <p class="last-message" :class="{ 'unread-text': conv.hasUnread }">
              {{ getLastMessageText(conv.id) }}
            </p>

            <span v-if="conv.hasUnread" class="unread-badge"></span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { OnyxAvatar } from 'sit-onyx'
import { useMessageStore } from '@/stores/messages'

const messageStore = useMessageStore()
const searchQuery = ref('')

const filteredConversations = computed(() => {
  if (!searchQuery.value.trim()) return messageStore.conversations
  return messageStore.conversations.filter(c => 
    c.partnerName?.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const getLastMessageText = (conversationId: number) => {
  const msgs = messageStore.messagesByConversation[conversationId]
  if (!msgs || msgs.length === 0) return 'No messages yet'
  
  const lastMsg = msgs[msgs.length - 1]
  return lastMsg?.content ?? 'No messages yet'
}
</script>

<style scoped>
.sidebar-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #ffffff;
}

.search-header {
  padding: 1rem;
  border-bottom: 1px solid #f3f4f6;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  width: 1rem;
  height: 1rem;
  color: #9ca3af;
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 0.5rem 0.75rem 0.5rem 2.25rem;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 0.875rem;
  color: #111827;
  outline: none;
  transition: all 0.15s ease;
}

.search-input:focus {
  background-color: #ffffff;
  border-color: #2563eb;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.15);
}

.conversation-list {
  flex: 1;
  overflow-y: auto;
}

.no-conversations {
  padding: 2rem 1rem;
  text-align: center;
  font-size: 0.875rem;
  color: #9ca3af;
}

.conversation-item {
  display: flex;
  gap: 0.75rem;
  padding: 1rem;
  border-bottom: 1px solid #f9fafb;
  cursor: pointer;
  transition: background-color 0.15s ease;
  position: relative;
}

.conversation-item:hover {
  background-color: #f9fafb;
}

.conversation-item.active {
  background-color: #eff6ff;
  border-left: 4px solid #2563eb;
}

.avatar-wrapper {
  flex-shrink: 0;
  display: flex;
  align-items: center;
}

.details-wrapper {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.top-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 0.25rem;
}

.partner-name {
  margin: 0;
  font-size: 0.875rem;
  color: #374151;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bottom-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.5rem;
}

.last-message {
  margin: 0;
  font-size: 0.8125rem;
  color: #6b7280;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.unread-text {
  font-weight: 700 !important;
  color: #111827 !important;
}

.unread-badge {
  width: 0.625rem;
  height: 0.625rem;
  background-color: #2563eb;
  border-radius: 50%;
  flex-shrink: 0;
}
</style>