<script setup>
    import { ref } from 'vue';
    import { useAuthStore } from '@/stores/auth.js';
    const auth = useAuthStore();

    const pickedNums = ref([]);

    const optionsNumPicks = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 40];
    const optionsNumGames = [1, 2, 3, 4, 5, 10, 20, 50, 100];

    const numBg = (row) => {
        const colors = [
            'bg-[#CA4538]',
            'bg-[#417CBF]',
            'bg-[#337B3D]',
            'bg-[#EF8E3C]',
            'bg-[#E94F83]',
            'bg-[#EB582D]',
            'bg-[#7A8486]',
            'bg-[#7E45C5]',
        ];

        return colors[row % colors.length];
    };

</script>

<template>
    <div class="absolute left-0 top-0 w-screen h-screen  flex">
        <div class="absolute left-0 top-0 w-screen h-screen bg-black opacity-60"></div>
        <div class="relative m-auto p-5 flex flex-col w-fit bg-white opacity-100 rounded-lg">
            <div class="mx-auto flex cursor-pointer rounded-lg mb-2">
                <p class="text-center select-none m-auto text-xl font-semibold">Select your Numbers</p>
            </div>
            <div v-if="optionsNumPicks.includes(pickedNums.length)" class="mx-auto flex cursor-pointer rounded-lg mb-2">
                <p class="text-center select-none m-auto text-sm">Selected: {{pickedNums.length}}</p>
            </div>
            <div v-else class="mx-auto flex cursor-pointer rounded-lg mb-2">
                <p class="text-center select-none m-auto text-sm text-[#CA4538]">Selected: {{pickedNums.length}}</p>
            </div>
            <p class="mx-auto w-2/3 sm:w-3/4 md:w-full  text-slate-400 text-xs md:text-md mb-5">
                <b>NOTE:</b>
                You will need to select 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20 or 40 numbers to play.
            </p>
            <div class="m-auto">
                <div class="flex gap-1 sm:gap-2 mb-2" v-for="y in [0, 1, 2, 3, 4, 5, 6, 7]">
                    <div v-for="x in [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]" >
                        <div v-if="pickedNums.includes((10 * y) + x)"
                            @click="pickedNums.splice(pickedNums.indexOf((10 * y) + x), 1)"
                            class="w-6 sm:w-8 md:w-10 text-center h-8 sm:h-10 md:h-14 flex rounded hover:opacity-90 cursor-pointer"
                            :class="numBg(y)"
                        >
                            <p class="m-auto select-none text-white text-xs md:text-md">{{ (10 * y) + x }}</p>
                        </div>
                        <div v-else 
                            @click="pickedNums.push((10 * y) + x)"
                            class="w-6 sm:w-8 md:w-10 text-center h-8 sm:h-10 md:h-14 flex bg-slate-50 hover:bg-slate-200 cursor-pointer rounded" 
                        >
                            <p class="m-auto select-none text-xs md:text-md">{{ (10 * y) + x }}</p>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Number of Games to play -->
            <div class="flex flex-col w-2/3 mx-auto">
                <p class="mx-auto mt-2 mb-1">No. Games:</p>
                <div class="flex mx-auto mb-2 gap-2 flex-wrap">
                    <div v-for="num in optionsNumGames"
                        class="rounded-full w-8 md:w-10 h-8 md:h-10 flex justify-center items-center bg-slate-50 hover:bg-slate-200 cursor-pointer"
                    >
                    <p class="select-none text-xs md:text-md">
                        {{ num }}
                    </p>
                    </div>
                </div>
            </div>
            
            <!-- Button to Place and do the Callback -->
            <div v-if="optionsNumPicks.includes(pickedNums.length)"
                class="mx-auto w-2/5 flex bg-[#337B3D] hover:opacity-90 cursor-pointer text-white rounded-lg h-10"
            >
                <p class="text-center select-none m-auto">Place Picks</p>
            </div>
            <div v-else
                class="mx-auto w-2/5 flex bg-slate-400 cursor-pointer text-white rounded-lg h-10"
            >
                <p class="text-center select-none m-auto">Place Picks</p>
            </div>
        </div>
    </div>

</template>