<template>
  <footer class="site-footer">
    <div class="footer-content">
      
      <div class="footer-top">
        <div class="footer-column brand-column">
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

.logo-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-text {
  font-size: 1.25rem;
  font-weight: 800;
  color: #111;
  line-height: 1;
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