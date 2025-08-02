<script lang="ts">
    import type { TournamentFormData, Round } from '$lib/types/tournament';
    import { calculateTotalRounds } from '$lib/utils/tournamentUtils';
    
    // Props
    export let formData: TournamentFormData;
    
    // Event dispatcher for form updates
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher<{
        addRound: void;
        removeRound: { index: number };
        updateRound: { index: number; field: string; value: number };
    }>();
    
    // Helper functions to dispatch events
    function handleAddRound() {
        dispatch('addRound');
    }
    
    function handleRemoveRound(index: number) {
        dispatch('removeRound', { index });
    }
    
    function updateRoundField(index: number, field: string, value: number) {
        dispatch('updateRound', { index, field, value });
    }
</script>

<div class="space-y-6">
    <div class="flex justify-between items-center">
        <h3 class="font-medium text-lg">Tournament Rounds</h3>
        <span class="badge badge-primary">{calculateTotalRounds(formData)} Rounds</span>
    </div>

    {#each formData.rounds as round, index}
        <div class="bg-base-100 border border-base-300 rounded-lg p-4">
            <div class="flex justify-between items-center mb-4">
                <h4 class="font-medium">{round.name}</h4>
                {#if index > 0}
                    <button
                        type="button"
                        class="btn btn-sm btn-error btn-outline"
                        onclick={() => handleRemoveRound(index)}
                    >
                        Remove
                    </button>
                {/if}
            </div>

            <!-- Round configuration in a single line -->
            <div class="flex flex-col space-y-4">
                <!-- Group information display -->
                <div class="flex items-center justify-between bg-base-200 p-3 rounded-lg">
                    <div class="flex items-center">
                        <div class="badge badge-primary mr-2">{round.groupCount}</div>
                        <span class="text-sm font-medium">Groups</span>
                    </div>
                    <div class="text-sm text-base-content/70">
                        Total: {round.groupCount * round.playersPerGroup} players
                    </div>
                </div>

                <!-- All inputs in a single line -->
                <div class="flex flex-wrap gap-3">
                    <!-- Players per Group -->
                    <div class="form-control flex-1 min-w-[150px]">
                        <label class="label" for="playersPerGroup">
                            <span class="label-text">Players per Group</span>
                        </label>
                        <input
                            type="number"
                            class="input input-bordered w-full"
                            min="2"
                            value={round.playersPerGroup}
                            oninput={(e) => updateRoundField(index, 'playersPerGroup', parseInt(e.currentTarget.value) || 2)}
                        />
                    </div>

                    <!-- Matches per Group -->
                    <div class="form-control flex-1 min-w-[150px]">
                        <label class="label" for="matchesPerGroup">
                            <span class="label-text">Matches per Group</span>
                        </label>
                        <input
                            type="number"
                            class="input input-bordered w-full"
                            min="1"
                            value={round.matchesPerGroup}
                            oninput={(e) => updateRoundField(index, 'matchesPerGroup', parseInt(e.currentTarget.value) || 1)}
                        />
                    </div>

                    <!-- Advancing Players per Group -->
                    <div class="form-control flex-1 min-w-[150px]">
                        <label class="label" for="advancingPlayersPerGroup">
                            <span class="label-text">Advancing Players</span>
                        </label>
                        <input
                            type="number"
                            class="input input-bordered w-full"
                            min="1"
                            max={round.playersPerGroup - 1}
                            value={round.advancingPlayersPerGroup}
                            oninput={(e) => updateRoundField(index, 'advancingPlayersPerGroup', parseInt(e.currentTarget.value) || 1)}
                        />
                    </div>

                    <!-- Concurrent Groups -->
                    <div class="form-control flex-1 min-w-[150px]">
                        <label class="label" for="concurrentGroups">
                            <span class="label-text">Concurrent Groups</span>
                        </label>
                        <input
                            type="number"
                            class="input input-bordered w-full"
                            min="1"
                            max={round.groupCount}
                            value={round.concurrentGroups}
                            oninput={(e) => updateRoundField(index, 'concurrentGroups', parseInt(e.currentTarget.value) || 1)}
                        />
                        <label class="label" for="concurrentGroups">
                            <span class="label-text-alt text-base-content/70">Groups that can play at the same time</span>
                        </label>
                    </div>
                </div>
            </div>
        </div>
    {/each}

    <!-- Add Round Button -->
    <div class="flex justify-center mt-4">
        <button
            type="button"
            class="btn btn-outline btn-primary"
            onclick={handleAddRound}
        >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
            </svg>
            Add Round
        </button>
    </div>
</div>