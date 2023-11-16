<template>
    <div class="flex w-screen h-screen">
      <canvas ref="gameCanvas" class="m-auto w-full xl:w-fit"></canvas>
    </div>
  </template>
  
  <script setup>
    import { useGameStore } from '@/stores/game.js';
    import { GameGrid } from '@/models/game.js';
    import { ref, onMounted } from 'vue';

    const gameStore = useGameStore();
    const picks = ref([]);
    const gameId = ref(0);
    const nextGame = ref(0);
    const gameCanvas = ref(null);
    
    const grid = new GameGrid(20, 20);

    gameStore.$subscribe((mutation, state) => {
        picks.value = state.picks;
        nextGame.value = state.nextGameTime;
        gameId.value = state.gameId;
    });

    onMounted(() => {
        const canvas = gameCanvas.value;
        const ctx = canvas.getContext('2d');

        // Set the canvas width and height to occupy full window
        ctx.canvas.width = 760;
        ctx.canvas.height = 380;

        // Start the animation loop
        drawLoop();
    });

    const drawLoop = () => {
        const ctx = gameCanvas.value.getContext('2d');
        grid.update(gameStore.gameId, gameStore.picks, gameStore.nextGameTime);

        // Clear the canvas
        ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
        grid.render(ctx);

        // Request the next frame
        requestAnimationFrame(drawLoop);
    };

</script>