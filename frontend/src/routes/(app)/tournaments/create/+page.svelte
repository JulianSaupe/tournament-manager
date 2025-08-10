<script lang="ts">
    import {goto} from '$app/navigation';
    import {
        resetTournamentForm,
        tournamentForm,
        tournamentFormErrors,
        tournamentFormValid,
        validateTournamentForm
    } from '$lib/stores/tournamentForm';

    // Import components
    import TournamentDetails from '$lib/components/tournaments/TournamentDetails.svelte';
    import TournamentStructure from '$lib/components/tournaments/TournamentStructure.svelte';
    // import TournamentVisualization from '$lib/components/tournaments/TournamentVisualization.svelte';

    let isSubmitting = false;
    let submitError = '';

    // Form submission
    async function handleSubmit(event: Event) {
        event.preventDefault();

        // Get current form data
        const currentFormData = $tournamentForm;

        // Validate form data
        const validation = validateTournamentForm(currentFormData);

        // Update error store
        tournamentFormErrors.set(validation.errors);
        tournamentFormValid.set(validation.isValid);

        if (!validation.isValid) {
            submitError = 'Please fix the validation errors before submitting.';
            return;
        }

        // Submit to backend
        try {
            isSubmitting = true;
            submitError = '';

            const response = await fetch('http://localhost:3000/api/tournament', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(currentFormData)
            });

            console.log('Response:', response);

            if (!response.ok) {
                const errorData = await response.text();
                throw new Error(`Failed to create tournament: ${errorData}`);
            }

            const createdTournament = await response.json();
            console.log('Tournament created successfully:', createdTournament);

            // Reset form and navigate to tournaments list
            resetTournamentForm();
            await goto('/');

        } catch (error) {
            console.error('Error creating tournament:', error);
            submitError = error instanceof Error ? error.message : 'An error occurred while creating the tournament';
        } finally {
            isSubmitting = false;
        }
    }

    // Cancel and go back to home page
    function handleCancel() {
        resetTournamentForm();
        goto('/');
    }
</script>

<div class="container mx-auto">
    <h1 class="mb-6 text-2xl">Create New Tournament</h1>

    <form onsubmit={handleSubmit} class="space-y-6">
        <!-- Display form-level error messages -->
        {#if submitError}
            <div class="alert alert-error">
                <span>{submitError}</span>
            </div>
        {/if}

        <!-- Display general validation errors -->
        {#if Object.keys($tournamentFormErrors).length > 0 && !$tournamentFormValid}
            <div class="alert alert-warning">
                <span>Please fix the following errors:</span>
                <ul class="list-disc list-inside mt-2">
                    {#each Object.entries($tournamentFormErrors) as [field, error]}
                        <li class="text-sm">{field}: {error}</li>
                    {/each}
                </ul>
            </div>
        {/if}

        <!-- Tournament Details Component -->
        <TournamentDetails/>

        <!-- Tournament Structure Component -->
        <TournamentStructure/>

        <!-- Tournament Visualization Component -->
        <!-- <TournamentVisualization /> -->

        <!-- Form Actions -->
        <div class="mt-8 card-actions justify-end">
            <button
                    type="button"
                    class="btn btn-ghost"
                    onclick={handleCancel}
                    disabled={isSubmitting}
            >
                Cancel
            </button>
            <button
                    type="submit"
                    class="btn btn-primary"
                    disabled={isSubmitting}
                    class:loading={isSubmitting}
            >
                {isSubmitting ? 'Creating...' : 'Create Tournament'}
            </button>
        </div>
    </form>
</div>
