<template>
  <OnyxCard class="property-card">
    <div class="card-image-wrapper" @click="emit('view', property.id)">
      <OnyxImage
        :height="256"
        :width="600"
        :src="imageUrl"
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

        <OnyxButton v-if="isEdit"
          mode="default"
          label="Edit"
          :icon="iconEdit"
          class="edit-btn"
          @click="emit('edit', property.id)"
        />
      </div>
    </div>
  </OnyxCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { OnyxCard, OnyxImage, OnyxButton } from 'sit-onyx'
import { iconEdit } from '@sit-onyx/icons'
import { Property } from '@/generated/proto/property-api'

const DEFAULT_IMAGE = 'https://images.unsplash.com/photo-1518780664697-55e3ad937233?q=80&w=1000'

const props = defineProps<{
  property: Property
  isLoading: boolean
  isEdit: boolean
}>()

const emit = defineEmits<{
  (e: 'view', id: number): void
  (e: 'edit', id: number): void
}>()

const imageUrl = computed(() => props.property.imageUrls[0] || DEFAULT_IMAGE)
</script>

<style scoped>
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
</style>