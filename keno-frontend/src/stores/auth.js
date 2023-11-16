import { defineStore } from 'pinia';
import axios from 'axios';

// Get VITE variables for the client ID and redirect URI
const CLIENT_ID = import.meta.env.VITE_CLIENT_ID;
const REDIRECT_URI = import.meta.env.VITE_REDIRECT_URI;
const AUTH_LOCAL_STORAGE_KEY = 'keno_auth';

export const useAuthStore = defineStore({
    // unique id of the store across your application
    id: 'auth',
  
    state: () => ({
        // Access Token Data
        accessToken: JSON.parse(localStorage.getItem(`${AUTH_LOCAL_STORAGE_KEY}.accessToken`)),
        expiresAt: JSON.parse(localStorage.getItem(`${AUTH_LOCAL_STORAGE_KEY}.expiresAt`)),

        // User Data
        userId: JSON.parse(localStorage.getItem(`${AUTH_LOCAL_STORAGE_KEY}.userId`)),
        avatar: JSON.parse(localStorage.getItem(`${AUTH_LOCAL_STORAGE_KEY}.avatar`)),
        username: JSON.parse(localStorage.getItem(`${AUTH_LOCAL_STORAGE_KEY}.username`)),
        globalName: JSON.parse(localStorage.getItem(`${AUTH_LOCAL_STORAGE_KEY}.globalName`)),
    }),
  
    actions: {
        // Set the access token and fetch the user data. Used by the callback
        // page after the user logs in eg. /callback#access_token=ACCESS_TOKEN&expires_in=EXPIRES_IN
        setToken(accessToken, expiresIn) {
            // Save the access token and expiration date to state
            this.accessToken = accessToken;
            this.expiresAt = Date.now() + expiresIn * 1000;
            this.stateToLocalStorage();

            // Fetch the user data from the Discord API.
            this.fetchUserData();
        },

        // To keep the auth state persistent across page refreshes, we store
        // the auth state in local storage. This function is called whenever
        // the auth state is updated.
        stateToLocalStorage() {
            localStorage.setItem(`${AUTH_LOCAL_STORAGE_KEY}.accessToken`, JSON.stringify(this.accessToken));
            localStorage.setItem(`${AUTH_LOCAL_STORAGE_KEY}.expiresAt`, JSON.stringify(this.expiresAt));
            localStorage.setItem(`${AUTH_LOCAL_STORAGE_KEY}.userId`, JSON.stringify(this.userId));
            localStorage.setItem(`${AUTH_LOCAL_STORAGE_KEY}.avatar`, JSON.stringify(this.avatar));
            localStorage.setItem(`${AUTH_LOCAL_STORAGE_KEY}.username`, JSON.stringify(this.username));
            localStorage.setItem(`${AUTH_LOCAL_STORAGE_KEY}.globalName`, JSON.stringify(this.globalName));
        },

        // Fetch the user data from the Discord API.
        fetchUserData() {
            axios.get('https://discord.com/api/users/@me', {headers: { Authorization: `Bearer ${this.accessToken}` }})
                .then((response) => {
                    console.log(response);
                    this.userId = response.data.id;
                    this.avatar = response.data.avatar;
                    this.username = response.data.username;
                    this.globalName = response.data.global_name;
                    
                    this.stateToLocalStorage();
                })  
                .catch((error) => {
                    console.error(error);
                }
            );
        },

        // Check if the access token is expired and redirect to the Discord 
        // login page if it is to get a new access token.
        checkToken() {

            // Redirect to the Discord login page if auth data is not available.
            if (Date.now() > this.expiresAt) {                
                const encoded_redirect = encodeURIComponent(REDIRECT_URI);
                window.location.href = `https://discord.com/api/oauth2/authorize?client_id=${CLIENT_ID}&redirect_uri=${encoded_redirect}&response_type=token&scope=identify%20email`;
                return
            }

            // Force a refresh of the user data if it's not available.
            if (this.userId === null) {
                this.fetchUserData();
            }
        },
    },
  });