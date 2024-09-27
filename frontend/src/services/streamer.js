import log from '@/services/logger.js';
import { useGameStore } from '@/stores/game';

export default class StreamerService {
    /**
     * Creates an instance of StreamerService and initiates a WebSocket connection.
     * @param {string} url - The WebSocket URL to connect to.
     * @param {number} maxRetries - The maximum number of reconnection attempts. Defaults to 5.
     * @param {number} initialBackoff - The initial delay (in ms) before retrying a connection, used in exponential backoff. Defaults to 1000 ms.
     */
    constructor(url, maxRetries = 5, initialBackoff = 1000) {
        this.url = url;
        this.socket = null;
        this.maxRetries = maxRetries;
        this.retryCount = 0;
        this.initialBackoff = initialBackoff;
        this.backoff = initialBackoff;

        this.gameStore = useGameStore();
        this.connect();
    }

    /**
     * Establishes the WebSocket connection and sets up event handlers for connection events.
     */
    connect() {
        log('INFO', 'Connecting to WebSocket', { url: this.url });
        this.socket = new WebSocket(this.url);
        this.socket.onopen = this.handleOpen.bind(this);
        this.socket.onmessage = this.handleMessage.bind(this);
        this.socket.onerror = this.handleError.bind(this);
        this.socket.onclose = this.handleClose.bind(this);
    }

    /**
     * Handles the WebSocket 'open' event, resetting retry count and backoff when a connection is successfully established.
     * @param {Event} event - The WebSocket open event.
     */
    handleOpen(event) {
        log('INFO', 'WebSocket connection opened', { url: this.url, retryCount: this.retryCount });
        this.retryCount = 0;  // Reset retry count on successful connection
        this.backoff = this.initialBackoff;  // Reset backoff to initial value
    }

    /**
     * Handles the WebSocket 'message' event, processing incoming data and updating the game state in the store.
     * @param {MessageEvent} event - The WebSocket message event containing data from the server.
     */
    handleMessage(event) {
        try {
            const message = JSON.parse(event.data);
            switch (message.type) {
                case 'NEW':
                    this.gameStore.setNewGame(message.body);
                    break;
                case 'PIC':
                    this.gameStore.addPick(message.body.pick);
                    break;
                case 'CUR':
                    this.gameStore.setCurrentGame(message.body);
                    break;
                default:
                    log('ERROR', 'Unknown message type received', { type: message.type });
            } 
        } catch (error) {
            log('ERROR', 'Failed to parse WebSocket message', { data: event.data, error: error.message });
        }
    }

    /**
     * Handles the WebSocket 'error' event by logging the error details.
     * @param {Event} event - The WebSocket error event.
     */
    handleError(event) {
        log('ERROR', 'WebSocket encountered an error', {
            message: event.message,
            stack: event.error ? event.error.stack : null
        });
    }

    /**
     * Handles the WebSocket 'close' event and initiates reconnection attempts with exponential backoff.
     * @param {CloseEvent} event - The WebSocket close event.
     */
    handleClose(event) {
        log('WARN', 'WebSocket connection closed', { url: this.url, reason: event.reason });
        this.attemptReconnect();
    }

    /**
     * Attempts to reconnect to the WebSocket server using exponential backoff. Stops after maxRetries.
     */
    attemptReconnect() {
        if (this.retryCount < this.maxRetries) {
            this.retryCount++;
            const reconnectDelay = this.backoff * Math.pow(2, this.retryCount);  // Exponential backoff
            log('WARN', 'Attempting to reconnect', { retryCount: this.retryCount, delay: reconnectDelay });

            setTimeout(() => {
                this.connect();
            }, reconnectDelay);
        } else {
            log('ERROR', 'Max reconnect attempts reached. Giving up.', { retryCount: this.retryCount });
        }
    }

    /**
     * Sends a message to the WebSocket server if the connection is open.
     * @param {Object} message - The message to send, which will be stringified before sending.
     */
    sendMessage(message) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(message));
        } else {
            log('WARN', 'Cannot send message. WebSocket is not open.', { message });
        }
    }

    /**
     * Closes the WebSocket connection manually and logs the closure.
     */
    close() {
        if (this.socket) {
            this.socket.close();
            log('INFO', 'WebSocket connection closed manually.', { url: this.url });
        }
    }
}