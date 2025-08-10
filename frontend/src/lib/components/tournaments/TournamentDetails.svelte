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
                        value={$tournamentForm.name}
                        oninput={(e) => updateName(e.currentTarget.value)}
                        class="input-bordered input w-full {$tournamentFormErrors.name ? 'input-error' : ''}"
                        placeholder="Enter tournament name"
                />
                {#if $tournamentFormErrors.name}
                    <label class="label" for="name">
                        <span class="label-text-alt text-error">{$tournamentFormErrors.name}</span>
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
                        value={$tournamentForm.playerCount}
                        oninput={(e) => updatePlayerCount(parseInt(e.currentTarget.value) || 0)}
                        min="1"
                        class="input-bordered input w-full {$tournamentFormErrors.playerCount
						? 'input-error'
						: ''}"
                        placeholder="Enter number of players"
                />
                {#if $tournamentFormErrors.playerCount}
                    <label class="label" for="playerCount">
                        <span class="label-text-alt text-error">{$tournamentFormErrors.playerCount}</span>
                    </label>
                {/if}
            </div>

            <!-- Tournament Description -->
            <div class="form-control w-full col-span-2">
                <label for="description" class="label">
                    <span class="label-text">Description</span>
                </label>
                <textarea
                        id="description"
                        value={$tournamentForm.description}
                        oninput={(e) => updateDescription(e.currentTarget.value)}
                        class="textarea textarea-bordered w-full {$tournamentFormErrors.description ? 'textarea-error' : ''}"
                        placeholder="Enter tournament description"
                        rows="3"
                ></textarea>
                {#if $tournamentFormErrors.description}
                    <label class="label" for="description">
                        <span class="label-text-alt text-error">{$tournamentFormErrors.description}</span>
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
                        value={$tournamentForm.startDate}
                        oninput={(e) => updateStartDate(e.currentTarget.value)}
                        class="input-bordered input w-full {$tournamentFormErrors.startDate
						? 'input-error'
						: ''}"
                />
                {#if $tournamentFormErrors.startDate}
                    <label class="label" for="startDate">
                        <span class="label-text-alt text-error">{$tournamentFormErrors.startDate}</span>
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
                        value={$tournamentForm.endDate}
                        oninput={(e) => updateEndDate(e.currentTarget.value)}
                        class="input-bordered input w-full {$tournamentFormErrors.endDate ? 'input-error' : ''}"
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
