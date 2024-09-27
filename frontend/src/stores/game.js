import { defineStore } from 'pinia';

export const useGameStore = defineStore({
    // unique id of the store across your application
    id: 'game',
  
    state: () => ({
      gameId: null,
      nextGameTime: null,
      currentGameStartTime: null,
      currentGameEndTime: null,
      picks: new Set([]),
    }),
  
    actions: {
      setNewGame(game) {
        this.gameId = game.gameId;
        this.nextGameTime = game.nextGameTime;
        this.currentGameStartTime = game.currentGameStartTime;
        this.currentGameEndTime = game.currentGameEndTime;
        this.picks = new Set([]);
      },
  
      addPick(pick) {
        this.picks.add(pick);
      },
  
      setCurrentGame(game) {
        this.gameId = game.gameId;
        this.nextGameTime = game.nextGameTime;
        this.currentGameStartTime = game.currentGameStartTime;
        this.currentGameEndTime = game.currentGameEndTime;
        this.picks = new Set(game.picks);
      },
    },
  });