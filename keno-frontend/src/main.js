import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { DiscordSDK } from "@discord/embedded-app-sdk";

import './style.css'

import App from './App.vue'
import router from './router';

const pinia = createPinia()
const app = createApp(App)

const windowLoc = window.location;

// Handle if running in Discord, to start the SDK.
if (windowLoc.hostname.includes("discordsays")) {
   const discordSdk = new DiscordSDK(import.meta.env.VITE_DISCORD_CLIENT_ID);
   
   async function setupDiscordSdk() {
      await discordSdk.ready();
   }
   
   setupDiscordSdk().then(() => {
      console.log("Discord SDK is ready");
   });
}

app.use(pinia)
   .use(router)
   .mount('#app')

