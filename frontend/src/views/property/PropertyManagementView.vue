<template>
  <div class="properties-page">
    <header class="header-section">
      <div class="header-content">
        <div class="header-main">
          <div class="header-text">
            <h1 class="page-title">My Properties</h1>
            <p class="page-description">
              Manage your listed rentals, edit details, and track performance.
            </p>
          </div>
          <OnyxButton
            mode="default"
            label="Add New Property"
            class="add-property-btn"
            @click="handleAddNewProperty"
            :icon="iconPlus"
          />
        </div>

        <div class="stats-grid">
          <div class="stat-card">
            <span class="stat-label">ACTIVE LISTINGS</span>
            <div class="stat-value">{{ propertyCount }}</div>
          </div>
        </div>
      </div>
    </header>

    <main class="content-section">
      <!-- Loading State -->
      <div v-if="isLoading && properties.length === 0" class="state-container">
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
            <div class="card-image-wrapper" @click="handleViewProperty(property.id)">
              <OnyxImage
                :height="256"
                :width="600"
                :src="getPropertyImage(property)"
                :alt="property.name"
                class="property-image"
                :skeleton="isLoading"
              />
            </div>

            <div class="card-body">
              <div class="card-info">
                <h3 class="property-name">{{ property.name }}</h3>
                <p class="property-address">{{ property.address }}</p>
              </div>

              <div class="card-footer">
                <div class="property-price">
                  <span class="price-amount">${{ property.price }} </span>
                  &nbsp;
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
import { OnyxButton, OnyxCard, OnyxIconButton, OnyxImage, useToast } from 'sit-onyx'
import { iconEdit, iconChevronLeft, iconChevronRight, iconPlus } from '@sit-onyx/icons'
import { useAuthStore } from '@/stores/auth'
import api from '@/utils/api'
import {
  CountPropertiesRequest,
  CountPropertiesResponse,
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
const propertyCount = ref(0)

// Pagination State
const currentPage = ref(0)
const nextPageToken = ref<string>('')
const pageTokens = ref<string[]>([''])

const getPropertyImage = (property: Property): string => {
  return property.imageUrls[0] || DEFAULT_IMAGE
}

const fetchProperties = async (pageTokenToRequest: string) => {


  isLoading.value = true

  try {
    const request = ListPropertiesRequest.create({
      pageSize: PAGE_SIZE,
      nextPageToken: pageTokenToRequest,
      filters: [{ owner: { value: authStore.user?.id} }],
    })

    const response = await api.post<ListPropertiesResponse>('/api/propertyList', request)

    if (response.data) {
      properties.value = response.data.properties || []
      nextPageToken.value = response.data.nextPageToken || ''
    }
  } catch (error: any) {
    console.error('Failed to fetch properties:', error)
    toast.show({
      headline: 'Unable to load properties.',
      description: error.message || 'Unable to load properties. Please try again.',
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
  const tokenToUse = pageTokens.value[currentPage.value] || ''
  if (tokenToUse === '') currentPage.value = 0
  await fetchProperties(tokenToUse)
}

const handleAddNewProperty = () => {
  router.push('/properties/create')
}

const handleEditProperty = (propertyId: number) => {
  router.push(`/properties/edit/${propertyId}`)
}

const handleViewProperty = (propertyId: number) => {
  // Adjust the route path if your details page URL is different
  router.push(`/properties/view/${propertyId}`)
}

const fetchPropertyCount = async () => {
  const request = CountPropertiesRequest.create({
    filters: [{ owner: { value: authStore.user?.id } }],
  })
  isLoading.value = true
  try {
    const response = await api.post<CountPropertiesResponse>('/api/propertyCount', request)
    if (response.data) {
      propertyCount.value = response.data.count
    }
  } catch (error: any) {
    console.error('Failed to fetch property count', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchPropertyCount()
  fetchProperties('')
})
</script>

<style scoped>
.header-section {
  background-color: #ffffff;
  padding: 2rem 1.5rem 1.5rem 1.5rem;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

.header-main {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1.5rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: #111827;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.page-description {
  font-size: 0.9375rem;
  color: #888888;
  margin: 0;
}

:deep(.add-property-btn) {
  background-color: #1e40af !important;
  color: #ffffff !important;
  border-radius: 12px !important;
  padding: 0.65rem 1.25rem !important;
  font-weight: 600 !important;
  border: none !important;
}

:deep(.add-property-btn:hover) {
  background-color: #1d4ed8 !important;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
}

.stat-card {
  background-color: #ffffff;
  color: #111827;
  border-radius: 16px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.stat-label {
  font-size: 0.75rem;
  font-weight: 700;
  color: #6b7280;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 800;
  color: #111827;
}

.stat-sub {
  font-size: 1.25rem;
  font-weight: 500;
  color: #6b7280;
}

.properties-page {
  min-height: 100vh;
  background-color: #fafafa;
  color: #111827;
  font-family: var(--onyx-font-family-base, system-ui, sans-serif);
}

.content-section {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem 4rem 1rem;
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

.properties-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
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
  justify-content: center;
  display: flex;
  cursor: pointer;
  background-color: #f3f4f6;
}

.property-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-body {
  padding: 0.9rem 1rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  flex: 1;
}

.card-info {
  margin-bottom: 1rem;
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
  padding-top: 1rem;
  margin-top: 0.5rem;
  border-top: 1px solid #d5d4d1;
}

.property-price {
  display: flex;
  align-items: baseline;
}

.price-amount {
  font-size: 1rem;
  font-weight: 800;
  color: #111827;
}

.price-unit {
  font-size: 1rem;
  color: #6b7280;
}

:deep(.edit-btn) {
  border-radius: 10px !important;
  border: 1px solid #dbeafe !important;
  background-color: #eff6ff !important;
  color: #1e40af !important;
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
