<template>
  <div class="owner-card">
    <div class="owner-info">
      <OnyxAvatar size="48px" :fullName="owner.userName" />
      <div class="owner-text">
        <div class="owner-title">Hosted by {{ owner.firstName }} {{ owner.lastName }}</div>
      </div>
    </div>
    <OnyxButton
      type="button"
      mode="outline"
      label="Ask a question"
      class="ask-question-btn"
    />
  </div>
</template>

<script setup lang="ts">
import type { UserProfile } from '@/generated/proto/user-api';
import { OnyxAvatar, OnyxButton } from 'sit-onyx'

defineProps<{
  owner: UserProfile
}>()
</script>

<style scoped>
.owner-card {
  border: 1px solid #1f2937;
  background-color: rgba(255, 255, 255, 0.02);
  border-radius: 16px;
  padding: 1.25rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.owner-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.owner-title {
  font-weight: 700;
  font-size: 0.95rem;
}


:deep(.ask-question-btn) {
  border-radius: 9999px !important;
  background-color: #ffffff !important;
  color: #111827 !important;
  border: none !important;
  font-weight: 700 !important;
}
</style>