<script setup>
import { useAuthStore } from '@/stores/auth.js';
import router from '@/router.js';


const auth = useAuthStore();

const hash = window.location.hash.substring(1);
const params = new URLSearchParams(hash);
const accessToken = params.get('access_token');
const expiresIn = params.get('expires_in');

if (accessToken && expiresIn) {
    auth.setToken(accessToken, Number(expiresIn));
    router.push("/");
}

</script>

<template>
    <div>Redirecting...</div>
</template>