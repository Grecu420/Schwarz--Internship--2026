<template>
  <div class="owner-card">
    <div class="owner-info">
      <OnyxAvatar size="48px" :fullName="owner.userName" />
      <div class="owner-text">
        <div class="owner-title">Hosted by {{ owner.firstName }} {{ owner.lastName }}</div>
      </div>
    </div>
    <OnyxButton v-if="owner.id !== authStore.user?.id"
      type="button"
      mode="outline"
      label="Ask a question"
      class="ask-question-btn"
      @click="askQuestion"
    />
  </div>
</template>

<script setup lang="ts">
import {
  CreateConversationRequest,
  CreateConversationResponse,
} from '@/generated/proto/conversation-api'
import type { UserProfile } from '@/generated/proto/user-api'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import { Axios, AxiosError } from 'axios'
import { OnyxAvatar, OnyxButton, useToast } from 'sit-onyx'
import { useRouter } from 'vue-router'

const props = defineProps<{
  owner: UserProfile
}>()

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const askQuestion = async () => {
  // create a conversation

  let conversationExists = false

  try {
    const userID = authStore.user?.id
    const ownerID = props.owner.id

    const request = CreateConversationRequest.create({
      user1Id: userID,
      user2Id: ownerID,
    })

    const response = await api.post<CreateConversationResponse>('/api/conversations', request)
    conversationExists = true
  } catch (error: any) {
    console.log(error)
    if (error instanceof AxiosError) {
      console.log(error)
    }

    if (error.status === 409) {
      // conversation exists
      conversationExists = true
    } else {

      toast.show({
      headline: 'Conversation Error',
      description: error.message || 'Failed to create conversation.',
      color: 'danger',
    })
    }

    
  }
  if (conversationExists) {
    router.push(`/conversations`)
  }
}
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
