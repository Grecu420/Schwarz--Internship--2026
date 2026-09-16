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
      if (this.userCache[partnerId]) return this.userCache[partnerId]

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
          
          conv.hasUnread = msgs.some((msg) => msg.senderId !== this.currentUserId && !msg.isRead)
        })

        await Promise.all(fetchMessagePromises)

        if (this.activeConversationId) {
          await this.markConversationAsRead(this.activeConversationId)
        }
      } catch (error) {
        throw new Error('Failed to load initial conversations.')
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
        throw new Error('Failed to fetch messages.')
      }
    },

    async pollUpdates() {
      if (!this.currentUserId) return

      try {
        const response = await api.get<ListConversationsResponse>(
          `/api/conversations?userId=${this.currentUserId}`,
        )
        const serverConvs = response.data.conversations || []

        for (const srvConv of serverConvs) {
          const localConvIndex = this.conversations.findIndex((c) => c.id === srvConv.id)
          const localConv = localConvIndex !== -1 ? this.conversations[localConvIndex] : undefined

          let localTime = 0
          if (localConv && localConv.updatedAt) {
             localTime = new Date(localConv.updatedAt).getTime()
          }
          
          let serverTime = 0
          if (srvConv.updatedAt) {
             serverTime = new Date(srvConv.updatedAt).getTime()
          }

          const isUpdated = serverTime > localTime

          if (!localConv || isUpdated) {
            const existingMessages = this.messagesByConversation[srvConv.id] || []
            
            const lastMsg = existingMessages[existingMessages.length - 1]
            const lastMsgId = lastMsg?.id ?? 0

            const msgResponse = await api.get<ListMessagesResponse>(
              `/api/message?conversationId=${srvConv.id}&lastMessageId=${lastMsgId}`,
            )
            const newMsgs = msgResponse.data.messages || []

            if (localConvIndex === -1) {
              this.messagesByConversation[srvConv.id] = newMsgs
              
              const partnerId = srvConv.user1Id === this.currentUserId ? srvConv.user2Id : srvConv.user1Id
              const partnerDetails = await this.getPartnerDetails(partnerId)
              
              this.conversations.unshift({
                ...srvConv,
                hasUnread: newMsgs.some((msg) => msg.senderId !== this.currentUserId && !msg.isRead),
                partnerName: partnerDetails.name,
                partnerUserName: partnerDetails.userName,
                partnerAvatar: partnerDetails.avatar,
              })
            } else if (newMsgs.length > 0) {
              const msgs = this.messagesByConversation[srvConv.id] || []
              msgs.push(...newMsgs)
              this.messagesByConversation[srvConv.id] = msgs
              
              this.conversations[localConvIndex]!.updatedAt = srvConv.updatedAt 
              this.updateUnreadStatus(srvConv.id)
              if (this.activeConversationId === srvConv.id) {
                this.markConversationAsRead(srvConv.id)
              }
            }
          }
        }
        
        if (this.activeConversationId) {
          const activeMsgs = this.messagesByConversation[this.activeConversationId] || []
          
          const waitingForReadReceipt = activeMsgs.some(
            (msg) => msg.senderId === this.currentUserId && !msg.isRead
          )

          if (waitingForReadReceipt) {
            await this.fetchMessages(this.activeConversationId)
          }
        }
        
      } catch (error) {
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
          
          const msgs = this.messagesByConversation[conversationId] || []
          msgs.push(createdMsg)
          this.messagesByConversation[conversationId] = msgs
          
          const convIndex = this.conversations.findIndex((c) => c.id === conversationId)
          if (convIndex !== -1) {
             this.conversations[convIndex]!.updatedAt = createdMsg.createdAt 
          }
        }
      } catch (error) {
        throw new Error('Failed to send message.')
      }
    },

    async markConversationAsRead(conversationId: number) {
      if (!this.currentUserId) return
      const messages = this.messagesByConversation[conversationId] || []
      const unreadMessages = messages.filter((msg) => msg.senderId !== this.currentUserId && !msg.isRead)

      if (unreadMessages.length === 0) return

      unreadMessages.forEach((msg) => { msg.isRead = true })
      this.updateUnreadStatus(conversationId)

      try {
        const updatePromises = unreadMessages.map((msg) => {
          const validCreatedAt = msg.createdAt ? new Date(msg.createdAt) : new Date()
          const request = UpdateMessageRequest.create({
            message: { ...msg, isRead: true, createdAt: validCreatedAt },
            fieldMask: ['isRead'],
          })
          return api.patch<UpdateMessageResponse>('/api/message', UpdateMessageRequest.toJSON(request))
        })
        await Promise.all(updatePromises)
      } catch (error) {
      }
    },

    updateUnreadStatus(conversationId: number) {
      if (!this.currentUserId) return
      const messages = this.messagesByConversation[conversationId] || []
      
      const convIndex = this.conversations.findIndex((c) => c.id === conversationId)
      if (convIndex !== -1) {
        this.conversations[convIndex]!.hasUnread = messages.some((msg) => msg.senderId !== this.currentUserId && !msg.isRead)
      }
    },
  },
})