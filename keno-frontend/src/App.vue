<script setup>
import StreamerService from './services/streamer.js'

const windowLoc = window.location

var socketEndpoint = "ws://localhost:8080/api/v1/ws"
// var socketEndpoint = "wss://tabo.tabdiscord.com/api/v1/ws"

// Override the socketEndpoint if the hostname is discordsays
if (windowLoc.hostname.includes("discordsays")) {
  socketEndpoint = "wss://" + windowLoc.hostname + "/api/v1/ws"
} 

// const streamerService = new StreamerService('wss://tabo.tabdiscord.com/api/v1/ws');
const streamerService = new StreamerService(socketEndpoint);

</script>

<template>
  <router-view></router-view>
</template>
