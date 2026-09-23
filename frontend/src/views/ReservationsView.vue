<template>
  <div class="reservations-page">
    <div class="main-container">
      <header class="page-header">
        <div class="header-top-row">
          <h1 class="page-title">Reservations</h1>

          <div v-if="activeTab === 'requests'" class="filter-dropdown">
            <OnyxSelect
              v-model="selectedFilter"
              :options="filterOptions"
              label="Filter requests"
              list-label="Filter options"
              hide-label
              density="compact"
              :with-clear="false"
              class="custom-onyx-select"
            />
          </div>
        </div>

        <p class="page-subtitle">Manage your bookings and guest requests in one place.</p>
      </header>

      <div class="tabs-and-filters">
        <div class="tabs-container">
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'bookings' }"
            @click="activeTab = 'bookings'"
          >
            My Bookings
          </button>
          <button
            class="tab-btn"
            :class="{ active: activeTab === 'requests' }"
            @click="activeTab = 'requests'"
          >
            Property Requests
          </button>
        </div>
      </div>

      <div class="content-section">
        <div v-if="activeTab === 'bookings'" class="cards-list">
          <GuestBookingCard
            v-for="booking in guestBookings"
            :key="booking.id"
            :booking="booking"
            @action="handleGuestAction(booking.id)"
          />
        </div>

        <div v-else-if="activeTab === 'requests'" class="cards-list">
          <HostRequestCard
            v-for="request in filteredHostRequests"
            :key="request.id"
            :request="request"
            @approve="handleApprove(request.id)"
            @reject="handleReject(request.id)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import GuestBookingCard from '@/components/reservations/GuestBookingCard.vue'
import HostRequestCard from '@/components/reservations/HostRequestCard.vue'
import { useToast, OnyxSelect, type SelectOption } from 'sit-onyx'
import {
  ListReservationsRequest,
  ListReservationsResponse,
  ReservationStatus,
  reservationStatusFromJSON,
  UpdateReservationRequest,
  UpdateReservationResponse,
} from '@/generated/proto/reservation-api'
import type { GetPropertyResponse } from '@/generated/proto/property-api'
import type { GetUserProfileResponse } from '@/generated/proto/user-api'

const auth = useAuthStore()
const toast = useToast()

const activeTab = ref<'bookings' | 'requests'>('bookings')
const isLoading = ref(false)

const guestBookings = ref<any[]>([])
const hostRequests = ref<any[]>([])

const selectedFilter = ref<'pending' | 'all'>('pending')

const pendingCount = computed(() => {
  return hostRequests.value.filter((r) => r.status === ReservationStatus.RESERVATION_STATUS_PENDING)
    .length
})

const filterOptions = computed<SelectOption<'pending' | 'all'>[]>(() => [
  {
    value: 'pending',
    label: `Show: Pending Requests (${pendingCount.value})`,
  },
  {
    value: 'all',
    label: 'Show: All',
  },
])

const filteredHostRequests = computed(() => {
  if (selectedFilter.value === 'pending') {
    return hostRequests.value.filter(
      (r) => r.status === ReservationStatus.RESERVATION_STATUS_PENDING,
    )
  }
  return hostRequests.value
})

const calculateNights = (start: string, end: string) => {
  const diffTime = Math.abs(new Date(end).getTime() - new Date(start).getTime())
  return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
}

const formatDateShort = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

const formatDateLong = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

const fetchGuestBookings = async () => {
  if (!auth.user?.id) return
  isLoading.value = true

  try {
    const request = ListReservationsRequest.create({
      nextPageToken: '',
      pageSize: 50,
      filters: [{ userId: { value: auth.user.id } }],
    })

    const res = await api.post<ListReservationsResponse>('/api/reservationsList', request)
    const reservations = res.data.reservations || []

    const populatedBookings = await Promise.all(
      reservations.map(async (r) => {
        const propRes = await api.get<GetPropertyResponse>(`/api/property?id=${r.propertyId}`)
        const property = propRes.data.property

        const checkIn = r.checkInDate || (r as any).check_in_date
        const checkOut = r.checkOutDate || (r as any).check_out_date
        const nights = calculateNights(checkIn, checkOut)

        const parsedStatus = reservationStatusFromJSON(r.status)

        const images = property?.imageUrls || (property as any)?.image_urls || []
        const firstImage = images.length > 0 ? images[0] : ''

        return {
          id: r.id,
          imageUrl: firstImage,
          propertyName: property?.name || 'Unknown Property',
          location: property?.address || 'Unknown Location',
          checkIn: formatDateLong(checkIn),
          checkOut: formatDateLong(checkOut),
          nights,
          price: (property?.price || 0) * nights,
          status: parsedStatus,
        }
      }),
    )

    guestBookings.value = populatedBookings
  } catch (error: any) {
    toast.show({ headline: 'Error loading bookings', color: 'danger' })
  } finally {
    isLoading.value = false
  }
}

const fetchHostRequests = async () => {
  if (!auth.user?.id) return
  isLoading.value = true

  try {
    const request = ListReservationsRequest.create({
      nextPageToken: '',
      pageSize: 50,
      filters: [{ ownerId: { value: auth.user.id } }],
    })

    const res = await api.post<ListReservationsResponse>('/api/reservationsList', request)
    const reservations = res.data.reservations || []

    const populatedRequests = await Promise.all(
      reservations.map(async (r) => {
        const [propRes, userRes] = await Promise.all([
          api.get<GetPropertyResponse>(`/api/property?id=${r.propertyId}`),
          api.get<GetUserProfileResponse>(`/api/user/profile?id=${r.userId}`),
        ])

        const property = propRes.data.property
        const guest = userRes.data.user

        const checkIn = r.checkInDate || (r as any).check_in_date
        const checkOut = r.checkOutDate || (r as any).check_out_date
        const nights = calculateNights(checkIn, checkOut)

        const parsedStatus = reservationStatusFromJSON(r.status)

        const images = property?.imageUrls || (property as any)?.image_urls || []
        const firstImage = images.length > 0 ? images[0] : ''

        return {
          id: r.id,
          imageUrl: firstImage,
          propertyName: property?.name || 'Unknown Property',
          guest: {
            name: `${guest?.firstName} ${guest?.lastName}`,
            avatar: guest?.profileImageUrl || '',
          },
          checkIn: formatDateShort(checkIn),
          checkOut: formatDateShort(checkOut),
          payout: (property?.price || 0) * nights,
          timeStatus: `Created on ${formatDateShort((r.createdAt as unknown as string) || (r as any).created_at)}`,
          status: parsedStatus,
        }
      }),
    )

    hostRequests.value = populatedRequests
  } catch (error: any) {
    toast.show({ headline: 'Error loading requests', color: 'danger' })
  } finally {
    isLoading.value = false
  }
}

watch(activeTab, (newTab) => {
  if (newTab === 'bookings') {
    fetchGuestBookings()
  } else if (newTab === 'requests') {
    fetchHostRequests()
  }
})

onMounted(() => {
  fetchGuestBookings()
})

const handleGuestAction = async (id: number) => {
  try {
    await api.delete(`/api/reservation?id=${id}`)

    toast.show({
      headline: 'Reservation cancelled successfully',
      color: 'success',
    })

    await fetchGuestBookings()
  } catch (error: any) {
    console.error('Error cancelling reservation:', error)
    toast.show({
      headline: 'Failed to cancel reservation',
      color: 'danger',
    })
  }
}

const handleApprove = async (id: number) => {
  try {
    const request = UpdateReservationRequest.create({
      reservation: {
        id: id,
        status: ReservationStatus.RESERVATION_STATUS_CONFIRMED,
      },
      fieldMask: ['status'],
    })

    const jsonBody = UpdateReservationRequest.toJSON(request)

    const response = await api.patch<UpdateReservationResponse>('/api/reservation', jsonBody)

    if (response.data.reservation) {
      toast.show({
        headline: 'Reservation Approved',
        description: 'The booking request has been confirmed.',
        color: 'success',
      })
    }

    await fetchHostRequests()
  } catch (error: any) {
    console.error('Error approving reservation:', error)
    toast.show({
      headline: 'Error',
      description: 'Failed to approve reservation.',
      color: 'danger',
    })
  }
}

const handleReject = async (id: number) => {
  try {
    const request = UpdateReservationRequest.create({
      reservation: {
        id: id,
        status: ReservationStatus.RESERVATION_STATUS_CANCELLED,
      },
      fieldMask: ['status'],
    })

    const jsonBody = UpdateReservationRequest.toJSON(request)

    await api.patch<UpdateReservationResponse>('/api/reservation', jsonBody)

    toast.show({
      headline: 'Reservation Rejected',
      description: 'The booking request has been rejected.',
      color: 'neutral',
    })

    await fetchHostRequests()
  } catch (error: any) {
    console.error('Error rejecting reservation:', error)
    toast.show({
      headline: 'Error',
      description: 'Failed to reject reservation.',
      color: 'danger',
    })
  }
}
</script>

<style scoped>
.reservations-page {
  padding: 2.5rem 1.5rem;
  background: #ffffff;
  min-height: 100vh;
}

.main-container {
  max-width: 1140px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.header-top-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 2rem;
  font-weight: 800;
  color: #111827;
  margin: 0;
}

.page-subtitle {
  color: #6b7280;
  margin: 0.5rem 0 0;
}

.tabs-and-filters {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 2rem;
}

.tabs-container {
  display: flex;
  gap: 2rem;
}

.tab-btn {
  background: none;
  border: none;
  padding: 0.75rem 0;
  font-size: 1rem;
  color: #6b7280;
  cursor: pointer;
  font-weight: 500;
  border-bottom: 2px solid transparent;
}

.tab-btn.active {
  color: #2563eb;
  font-weight: 700;
  border-bottom-color: #2563eb;
}

.filter-dropdown {
  width: 230px;
}

:deep(.custom-onyx-select) {
  width: 100% !important;
}

:deep(.custom-onyx-select .onyx-form-element-v2__wrapper) {
  background-color: #ffffff !important;
  border: 1px solid #e5e7eb !important;
  border-radius: 10px !important;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03) !important;
  padding: 0.2rem 0.5rem !important;
  transition: all 0.15s ease !important;
  height: 36px !important;
  min-height: 36px !important;
  display: flex !important;
  align-items: center !important;
}

:deep(.custom-onyx-select:hover .onyx-form-element-v2__wrapper) {
  border-color: #d1d5db !important;
  background-color: #f9fafb !important;
}

:deep(.custom-onyx-select input),
:deep(.custom-onyx-select .onyx-select__value) {
  font-size: 0.75rem !important; 
  font-weight: 600 !important;
  color: #374151 !important;
  cursor: pointer !important;
}

:deep(.custom-onyx-select .onyx-icon) {
  color: #6b7280 !important;
  width: 14px !important;
  height: 14px !important;
}

:deep(.custom-onyx-select .onyx-icon),
:deep(.custom-onyx-select button) {
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  color: #6b7280 !important;
  width: 14px !important;
  height: 14px !important;
  margin: 0 !important;
}

.cards-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
</style>

<style>
.onyx-options-menu,
.onyx-listbox,
.onyx-select-dropdown {
  border-radius: 10px !important;
  border: 1px solid #e5e7eb !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08) !important;
  padding: 0.25rem !important;
  background-color: #ffffff !important;
}

.onyx-options-menu-item,
.onyx-listbox-option,
.onyx-select-option {
  border-radius: 6px !important;
  font-size: 0.75rem !important; 
  font-weight: 500 !important;
  color: #374151 !important;
  padding: 0.4rem 0.6rem !important;
}

.onyx-options-menu-item:hover,
.onyx-listbox-option:hover,
.onyx-listbox-option--active,
.onyx-select-option:hover {
  background-color: #f3f4f6 !important;
  color: #111827 !important;
}

.onyx-options-menu-item--selected,
.onyx-listbox-option--selected,
.onyx-select-option[aria-selected="true"],
[aria-selected="true"].onyx-listbox-option {
  background-color: #eff6ff !important;
  color: #1d4ed8 !important;
  font-weight: 700 !important;
}
</style>