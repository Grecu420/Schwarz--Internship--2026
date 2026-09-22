<template>
  <div class="properties-page">
    <main class="content-section">
      <!-- Loading State -->
      <div v-if="isLoading" class="state-container">
        <p>Loading properties...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="properties.length === 0" class="state-container">
        <h3>No properties found</h3>
        <p>You haven't listed any properties yet.</p>
      </div>

      <!-- Property Grid -->
      <div v-else class="properties-container">
        <div class="properties-grid">
          <OnyxCard v-for="property in properties" :key="property.id" class="property-card">
            <div class="card-image-wrapper">
              <img :src="getPropertyImage(property)" :alt="property.name" class="property-image" />
            </div>

            <div class="card-body">
              <div class="card-info">
                <h3 class="property-name">{{ property.name }}</h3>
                <p class="property-address">{{ property.address }}</p>
              </div>

              <div class="card-footer">
                <div class="property-price">
                  <span class="price-amount">${{ property.price }}</span>
                  <span class="price-unit"> / night</span>
                </div>

                <OnyxButton
                  mode="default"
                  label="Edit"
                  :icon="iconEdit"
                  class="edit-btn"
                  @click="handleEditProperty(property.id)"
                />
              </div>
            </div>
          </OnyxCard>
        </div>

        <!-- Pagination Controls -->
        <div class="pagination-controls">
          <OnyxIconButton
            mode="outline"
            label="Previous Page"
            :icon="iconChevronLeft"
            :disabled="currentPage === 0 || isLoading"
            class="pagination-btn"
            @click="goToPreviousPage"
          />

          <span class="page-indicator">Page {{ currentPage + 1 }}</span>

          <OnyxIconButton
            mode="outline"
            label="Next Page"
            :icon="iconChevronRight"
            :disabled="!nextPageToken || isLoading"
            class="pagination-btn"
            @click="goToNextPage"
          />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { OnyxButton, OnyxCard, OnyxIconButton, useToast } from 'sit-onyx'
import { iconEdit, iconChevronLeft, iconChevronRight } from '@sit-onyx/icons'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import {
  ListPropertiesRequest,
  ListPropertiesResponse,
  Property,
} from '@/generated/proto/property-api'

const DEFAULT_IMAGE = 'https://images.unsplash.com/photo-1518780664697-55e3ad937233?q=80&w=1000'
const PAGE_SIZE = 6

const router = useRouter()
const toast = useToast()
const authStore = useAuthStore()

const properties = ref<Property[]>([])
const isLoading = ref<boolean>(false)

// Pagination State
const currentPage = ref(0)
const nextPageToken = ref<string>('')
const pageTokens = ref<string[]>([''])

const getPropertyImage = (property: Property): string => {
  return property.imageUrls[0] || DEFAULT_IMAGE
}

const fetchProperties = async (pageTokenToRequest: string) => {
  if (!authStore.user?.id) {
    toast.show({
      headline: 'Authentication Error',
      description: 'You must be logged in to view your properties.',
      color: 'danger',
    })
    return
  }

  isLoading.value = true

  try {
    const request = ListPropertiesRequest.create({
      pageSize: PAGE_SIZE,
      nextPageToken: pageTokenToRequest,
      filters: [
        {owner: {value: authStore.user.id}}
      ]
    })

    const jsonBody = ListPropertiesRequest.toJSON(request)
    const response = await api.post<ListPropertiesResponse>('/api/propertyList', jsonBody)

    if (response.data) {
      properties.value = response.data.properties || []
      nextPageToken.value = response.data.nextPageToken || ''
    }
  } catch (error) {
    console.error('Failed to fetch properties:', error)
    toast.show({
      headline: 'Error',
      description: 'Unable to load properties. Please try again.',
      color: 'danger',
    })
  } finally {
    isLoading.value = false
  }
}

const goToNextPage = async () => {
  if (!nextPageToken.value) return

  const tokenToUse = nextPageToken.value
  currentPage.value++

  if (pageTokens.value.length <= currentPage.value) {
    pageTokens.value.push(tokenToUse)
  } else {
    pageTokens.value[currentPage.value] = tokenToUse
  }

  await fetchProperties(tokenToUse)
}

const goToPreviousPage = async () => {
  if (currentPage.value === 0) return

  currentPage.value--
  const tokenToUse = pageTokens.value[currentPage.value]
  if (!tokenToUse) {
    currentPage.value = 0
    return
  }
  await fetchProperties(tokenToUse)
}

const handleEditProperty = (propertyId: string | number) => {
  router.push(`/properties/edit/${propertyId}`)
}

onMounted(() => {
  fetchProperties('')
})
</script>

<style scoped>
.properties-page {
  min-height: 100vh;
  background-color: #fafafa;
  color: #111827;
  font-family: var(--onyx-font-family-base, system-ui, sans-serif);
}

.content-section {
  max-width: 1200px;
  margin: 0 auto;
  padding: 3rem 1.5rem 5rem 1.5rem;
}

.state-container {
  text-align: center;
  padding: 4rem 2rem;
  color: #6b7280;
}

.state-container h3 {
  color: #111827;
  margin-bottom: 0.5rem;
}

.properties-container {
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

/* Property Cards Grid */
.properties-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.75rem;
}

:deep(.property-card) {
  padding: 0 !important;
  border-radius: 20px !important;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  border: 1px solid #e5e7eb !important;
  background-color: #ffffff !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03) !important;
}

.card-image-wrapper {
  width: 100%;
  height: 220px;
  overflow: hidden;
  background-color: #f3f4f6;
}

.property-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-body {
  padding: 1.25rem 1.5rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  flex: 1;
}

.card-info {
  margin-bottom: 1.25rem;
}

.property-name {
  font-size: 1.05rem;
  font-weight: 700;
  color: #111827;
  margin: 0 0 0.25rem 0;
  line-height: 1.3;
}

.property-address {
  font-size: 0.85rem;
  color: #6b7280;
  margin: 0;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 0.5rem;
}

.property-price {
  display: flex;
  align-items: baseline;
}

.price-amount {
  font-size: 1.125rem;
  font-weight: 800;
  color: #111827;
}

.price-unit {
  font-size: 0.85rem;
  color: #6b7280;
}

:deep(.edit-btn) {
  border-radius: 10px !important;
  border: 1px solid #dbeafe !important;
  background-color: #eff6ff !important;
  color: #2563eb !important;
  font-weight: 600 !important;
  padding: 0.4rem 0.9rem !important;
  font-size: 0.85rem !important;
}

:deep(.edit-btn:hover) {
  background-color: #dbeafe !important;
  border-color: #bfdbfe !important;
}

/* Pagination Controls */
.pagination-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1.5rem;
}

.page-indicator {
  color: #374151;
  font-size: 0.875rem;
  font-weight: 600;
}

:deep(.pagination-btn) {
  border-radius: 50% !important;
  width: 40px !important;
  height: 40px !important;
  padding: 0 !important;
  display: flex !important;
  justify-content: center !important;
  align-items: center !important;
  border: 1px solid #e5e7eb !important;
  background-color: #ffffff !important;
  color: #374151 !important;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.03) !important;
}

:deep(.pagination-btn:not(:disabled):hover) {
  background-color: #f9fafb !important;
  border-color: #d1d5db !important;
  color: #111827 !important;
}

:deep(.pagination-btn:disabled) {
  opacity: 0.4;
  cursor: not-allowed;
  background-color: #f3f4f6 !important;
  border-color: #e5e7eb !important;
}

@media (max-width: 1024px) {
  .properties-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .properties-grid {
    grid-template-columns: 1fr;
  }
}
</style>
