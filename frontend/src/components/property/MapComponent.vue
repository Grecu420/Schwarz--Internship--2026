<template>
  <div style="height: 50vh; width: 100%">
    <l-map ref="map" v-model:zoom="zoom" :center="initialCenter">
      <l-tile-layer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        layer-type="base"
        name="OpenStreetMap"
      />
      <l-marker v-model:lat-lng="marker" :draggable="!loading" @moveend="dropMarker">
        <l-tooltip> {{ loading ? "Loading..." : "Drag to set location" }}</l-tooltip>
        <l-icon :icon-url="iconUrl" :icon-size="iconSize" />
      </l-marker>
    </l-map>
  </div>
</template>

<script setup lang="ts">
import { LMap, LTileLayer, LMarker, LIcon, LTooltip } from '@vue-leaflet/vue-leaflet'
import icon from '@/../public/favicon.ico'
import { ref } from 'vue'

interface AddressDetails {
  house_number?: string
  road?: string
  neighbourhood?: string
  suburb?: string
  city?: string
  town?: string
  village?: string
  state?: string
  postcode?: string
  country?: string
  country_code?: string
}

interface NominatimResponse {
  place_id: number
  licence: string
  osm_type: string
  osm_id: number
  lat: string
  lon: string
  display_name: string
  address?: AddressDetails
  boundingbox: [string, string, string, string]
}

export interface LocationPayload {
  location: [number, number]
  address: string
}

interface Props {
  initialZoom?: number
  initialCenter?: [number, number]
}

const props = withDefaults(defineProps<Props>(), {
  initialZoom: 15,
  initialCenter: () => [47.41322, -1.219482],
})

const emit = defineEmits<{
  (e: 'change', payload: LocationPayload): void
}>()

const zoom = ref(props.initialZoom)
const marker = ref({ lat: props.initialCenter[0], lng: props.initialCenter[1] })

const iconWidth = 25
const iconHeight = 40
const iconUrl = icon
const iconSize: [number, number] = [iconWidth, iconHeight]

const address = ref<string>('')
const loading = ref<boolean>(false)
const error = ref<string | null>(null)

async function dropMarker() {
  let lat = marker.value.lat
  let lng = marker.value.lng

  await fetchAddress(lat, lng)

  emit('change', {
    location: [lat, lng],
    address: address.value,
  })
}

async function fetchAddress(lat: number, lon: number): Promise<void> {
  loading.value = true
  error.value = null

  try {
    const url = `https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=${lat}&lon=${lon}`
    const response = await fetch(url, {
      headers: {
        'User-Agent': 'VueLeafletApp/1.0',
      },
    })

    if (!response.ok) {
      throw new Error(`Geocoding failed with status: ${response.status}`)
    }

    const data: NominatimResponse = await response.json()
    address.value = data.display_name || 'No address found'
  } catch (err) {
    if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = 'An unknown error occurred'
    }
  } finally {
    loading.value = false
  }
}
</script>
