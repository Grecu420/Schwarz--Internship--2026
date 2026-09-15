import { defineStore } from 'pinia'
import api from '@/utils/api'
import { useAuthStore } from '@/stores/auth'

import {
  Message,
  CreateMessageRequest,
  CreateMessageResponse,
  ListMessagesResponse,
  UpdateMessageRequest,
  UpdateMessageResponse,
} from '@/generated/proto/message-api'

import { Conversation, ListConversationsResponse } from '@/generated/proto/conversation-api'

import { GetUserProfileResponse } from '@/generated/proto/user-api'

export interface UIConversation extends Conversation {
  partnerName?: string
  partnerUserName?: string
  partnerAvatar?: string
  propertyName?: string
  hasUnread?: boolean
}

export const useMessageStore = defineStore('messages', {
  state: () => ({
    activeConversationId: null as number | null,
    conversations: [] as UIConversation[],
    messagesByConversation: {} as Record<number, Message[]>,
    isLoading: false,

    userCache: {} as Record<number, { name: string; userName: string; avatar: string }>,
  }),

  getters: {
    currentUserId: (): number | null => {
      const authStore = useAuthStore()
      return authStore.user?.id ?? null
    },

    activeMessages: (state): Message[] => {
      if (!state.activeConversationId) return []
      return state.messagesByConversation[state.activeConversationId] || []
    },

    activeConversationDetails: (state): UIConversation | undefined => {
      return state.conversations.find((c) => c.id === state.activeConversationId)
    },
  },

  actions: {
    async setActiveConversation(conversationId: number) {
      this.activeConversationId = conversationId

      if (!this.messagesByConversation[conversationId]) {
        await this.fetchMessages(conversationId)
      } else {
        await this.markConversationAsRead(conversationId)
      }
    },

    async getPartnerDetails(partnerId: number) {
      if (this.userCache[partnerId]) {
        return this.userCache[partnerId]
      }

      try {
        const response = await api.get<GetUserProfileResponse>(`/api/user/profile?id=${partnerId}`)
        const user = response.data.user

        const fullName = `${user?.firstName || ''} ${user?.lastName || ''}`.trim()

        const details = {
          name: fullName || user?.userName || `User #${partnerId}`,
          userName: user?.userName || `user_${partnerId}`,
          avatar: user?.profileImageUrl || '',
        }

        this.userCache[partnerId] = details
        return details
      } catch (error) {
        console.error(`Failed to fetch user ${partnerId}:`, error)
        return { name: `User #${partnerId}`, userName: `user_${partnerId}`, avatar: '' }
      }
    },

    async fetchConversations() {
      if (!this.currentUserId) return

      this.isLoading = true
      try {
        const response = await api.get<ListConversationsResponse>(
          `/api/conversations?userId=${this.currentUserId}`,
        )

        const rawConversations = response.data.conversations || []

        this.conversations = await Promise.all(
          rawConversations.map(async (conv) => {
            const partnerId = conv.user1Id === this.currentUserId ? conv.user2Id : conv.user1Id
            const partnerDetails = await this.getPartnerDetails(partnerId)

            return {
              ...conv,
              hasUnread: false,
              partnerName: partnerDetails.name,
              partnerUserName: partnerDetails.userName,
              partnerAvatar: partnerDetails.avatar,
            }
          }),
        )

        const fetchMessagePromises = this.conversations.map(async (conv) => {
          const msgResponse = await api.get<ListMessagesResponse>(
            `/api/message?conversationId=${conv.id}`,
          )
          const msgs = msgResponse.data.messages || []
          this.messagesByConversation[conv.id] = msgs

          const hasUnread = msgs.some((msg) => msg.senderId !== this.currentUserId && !msg.isRead)
          conv.hasUnread = hasUnread
        })

        await Promise.all(fetchMessagePromises)

        if (this.activeConversationId) {
          await this.markConversationAsRead(this.activeConversationId)
        }
      } catch (error) {
        console.error('Failed to fetch conversations:', error)
      } finally {
        this.isLoading = false
      }
    },

    async fetchMessages(conversationId: number) {
      try {
        const response = await api.get<ListMessagesResponse>(
          `/api/message?conversationId=${conversationId}`,
        )

        this.messagesByConversation[conversationId] = response.data.messages || []

        this.updateUnreadStatus(conversationId)
        if (this.activeConversationId === conversationId) {
          await this.markConversationAsRead(conversationId)
        }
      } catch (error) {
        console.error(`Failed to fetch messages for conv ${conversationId}:`, error)
      }
    },

    async pollUpdates() {
      if (!this.currentUserId) return

      if (this.activeConversationId) {
        await this.fetchMessages(this.activeConversationId)
      }

      for (const conv of this.conversations) {
        if (conv.id !== this.activeConversationId) {
          const response = await api.get<ListMessagesResponse>(
            `/api/message?conversationId=${conv.id}`,
          )
          this.messagesByConversation[conv.id] = response.data.messages || []
          this.updateUnreadStatus(conv.id)
        }
      }
    },

    async sendMessage(conversationId: number, content: string) {
      if (!this.currentUserId || !content.trim()) return

      try {
        const request = CreateMessageRequest.create({
          conversationId,
          senderId: this.currentUserId,
          content: content.trim(),
        })

        const jsonBody = CreateMessageRequest.toJSON(request)
        const response = await api.post<CreateMessageResponse>('/api/message', jsonBody)

        if (response.data.message) {
          const createdMsg = response.data.message

          if (!this.messagesByConversation[conversationId]) {
            this.messagesByConversation[conversationId] = []
          }
          this.messagesByConversation[conversationId].push(createdMsg)
        }
      } catch (error) {
        console.error('Failed to send message:', error)
      }
    },

    async markConversationAsRead(conversationId: number) {
      if (!this.currentUserId) return

      const messages = this.messagesByConversation[conversationId] || []

      const unreadMessages = messages.filter(
        (msg) => msg.senderId !== this.currentUserId && !msg.isRead,
      )

      if (unreadMessages.length === 0) return

      unreadMessages.forEach((msg) => {
        msg.isRead = true
      })
      this.updateUnreadStatus(conversationId)

      try {
        const updatePromises = unreadMessages.map((msg) => {
          const validCreatedAt = msg.createdAt ? new Date(msg.createdAt) : new Date()

          const request = UpdateMessageRequest.create({
            message: {
              ...msg,
              isRead: true,
              createdAt: validCreatedAt,
            },
            fieldMask: ['isRead'],
          })

          const jsonBody = UpdateMessageRequest.toJSON(request)
          return api.patch<UpdateMessageResponse>('/api/message', jsonBody)
        })

        await Promise.all(updatePromises)
      } catch (error) {
        console.error('Failed to update read status on backend:', error)
      }
    },

    updateUnreadStatus(conversationId: number) {
      if (!this.currentUserId) return
      const messages = this.messagesByConversation[conversationId] || []
      const hasUnread = messages.some((msg) => msg.senderId !== this.currentUserId && !msg.isRead)

      const conv = this.conversations.find((c) => c.id === conversationId)
      if (conv) {
        conv.hasUnread = hasUnread
      }
    },
  },
})