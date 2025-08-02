<script lang="ts">
    import {goto} from '$app/navigation';
    import {mockTournaments} from '$lib/mockData';
    import type {TournamentFormData, TournamentFormErrors} from '$lib/types/tournament';
    import {updateGroupCounts, validateForm} from '$lib/utils/tournamentUtils';

    // Import components
    import TournamentDetails from '$lib/components/tournaments/TournamentDetails.svelte';
    import TournamentStructure from '$lib/components/tournaments/TournamentStructure.svelte';
    import TournamentVisualization from '$lib/components/tournaments/TournamentVisualization.svelte';

    // Form data using $state rune
    let formData = $state<TournamentFormData>({
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
    });

    // Form validation using $state rune
    let errors = $state<TournamentFormErrors>({
        name: '',
        startDate: '',
        endDate: '',
        playerCount: ''
    });

    let updatedFormData = $derived(updateGroupCounts(formData));

    $effect(() => {
        formData = updatedFormData;
    });

    // Simple callback function for form updates
    function handleFormUpdate(field: string, value: string | number | boolean) {
        (formData as any)[field] = value;

        // Clear error for this field
        if (errors[field as keyof TournamentFormErrors]) {
            errors[field as keyof TournamentFormErrors] = '';
        }
    }

    // Simple callback function for rounds update
    function handleRoundsUpdate(rounds: any[]) {
        formData.rounds = rounds;
    }

    // Form submission
    function handleSubmit(event: Event) {
        event.preventDefault();
        const validation = validateForm(formData);

        if (validation.isValid) {
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
                onUpdate={handleFormUpdate}
        />

        <!-- Tournament Structure Component -->
        <TournamentStructure
                {formData}
                onUpdate={handleFormUpdate}
                onUpdateRounds={handleRoundsUpdate}
        />

        <!-- Tournament Visualization Component -->
        <TournamentVisualization
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