<script setup>
import KenoCanvas from './components/KenoCanvas.vue'
import StreamerService from './services/streamer.js'

const streamerService = new StreamerService('wss://tabo.tabdiscord.com/api/v1/ws');
// const streamerService = new StreamerService('ws://localhost:8080/api/v1/ws');

</script>

<template>
  <div class="w-screen h-screen flex">
    <KenoCanvas />
  </div>
</template>
