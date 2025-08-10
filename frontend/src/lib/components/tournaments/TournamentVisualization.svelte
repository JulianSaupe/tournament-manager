<script lang="ts">
	import { tournamentForm } from '$lib/stores/tournamentForm';

	// Use the store directly with $ syntax for reactive access
	const playerCount = $derived($tournamentForm.playerCount);
</script>

<div class="card bg-base-100 shadow-sm">
	<div class="card-body">
		<div class="mb-4 card-title flex items-center justify-between">
			<h2 class="card-title">Tournament Progression</h2>
			{#if playerCount > 0}
				<div class="badge badge-primary">{playerCount} Players</div>
			{/if}
		</div>

		<div class="mb-6 rounded-lg border-l-4 border-primary bg-base-200 p-4">
			<p class="text-sm">
				<span class="font-medium">How to use:</span> Configure your tournament structure by adding rounds
				and setting the number of players per group and advancing players. The visualization below will
				update automatically to show how players progress through the tournament.
			</p>
		</div>

		{#if playerCount > 0}{:else}
			<div class="alert alert-info shadow-sm">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					class="h-6 w-6 shrink-0 stroke-current"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
					></path>
				</svg>
				<span>Enter the number of players to see the tournament progression visualization.</span>
			</div>
		{/if}
	</div>
</div>
