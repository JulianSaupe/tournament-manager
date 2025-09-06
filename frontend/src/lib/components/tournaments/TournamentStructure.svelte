<script lang="ts">
    import type {Round} from '$lib/types/tournament/tournament';
    import {tournamentForm} from '$lib/stores/tournamentForm';
    import type {TournamentRoundData} from '$lib/validation/tournamentSchema';

    // Simplified state - just rounds and auto-calculation mode
    let rounds: Round[] = $state([
        {
            name: 'Round 1',
            groupCount: 1,
            playersPerGroup: 8,
            matchesPerGroup: 1,
            advancingPlayersPerGroup: 1,
            concurrentGroups: 1
        }
    ]);

    let autoCalculate: boolean = $state(true);

    // Derived calculations
    const totalRounds = $derived(rounds.length || 0);
    const totalMatches = $derived(calculateTotalMatches());
    const validationStatus = $derived(checkValidation());

    // Transform Round objects to backend-expected TournamentRoundData format
    function transformRoundToTournamentRoundData(round: Round): TournamentRoundData {
        return {
            name: round.name,
            matchCount: round.matchesPerGroup,
            playerAdvancementCount: round.advancingPlayersPerGroup,
            groupSize: round.playersPerGroup,
            groupCount: round.groupCount,
            concurrentGroupCount: round.concurrentGroups
        };
    }

    // Auto-detect final round and auto-recalculate group counts
    $effect(() => {
        if (autoCalculate && rounds.length > 0) {
            const updatedRounds = [...rounds];
            let hasChanges = false;
            
            // First round: calculate groups based on total players
            if (updatedRounds[0]) {
                const newGroupCount = Math.ceil($tournamentForm.playerCount / updatedRounds[0].playersPerGroup);
                if (updatedRounds[0].groupCount !== newGroupCount) {
                    updatedRounds[0].groupCount = newGroupCount;
                    hasChanges = true;
                }
            }
            
            // Subsequent rounds: calculate groups based on advancing players from previous round
            for (let i = 1; i < updatedRounds.length; i++) {
                const prevRound = updatedRounds[i - 1];
                const currentRound = updatedRounds[i];
                const advancingPlayers = prevRound.groupCount * prevRound.advancingPlayersPerGroup;
                const newGroupCount = Math.ceil(advancingPlayers / currentRound.playersPerGroup);
                
                if (updatedRounds[i].groupCount !== newGroupCount) {
                    updatedRounds[i].groupCount = newGroupCount;
                    hasChanges = true;
                }
            }

            // Auto-detect final round - only the LAST round with 1 group should be "Finals"
            // Also set advancing players to 1 (winner) for final round
            for (let i = 0; i < updatedRounds.length; i++) {
                const isLastRound = i === updatedRounds.length - 1;
                const isFinalRound = updatedRounds[i].groupCount === 1 && isLastRound;
                
                if (isFinalRound && !updatedRounds[i].name.toLowerCase().includes('final')) {
                    updatedRounds[i].name = 'Finals';
                    updatedRounds[i].advancingPlayersPerGroup = 1; // Winner takes all
                    hasChanges = true;
                } else if (!isFinalRound && updatedRounds[i].name.toLowerCase().includes('final')) {
                    // Remove "Finals" from non-final rounds
                    updatedRounds[i].name = `Round ${i + 1}`;
                    hasChanges = true;
                }
            }
            
            if (hasChanges) {
                rounds = updatedRounds;
            }
        }
    });

    // Synchronize rounds data with form store
    $effect(() => {
        $tournamentForm.rounds = rounds.map(transformRoundToTournamentRoundData);
    });

    // Helper to check if a round is the final round
    function isFinalRound(index: number): boolean {
        return index === rounds.length - 1 && rounds[index].groupCount === 1;
    }

    // Calculation helpers
    function calculateTotalMatches(): number {
        return rounds.reduce((total, round) => {
            return total + (round.groupCount * round.matchesPerGroup);
        }, 0);
    }

    function checkValidation(): 'good' | 'warning' | 'error' {
        // Check if final round has exactly 1 group
        const lastRound = rounds[rounds.length - 1];
        if (!lastRound || lastRound.groupCount !== 1) {
            return 'error';
        }

        // Check player progression
        for (let i = 0; i < rounds.length; i++) {
            const round = rounds[i];
            const prevRound = i > 0 ? rounds[i - 1] : null;
            
            if (!$tournamentForm.allowUnderfilledGroups) {
                const expectedPlayers = prevRound 
                    ? prevRound.groupCount * prevRound.advancingPlayersPerGroup
                    : $tournamentForm.playerCount;
                
                const actualPlayers = round.groupCount * round.playersPerGroup;
                
                if (actualPlayers !== expectedPlayers) {
                    return 'error';
                }
            }

            if (round.advancingPlayersPerGroup >= round.playersPerGroup) {
                return 'error';
            }
        }

        return 'good';
    }

    function getValidationMessage(round: Round, index: number): string {
        const prevRound = index > 0 ? rounds[index - 1] : null;
        const expectedPlayers = prevRound 
            ? prevRound.groupCount * prevRound.advancingPlayersPerGroup
            : $tournamentForm.playerCount;
        
        const actualPlayers = round.groupCount * round.playersPerGroup;
        
        if (actualPlayers !== expectedPlayers && !$tournamentForm.allowUnderfilledGroups) {
            return `Expected ${expectedPlayers} players, but configured for ${actualPlayers}. Click "Quick Fix" to adjust.`;
        }
        return '';
    }

    // Action handlers
    function addRound(): void {
        const lastRound = rounds[rounds.length - 1];
        const advancingPlayers = lastRound.advancingPlayersPerGroup * lastRound.groupCount;
        
        // Smart defaults for new round
        const playersPerGroup = Math.min(8, Math.max(4, advancingPlayers));
        const groupCount = Math.ceil(advancingPlayers / playersPerGroup);
        
        const newRound: Round = {
            name: `Round ${rounds.length + 1}`,
            groupCount: groupCount,
            playersPerGroup: playersPerGroup,
            matchesPerGroup: 1,
            advancingPlayersPerGroup: Math.max(1, Math.floor(playersPerGroup / 2)),
            concurrentGroups: Math.min(groupCount, 2)
        };
        
        rounds = [...rounds, newRound];
    }

    function removeRound(index: number): void {
        if (rounds.length > 1) {
            rounds = rounds.filter((_, i) => i !== index);
        }
    }

    function quickFix(index: number): void {
        const updatedRounds = [...rounds];
        const round = updatedRounds[index];
        const prevRound = index > 0 ? updatedRounds[index - 1] : null;
        const expectedPlayers = prevRound 
            ? prevRound.groupCount * prevRound.advancingPlayersPerGroup
            : $tournamentForm.playerCount;

        // Try to fix by adjusting group count first
        const newGroupCount = Math.ceil(expectedPlayers / round.playersPerGroup);
        updatedRounds[index] = {
            ...round,
            groupCount: newGroupCount
        };
        
        rounds = updatedRounds;
    }

    function updateRoundField(index: number, field: keyof Round, value: number): void {
        if (index < 0 || index >= rounds.length || value < 0) return;

        const updatedRounds = [...rounds];
        updatedRounds[index] = {
            ...updatedRounds[index],
            [field]: value
        };

        // Enforce constraints
        if (field === 'advancingPlayersPerGroup') {
            const maxAdvancing = isFinalRound(index) ? 1 : updatedRounds[index].playersPerGroup - 1;
            updatedRounds[index].advancingPlayersPerGroup = Math.min(value, maxAdvancing);
        }

        if (field === 'concurrentGroups') {
            updatedRounds[index].concurrentGroups = Math.min(value, updatedRounds[index].groupCount);
        }

        rounds = updatedRounds;
    }
</script>

<div class="card bg-base-100 shadow-sm">
    <div class="card-body">
        <div class="mb-6 flex items-center justify-between">
            <h2 class="text-xl font-semibold">Tournament Structure</h2>
            <span class="badge badge-primary">{totalRounds} Rounds</span>
        </div>

        <!-- Auto-Calculate Toggle -->
        <div class="mb-6 flex items-center gap-4">
            <label class="cursor-pointer label">
                <span class="label-text mr-3">Auto-calculate group counts</span>
                <input type="checkbox" class="toggle toggle-primary" bind:checked={autoCalculate} />
            </label>
            <div class="text-sm text-base-content/70">
                When enabled, group counts are automatically calculated based on player count and advancement
            </div>
        </div>

        <!-- Tournament Overview Stats -->
        <div class="mb-6 grid grid-cols-3 gap-4">
            <div class="stat bg-base-200 rounded-lg p-3">
                <div class="stat-title text-xs">Total Matches</div>
                <div class="stat-value text-xl">{totalMatches}</div>
            </div>
            <div class="stat bg-base-200 rounded-lg p-3">
                <div class="stat-title text-xs">Final Groups</div>
                <div class="stat-value text-xl">{rounds[rounds.length - 1]?.groupCount || 0}</div>
            </div>
            <div class="stat bg-base-200 rounded-lg p-3">
                <div class="stat-title text-xs">Status</div>
                <div class="stat-value text-lg {validationStatus === 'good' ? 'text-success' : validationStatus === 'warning' ? 'text-warning' : 'text-error'}">
                    {validationStatus === 'good' ? '✓ Valid' : '✗ Issues'}
                </div>
            </div>
        </div>

        <!-- Player Flow Visualization -->
        <div class="mb-6 overflow-x-auto">
            <div class="flex items-center gap-3 min-w-fit p-2">
                {#each rounds as round, index}
                    <div class="flex flex-col items-center">
                        <div class="badge badge-sm mb-2">{round.name}</div>
                        <div class="bg-base-200 rounded-lg p-3 min-w-[100px] text-center {getValidationMessage(round, index) ? 'border-2 border-error' : ''}">
                            <div class="text-lg font-bold">{round.groupCount}</div>
                            <div class="text-xs text-base-content/70">groups</div>
                            <div class="text-sm mt-1">{round.groupCount * round.playersPerGroup} players</div>
                        </div>
                    </div>
                    
                    {#if index < rounds.length - 1}
                        <div class="flex flex-col items-center px-2">
                            <div class="text-xs text-base-content/60 mb-1 text-center">
                                {round.advancingPlayersPerGroup * round.groupCount} advance
                            </div>
                            <svg class="w-6 h-4 text-base-content/40" fill="currentColor" viewBox="0 0 20 10">
                                <path d="M12 1l6 4-6 4V6H2V4h10V1z"/>
                            </svg>
                        </div>
                    {/if}
                {/each}
            </div>
        </div>

        <!-- Round Configuration -->
        <div class="space-y-6">
            {#each rounds as round, index (index)}
                <div class="card bg-base-50 border {getValidationMessage(round, index) ? 'border-error' : 'border-base-300'}">
                    <div class="card-body p-4">
                        <!-- Round Header -->
                        <div class="flex items-center justify-between mb-4">
                            <div class="flex items-center gap-3">
                                <h4 class="font-medium text-lg">{round.name}</h4>
                                {#if getValidationMessage(round, index)}
                                    <div class="tooltip tooltip-error" data-tip={getValidationMessage(round, index)}>
                                        <div class="badge badge-error badge-sm">!</div>
                                    </div>
                                {/if}
                            </div>
                            
                            <div class="flex items-center gap-2">
                                {#if getValidationMessage(round, index)}
                                    <button type="button" class="btn btn-xs btn-success" onclick={() => quickFix(index)}>
                                        Quick Fix
                                    </button>
                                {/if}
                                
                                {#if rounds.length > 1}
                                    <button type="button" class="btn btn-xs btn-error" onclick={() => removeRound(index)}>
                                        Remove
                                    </button>
                                {/if}
                            </div>
                        </div>

                        <!-- Labels Row -->
                        <div class="grid grid-cols-5 gap-4 mb-2">
                            <div class="text-sm font-medium text-base-content/80">Players per Group</div>
                            <div class="text-sm font-medium text-base-content/80">Matches per Group</div>
                            <div class="text-sm font-medium text-base-content/80">Advancing per Group</div>
                            <div class="text-sm font-medium text-base-content/80">Concurrent Groups</div>
                            <div class="text-sm font-medium text-base-content/80">Summary</div>
                        </div>

                        <!-- Inputs Row -->
                        <div class="grid grid-cols-5 gap-4">
                            <!-- Players per Group -->
                            <div class="flex flex-col">
                                <div class="join">
                                    <input 
                                        id="playersPerGroup-{index}"
                                        type="number" 
                                        class="input input-bordered input-sm join-item flex-1"
                                        min="2"
                                        value={round.playersPerGroup}
                                        oninput={(e) => updateRoundField(index, 'playersPerGroup', parseInt(e.currentTarget.value) || 2)} />
                                    
                                    <button type="button" class="btn btn-sm join-item px-2" 
                                            onclick={() => updateRoundField(index, 'playersPerGroup', Math.max(2, round.playersPerGroup - 1))}>
                                        −
                                    </button>
                                    <button type="button" class="btn btn-sm join-item px-2" 
                                            onclick={() => updateRoundField(index, 'playersPerGroup', round.playersPerGroup + 1)}>
                                        +
                                    </button>
                                </div>
                                <div class="text-xs text-base-content/60 mt-1">
                                    = {round.groupCount} groups
                                </div>
                            </div>

                            <!-- Matches per Group -->
                            <div>
                                <input 
                                    id="matchesPerGroup-{index}"
                                    type="number" 
                                    class="input input-bordered input-sm w-full"
                                    min="1"
                                    value={round.matchesPerGroup}
                                    oninput={(e) => updateRoundField(index, 'matchesPerGroup', parseInt(e.currentTarget.value) || 1)} />
                            </div>

                            <!-- Advancing Players -->
                            <div>
                                {#if isFinalRound(index)}
                                    <!-- Final round: winner only -->
                                    <input 
                                        id="advancingPlayers-{index}"
                                        type="number" 
                                        class="input input-bordered input-sm w-full"
                                        value="1"
                                        disabled />
                                    <div class="text-xs text-base-content/60 mt-1">
                                        winner only
                                    </div>
                                {:else}
                                    <!-- Regular rounds -->
                                    <input 
                                        id="advancingPlayers-{index}"
                                        type="number" 
                                        class="input input-bordered input-sm w-full"
                                        min="0"
                                        max={round.playersPerGroup - 1}
                                        value={round.advancingPlayersPerGroup}
                                        oninput={(e) => updateRoundField(index, 'advancingPlayersPerGroup', parseInt(e.currentTarget.value) || 0)} />
                                    <div class="text-xs text-base-content/60 mt-1">
                                        max {round.playersPerGroup - 1}
                                    </div>
                                {/if}
                            </div>

                            <!-- Concurrent Groups -->
                            <div>
                                {#if round.groupCount === 1}
                                    <!-- Single group: concurrent input is obsolete -->
                                    <input 
                                        id="concurrentGroups-{index}"
                                        type="number" 
                                        class="input input-bordered input-sm w-full"
                                        value="1"
                                        disabled />
                                    <div class="text-xs text-base-content/60 mt-1">
                                        single group
                                    </div>
                                {:else}
                                    <!-- Multiple groups -->
                                    <input 
                                        id="concurrentGroups-{index}"
                                        type="number" 
                                        class="input input-bordered input-sm w-full"
                                        min="1"
                                        max={round.groupCount}
                                        value={round.concurrentGroups}
                                        oninput={(e) => updateRoundField(index, 'concurrentGroups', parseInt(e.currentTarget.value) || 1)} />
                                    <div class="text-xs text-base-content/60 mt-1">
                                        max {round.groupCount}
                                    </div>
                                {/if}
                            </div>

                            <!-- Simplified Round Summary -->
                            <div class="bg-base-200 rounded p-3 text-center">
                                <div class="text-sm font-semibold">{round.groupCount * round.matchesPerGroup} matches</div>
                                {#if !isFinalRound(index)}
                                    <div class="text-xs text-base-content/70 mt-1">
                                        {round.advancingPlayersPerGroup * round.groupCount} advance
                                    </div>
                                {:else}
                                    <div class="text-xs text-base-content/70 mt-1">
                                        1 winner
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>
            {/each}
        </div>

        <!-- Add Round Button -->
        <div class="mt-6 text-center">
            <button type="button" class="btn btn-primary" onclick={() => addRound()}>
                <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                Add Round
            </button>
        </div>
    </div>
</div>
