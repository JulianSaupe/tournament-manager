<script lang="ts">
    import type {TournamentFormData, TournamentFormErrors} from '$lib/types/tournament';

    const {formData, errors, onUpdate} = $props<{
        formData: TournamentFormData;
        errors: TournamentFormErrors;
        onUpdate: (field: string, value: string | number | boolean) => void;
    }>();
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
                        value={formData.name}
                        oninput={(e) => onUpdate('name', e.currentTarget.value)}
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
                        value={formData.playerCount}
                        oninput={(e) => onUpdate('playerCount', parseInt(e.currentTarget.value) || 0)}
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
                        value={formData.startDate}
                        oninput={(e) => onUpdate('startDate', e.currentTarget.value)}
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
                        value={formData.endDate}
                        oninput={(e) => onUpdate('endDate', e.currentTarget.value)}
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