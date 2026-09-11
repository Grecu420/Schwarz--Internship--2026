<template>
  <div class="profile-page-wrapper">
    <div class="header-container">
      <span class="category-tag">ACCOUNT SETTINGS</span>
      <h1 class="page-title">Edit Profile</h1>
      <p class="page-description">
        Keep your personal details up to date to ensure seamless bookings and smooth communication with hosts.
      </p>
    </div>

    <ProfileForm />
  </div>
</template>

<script setup lang="ts">
import ProfileForm from '@/components/ProfileForm.vue'
</script>

<style scoped>
.profile-page-wrapper {
  min-height: calc(100vh - 80px);
  background-color: #fafafa;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem 1.5rem 5rem 1.5rem;
}

.header-container {
  text-align: center;
  max-width: 520px;
  margin-bottom: 2rem;
}

.category-tag {
  font-size: 0.75rem;
  font-weight: 800;
  color: #2563eb;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.page-title {
  font-size: 2rem;
  font-weight: 800;
  color: #111827;
  margin: 0.25rem 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.page-description {
  font-size: 0.875rem;
  color: #6b7280;
  line-height: 1.5;
  margin: 0;
}
</style>