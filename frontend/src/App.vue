<script setup>
import StreamerService from '@/services/streamer.js'



// const streamerService = new StreamerService('wss://tabo.tabdiscord.com/api/v1/ws');
const streamerService = new StreamerService(import.meta.env.VITE_TABO_BACKEND_WEBSOCKET);

</script>

<template>
  <router-view></router-view>
</template>
