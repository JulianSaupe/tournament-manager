<script lang="ts">
    import type { PhaseVisualizationData } from '$lib/types/tournament';
    
    // Props
    export let visualizationData: PhaseVisualizationData[] = [];
    export let playerCount: number = 0;
</script>

<div class="card bg-base-100 shadow-sm">
    <div class="card-body">
        <div class="card-title mb-4 flex justify-between items-center">
            <h2 class="card-title">Tournament Progression</h2>
            {#if playerCount > 0}
                <div class="badge badge-primary">{playerCount} Players</div>
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

        {#if playerCount > 0}
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