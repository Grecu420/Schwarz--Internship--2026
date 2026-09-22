<template>
  <div class="property-view-wrapper">
    <div v-if="!isLoading && property" class="main-container">
      <!-- Image Gallery -->
      <section v-if="property.imageUrls?.length" class="image-scroll-container">
        <img
          v-for="(url, index) in property.imageUrls"
          :key="index"
          :src="url"
          :alt="`${property.name} photo ${index + 1}`"
          class="scroll-image"
        />
      </section>

      <!-- Property Details & Reservation Grid -->
      <section class="property-details-grid">
        <!-- Left Column: Property Data -->
        <div class="details-column">
          <div class="title-header-block">
            <h1 class="property-title">{{ property.name }}</h1>
            <p class="property-address">{{ property.address }}</p>
          </div>

          <hr class="divider" />

          <div class="about-section">
            <h2 class="section-heading">About this space</h2>
            <p class="description-text">{{ property.description }}</p>
          </div>

          <hr class="divider" />

          <!-- Owner Information Component -->
          <PropertyOwnerCard v-if="owner" :owner="owner" />
        </div>

        <!-- Right Column: Reservation Component -->
        <div v-if="isBookable" class="booking-column">
          <ReservationCard
            :price="property.price"
            :existing-reservations="existingReservations"
            @submit="handleReservation"
          />
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { GetPropertyResponse, Property } from '@/generated/proto/property-api'
import PropertyOwnerCard from '@/components/property/PropertyOwnerCard.vue'
import ReservationCard from '@/components/property/ReservationCard.vue'
import api from '@/utils/api'
import { onMounted, ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast, OnyxImage } from 'sit-onyx'
import { GetUserProfileResponse, UserProfile } from '@/generated/proto/user-api'
import {
  CreateReservationRequest,
  CreateReservationResponse,
  ListReservationsRequest,
  ListReservationsResponse,
  Reservation,
  ReservationStatus,
} from '@/generated/proto/reservation-api'
import { useAuthStore } from '@/stores/auth'
import type { AxiosResponse } from 'axios'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const auth = useAuthStore()

const property = ref<Property | null>(null)
const owner = ref<UserProfile | null>(null)
const isLoading = ref(true)
const existingReservations = ref<[string, string][]>([])

const isBookable = computed(() => {
  return property.value && owner.value?.id !== auth.user?.id
})

const goBack = () => {
  router.push('/properties')
}

onMounted(async () => {
  const id = route.params.id
  if (!id) {
    goBack()
    return
  }

  isLoading.value = true
  try {
    // Load property
    const response = await api.get<GetPropertyResponse>(`/api/property?id=${id}`)
    property.value = response.data.property ?? null
    if (property.value === null) {
      goBack()
      return
    }

    const ownerPromise = api.get<GetUserProfileResponse>(
      `/api/user/profile?id=${property.value.userId}`,
    )

    let reservationsPromise: Promise<AxiosResponse<ListReservationsResponse>> | null = null
    if (auth.user?.id !== property.value.userId) {
      const reservationRequest = ListReservationsRequest.create({
        nextPageToken: '',
        pageSize: 1000,
        filters: [
          { propertyId: { value: property.value.id } },
          { status: { value: ReservationStatus.RESERVATION_STATUS_CONFIRMED } },
        ],
      })
      reservationsPromise = api.post<ListReservationsResponse>(
        '/api/reservationsList',
        reservationRequest,
      )
    }

    const [ownerRes, reservationsRes] = await Promise.all([ownerPromise, reservationsPromise])

    owner.value = ownerRes.data.user ?? null

    if (reservationsRes) {
      existingReservations.value = (reservationsRes.data.reservations ?? []).map(
        (r: Reservation) => [r.checkInDate, r.checkOutDate],
      )
    }
  } catch (error: any) {
    toast.show({
      headline: 'Failed to load property details',
      description: error?.message || 'Failed to send property request.',
      color: 'danger',
    })
    goBack()
  } finally {
    isLoading.value = false
  }
})

const handleReservation = async (reservation: { start: string; end: string }) => {
  if (!property.value) return

  const request = CreateReservationRequest.create({
    reservation: {
      propertyId: property.value.id,
      ownerId: property.value.userId,
      userId: auth.user?.id,
      checkInDate: reservation.start,
      checkOutDate: reservation.end,
    },
  })

  try {
    const response = await api.post<CreateReservationResponse>('/api/reservation', request)
    router.push('/reservations')
    toast.show({
      headline: 'Reservation submitted',
      description: `Start: ${reservation.start}; End: ${reservation.end}`,
      color: 'neutral',
    })
  } catch (error: any) {
    toast.show({
      headline: 'Failed to create reservation',
      description: error?.message || 'Failed to send reservation request.',
      color: 'danger',
    })
  }
}
</script>

<style scoped>
.property-view-wrapper {
  width: 100%;
  min-height: 100vh;
  background-color: #ffffff;
  padding: 2rem 1.5rem 5rem 1.5rem;
}

.main-container {
  max-width: 1140px;
  margin: 0 auto;
}

.image-scroll-container {
  display: flex;
  gap: 1rem;
  overflow-x: scroll;
  scroll-snap-type: x proximity;
  padding-bottom: 1rem;
  margin-bottom: 2.5rem;
  scrollbar-width: thin;
  scrollbar-color: #515b69 transparent;
}

.image-scroll-container::-webkit-scrollbar {
  height: 8px;
}

.image-scroll-container::-webkit-scrollbar-track {
  background-color: #f1f1f1;
  background-clip: padding-box;
  border-left: 2rem solid transparent;
  border-right: 2rem solid transparent;
  border-radius: 8px;
}

.image-scroll-container::-webkit-scrollbar-thumb {
  background-color: #6b7280;
  border-radius: 4px;
}

.scroll-image {
  flex: 0 0 auto;
  width: 70%;
  max-width: 640px;
  height: 380px;
  object-fit: cover;
  border-radius: 16px;
  scroll-snap-align: start end;
}

/* Grid Layout */
.property-details-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 3.5rem;
  align-items: start;
}

.booking-column {
  position: sticky;
  top: 2rem;
}

@media (max-width: 900px) {
  .property-details-grid {
    grid-template-columns: 1fr;
  }

  .booking-column {
    position: static;
  }
}

.property-title {
  font-size: 2.25rem;
  font-weight: 800;
  margin: 0 0 0.5rem 0;
}

.property-address {
  font-size: 1rem;
  color: #9ca3af;
  margin: 0;
}

.divider {
  border: none;
  border-top: 1px solid #e6e5e3;
  margin: 1.75rem 0;
}

.section-heading {
  font-size: 1.25rem;
  font-weight: 700;
  margin: 0 0 0.875rem 0;
}

.description-text {
  font-size: 0.9375rem;
  color: #2c2c2c;
  line-height: 1.65;
  margin: 0;
  white-space: pre-line;
}
</style>
