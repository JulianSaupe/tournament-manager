<script lang="ts">
    import type { TournamentFormData } from '$lib/types/tournament';
    
    // Props
    export let formData: TournamentFormData;
    
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

<!-- Group Phase Toggle -->
<div class="form-control bg-base-200 p-4 rounded-lg mb-4">
    <div class="flex items-left mb-1">
        <span class="label-text font-medium text-base">Enable Group Phase</span>
        <input 
            type="checkbox" 
            class="toggle toggle-primary ms-3" 
            checked={formData.groupPhase}
            on:change={(e) => updateFormData('groupPhase', e.currentTarget.checked)}
        />
    </div>
    <p class="text-xs text-base-content/70">
        Organize players into groups for the initial tournament phase
    </p>
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
                    value={formData.groupSize}
                    on:input={(e) => updateFormData('groupSize', parseInt(e.currentTarget.value) || 4)}
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
                        checked={formData.allowPartiallyFilledGroups}
                        on:change={(e) => updateFormData('allowPartiallyFilledGroups', e.currentTarget.checked)}
                    />
                </label>
                <label class="label mt-[-8px]" for="allowPartiallyFilledGroups">
                    <span class="label-text-alt text-base-content/70">If disabled, all groups will have the same number of players</span>
                </label>
            </div>
        </div>
    </div>
{/if}