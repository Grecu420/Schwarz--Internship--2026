<template>
  <div class="chat-area-container">
    
    <header v-if="activeChat" class="chat-header">
      <div class="partner-info">
        <OnyxAvatar 
          :src="activeChat.partnerAvatar" 
          :fullName="activeChat.partnerName || 'User'" 
        />
        <div class="partner-details">
          <h2 class="partner-name">{{ activeChat.partnerName }}</h2>
          <p class="partner-username">{{activeChat.partnerUserName}}</p>
        </div>
      </div>
    </header>

    <div ref="messagesFeedRef" class="messages-feed">
      <template v-for="(msg, index) in messages" :key="msg.id">
        
        <div v-if="isFirstUnread(msg, index)" class="unread-divider">
          <div class="divider-line"></div>
          <span class="divider-text">Mesaje noi</span>
          <div class="divider-line"></div>
        </div>

        <MessageBubble 
          :message="msg"
          :is-mine="msg.senderId === currentUserId"
          :partner-avatar="activeChat?.partnerAvatar"
          :partner-name="activeChat?.partnerName"
        />
        
      </template>
    </div>

    <footer class="chat-input-area">
      <form @submit.prevent="handleSendMessage" class="input-form">
        <input 
          v-model="newMessageText"
          type="text" 
          placeholder="Type a message..." 
          class="message-input"
        />
        <button 
          type="submit" 
          class="send-btn"
          :disabled="!newMessageText.trim()"
        >
          <OnyxIcon :icon="iconSend" class="send-icon" />
        </button>
      </form>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { OnyxAvatar, OnyxIcon, useToast } from 'sit-onyx' 
import { iconSend } from "@sit-onyx/icons";

import { useMessageStore } from '@/stores/messages'
import MessageBubble from './MessageBubble.vue';
import type { Message } from '@/generated/proto/message-api'

const props = defineProps<{
  conversationId: number
}>()

const messageStore = useMessageStore()
const toast = useToast() 

const newMessageText = ref('')
const messagesFeedRef = ref<HTMLDivElement | null>(null)

const currentUserId = computed(() => messageStore.currentUserId)
const activeChat = computed(() => messageStore.activeConversationDetails)
const messages = computed(() => messageStore.activeMessages)

const isFirstUnread = (msg: Message, index: number): boolean => {
  if (msg.senderId === currentUserId.value || msg.isRead) return false
  const prevMsg = messages.value[index - 1]
  return !prevMsg || prevMsg.isRead || prevMsg.senderId === currentUserId.value
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesFeedRef.value) {
    messagesFeedRef.value.scrollTop = messagesFeedRef.value.scrollHeight
  }
}

const handleSendMessage = async () => {
  const text = newMessageText.value.trim()
  if (!text) return
  
  try {
    await messageStore.sendMessage(props.conversationId, text)
    newMessageText.value = ''
    scrollToBottom()
  } catch (error: any) {
    toast.show({
      headline: 'Error',
      description: error.message || 'Could not send message. Please try again.',
      color: 'danger'
    })
  }
}

watch(
  () => messages.value.length,
  () => {
    scrollToBottom()
  },
  { immediate: true }
)
</script>

<style scoped>
.chat-area-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #f9fafb;
}

.chat-header {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e5e7eb;
  background-color: #ffffff;
}

.partner-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.partner-name {
  margin: 0;
  font-size: 1rem;
  font-weight: 700;
  color: #111827;
}

.partner-username {
  margin: 0;
  font-size: 0.75rem;
  color: #6b7280;
}

.messages-feed {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
}

.unread-divider {
  display: flex;
  align-items: center;
  margin: 1.5rem 0;
}

.divider-line {
  flex: 1;
  border-top: 1px solid #bfdbfe;
}

.divider-text {
  padding: 0 1rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: #2563eb;
  background-color: #f9fafb;
}

.chat-input-area {
  padding: 1rem;
  background-color: #ffffff;
  border-top: 1px solid #e5e7eb;
}

.input-form {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.message-input {
  flex: 1;
  padding: 0.75rem 1.25rem;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 9999px;
  font-size: 0.875rem;
  outline: none;
  transition: all 0.15s ease;
}

.message-input:focus {
  background-color: #ffffff;
  border-color: #2563eb;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.15);
}

.send-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.75rem;
  height: 2.75rem;
  background-color: #2563eb;
  color: #ffffff;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.send-btn:hover:not(:disabled) {
  background-color: #1d4ed8;
}

.send-btn:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
  opacity: 0.6;
}

.send-icon {
  width: 1.25rem;
  height: 1.25rem;
  transform: translateX(3px); 
}
</style>