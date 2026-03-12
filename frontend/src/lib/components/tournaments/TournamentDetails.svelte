<script lang="ts">
    import {tournamentForm, tournamentFormErrors} from '$lib/stores/tournamentForm';

    function updateName(name: string): void {
        $tournamentForm.name = name;
    }

    function updatePlayerCount(count: number): void {
        $tournamentForm.playerCount = count;
    }

    function updateStartDate(date: string): void {
        $tournamentForm.startDate = date;
    }

    function updateEndDate(date: string): void {
        $tournamentForm.endDate = date;
    }

    function updateDescription(description: string): void {
        $tournamentForm.description = description;
    }
</script>

<div class="card bg-base-100 shadow-sm">
    <div class="card-body">
        <div class="card-title">
            <h2 class="card-title">TournamentProvider Details</h2>
        </div>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 sm:gap-6">
            <!-- TournamentProvider Name -->
            <div class="form-control w-full">
                <label class="label" for="name">
                    <span class="label-text">TournamentProvider Name</span>
                </label>
                <input
                        class="input-bordered input w-full {$tournamentFormErrors.name ? 'input-error' : ''}"
                        id="name"
                        oninput={(e) => updateName(e.currentTarget.value)}
                        placeholder="Enter tournament name"
                        type="text"
                        value={$tournamentForm.name}
                />
                {#if $tournamentFormErrors.name}
                    <label class="label" for="name">
                        <span class="label-text-alt text-error">{$tournamentFormErrors.name}</span>
                    </label>
                {/if}
            </div>

            <!-- Number of Players -->
            <div class="form-control w-full">
                <label class="label" for="playerCount">
                    <span class="label-text">Number of Players</span>
                </label>
                <input
                        class="input-bordered input w-full {$tournamentFormErrors.playerCount
						? 'input-error'
						: ''}"
                        id="playerCount"
                        min="1"
                        oninput={(e) => updatePlayerCount(parseInt(e.currentTarget.value) || 0)}
                        placeholder="Enter number of players"
                        type="number"
                        value={$tournamentForm.playerCount}
                />
                {#if $tournamentFormErrors.playerCount}
                    <label class="label" for="playerCount">
                        <span class="label-text-alt text-error">{$tournamentFormErrors.playerCount}</span>
                    </label>
                {/if}
            </div>

            <!-- TournamentProvider Description -->
            <div class="form-control w-full sm:col-span-2">
                <label class="label" for="description">
                    <span class="label-text">Description</span>
                </label>
                <textarea
                        class="textarea-bordered textarea w-full {$tournamentFormErrors.description
						? 'textarea-error'
						: ''}"
                        id="description"
                        oninput={(e) => updateDescription(e.currentTarget.value)}
                        placeholder="Enter tournament description"
                        rows="3"
                >{$tournamentForm.description}</textarea>
                {#if $tournamentFormErrors.description}
                    <label class="label" for="description">
                        <span class="label-text-alt text-error">{$tournamentFormErrors.description}</span>
                    </label>
                {/if}
            </div>

            <!-- Start Date -->
            <div class="form-control w-full">
                <label class="label" for="startDate">
                    <span class="label-text">Start Date</span>
                </label>
                <input
                        class="input-bordered input w-full {$tournamentFormErrors.startDate ? 'input-error' : ''}"
                        id="startDate"
                        oninput={(e) => updateStartDate(e.currentTarget.value)}
                        type="date"
                        value={$tournamentForm.startDate}
                />
                {#if $tournamentFormErrors.startDate}
                    <label class="label" for="startDate">
                        <span class="label-text-alt text-error">{$tournamentFormErrors.startDate}</span>
                    </label>
                {/if}
            </div>

            <!-- End Date -->
            <div class="form-control w-full">
                <label class="label" for="endDate">
                    <span class="label-text">End Date</span>
                </label>
                <input
                        class="input-bordered input w-full {$tournamentFormErrors.endDate ? 'input-error' : ''}"
                        id="endDate"
                        oninput={(e) => updateEndDate(e.currentTarget.value)}
                        type="date"
                        value={$tournamentForm.endDate}
                />
                {#if $tournamentFormErrors.endDate}
                    <label class="label" for="endDate">
                        <span class="label-text-alt text-error">{$tournamentFormErrors.endDate}</span>
                    </label>
                {/if}
            </div>
        </div>
    </div>
</div>
