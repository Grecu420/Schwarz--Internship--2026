<template>
  <div class="empty-state-container">
    <div class="empty-state-content">
      <div class="icon-wrapper">
        <OnyxIcon :icon="iconMessageDots" class="chat-icon" />
      </div>

      <h3 class="empty-title">Your Messages</h3>
      <p class="empty-description">
        Select a conversation from the sidebar on the left to start chatting or view message history.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { OnyxIcon } from 'sit-onyx'
import { iconMessageDots } from '@sit-onyx/icons'
</script>

<style scoped>
.empty-state-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 2rem;
  background-color: #f9fafb;
}

.empty-state-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  max-width: 320px;
}

.icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 4rem;
  height: 4rem;
  background-color: #eff6ff;
  border-radius: 50%;
  margin-bottom: 1.25rem;
}

:deep(.chat-icon) {
  width: 2.25rem;
  height: 2.25rem;
  color: #2563eb;
}

.empty-title {
  margin: 0 0 0.5rem 0;
  font-size: 1.125rem;
  font-weight: 700;
  color: #111827;
}

.empty-description {
  margin: 0;
  font-size: 0.875rem;
  color: #6b7280;
  line-height: 1.5;
}
</style>