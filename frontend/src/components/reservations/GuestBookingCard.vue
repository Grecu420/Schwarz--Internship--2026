<template>
  <div class="booking-card">
    <div class="image-wrapper">
      <img :src="booking.imageUrl" alt="Property" class="property-img" />
    </div>

    <div class="card-content">
      <div class="top-row">
        <div class="title-block">
          <h3>{{ booking.propertyName }}</h3>
          <p class="location">{{ booking.location }}</p>
        </div>
        <!-- OnyxBadge cu styling subtil -->
        <OnyxBadge
          :color="statusDetails.color"
          class="custom-onyx-badge"
          :class="statusDetails.customClass"
        >
          {{ statusDetails.text }}
        </OnyxBadge>
      </div>

      <div class="bottom-row">
        <div class="dates-grid">
          <div class="date-item">
            <span class="label">CHECK-IN</span>
            <span class="value">{{ booking.checkIn }}</span>
          </div>
          <div class="date-item">
            <span class="label">CHECK-OUT</span>
            <span class="value">{{ booking.checkOut }}</span>
          </div>
        </div>

        <div class="price-and-action">
          <div class="price-block">
            <span class="nights">{{ booking.nights }} nights total</span>
            <span class="price">${{ booking.price.toLocaleString() }}</span>
          </div>

          <!-- OnyxButton deschide modalul de alertă -->
          <OnyxButton
            v-if="actionLabel"
            :label="actionLabel"
            color="danger"
            mode="outline"
            class="custom-cancel-btn"
            @click="alertOpen = true"
          />
          <span v-else class="no-action-text">Booking Declined</span>
        </div>
      </div>
    </div>

    <OnyxAlertModal
      v-model:open="alertOpen"
      density="cozy"
      class="custom-alert-modal"
      :icon="{ icon: iconCircleAttention, color: 'danger' }"
      label="Cancel Reservation"
    >
      <p class="modal-description">
        Are you sure that you want to cancel your booking for <strong>{{ booking.propertyName }}</strong>? This action cannot be reverted.
      </p>

      <template #actions>
        <OnyxButton
          label="Keep Booking"
          color="neutral"
          mode="plain"
          class="modal-btn-plain"
          autofocus
          @click="alertOpen = false"
        />
        <OnyxButton
          label="Yes, Cancel"
          color="danger"
          class="modal-btn-danger"
          @click="confirmCancel"
        />
      </template>
    </OnyxAlertModal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { OnyxBadge, OnyxButton, OnyxAlertModal } from 'sit-onyx'
import { iconCircleAttention } from '@sit-onyx/icons'
import { ReservationStatus } from '@/generated/proto/reservation-api'

export interface GuestBookingData {
  id: number
  imageUrl: string
  propertyName: string
  location: string
  checkIn: string
  checkOut: string
  nights: number
  price: number
  status: ReservationStatus 
}

const props = defineProps<{
  booking: GuestBookingData
}>()

const emit = defineEmits(['action'])

const alertOpen = ref(false)

const confirmCancel = () => {
  alertOpen.value = false
  emit('action')
}

const statusDetails = computed(() => {
  switch (props.booking.status) {
    case ReservationStatus.RESERVATION_STATUS_PENDING:
      return { 
        text: 'PENDING HOST APPROVAL', 
        color: 'primary' as const, 
        customClass: 'badge-blue' 
      }
    case ReservationStatus.RESERVATION_STATUS_CONFIRMED:
      return { 
        text: 'CONFIRMED • UPCOMING', 
        color: 'success' as const, 
        customClass: 'badge-green' 
      }
    case ReservationStatus.RESERVATION_STATUS_CANCELLED:
      return { 
        text: 'REJECTED', 
        color: 'danger' as const, 
        customClass: 'badge-red' 
      }
    default:
      return { 
        text: 'UNKNOWN', 
        color: 'warning' as const, 
        customClass: '' 
      }
  }
})

const actionLabel = computed(() => {
  if (props.booking.status === ReservationStatus.RESERVATION_STATUS_CANCELLED) {
    return null
  }
  return props.booking.status === ReservationStatus.RESERVATION_STATUS_PENDING 
    ? 'Withdraw Request' 
    : 'Cancel Reservation'
})
</script>

<style scoped>
.booking-card {
  position: relative;
  display: flex;
  gap: 1.5rem;
  padding: 1.5rem;
  border: 1px solid #f3f4f6;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 1px 3px rgba(0,0,0,0.02);
}

.image-wrapper {
  flex-shrink: 0;
  width: 250px;
  height: 150px;
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

.top-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

h3 {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 700;
  color: #111827;
}

.location {
  margin: 0.25rem 0 0;
  color: #6b7280;
  font-size: 0.875rem;
}

.bottom-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  padding-top: 1.25rem;
  border-top: 1px solid #f3f4f6; 
}

.dates-grid {
  display: flex;
  gap: 3rem;
}

.date-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.label {
  font-size: 0.65rem;
  font-weight: 700;
  color: #9ca3af;
  text-transform: uppercase;
}

.value {
  font-size: 0.875rem;
  font-weight: 600;
  color: #111827;
}

.price-and-action {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.price-block {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  line-height: 1.3;
}

.nights {
  font-size: 0.8125rem;
  color: #6b7280;
}

.price {
  font-size: 1.125rem;
  font-weight: 800;
  color: #111827;
}

.no-action-text {
  font-size: 0.875rem;
  font-weight: 500;
  color: #9ca3af;
}

:deep(.custom-onyx-badge) {
  font-size: 0.7rem !important;
  padding: 0.15rem 0.6rem !important;
  font-weight: 700 !important;
  border-radius: 6px !important;
}

:deep(.badge-blue) {
  background-color: #eff6ff !important;
  color: #1e40af !important;
  border: none !important; 
}

:deep(.badge-green) {
  background-color: #dcfce7 !important;
  color: #166534 !important;
  border: none !important;
}

:deep(.badge-red) {
  background-color: #fee2e2 !important;
  color: #991b1b !important;
  border: none !important;
}

:deep(.custom-cancel-btn) {
  border-radius: 8px !important;
  padding: 0.4rem 1rem !important;
  height: 35px !important;
  background-color: #fff !important;
  color: #dc2626 !important;
  border: 1px solid #fca5a5 !important;
}

:deep(.custom-cancel-btn:hover) {
  background-color: #fef2f2 !important; 
  border-color: #ef4444 !important; 
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
  background-color: #fff !important;
  color: #dc2626 !important;
}

:deep(.modal-btn-danger:hover) {
  background-color: #fef2f2 !important; 
  border-color: #ef4444 !important; 
}
</style>