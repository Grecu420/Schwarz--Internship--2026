<template>
  <footer class="site-footer">
    <div class="footer-content">
      
      <div class="footer-top">
        <div class="footer-column brand-column">
          <Logo/>
          <p class="brand-description">
            Discover a warmer, more connected way to travel. From mountain hideaways to seaside lofts, find your next cozy home with us.
          </p>
        </div>

        <div class="footer-column">
          <h4>SUPPORT</h4>
          <nav>
            <router-link to="#">Help Center</router-link>
            <router-link to="#">Safety Information</router-link>
            <router-link to="#">Cancellation Options</router-link>
          </nav>
        </div>

        <div class="footer-column">
          <h4>COMMUNITY</h4>
          <nav>
            <router-link to="#">StayNest Blog</router-link>
            <router-link to="#">Cozy Forum</router-link>
            <router-link to="#">Host Resources</router-link>
          </nav>
        </div>

        <div class="footer-column">
          <h4>HOSTING</h4>
          <nav>
            <router-link to="#">Try Hosting</router-link>
            <router-link to="#">StayCover for Hosts</router-link>
            <router-link to="#">Hosting Tools</router-link>
          </nav>
        </div>
      </div>

      <div class="footer-bottom">
        <div class="copyright">
          © 2026 StayNest, Inc. All rights reserved.
        </div>
        <div class="legal-links">
          <router-link to="#">Privacy Policy</router-link>
          <router-link to="#">Terms of Service</router-link>
        </div>
      </div>

    </div>
  </footer>
</template>

<script setup lang="ts">
import Logo from '@/components/Logo.vue' 
</script>

<style scoped>
.site-footer {
  background-color: #ffffff;
  border-top: 1px solid #ebebeb;
  padding: 4rem 3rem 1.5rem;
}

.footer-content {
  max-width: 1440px;
  margin: 0 auto;
}

.footer-top {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 2rem;
  margin-bottom: 3rem;
}

.brand-description {
  font-size: 0.875rem;
  line-height: 1.5;
  color: #717171;
  margin-top: 1.5rem;
  max-width: 260px;
  font-weight: 400;
}

.footer-column h4 {
  font-size: 0.75rem;
  font-weight: 700;
  color: #222222;
  margin-top: 0;
  margin-bottom: 1.25rem;
}

.footer-column nav {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.footer-column nav a {
  text-decoration: none;
  color: #717171;
  font-size: 0.875rem;
  font-weight: 400;
  transition: color 0.2s ease;
}

.footer-column nav a:hover {
  color: #222222;
  text-decoration: underline;
}

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 1.5rem;
  border-top: 1px solid #ebebeb;
  font-size: 0.875rem;
  color: #717171;
}

.legal-links {
  display: flex;
  gap: 1.5rem;
}

.legal-links a {
  text-decoration: none;
  color: #717171;
  font-weight: 400;
  transition: color 0.2s ease;
}

.legal-links a:hover {
  color: #222222;
  text-decoration: underline;
}

@media (max-width: 900px) {
  .footer-top {
    grid-template-columns: 1fr;
    gap: 2.5rem;
  }
  .footer-bottom {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
}
</style>