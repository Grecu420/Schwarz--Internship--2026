<template>
  <div class="properties-page">
    <header
      class="header-section"
      :style="{
        backgroundImage: `linear-gradient(rgba(0, 0, 0, 0.45), rgba(0, 0, 0, 0.45)), url(${SPLASH_URL})`,
      }"
    >
      <div class="header-content">
        <div class="header-main">
          <div class="header-text">
            <h1 class="page-title">Find your next stay</h1>
            <p class="page-description">
              Discover charming lofts, forest cabins, and seaside villas hand-picked for comfort
            </p>
          </div>
        </div>

        <PropertySearchBar
          :price="[0, 10000]"
          :cities="cities"
          :sort-options="sortOptions"
          v-model="filters"
          @submit="handleSubmit"
          class="search-bar"
        />
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
      </div>

      <!-- Property Grid -->
      <div v-else class="properties-container">
        <h3 v-if="propertyCount > 0">Found {{ propertyCount }} results</h3>
        <div class="properties-grid">
          <PropertyCard
            v-for="property in properties"
            :key="property.id"
            :is-edit="false"
            :is-loading="isLoading"
            :property="property"
            @view="handleViewProperty"
          />
        </div>

        <!-- Pagination Controls -->
        <PaginationControls
          :current-page="currentPage"
          :is-next-disabled="!nextPageToken"
          :is-prev-disabled="currentPage === 0"
          :is-loading="isLoading"
          @next-page="goToNextPage"
          @previous-page="goToPreviousPage"
        />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'sit-onyx'
import PropertySearchBar from '@/components/property/PropertySearchBar.vue'
import { type FilterState } from '@/components/property/PropertySearchBar.vue'
import PropertyCard from '@/components/property/PropertyCard.vue'
import PaginationControls from '@/components/property/PaginationControls.vue'
import api from '@/utils/api'
import { getCitiesByCountry, type City } from '@/utils/cities'
import { Property } from '@/generated/proto/property-api'
import {
  CountPropertiesRequest,
  CountPropertiesResponse,
  ListPropertiesRequest,
  ListPropertiesResponse,
  SortType,
  type ListPropertiesFiltersOneOf,
} from '@/generated/proto/property-list-api'
import { outputFormatter } from '@/utils/dates'

const SPLASH_URL =
  'https://images.unsplash.com/photo-1470770841072-f978cf4d019e?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D'
const PAGE_SIZE = 6

const router = useRouter()
const toast = useToast()

const properties = ref<Property[]>([])
const isLoading = ref<boolean>(false)
const propertyCount = ref(0)

// Pagination State
const currentPage = ref(0)
const nextPageToken = ref<string>('')
const pageTokens = ref<string[]>([''])

const filters = ref<FilterState>({})
const cities = ref<City[]>([])

const sortOptions = [
  { value: 'name', label: 'Name' },
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
]

const getFilterList = () => {
  const reqFilters: ListPropertiesFiltersOneOf[] = []
  const filterState = filters.value
  // Name filter
  if (filterState.search) {
    const name = filterState.search
    reqFilters.push({ name: { value: name } })
  }

  // Location filter
  if (filterState.location) {
    const city = cities.value.find((city) => city.name === filterState.location)
    const lat = city?.latitude
    const long = city?.longitude
    const radius = filterState.radius
    if (!!lat && !!long && !!radius) {
      reqFilters.push({ location: { center: { lat, long }, radius } })
    }
  }

  if (filterState.price) {
    const [min, max] = filterState.price

    reqFilters.push({ priceRange: { min, max } })
  }

  if (filterState.dates) {
    const checkIn = filterState.dates.start
    const checkOut = filterState.dates.end
    if (checkIn && checkOut) {
      const checkInDate = outputFormatter.format(checkIn)
      const checkOutDate = outputFormatter.format(checkOut)
      reqFilters.push({ dateInterval: { checkInDate, checkOutDate } })
    }
  }

  return reqFilters
}

const fetchProperties = async (pageTokenToRequest: string) => {
  isLoading.value = true

  try {
    let sortType = SortType.SORT_TYPE_UNSPECIFIED

    switch (filters.value.sort) {
      case 'name':
        sortType = SortType.SORT_TYPE_NAME
        break
      case 'price_asc':
        sortType = SortType.SORT_TYPE_PRICE_ASC
        break
      case 'price_desc':
        sortType = SortType.SORT_TYPE_PRICE_DESC
        break
      default:
        sortType = SortType.SORT_TYPE_UNSPECIFIED
    }
    const request = ListPropertiesRequest.create({
      pageSize: PAGE_SIZE,
      nextPageToken: pageTokenToRequest,
      filters: getFilterList(),
      sortType: sortType,
    })

    const response = await api.post<ListPropertiesResponse>('/api/propertyList', request)
    console.log(response)
    if (response.data) {
      properties.value = response.data.properties || []
      nextPageToken.value = response.data.nextPageToken || ''
    }
    console.log(properties)
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

const handleViewProperty = (propertyId: number) => {
  router.push(`/properties/view/${propertyId}`)
}

const fetchPropertyCount = async () => {
  const request = CountPropertiesRequest.create({
    filters: getFilterList(),
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

const fetchCities = async () => {
  cities.value = await getCitiesByCountry('RO')
  console.log(cities)
}

const resetPage = () => {
  currentPage.value = 0
  nextPageToken.value = ''
  pageTokens.value = ['']

  fetchPropertyCount()
  fetchProperties('')
}

const handleSubmit = (newFilters: FilterState) => {
  console.log(newFilters)
  filters.value = newFilters

  resetPage()
}

onMounted(() => {
  // fetchPropertyCount()
  fetchProperties('')
  fetchCities()
})
</script>

<style scoped>
.properties-page {
  min-height: 100vh;
  background-color: #fafafa;
  color: #111827;
  font-family: var(--onyx-font-family-base, system-ui, sans-serif);
}

.header-section {
  background-size: cover;
  background-position: top center;
  background-repeat: no-repeat;
  background-attachment: fixed;
  padding: 4rem 1.5rem 3rem 1.5rem;
}

.header-content {
  max-width: 600px;
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
  color: #ffffff;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.page-description {
  font-size: 0.9375rem;
  color: #f3f4f6;
  margin: 0;
}

.search-bar {
  max-width: 600px;
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
  gap: 1.5rem;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
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
</style>
