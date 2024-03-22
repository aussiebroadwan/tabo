<script setup>
import { useAuthStore } from '@/stores/auth.js';

// Components
import KenoCanvas from '@/components/KenoCanvas.vue'

const auth = useAuthStore();
// auth.checkToken("/");

</script>

<template>
    <div class="w-screen h-screen flex flex-col">
        <KenoCanvas />
    </div>
</template>
