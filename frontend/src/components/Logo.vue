<template>
  <div class="logo-container">
    <svg width="32" height="28" viewBox="0 0 32 28" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="32" height="28" rx="6" fill="#1E40AF"/>
      <g clip-path="url(#clip0_3_13)">
        <path d="M17.75 19.2495V14.5829C17.75 14.4282 17.6885 14.2798 17.5791 14.1704C17.4697 14.061 17.3214 13.9995 17.1667 13.9995H14.8333C14.6786 13.9995 14.5303 14.061 14.4209 14.1704C14.3115 14.2798 14.25 14.4282 14.25 14.5829V19.2495M10.75 12.8331C10.75 12.6634 10.7869 12.4958 10.8584 12.3418C10.9298 12.1879 11.034 12.0514 11.1636 11.9418L15.2469 8.44181C15.4575 8.26384 15.7243 8.1662 16 8.1662C16.2757 8.1662 16.5425 8.26384 16.7531 8.44181L20.8364 11.9418C20.966 12.0514 21.0702 12.1879 21.1416 12.3418C21.2131 12.4958 21.25 12.6634 21.25 12.8331V18.0831C21.25 18.3926 21.1271 18.6893 20.9083 18.9081C20.6895 19.1269 20.3928 19.2498 20.0833 19.2498H11.9167C11.6072 19.2498 11.3105 19.1269 11.0917 18.9081C10.8729 18.6893 10.75 18.3926 10.75 18.0831V12.8331Z" stroke="white" stroke-width="2" stroke-linecap="round"/>
      </g>
      <defs>
        <clipPath id="clip0_3_13">
          <rect width="14" height="14" fill="white" transform="translate(9 7)"/>
        </clipPath>
      </defs>
    </svg>
    <span class="logo-text">StayNest</span>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.logo-container {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;  
  transition: opacity 0.2s ease;
}

.logo-text {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1e40af; 
}
</style>