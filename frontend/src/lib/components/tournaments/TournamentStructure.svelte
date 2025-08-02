<script lang="ts">
    import type { TournamentFormData } from '$lib/types/tournament';
    import { calculateTotalRounds, addRound, removeRound } from '$lib/utils/tournamentUtils';
    import GroupPhaseSettings from './GroupPhaseSettings.svelte';
    import RoundConfiguration from './RoundConfiguration.svelte';
    
    // Props
    export let formData: TournamentFormData;
    
    // Event dispatcher for form updates
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher<{
        update: { field: string; value: string | number | boolean };
        updateRounds: { rounds: any[] };
    }>();
    
    // Helper function to update form data
    function updateFormData(field: string, value: string | number | boolean) {
        dispatch('update', { field, value });
    }
    
    // Handle group phase settings updates
    function handleGroupPhaseUpdate(event: CustomEvent<{ field: string; value: string | number | boolean }>) {
        updateFormData(event.detail.field, event.detail.value);
    }
    
    // Handle round operations
    function handleAddRound() {
        const updatedFormData = addRound(formData);
        dispatch('updateRounds', { rounds: updatedFormData.rounds });
    }
    
    function handleRemoveRound(event: CustomEvent<{ index: number }>) {
        const updatedFormData = removeRound(formData, event.detail.index);
        dispatch('updateRounds', { rounds: updatedFormData.rounds });
    }
    
    function handleUpdateRound(event: CustomEvent<{ index: number; field: string; value: number }>) {
        const { index, field, value } = event.detail;
        const updatedRounds = [...formData.rounds];
        updatedRounds[index] = { ...updatedRounds[index], [field]: value };
        dispatch('updateRounds', { rounds: updatedRounds });
    }
</script>

<div class="card bg-base-100 shadow-sm">
    <div class="card-body">
        <div class="card-title mb-4 flex justify-between items-center">
            <h2 class="card-title">Tournament Structure</h2>
            <span class="badge badge-primary">{calculateTotalRounds(formData)} Rounds</span>
        </div>

        <!-- Group Phase Settings -->
        <GroupPhaseSettings 
            {formData} 
            on:update={handleGroupPhaseUpdate} 
        />

        <!-- Rounds Configuration -->
        <RoundConfiguration 
            {formData} 
            on:addRound={handleAddRound}
            on:removeRound={handleRemoveRound}
            on:updateRound={handleUpdateRound}
        />
    </div>
</div>