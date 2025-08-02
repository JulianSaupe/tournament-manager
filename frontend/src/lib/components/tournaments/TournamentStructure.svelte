<script lang="ts">
    import type { TournamentFormData } from '$lib/types/tournament';
    import { calculateTotalRounds, addRound, removeRound } from '$lib/utils/tournamentUtils';
    import GroupPhaseSettings from './GroupPhaseSettings.svelte';
    import RoundConfiguration from './RoundConfiguration.svelte';
    
    // Props
    export let formData: TournamentFormData;
    
    // Callback props instead of event dispatcher
    export let onUpdate = (field: string, value: string | number | boolean) => {};
    export let onUpdateRounds = (rounds: any[]) => {};
    
    // Helper function to update form data
    function updateFormData(field: string, value: string | number | boolean) {
        onUpdate(field, value);
    }
    
    // Handle group phase settings updates
    function handleGroupPhaseUpdate(event: CustomEvent<{ field: string; value: string | number | boolean }>) {
        updateFormData(event.detail.field, event.detail.value);
    }
    
    // Handle round operations
    function handleAddRound() {
        const updatedFormData = addRound(formData);
        onUpdateRounds(updatedFormData.rounds);
    }
    
    function handleRemoveRound(event: CustomEvent<{ index: number }> | { detail: { index: number } }) {
        const index = 'detail' in event ? event.detail.index : event.detail.index;
        const updatedFormData = removeRound(formData, index);
        onUpdateRounds(updatedFormData.rounds);
    }
    
    function handleUpdateRound(event: CustomEvent<{ index: number; field: string; value: number }> | { detail: { index: number; field: string; value: number } }) {
        const { index, field, value } = 'detail' in event ? event.detail : event.detail;
        const updatedRounds = [...formData.rounds];
        updatedRounds[index] = { ...updatedRounds[index], [field]: value };
        onUpdateRounds(updatedRounds);
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
            onUpdate={(field, value) => handleGroupPhaseUpdate({ detail: { field, value } })} 
        />

        <!-- Rounds Configuration -->
        <RoundConfiguration 
            {formData} 
            onAddRound={handleAddRound}
            onRemoveRound={(index) => handleRemoveRound({ detail: { index } })}
            onUpdateRound={(index, field, value) => handleUpdateRound({ detail: { index, field, value } })}
        />
    </div>
</div>