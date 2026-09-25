<template>
  <div class="request-card">
    <div class="image-wrapper">
      <img :src="request.imageUrl" alt="Property" class="property-img" />
    </div>

    <div class="card-content">
      <div class="main-info">
        <div class="property-guest-col">
          <h3>{{ request.propertyName }}</h3>
          <div class="guest-info">
            <OnyxAvatar
              :src="request.guest.avatar"
              :fullName="request.guest.name"
              size="24px"
            />
            <span class="guest-name">{{ request.guest.name }}</span>
          </div>
        </div>

        <div class="dates-col">
          <div class="date-group">
            <span class="label">CHECK-IN</span>
            <span class="value">{{ request.checkIn }}</span>
          </div>
          <div class="date-group">
            <span class="label">CHECK-OUT</span>
            <span class="value">{{ request.checkOut }}</span>
          </div>
        </div>

        <div class="payout-col">
          <span class="label">YOUR PAYOUT</span>
          <span class="payout-amount">${{ request.payout.toLocaleString() }}</span>
        </div>
      </div>

      <div class="footer-row">
        <span class="time-status">{{ request.timeStatus }}</span>

        <div v-if="isPending" class="actions">
          <OnyxButton
            label="Reject"
            color="danger"
            mode="outline"
            class="custom-cancel-btn"
            @click="rejectModalOpen = true"
          />
          <OnyxButton
            label="Approve Request"
            color="primary"
            class="custom-approve-btn"
            @click="$emit('approve')"
          />
        </div>

        <OnyxBadge
          v-else
          :color="statusDetails.color"
          class="custom-onyx-badge"
          :class="statusDetails.customClass"
        >
          {{ statusDetails.text }}
        </OnyxBadge>
      </div>
    </div>

    <OnyxAlertModal
      v-model:open="rejectModalOpen"
      density="cozy"
      class="custom-alert-modal"
      :icon="{ icon: iconCircleAttention, color: 'danger' }"
      label="Reject Request"
    >
      <p class="modal-description">
        Are you sure you want to reject the booking request from <strong>{{ request.guest.name }}</strong>?
      </p>

      <template #actions>
        <OnyxButton
          label="Cancel"
          color="neutral"
          mode="plain"
          class="modal-btn-plain"
          autofocus
          @click="rejectModalOpen = false"
        />
        <OnyxButton
          label="Reject Request"
          color="danger"
          class="modal-btn-danger"
          @click="confirmReject"
        />
      </template>
    </OnyxAlertModal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { OnyxAvatar, OnyxButton, OnyxBadge, OnyxAlertModal } from 'sit-onyx'
import { iconCircleAttention } from '@sit-onyx/icons'
import { ReservationStatus } from '@/generated/proto/reservation-api'

export interface GuestInfo {
  name: string
  avatar: string
}

export interface HostRequestData {
  id: number
  imageUrl: string
  propertyName: string
  guest: GuestInfo
  checkIn: string
  checkOut: string
  payout: number
  timeStatus: string
  status: ReservationStatus
}

const props = defineProps<{
  request: HostRequestData
}>()

const emit = defineEmits(['approve', 'reject'])

const rejectModalOpen = ref(false)

const confirmReject = () => {
  rejectModalOpen.value = false
  emit('reject')
}

const isPending = computed(() => 
  props.request.status === ReservationStatus.RESERVATION_STATUS_PENDING
)

const statusDetails = computed(() => {
  if (props.request.status === ReservationStatus.RESERVATION_STATUS_CONFIRMED) {
    return { text: 'APPROVED', color: 'success' as const, customClass: 'badge-green' }
  }
  if (props.request.status === ReservationStatus.RESERVATION_STATUS_CANCELLED) {
    return { text: 'REJECTED / CANCELLED', color: 'danger' as const, customClass: 'badge-red' }
  }
  return { text: 'PENDING', color: 'warning' as const, customClass: 'badge-blue' }
})
</script>

<style scoped>
.request-card {
  position: relative;
  display: flex;
  gap: 1.75rem;
  padding: 1.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  background: #ffffff;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
}

.image-wrapper {
  flex-shrink: 0;
  width: 260px;
  height: 145px;
  border-radius: 12px;
  overflow: hidden;
}

.property-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.main-info {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.125rem;
  font-weight: 700;
  color: #111827;
  letter-spacing: -0.01em;
}

.guest-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: #6b7280;
}

.guest-name {
  font-weight: 600;
  color: #374151;
}

.dates-col {
  display: flex;
  gap: 3.5rem;
}

.date-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.label {
  font-size: 0.6875rem;
  font-weight: 700;
  color: #9ca3af;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.value {
  font-size: 0.875rem;
  font-weight: 600;
  color: #111827;
}

.payout-col {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  line-height: 1.35;
}

.payout-amount {
  font-size: 1.25rem;
  font-weight: 800;
  color: #10b981;
}

.footer-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  padding-top: 1.25rem;
  border-top: 1px solid #f3f4f6;
}

.time-status {
  font-size: 0.8125rem;
  color: #6b7280;
}

.actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

:deep(.custom-onyx-badge) {
  font-size: 0.6875rem !important;
  padding: 0.35rem 0.75rem !important;
  font-weight: 800 !important;
  border-radius: 8px !important;
  letter-spacing: 0.03em;
}

:deep(.badge-green) {
  background-color: #d1fae5 !important;
  color: #059669 !important;
  border: none !important;
}

:deep(.badge-red) {
  background-color: #fee2e2 !important;
  color: #ef4444 !important;
  border: none !important;
}

:deep(.badge-blue) {
  background-color: #e0f2fe !important;
  color: #0284c7 !important;
  border: none !important;
}

:deep(.custom-cancel-btn) {
  border-radius: 8px !important;
  font-size: 0.7rem !important;
  font-weight: 600 !important;
  letter-spacing: 0.01em !important;
  padding: 0.25rem 0.75rem !important;
  min-height: unset !important;
  height: 32px !important;
  background-color: #ffffff !important;
  color: #dc2626 !important;
  border: 1px solid #fca5a5 !important;
}

:deep(.custom-cancel-btn:hover) {
  background-color: #fef2f2 !important;
  border-color: #ef4444 !important;
}

:deep(.custom-approve-btn) {
  border-radius: 8px !important;
  font-size: 0.7rem !important;
  font-weight: 600 !important;
  letter-spacing: 0.01em !important;
  padding: 0.25rem 0.75rem !important;
  min-height: unset !important;
  height: 32px !important;
  background-color: #1e40af !important;
  border-color: #1e40af !important;
}

:deep(.custom-approve-btn:hover) {
  background-color: #1d4ed8 !important;
}

.modal-description {
  font-size: 0.875rem;
  color: #4b5563;
  line-height: 1.4;
  margin: 0;
}

:deep(.modal-btn-plain) {
  border: none !important;
  border-radius: 8px !important;
  font-size: 0.875rem !important;
  padding: 0.4rem 1rem !important;
  color: #4b5563 !important;
  background-color: transparent !important;
}

:deep(.modal-btn-plain:hover) {
  background-color: #f3f4f6 !important;
  color: #111827 !important;
}

:deep(.modal-btn-danger) {
  border-radius: 8px !important;
  font-size: 0.875rem !important;
  padding: 0.4rem 1rem !important;
}

:deep(.modal-btn-danger:hover) {
  background-color: #fef2f2 !important; 
  border-color: #ef4444 !important; 
}
</style>