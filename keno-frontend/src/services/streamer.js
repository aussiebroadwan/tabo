import { useGameStore } from '@/stores/game';

export default class StreamerService {
    constructor(url) {
        this.socket = new WebSocket(url);
        this.socket.onmessage = this.handleMessage.bind(this);
        this.socket.onerror= function (evt) {console.log("onerror", evt);}
        this.socket.onclose= function (evt) {console.log("onclose",evt);}
        this.gameStore = useGameStore();
    }

    handleMessage(event) {
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
                console.error('Unknown message type:', message.type);
        }
    }
}