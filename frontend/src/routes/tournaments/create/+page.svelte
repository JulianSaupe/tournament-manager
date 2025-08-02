<script lang="ts">
    import { goto } from '$app/navigation';
    import { mockTournaments } from '$lib/mockData';
    import type { TournamentFormData, TournamentFormErrors, PhaseVisualizationData } from '$lib/types/tournament';
    import { validateForm, generateVisualizationData, updateGroupCounts } from '$lib/utils/tournamentUtils';
    
    // Import components
    import TournamentDetails from '$lib/components/tournaments/TournamentDetails.svelte';
    import TournamentStructure from '$lib/components/tournaments/TournamentStructure.svelte';
    import TournamentVisualization from '$lib/components/tournaments/TournamentVisualization.svelte';

    // Form data
    let formData: TournamentFormData = {
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
        ]
    };

    // Form validation
    let errors: TournamentFormErrors = {
        name: '',
        startDate: '',
        endDate: '',
        playerCount: ''
    };

    // Visualization data
    let visualizationData: PhaseVisualizationData[] = [];

    // Update form data field
    function handleFormUpdate(event: CustomEvent<{ field: string; value: string | number | boolean }>) {
        const { field, value } = event.detail;
        formData = { ...formData, [field]: value };
    }

    // Update rounds
    function handleRoundsUpdate(event: CustomEvent<{ rounds: any[] }>) {
        formData = { ...formData, rounds: event.detail.rounds };
    }

    // Reactive statements to calculate group counts and update visualization
    $: {
        if (formData.playerCount > 0) {
            formData = updateGroupCounts(formData);
        }
    }

    // Reactive statement to update visualization when form data changes
    $: visualizationData = generateVisualizationData(formData);

    // Form submission
    function handleSubmit(event: Event) {
        event.preventDefault();
        const validation = validateForm(formData);
        
        if (validation.isValid) {
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
        } else {
            errors = validation.errors;
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
        <!-- Tournament Details Component -->
        <TournamentDetails 
            {formData} 
            {errors} 
            on:update={handleFormUpdate} 
        />

        <!-- Tournament Structure Component -->
        <TournamentStructure 
            {formData} 
            on:update={handleFormUpdate}
            on:updateRounds={handleRoundsUpdate}
        />

        <!-- Tournament Visualization Component -->
        <TournamentVisualization 
            {visualizationData} 
            playerCount={formData.playerCount} 
        />

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