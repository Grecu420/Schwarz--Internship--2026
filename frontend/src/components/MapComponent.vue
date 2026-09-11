<template>
  <div style="height: 600px; width: 800px">
    <l-map ref="map" v-model:zoom="zoom" :center="[47.41322, -1.219482]">
      <l-tile-layer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        layer-type="base"
        name="OpenStreetMap"
      ></l-tile-layer>
      <l-control-layers />
      <l-marker :lat-lng="marker" draggable @moveend="console.log('moveend')">
        <l-tooltip> lol </l-tooltip>
      </l-marker>

      <l-marker :lat-lng="[47.41322, -1.219482]">
        <l-icon :icon-url="iconUrl" :icon-size="iconSize" />
      </l-marker>

      <l-marker :lat-lng="[50, 50]" draggable @moveend="console.log('moveend')">
        <l-popup> lol </l-popup>
      </l-marker>
    </l-map>
  </div>
</template>

<script setup lang="ts">
import {
  LMap,
  LTileLayer,
  LMarker,
  LControlLayers,
  LIcon,
  LPopup,
  LTooltip,
} from '@vue-leaflet/vue-leaflet'
import { ref } from 'vue'

// export default {
//   components: {
//     LMap,
//     LTileLayer,
//   },
//   data() {
//     return {
//       zoom: 2,
//     }
//   },
// }

const zoom = ref(2)
const iconWidth = 25
const iconHeight = 40
const iconUrl = `https://placekitten.com/${iconWidth}/${iconHeight}`
const iconSize = [iconWidth, iconHeight]

const marker = ref([0, 0])
</script>

<style scoped></style>
