<script lang="ts">
    import type { TournamentFormData, Round } from '$lib/types/tournament';
    import { calculateTotalRounds } from '$lib/utils/tournamentUtils';
    
    // Props
    export let formData: TournamentFormData;
    
    // Callback props instead of event dispatcher
    export let onAddRound = () => {};
    export let onRemoveRound = (index: number) => {};
    export let onUpdateRound = (index: number, field: string, value: number) => {};
    
    // Helper functions to call callbacks
    function handleAddRound() {
        onAddRound();
    }
    
    function handleRemoveRound(index: number) {
        onRemoveRound(index);
    }
    
    function updateRoundField(index: number, field: string, value: number) {
        onUpdateRound(index, field, value);
    }
    
    // Derived value for total rounds
    $: totalRounds = calculateTotalRounds(formData);
</script>

<div class="space-y-6">
    <div class="flex justify-between items-center">
        <h3 class="font-medium text-lg">Tournament Rounds</h3>
        <span class="badge badge-primary">{totalRounds} Rounds</span>
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
                        Total: {round.groupCount * (round.playersPerGroup || 0) || 0} players
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
                            oninput={(e) => {
                                const value = e.currentTarget.value;
                                if (value === '') {
                                    // Allow empty value during editing
                                    updateRoundField(index, 'playersPerGroup', null);
                                } else {
                                    updateRoundField(index, 'playersPerGroup', parseInt(value));
                                }
                            }}
                            onblur={(e) => {
                                // Apply fallback only when field loses focus
                                const value = parseInt(e.currentTarget.value);
                                if (isNaN(value) || value < 2) {
                                    updateRoundField(index, 'playersPerGroup', 2);
                                    e.currentTarget.value = '2';
                                }
                            }}
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
                            oninput={(e) => {
                                const value = e.currentTarget.value;
                                if (value === '') {
                                    // Allow empty value during editing
                                    updateRoundField(index, 'matchesPerGroup', null);
                                } else {
                                    updateRoundField(index, 'matchesPerGroup', parseInt(value));
                                }
                            }}
                            onblur={(e) => {
                                // Apply fallback only when field loses focus
                                const value = parseInt(e.currentTarget.value);
                                if (isNaN(value) || value < 1) {
                                    updateRoundField(index, 'matchesPerGroup', 1);
                                    e.currentTarget.value = '1';
                                }
                            }}
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
                            value={round.advancingPlayersPerGroup}
                            oninput={(e) => {
                                const value = e.currentTarget.value;
                                if (value === '') {
                                    // Allow empty value during editing
                                    updateRoundField(index, 'advancingPlayersPerGroup', null);
                                } else {
                                    updateRoundField(index, 'advancingPlayersPerGroup', parseInt(value));
                                }
                            }}
                            onblur={(e) => {
                                // Apply fallback only when field loses focus
                                const value = parseInt(e.currentTarget.value);
                                if (isNaN(value) || value < 1) {
                                    updateRoundField(index, 'advancingPlayersPerGroup', 1);
                                    e.currentTarget.value = '1';
                                } else if (round.playersPerGroup && value >= round.playersPerGroup) {
                                    // Ensure advancing players is less than total players
                                    const maxAdvancing = Math.max(1, round.playersPerGroup - 1);
                                    updateRoundField(index, 'advancingPlayersPerGroup', maxAdvancing);
                                    e.currentTarget.value = maxAdvancing.toString();
                                }
                            }}
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
                            oninput={(e) => {
                                const value = e.currentTarget.value;
                                if (value === '') {
                                    // Allow empty value during editing
                                    updateRoundField(index, 'concurrentGroups', null);
                                } else {
                                    updateRoundField(index, 'concurrentGroups', parseInt(value));
                                }
                            }}
                            onblur={(e) => {
                                // Apply fallback only when field loses focus
                                const value = parseInt(e.currentTarget.value);
                                if (isNaN(value) || value < 1) {
                                    updateRoundField(index, 'concurrentGroups', 1);
                                    e.currentTarget.value = '1';
                                } else if (value > round.groupCount) {
                                    // Ensure concurrent groups is not more than total groups
                                    updateRoundField(index, 'concurrentGroups', round.groupCount);
                                    e.currentTarget.value = round.groupCount.toString();
                                }
                            }}
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