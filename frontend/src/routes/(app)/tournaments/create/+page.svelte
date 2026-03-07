<script lang="ts">
	import { goto } from '$app/navigation';
	import { enhance } from '$app/forms';
	import {
		resetTournamentForm,
		tournamentForm,
		tournamentFormErrors,
		tournamentFormValid
	} from '$lib/stores/tournamentForm';

	import TournamentDetails from '$lib/components/tournaments/TournamentDetails.svelte';
	import TournamentStructure from '$lib/components/tournaments/TournamentStructure.svelte';

	let isSubmitting = false;
	let submitError = '';

	// Cancel and go back to home page
	function handleCancel() {
		resetTournamentForm();
		goto('/');
	}

	const handleEnhance = (_event: any) => {
		// pending
		isSubmitting = true;
		submitError = '';
		tournamentFormErrors.set({});
		tournamentFormValid.set(true);

		return async (payload: any) => {
			// result handler
			isSubmitting = false;
			const result = payload?.result;
			const update = payload?.update;

			if (result && result.type === 'failure') {
				const data: any = result.data || {};
				if (data.errors) tournamentFormErrors.set(data.errors);
				tournamentFormValid.set(false);
				submitError =
					data.message || data.backendError || 'An error occurred while creating the tournament';
			} else {
				resetTournamentForm();
			}

			if (typeof update === 'function') {
				await update();
			}
		};
	};
</script>

<div class="container mx-auto">
	<h1 class="mb-6 text-2xl">Create New Tournament</h1>

	<form class="space-y-6" method="POST" use:enhance={handleEnhance}>
		<input name="payload" type="hidden" value={JSON.stringify($tournamentForm)} />
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
				<ul class="mt-2 list-inside list-disc">
					{#each Object.entries($tournamentFormErrors) as [field, error]}
						<li class="text-sm">{field}: {error}</li>
					{/each}
				</ul>
			</div>
		{/if}

		<!-- Tournament Details Component -->
		<TournamentDetails />

		<!-- Tournament Structure Component -->
		<TournamentStructure />

		<!-- Form Actions -->
		<div class="mt-8 card-actions justify-end">
			<button class="btn btn-ghost" disabled={isSubmitting} on:click={handleCancel} type="button">
				Cancel
			</button>
			<button
				class="btn btn-primary"
				class:loading={isSubmitting}
				disabled={isSubmitting}
				type="submit"
			>
				{isSubmitting ? 'Creating...' : 'Create Tournament'}
			</button>
		</div>
	</form>
</div>
