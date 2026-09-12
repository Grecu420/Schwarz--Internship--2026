<template>
  <div class="page-container">
    <header class="page-header">
      <OnyxButton
        label="Back"
        :icon="iconChevronLeft"
        mode="outline"
        class="back-btn"
        @click="goBack"
      />

      <div class="header-titles">
        <h1 class="page-title">Create New Property</h1>
        <p class="page-subtitle">Fill in the details below to publish a new property listing.</p>
      </div>
    </header>

    <main class="content-card">
      <PropertyForm @submit="handleSuccess" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { OnyxButton } from 'sit-onyx'
import { iconChevronLeft } from '@sit-onyx/icons'
import PropertyForm from '@/components/PropertyForm.vue'
import type { Property } from '@/generated/proto/property-api'

const router = useRouter()

const goBack = () => {
  router.push('/properties')
}

const handleSuccess = (createdProperty: Property) => {
  // Option 1: Navigate back to main properties list

  console.log('here')
  console.log(createdProperty)

  router.push('/properties')

  // Option 2 (Alternative): Navigate directly to the new property details page
  // if (createdProperty?.id) {
  //   router.push(`/properties/${createdProperty.id}`)
  // }
}
</script>

<style scoped>
.page-container {
  max-width: 720px;
  margin: 2rem auto;
  padding: 0 1rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.page-header {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  align-items: flex-start;
}

.header-titles {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.page-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--onyx-color-neutral-900, #0f172a);
  margin: 0;
}

.page-subtitle {
  font-size: 0.95rem;
  color: var(--onyx-color-neutral-600, #475569);
  margin: 0;
}

.content-card {
  background: var(--onyx-color-surface-default, #ffffff);
  border: 1px solid var(--onyx-color-neutral-200, #e2e8f0);
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

@media (max-width: 640px) {
  .content-card {
    padding: 1.25rem;
  }
}
</style>
