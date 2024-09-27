<script setup>
import { useAuthStore } from '@/stores/auth.js';

// Components
import KenoCanvas from '@/components/KenoCanvas.vue'
import KenoCard from '@/components/KenoCard.vue'

const auth = useAuthStore();
auth.checkToken("/");

</script>

<template>
    <div class="w-screen h-screen flex flex-col">
        <div class="m-auto">
            <KenoCanvas />
            <div class="mx-5 flex justify-end">
                <div class="w-1/5 flex bg-[#337B3D] hover:opacity-90 cursor-pointer text-white rounded-lg text-center h-10 select-none font-semibold">
                    <p class="m-auto">
                        PLACE BET
                    </p>
                </div>
            </div>
            <!-- <KenoCard /> -->
        </div>
    </div>
</template>
