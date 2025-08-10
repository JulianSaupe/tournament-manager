<script lang="ts">
	import { goto } from '$app/navigation';

	// Import components
	import TournamentDetails from '$lib/components/tournaments/TournamentDetails.svelte';
	import TournamentStructure from '$lib/components/tournaments/TournamentStructure.svelte';
	// import TournamentVisualization from '$lib/components/tournaments/TournamentVisualization.svelte';

	// Form submission
	function handleSubmit(event: Event) {
		event.preventDefault();
	}

	// Cancel and go back to home page
	function handleCancel() {
		goto('/');
	}
</script>

<div class="container mx-auto">
	<h1 class="mb-6 text-2xl">Create New Tournament</h1>

	<form onsubmit={handleSubmit} class="space-y-6">
		<!-- Tournament Details Component -->
		<TournamentDetails />

		<!-- Tournament Structure Component -->
		<TournamentStructure />

		<!-- Tournament Visualization Component -->
		<!-- <TournamentVisualization /> -->

		<!-- Form Actions -->
		<div class="mt-8 card-actions justify-end">
			<button type="button" class="btn btn-ghost" onclick={handleCancel}> Cancel </button>
			<button type="submit" class="btn btn-primary"> Create Tournament </button>
		</div>
	</form>
</div>
