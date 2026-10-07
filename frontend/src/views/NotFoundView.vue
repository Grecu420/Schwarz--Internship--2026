<template>
  <div class="not-found-page">
    <div class="content-wrapper">
      
      <div class="text-section">
        <OnyxBadge 
          color="primary" 
          class="custom-onyx-badge badge-blue"
        >
          <span class="dot"></span>
          ERROR 404
        </OnyxBadge>
        
        <h1 class="title">
          Looks like this<br />
          trail ends here.
        </h1>
        
        <p class="description">
          The cabin you're looking for may have moved, or this path was never mapped. Let's get you back somewhere warm.
        </p>
        
        <div class="actions">
          <OnyxButton 
            label="Back to homepage" 
            :icon="iconHome" 
            mode="default"
            class="home-btn"
            @click="goHome" 
          />
        </div>
      </div>

      <div class="image-section">
        <img 
          src="https://images.unsplash.com/photo-1525113990976-399835c43838?q=80&w=1364&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D" 
          alt="Cabin in the woods" 
          class="hero-image" 
        />
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { OnyxButton, OnyxBadge } from 'sit-onyx'
import { iconHome } from "@sit-onyx/icons";

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.not-found-page {
  min-height: calc(100vh - 200px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  background-color: #ffffff;
  font-family: var(--onyx-font-family-base, system-ui, sans-serif);
}

.content-wrapper {
  max-width: 1200px;
  width: 100%;
  display: flex;
  align-items: center;
  gap: 4rem;
}

.text-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 1.5rem;
}

:deep(.custom-onyx-badge) {
  display: inline-flex !important;
  align-items: center !important;
  gap: 0.5rem !important;
  font-size: 0.875rem !important;
  padding: 0.25rem 0.75rem !important;
  font-weight: 700 !important;
  border-radius: 6px !important;
  letter-spacing: 0.05em !important;
}

:deep(.badge-blue) {
  background-color: #eff6ff !important;
  color: #1a41b5 !important;
  border: none !important; 
}

.dot {
  width: 8px;
  height: 8px;
  background-color: #1a41b5;
  border-radius: 50%;
  display: inline-block;
}

.title {
  font-size: 4rem;
  line-height: 1.1;
  font-weight: 900;
  color: #111827;
  margin: 0;
  letter-spacing: -0.02em;
}

.description {
  font-size: 1.125rem;
  line-height: 1.6;
  color: #6b7280;
  margin: 0;
  max-width: 90%;
}

.actions {
  display: flex;
  align-items: center;
  margin-top: 1rem;
}

:deep(.home-btn) {
  border-radius: 8px !important;
  padding: 0.75rem 1.5rem !important;
  font-weight: 600 !important;
}

.image-section {
  flex: 1.2;
  display: flex;
  justify-content: flex-end;
}

.hero-image {
  width: 100%;
  height: 650px;
  object-fit: cover;
  border-radius: 24px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

@media (max-width: 1024px) {
  .title {
    font-size: 3rem;
  }
  .hero-image {
    height: 500px;
  }
}

@media (max-width: 768px) {
  .content-wrapper {
    flex-direction: column;
    text-align: center;
    gap: 3rem;
  }
  
  .text-section {
    align-items: center;
  }
  
  .description {
    max-width: 100%;
  }
  
  .image-section {
    width: 100%;
    justify-content: center;
  }
  
  .hero-image {
    height: 400px;
  }
}
</style>