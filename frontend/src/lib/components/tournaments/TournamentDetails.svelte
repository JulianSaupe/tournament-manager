<script lang="ts">
    import type { TournamentFormData, TournamentFormErrors } from '$lib/types/tournament';
    
    // Props
    export let formData: TournamentFormData;
    export let errors: TournamentFormErrors;
    
    // Event dispatcher for form updates
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher<{
        update: { field: string; value: string | number | boolean };
    }>();
    
    // Helper function to update form data
    function updateFormData(field: string, value: string | number | boolean) {
        dispatch('update', { field, value });
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
                        value={formData.name}
                        on:input={(e) => updateFormData('name', e.currentTarget.value)}
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
                        on:input={(e) => updateFormData('playerCount', parseInt(e.currentTarget.value) || 0)}
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
                        on:input={(e) => updateFormData('startDate', e.currentTarget.value)}
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
                        on:input={(e) => updateFormData('endDate', e.currentTarget.value)}
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