<template>
  <div class="property-view-wrapper">
    <div v-if="!isLoading && property" class="main-container">
      <!-- Scrollable Image Gallery -->
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
        <!-- Left Column: Factual Property Data -->
        <div class="details-column">
          <div class="title-header-block">
            <h1 class="property-title">{{ property.name }}</h1>
            <p class="property-address">{{ property.address }}</p>
            <p v-if="property.location" class="property-coords">
              Lat: {{ property.location.lat }}, Long: {{ property.location.long }}
            </p>
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
        <div class="booking-column">
          <ReservationCard :price="property.price" :disabled-days="[]" />
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { GetPropertyResponse, Property } from '@/generated/proto/property-api'
import PropertyOwnerCard from '@/components/PropertyOwnerCard.vue'
import ReservationCard from '@/components/ReservationCard.vue'
import api from '@/utils/api'
import { onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'sit-onyx'
import { GetUserProfileResponse, User, UserProfile } from '@/generated/proto/user-api'

const router = useRouter()
const route = useRoute()
const toast = useToast()

const property = ref<Property | null>(null)
const owner = ref<UserProfile | null>(null)
const isLoading = ref(true)

const goBack = () => {
  router.push('/properties')
}

onMounted(async () => {
  const id = route.params.id
  isLoading.value = true
  try {
    const response = await api.get<GetPropertyResponse>(`/api/property?id=${id}`)
    property.value = response.data.property ?? null
    if (property.value === null) {
      goBack()
    }
    const response2 = await api.get<GetUserProfileResponse>(`/api/user/profile?id=${property.value?.userId}`)
    owner.value = response2.data.user ?? null

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

/* Scrollable Horizontal Image Gallery */
.image-scroll-container {
  display: flex;
  gap: 1rem;
  overflow-x: auto;
  scroll-snap-type: x mandatory;
  padding-bottom: 1rem;
  margin-bottom: 2.5rem;
  scrollbar-width: thin;
  scrollbar-color: #374151 #111827;
}

.image-scroll-container::-webkit-scrollbar {
  height: 8px;
}

.image-scroll-container::-webkit-scrollbar-thumb {
  background: #374151;
  border-radius: 4px;
}

.scroll-image {
  flex: 0 0 auto;
  width: 70%;
  max-width: 640px;
  height: 380px;
  object-fit: cover;
  border-radius: 16px;
  scroll-snap-align: start;
}

/* Grid Layout */
.property-details-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 3.5rem;
  align-items: start;
}

@media (max-width: 900px) {
  .property-details-grid {
    grid-template-columns: 1fr;
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

.property-coords {
  font-size: 0.8125rem;
  color: #6b7280;
  margin-top: 0.25rem;
}

.divider {
  border: none;
  border-top: 1px solid #1f2937;
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
