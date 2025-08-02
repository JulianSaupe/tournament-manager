<script lang="ts">
    import {goto} from '$app/navigation';
    import {mockTournaments} from '$lib/mockData';

    // Define interfaces for tournament structure
    interface Round {
        name: string;
        groupCount: number; // Will be calculated based on playerCount and playersPerGroup
        playersPerGroup: number;
        matchesPerGroup: number;
        advancingPlayersPerGroup: number;
        concurrentGroups: number; // Number of groups that can play concurrently
    }

    // Form data
    let formData = {
        name: '',
        startDate: '',
        endDate: '',
        playerCount: 0,
        groupPhase: false,
        allowPartiallyFilledGroups: false,
        groupSize: 4,
        rounds: [
            {
                name: 'Round 1',
                groupCount: 4,
                playersPerGroup: 4,
                matchesPerGroup: 6,
                advancingPlayersPerGroup: 2,
                concurrentGroups: 2
            }
        ] as Round[]
    };

    // Form validation
    let errors = {
        name: '',
        startDate: '',
        endDate: '',
        playerCount: ''
    };

    // Helper functions for rounds management
    function addRound() {
        const lastRound = formData.rounds[formData.rounds.length - 1];
        const newRoundNumber = formData.rounds.length + 1;

        // Calculate default values based on previous round
        const newPlayersPerGroup = lastRound.advancingPlayersPerGroup * 2; // Double the advancing players

        const newRound: Round = {
            name: `Round ${newRoundNumber}`,
            groupCount: 1, // Will be calculated automatically by reactive statement
            playersPerGroup: newPlayersPerGroup,
            matchesPerGroup: calculateDefaultMatches(newPlayersPerGroup),
            advancingPlayersPerGroup: Math.max(1, Math.floor(lastRound.advancingPlayersPerGroup / 2)), // Half the advancing players
            concurrentGroups: Math.max(1, Math.floor(lastRound.concurrentGroups / 2)) // Half the concurrent groups, minimum 1
        };

        formData.rounds = [...formData.rounds, newRound];
    }

    // Calculate default number of matches for a group
    function calculateDefaultMatches(players: number): number {
        // Round-robin tournament: each player plays against every other player once
        return players > 1 ? Math.floor((players * (players - 1)) / 2) : 0;
    }

    function removeRound(index: number) {
        if (formData.rounds.length > 1) {
            formData.rounds = formData.rounds.filter((_, i) => i !== index);
        }
    }

    function calculateTotalRounds() {
        return formData.rounds.length;
    }

    // Functions for tournament visualization
    function generateVisualizationData() {
        const data = [];
        let totalPlayers = formData.playerCount;

        // If group phase is enabled, use it as the starting point
        if (formData.groupPhase && formData.playerCount > 0) {
            const groupCount = Math.ceil(totalPlayers / formData.groupSize);
            const playersPerGroup = formData.groupSize;

            // Calculate how many players advance to the first round
            const advancingPlayers = groupCount * formData.rounds[0].advancingPlayersPerGroup;

            data.push({
                name: 'Group Phase',
                groupCount,
                playersPerGroup,
                totalPlayers,
                advancingPlayers,
                matchesPerGroup: calculateMatchesForGroup(playersPerGroup),
                advancingPlayersPerGroup: formData.rounds[0].advancingPlayersPerGroup,
                concurrentGroups: formData.groupSize // For group phase, default to group size as concurrent groups
            });

            totalPlayers = advancingPlayers;
        }

        // Helper function to calculate default matches for a group
        function calculateMatchesForGroup(players: number): number {
            // Default calculation: each player plays against every other player once
            return players > 1 ? (players * (players - 1)) / 2 : 0;
        }

        // Add data for each configured round
        formData.rounds.forEach((round, index) => {
            const advancingPlayers = index < formData.rounds.length - 1
                ? formData.rounds[index + 1].groupCount * formData.rounds[index + 1].playersPerGroup
                : 1; // Final winner

            data.push({
                name: round.name,
                groupCount: round.groupCount,
                playersPerGroup: round.playersPerGroup,
                totalPlayers,
                advancingPlayers,
                matchesPerGroup: round.matchesPerGroup,
                advancingPlayersPerGroup: round.advancingPlayersPerGroup,
                concurrentGroups: round.concurrentGroups
            });

            totalPlayers = advancingPlayers;
        });

        return data;
    }

    // Reactive statements to calculate group counts and update visualization
    $: {
        // Calculate group counts for each round based on player count and players per group
        if (formData.playerCount > 0) {
            let availablePlayers = formData.playerCount;

            // If group phase is enabled, calculate players for the first round
            if (formData.groupPhase) {
                // Calculate number of groups in the group phase
                const groupPhaseGroupCount = Math.ceil(formData.playerCount / formData.groupSize);

                // Calculate players advancing from group phase to first round
                availablePlayers = groupPhaseGroupCount * formData.rounds[0].advancingPlayersPerGroup;
            }

            // Update each round's group count
            formData.rounds.forEach((round, index) => {
                if (index > 0) {
                    // For rounds after the first, available players come from previous round
                    const prevRound = formData.rounds[index - 1];
                    availablePlayers = prevRound.groupCount * prevRound.advancingPlayersPerGroup;
                }

                // Calculate group count based on available players and players per group
                // Ensure we always have at least 1 group
                round.groupCount = Math.max(1, Math.ceil(availablePlayers / round.playersPerGroup));
            });
        }
    }

    // Reactive statement to update visualization when form data changes
    $: visualizationData = generateVisualizationData();

    // Form validation function
    function validateForm() {
        let isValid = true;

        // Reset errors
        errors = {
            name: '',
            startDate: '',
            endDate: '',
            playerCount: ''
        };

        // Validate name
        if (!formData.name.trim()) {
            errors.name = 'Tournament name is required';
            isValid = false;
        }

        // Validate start date
        if (!formData.startDate) {
            errors.startDate = 'Start date is required';
            isValid = false;
        }

        // Validate end date
        if (!formData.endDate) {
            errors.endDate = 'End date is required';
            isValid = false;
        } else if (formData.startDate && formData.endDate && new Date(formData.endDate) < new Date(formData.startDate)) {
            errors.endDate = 'End date must be after start date';
            isValid = false;
        }

        // Validate player count
        if (!formData.playerCount || formData.playerCount <= 0) {
            errors.playerCount = 'Number of players must be greater than 0';
            isValid = false;
        }

        // Validate rounds configuration
        if (formData.rounds.length === 0) {
            alert('At least one round is required');
            isValid = false;
        }

        // Check if each round has valid settings
        for (let i = 0; i < formData.rounds.length; i++) {
            const round = formData.rounds[i];

            if (round.groupCount <= 0) {
                alert(`Round ${i + 1}: Number of groups must be greater than 0`);
                isValid = false;
            }

            if (round.playersPerGroup <= 1) {
                alert(`Round ${i + 1}: Players per group must be at least 2`);
                isValid = false;
            }

            if (round.matchesPerGroup <= 0) {
                alert(`Round ${i + 1}: Matches per group must be greater than 0`);
                isValid = false;
            }

            if (round.advancingPlayersPerGroup <= 0 || round.advancingPlayersPerGroup >= round.playersPerGroup) {
                alert(`Round ${i + 1}: Advancing players must be between 1 and ${round.playersPerGroup - 1}`);
                isValid = false;
            }

            if (round.concurrentGroups <= 0) {
                alert(`Round ${i + 1}: Concurrent groups must be at least 1`);
                isValid = false;
            }

            if (round.concurrentGroups > round.groupCount) {
                alert(`Round ${i + 1}: Concurrent groups cannot exceed the total number of groups (${round.groupCount})`);
                isValid = false;
            }
        }

        return isValid;
    }

    // Form submission
    function handleSubmit(event: Event) {
        event.preventDefault();
        if (validateForm()) {
            // In a real application, this would send data to the backend
            // For now, we'll just log it and navigate back to the home page

            // Create a new tournament object
            const newTournament = {
                id: (mockTournaments.length + 1).toString(),
                name: formData.name,
                startDate: new Date(formData.startDate).toISOString(),
                endDate: new Date(formData.endDate).toISOString(),
                playerCount: formData.playerCount,
                status: 'draft',
                groupPhase: formData.groupPhase,
                groupSize: formData.groupPhase ? formData.groupSize : null,
                allowPartiallyFilledGroups: formData.groupPhase ? formData.allowPartiallyFilledGroups : false,
                rounds: formData.rounds.map(round => ({
                    name: round.name,
                    groupCount: round.groupCount,
                    playersPerGroup: round.playersPerGroup,
                    matchesPerGroup: round.matchesPerGroup,
                    advancingPlayersPerGroup: round.advancingPlayersPerGroup,
                    concurrentGroups: round.concurrentGroups
                }))
            };

            console.log('New tournament created:', newTournament);

            // Navigate back to the home page
            goto('/');
        }
    }

    // Cancel and go back to home page
    function handleCancel() {
        goto('/');
    }
</script>

<div class="container mx-auto">
    <h1 class="text-2xl mb-6">Create New Tournament</h1>

    <form onsubmit={handleSubmit} class="space-y-6">
        <div class="card bg-base-100 shadow-sm">
            <div class="card-body">
                <div class="card-title">
                    <h2 class="card-title">Tournament Details</h2>
                </div>
                <div class="grid grid-cols-2 gap-6">
                    <!-- Tournament Name -->
                    <div class="form-control w-full">
                        <label for="name" class="label">
                            <span class="label-text">Tournament Name</span>
                        </label>
                        <input
                                type="text"
                                id="name"
                                bind:value={formData.name}
                                class="input input-bordered w-full {errors.name ? 'input-error' : ''}"
                                placeholder="Enter tournament name"
                        />
                        {#if errors.name}
                            <label class="label" for="name">
                                <span class="label-text-alt text-error">{errors.name}</span>
                            </label>
                        {/if}
                    </div>

                    <!-- Number of Players -->
                    <div class="form-control w-full">
                        <label for="playerCount" class="label">
                            <span class="label-text">Number of Players</span>
                        </label>
                        <input
                                type="number"
                                id="playerCount"
                                bind:value={formData.playerCount}
                                min="1"
                                class="input input-bordered w-full {errors.playerCount ? 'input-error' : ''}"
                                placeholder="Enter number of players"
                        />
                        {#if errors.playerCount}
                            <label class="label" for="playerCount">
                                <span class="label-text-alt text-error">{errors.playerCount}</span>
                            </label>
                        {/if}
                    </div>

                    <!-- Start Date -->
                    <div class="form-control w-full">
                        <label for="startDate" class="label">
                            <span class="label-text">Start Date</span>
                        </label>
                        <input
                                type="date"
                                id="startDate"
                                bind:value={formData.startDate}
                                class="input input-bordered w-full {errors.startDate ? 'input-error' : ''}"
                        />
                        {#if errors.startDate}
                            <label class="label" for="startDate">
                                <span class="label-text-alt text-error">{errors.startDate}</span>
                            </label>
                        {/if}
                    </div>

                    <!-- End Date -->
                    <div class="form-control w-full">
                        <label for="endDate" class="label">
                            <span class="label-text">End Date</span>
                        </label>
                        <input
                                type="date"
                                id="endDate"
                                bind:value={formData.endDate}
                                class="input input-bordered w-full {errors.endDate ? 'input-error' : ''}"
                        />
                        {#if errors.endDate}
                            <label class="label" for="endDate">
                                <span class="label-text-alt text-error">{errors.endDate}</span>
                            </label>
                        {/if}
                    </div>
                </div>
            </div>
        </div>

        <div class="card bg-base-100 shadow-sm">
            <div class="card-body">
                <div class="card-title mb-4 flex justify-between items-center">
                    <h2 class="card-title">Tournament Structure</h2>
                    <span class="badge badge-primary">{calculateTotalRounds()} Rounds</span>
                </div>

                <!-- Group Phase Toggle -->
                <div class="form-control bg-base-200 p-4 rounded-lg mb-4">
                    <div class="flex items-left mb-1">
                        <span class="label-text font-medium text-base">Enable Group Phase</span>
                        <input type="checkbox" class="toggle toggle-primary ms-3" bind:checked={formData.groupPhase}/>
                    </div>
                    <p class="text-xs text-base-content/70">Organize players into groups for the initial tournament
                        phase</p>
                </div>

                <!-- Group Phase Settings (only visible when group phase is enabled) -->
                {#if formData.groupPhase}
                    <div class="bg-base-100 border border-base-300 rounded-lg p-4 ml-2 sm:ml-4 mb-6">
                        <h3 class="font-medium mb-4">Group Phase Settings</h3>

                        <div class="grid gap-4">
                            <!-- Group Size -->
                            <div class="form-control w-full">
                                <label for="groupSize" class="label">
                                    <span class="label-text">Group Size</span>
                                </label>
                                <input
                                        type="number"
                                        id="groupSize"
                                        class="input input-bordered w-full"
                                        placeholder="Enter group size"
                                        min="2"
                                        bind:value={formData.groupSize}
                                />
                                <label class="label" for="groupSize">
                                    <span class="label-text-alt text-base-content/70">Recommended: 4-6 players per group</span>
                                </label>
                            </div>

                            <!-- Allow Partially Filled Groups -->
                            <div class="form-control">
                                <label class="label cursor-pointer justify-between flex-wrap sm:flex-nowrap">
                                    <span class="label-text w-full sm:w-auto mb-2 sm:mb-0">Allow Partially Filled Groups</span>
                                    <input
                                            type="checkbox"
                                            class="checkbox"
                                            bind:checked={formData.allowPartiallyFilledGroups}
                                    />
                                </label>
                                <label class="label mt-[-8px]" for="allowPartiallyFilledGroups">
                                    <span class="label-text-alt text-base-content/70">If disabled, all groups will have the same number of players</span>
                                </label>
                            </div>
                        </div>
                    </div>
                {/if}

                <!-- Rounds Configuration -->
                <div class="space-y-6">
                    <h3 class="font-medium text-lg">Tournament Rounds</h3>

                    {#each formData.rounds as round, index}
                        <div class="bg-base-100 border border-base-300 rounded-lg p-4">
                            <div class="flex justify-between items-center mb-4">
                                <h4 class="font-medium">{round.name}</h4>
                                {#if index > 0}
                                    <button
                                            type="button"
                                            class="btn btn-sm btn-error btn-outline"
                                            onclick={() => removeRound(index)}
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
                                                bind:value={round.playersPerGroup}
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
                                                bind:value={round.matchesPerGroup}
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
                                                bind:value={round.advancingPlayersPerGroup}
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
                                                bind:value={round.concurrentGroups}
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
                                onclick={addRound}
                        >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 24 24"
                                 stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                      d="M12 4v16m8-8H4"/>
                            </svg>
                            Add Round
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Tournament Visualization -->
        <div class="card bg-base-100 shadow-sm">
            <div class="card-body">
                <div class="card-title mb-4 flex justify-between items-center">
                    <h2 class="card-title">Tournament Progression</h2>
                    {#if formData.playerCount > 0}
                        <div class="badge badge-primary">{formData.playerCount} Players</div>
                    {/if}
                </div>

                <div class="mb-6 p-4 bg-base-200 rounded-lg border-l-4 border-primary">
                    <p class="text-sm">
                        <span class="font-medium">How to use:</span> Configure your tournament structure by adding
                        rounds and setting the number of players per group and advancing players.
                        The visualization below will update automatically to show how players progress through the
                        tournament.
                    </p>
                </div>

                {#if formData.playerCount > 0}
                    <div class="overflow-x-auto">
                        <div class="tournament-tree pb-6">
                            <!-- Tournament Tree Visualization -->
                            <div class="flex flex-col items-center">
                                {#if visualizationData.length > 0}
                                    <!-- Calculate dynamic width based on maximum number of groups in any phase -->
                                    {@const maxGroups = Math.max(...visualizationData.map(phase => phase.groupCount))}
                                    {@const totalWidth = Math.max(900, maxGroups * 120)}
                                    <!-- Ensure minimum width of 900px -->
                                    {@const svgHeight = visualizationData.length * 200}
                                    <!-- SVG for the tournament tree -->
                                    <div class="w-full overflow-x-auto">
                                        <div class="min-w-[800px] md:min-w-0">
                                            <svg class="tournament-bracket" width="100%"
                                                 height="{svgHeight}"
                                                 viewBox="0 0 {totalWidth} {svgHeight}"
                                                 preserveAspectRatio="xMidYMid meet">

                                                <!-- Draw the tree structure -->
                                                {#each visualizationData.slice().reverse() as phase, reversedIndex}
                                                    {@const phaseIndex = visualizationData.length - 1 - reversedIndex}
                                                    {@const yPosition = 50 + reversedIndex * 150}
                                                    {@const numGroups = phase.groupCount}
                                                    {@const groupWidth = totalWidth / numGroups}

                                                    <!-- Phase Label with concurrent groups info -->
                                                    <text
                                                            x="{totalWidth / 2}"
                                                            y="{yPosition - 30}"
                                                            text-anchor="middle"
                                                            class="text-lg font-medium fill-primary"
                                                    >
                                                        {phase.name} ({phase.totalPlayers}
                                                        Players, {phase.concurrentGroups || 'All'} Concurrent)
                                                    </text>

                                                    <!-- Groups in this phase -->
                                                    {#each Array(numGroups) as _, groupIndex}
                                                        {@const xPosition = (groupIndex + 0.5) * groupWidth}

                                                        <!-- Group node -->
                                                        <g class="group-node"
                                                           transform="translate({xPosition}, {yPosition})">
                                                            <!-- Group box -->
                                                            <rect
                                                                    x="-50"
                                                                    y="-40"
                                                                    width="100"
                                                                    height="80"
                                                                    rx="8"
                                                                    class="fill-base-100 stroke-base-300 hover:stroke-primary"
                                                                    stroke-width="2"
                                                            />

                                                            <!-- Group label -->
                                                            <text
                                                                    y="-20"
                                                                    text-anchor="middle"
                                                                    class="font-medium text-base"
                                                            >
                                                                Group {groupIndex + 1}
                                                            </text>

                                                            <!-- Players info -->
                                                            <text
                                                                    y="0"
                                                                    text-anchor="middle"
                                                                    class="text-sm"
                                                            >
                                                                {phase.playersPerGroup} Players
                                                            </text>

                                                            <!-- Matches info -->
                                                            <text
                                                                    y="20"
                                                                    text-anchor="middle"
                                                                    class="text-sm"
                                                            >
                                                                {phase.matchesPerGroup} Matches
                                                            </text>

                                                            <!-- Advancing info for non-final rounds -->
                                                            {#if phaseIndex < visualizationData.length - 1}
                                                                <g class="advancing-indicator"
                                                                   transform="translate(0, 50)">
                                                                    <rect
                                                                            x="-45"
                                                                            y="-15"
                                                                            width="90"
                                                                            height="30"
                                                                            rx="4"
                                                                            class="fill-success/10 stroke-success"
                                                                            stroke-width="1"
                                                                    />
                                                                    <text
                                                                            y="5"
                                                                            text-anchor="middle"
                                                                            class="text-xs text-success"
                                                                    >
                                                                        Top {phase.advancingPlayersPerGroup} Advance
                                                                    </text>
                                                                </g>
                                                            {/if}
                                                        </g>

                                                        <!-- Connection lines to next phase (if not the final phase) -->
                                                        {#if reversedIndex > 0 && phaseIndex < visualizationData.length - 1}
                                                            {@const
                                                                nextPhase = visualizationData[visualizationData.length - reversedIndex]}
                                                            {@const nextNumGroups = nextPhase.groupCount}
                                                            {@const nextGroupWidth = totalWidth / nextNumGroups}
                                                            {@const nextYPosition = 50 + (reversedIndex - 1) * 150}

                                                            <!-- Calculate which group in the next phase this connects to -->
                                                            {@const
                                                                nextGroupIndex = Math.floor(groupIndex * nextNumGroups / numGroups)}
                                                            {@const
                                                                nextXPosition = (nextGroupIndex + 0.5) * nextGroupWidth}

                                                            <!-- Draw connection line with arrow -->
                                                            <path
                                                                    d="M {xPosition} {yPosition + 60} C {xPosition} {yPosition + 90}, {nextXPosition} {nextYPosition - 90}, {nextXPosition} {nextYPosition - 60}"
                                                                    fill="none"
                                                                    stroke="currentColor"
                                                                    stroke-width="2"
                                                                    class="text-primary"
                                                            />

                                                            <!-- Add arrow at the end of the path -->
                                                            <polygon
                                                                    points="{nextXPosition-5},{nextYPosition-65} {nextXPosition},{nextYPosition-60} {nextXPosition+5},{nextYPosition-65}"
                                                                    fill="currentColor"
                                                                    class="text-primary"
                                                            />

                                                            <!-- Add advancement info on the path -->
                                                            {@const midX = (xPosition + nextXPosition) / 2}
                                                            {@const midY = (yPosition + 60 + nextYPosition - 60) / 2}

                                                            <g class="advancement-label"
                                                               transform="translate({midX}, {midY})">
                                                                <rect
                                                                        x="-30"
                                                                        y="-10"
                                                                        width="60"
                                                                        height="20"
                                                                        rx="10"
                                                                        class="fill-base-100 stroke-primary"
                                                                        stroke-width="1"
                                                                        opacity="0.9"
                                                                />
                                                                <text
                                                                        y="4"
                                                                        text-anchor="middle"
                                                                        class="text-xs text-primary font-medium"
                                                                >
                                                                    Top {phase.advancingPlayersPerGroup}
                                                                </text>
                                                            </g>
                                                        {/if}
                                                    {/each}
                                                {/each}

                                                <!-- Connect the final winner to the last phase -->
                                                {#if visualizationData.length > 0}
                                                    {@const finalPhase = visualizationData[0]}
                                                    {@const finalYPosition = 50 + (visualizationData.length - 1) * 150}
                                                    {@const finalNumGroups = finalPhase.groupCount}
                                                    {@const finalGroupWidth = totalWidth / finalNumGroups}

                                                    <!-- If there's only one group in the final phase, connect directly -->
                                                    {#if finalNumGroups === 1}
                                                        {@const finalXPosition = finalGroupWidth / 2}
                                                        <path
                                                                d="M {totalWidth / 2} 20 L {totalWidth / 2} 50 L {finalXPosition} 50 L {finalXPosition} {finalYPosition - 60}"
                                                                fill="none"
                                                                stroke="currentColor"
                                                                stroke-width="2"
                                                                stroke-dasharray="4 2"
                                                                class="text-success"
                                                        />
                                                    {:else}
                                                        <!-- For multiple groups in final phase, connect to the middle -->
                                                        <path
                                                                d="M {totalWidth / 2} 20 L {totalWidth / 2} {finalYPosition - 100}"
                                                                fill="none"
                                                                stroke="currentColor"
                                                                stroke-width="2"
                                                                stroke-dasharray="4 2"
                                                                class="text-success"
                                                        />

                                                        <!-- Add text to explain -->
                                                        <text
                                                                x="{totalWidth / 2}"
                                                                y="{finalYPosition - 80}"
                                                                text-anchor="middle"
                                                                class="text-sm text-success"
                                                        >
                                                            Winner from final round
                                                        </text>
                                                    {/if}
                                                {/if}
                                            </svg>
                                        </div>
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <div class="mt-6 text-sm text-base-content/70 bg-base-200 p-3 rounded-lg">
                        <p class="italic flex items-center">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-info" viewBox="0 0 24 24"
                                 stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                      d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                            </svg>
                            This visualization is a simplified representation of your tournament structure. The actual
                            tournament may vary based on the final number of participants.
                        </p>
                    </div>
                {:else}
                    <div class="alert alert-info shadow-sm">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                             class="stroke-current shrink-0 w-6 h-6">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                        <span>Enter the number of players to see the tournament progression visualization.</span>
                    </div>
                {/if}
            </div>
        </div>

        <!-- Form Actions -->
        <div class="card-actions justify-end mt-8">
            <button type="button" class="btn btn-ghost" onclick={handleCancel}>
                Cancel
            </button>
            <button type="submit" class="btn btn-primary">
                Create Tournament
            </button>
        </div>
    </form>
</div>