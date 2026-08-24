<template>
    <div class="page-container">
        <p>If you see this, you have a valid token</p>

        <button @click="handleLogout">Logout</button>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';



const router = useRouter()
const handleLogout = (): void => {

    localStorage.removeItem("jwt_token")
    router.push({ name: 'login' });
};

</script>

<style scoped>
@import '../styles/global.css';
</style>