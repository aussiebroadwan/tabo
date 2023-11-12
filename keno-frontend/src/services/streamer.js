import { useGameStore } from '@/stores/game';

export default class StreamerService {
    constructor(url) {
        this.socket = new WebSocket(url);
        this.socket.onmessage = this.handleMessage.bind(this);
        this.gameStore = useGameStore();
    }

    handleMessage(event) {
        const message = JSON.parse(event.data);
        console.log("message", message)
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
                console.error('Unknown message type:', message.type);
        }
    }
}