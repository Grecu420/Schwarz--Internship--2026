<template>
  <div class="bubble-wrapper" :class="isMine ? 'mine' : 'partner'">
    <div v-if="!isMine" class="avatar-container">
      <OnyxAvatar 
        :src="partnerAvatar" 
        :fullName="partnerName || 'User'"
        size="32px"
      />
    </div>

    <div class="message-content-wrapper">
      <div class="bubble-box">
        {{ message.content }}
      </div>

      <div class="message-meta">
        <span>{{ formattedTime }}</span>
        <template v-if="isMine">
          <span class="dot-separator">•</span>
          <span>{{ message.isRead ? 'Read' : 'Sent' }}</span>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { OnyxAvatar } from 'sit-onyx'
import type { Message } from '@/generated/proto/message-api'

const props = defineProps<{
  message: Message
  isMine: boolean
  partnerAvatar?: string
  partnerName?: string
}>()

const formattedTime = computed(() => {
  if (!props.message.createdAt) return ''
  const date = new Date(props.message.createdAt)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
})
</script>

<style scoped>
.bubble-wrapper {
  display: flex;
  width: 100%;
  margin-bottom: 1rem;
}

.bubble-wrapper.mine {
  justify-content: flex-end;
}

.bubble-wrapper.partner {
  justify-content: flex-start;
}

.avatar-container {
  flex-shrink: 0;
  margin-right: 0.75rem;
  margin-top: auto;
  margin-bottom: 1.25rem;
  display: flex;
  align-items: center;
}

.message-content-wrapper {
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.mine .message-content-wrapper {
  align-items: flex-end;
}

.partner .message-content-wrapper {
  align-items: flex-start;
}

.bubble-box {
  padding: 0.75rem 1.25rem;
  border-radius: 16px;
  font-size: 0.875rem;
  line-height: 1.4;
  word-break: break-word;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.mine .bubble-box {
  background-color: #2563eb;
  color: #ffffff;
  border-bottom-right-radius: 4px;
}

.partner .bubble-box {
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  color: #111827;
  border-bottom-left-radius: 4px;
}

.message-meta {
  font-size: 0.6875rem;
  color: #9ca3af;
  margin-top: 0.25rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.dot-separator {
  user-select: none;
}
</style>