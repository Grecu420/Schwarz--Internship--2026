<script setup lang="ts">
import { RouterView } from 'vue-router'
import { OnyxToast } from 'sit-onyx';
</script>

<template>
  <RouterView />
  <OnyxToast />
</template>

<style>
:root,
:root,
html[data-onyx-theme] {
  --onyx-color-base-primary-500: #1a41b5 !important;
  --onyx-color-base-primary-600: #143393 !important;
  --onyx-color-base-primary-700: #0f2671 !important;
  --onyx-color-base-primary-100: #e8ecf8 !important;

  --onyx-color-brand-500: #1a41b5 !important;
  --onyx-color-brand-600: #143393 !important;
  --onyx-color-brand-700: #0f2671 !important;

  --onyx-color-component-focus-ring: #1a41b5 !important;
  --onyx-color-form-border-focus: #1a41b5 !important;
}

.onyx-button:hover {
  background-color: var(--onyx-color-base-primary-600) !important;
  border-color: var(--onyx-color-base-primary-600) !important;
}

.onyx-button:active {
  background-color: var(--onyx-color-base-primary-700) !important;
  border-color: var(--onyx-color-base-primary-700) !important;
}

.onyx-toast, 
.onyx-toast-container,
[data-onyx-toast-container] {
  top: auto !important;         
  left: auto !important;        
  bottom: 2rem !important;     
  right: 2rem !important;      
  transform: none !important;   
}
</style>